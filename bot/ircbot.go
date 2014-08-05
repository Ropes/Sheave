package bot

import (
	"container/heap"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/ropes/anagrams"
	"github.com/ropes/sheave/history"
	"github.com/ropes/sheave/parse"
	"github.com/thoj/go-ircevent"
)

//Events contain both hack/talk night Event structs for global access
type Events struct {
	hacknight parse.CalEvent
	talknight parse.CalEvent
}

//CalEvents variable holds talk/hack night calendar/event information
var CalEvents Events

//AM holds anagram maps keyed off of their characters sorted
var AM *anagrams.AnagramMap

//LoadCalendar reads in the Events from their JSON definitions and
//applies it to global 'events' variable for access
func LoadCalendar() {
	c := make(chan parse.CalEvent)
	go parse.ParseEvent("hacknights.json", c)
	go parse.ParseEvent("talknights.json", c)

	CalEvents.hacknight, CalEvents.talknight = <-c, <-c
}

//EventResponse creates a []string of useful information of an Event(struct)
//which will be sent to inquiring user.
func EventResponse(e parse.CalEvent, user string, etype string) []string {
	var resp []string
	msg := fmt.Sprintf("%s: Next %s: %s", user, etype, e.Localtime)
	resp = append(resp, msg)

	msg = fmt.Sprintf(">>> %s @ %s <<<", e.Topic, e.Location)
	resp = append(resp, msg)

	msg = fmt.Sprintf("Info: %s", e.Link)
	resp = append(resp, msg)
	return resp
}

//SendPrivMsgs broadcasts PRIVMSGs via the given client and channel.
//msgs []string; are the messages to be sent
func SendPrivMsgs(con *irc.Connection, channel string, msgs []string) {
	for _, msg := range msgs {
		con.Privmsg(channel, msg)
	}
}

//GopherHandler which responds with the next meeting type for the !nextmeetup command
func GopherHandler(e *irc.Event, con *irc.Connection) {
	channel := e.Arguments[0]

	if cmd := strings.Trim(e.Arguments[1], " "); cmd[0] == '!' && len(e.Arguments) >= 2 {
		LoadCalendar()
		switch cmd {
		case "!nextmeetup":
			msg := fmt.Sprintf("%s: %s", e.Nick, " TODO: meetingtime!")
			con.Privmsg(channel, msg)
		case "!nexttalk":
			msgs := EventResponse(CalEvents.talknight, e.Nick, "Talk Night")
			SendPrivMsgs(con, channel, msgs)
		case "!nexthack":
			msgs := EventResponse(CalEvents.hacknight, e.Nick, "Hack Night")
			SendPrivMsgs(con, channel, msgs)
		case "!sheavehelp":
			con.Privmsg(channel, "Sheavebot Cmds: !nextmeetup !nexttalk !nexthack")
		}
	}
}

//ChannelHistory is a map of channel histories recording recent user messages
var ChannelHistory map[string]*history.HistoryHeap

//ChannelHistorian takes channel messages and records them in a HistoryHeap to
//save them for potential anagramming!
func ChannelHistorian(e *irc.Event) {
	channel := e.Arguments[0]

	msg := strings.Split(strings.Trim(e.Arguments[1], " "), " ")

	hh := ChannelHistory[channel]
	//fmt.Printf("Msg received: %#v\n", msg)
	//fmt.Printf("ChanHist: %#v %#v\n", channel, hh)

	if msg != nil {
		heap.Push(hh, msg)
		hh.PrintDump()
		//hh.Add(msg)
	}
}

//parseMsg breaks apart a privmsg and returns a list of strings to be anagramed
func parseMsg(s string) []string {
	return strings.Split(strings.Trim(s, " "), " ")
}

func AnagramCmdParse(trimmed string) []string {
	re := regexp.MustCompile("([0-9]*)*([!]+)anagram[s]*.*")
	sentence := regexp.MustCompile(".*!anagram[s]* [(.+) ]*")

	if cmd := re.FindStringSubmatch(trimmed); len(cmd) == 3 {
		words := sentence.FindStringSubmatch(trimmed)
		if len(words) > 0 {
			return words
		}
		return cmd
	}
	return make([]string, 0)
}

//AnagramResponder returns previous user messages with text anagramed
func AnagramResponder(e *irc.Event, con *irc.Connection, logger *log.Logger) {
	channel := e.Arguments[0]
	trimmed := strings.Trim(e.Arguments[1], " ")
	re := regexp.MustCompile("([0-9]*)*([!]+)anagram[s]*")

	if cmd := re.FindStringSubmatch(trimmed); len(cmd) == 3 {
		//Parse command
		back := -1
		if cmd[2] != "" {
			back += len(cmd[2])
		}
		if cmd[1] != "" {
			b, _ := strconv.Atoi(cmd[1])
			back += b
		}

		hh := ChannelHistory[channel]
		//Pull previous msg from history
		if back >= hh.Len() && back < 20 {
			con.Privmsg(channel, "Sorry, sheave doesn't have history recorded that far back :(")
			return
		} else if back >= hh.Len() && back > 20 {
			con.Privmsg(channel, "sheave is limited to 20 previous messages per channel")
			return
		}

		logger.Printf("Anagraming %d back in channel: %#v\n%#v\n", back, channel, hh)
		x := hh.Hist(back)

		s := AM.AnagramSentence(x)
		logger.Printf("Anagraming: %#v -> %#v\n", x, s)

		con.Privmsg(channel, "Anagramed: "+strings.Join(s, " "))
	}
}

//IRCConnect initializes and runs the irc connection and adds the GopherHandler to its event loop for parsing messages
func IRCConnect(ircconfig parse.IRCConfig) {
	con := irc.IRC(ircconfig.UserName, ircconfig.UserName)
	con.Password = ircconfig.Passwd

	file, err := os.OpenFile("irc.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("Failed to open log file:%s", err))
	}
	irclog := log.New(file, "sheave:", log.Ldate|log.Ltime|log.Lshortfile)

	anagramfile, err := os.OpenFile("anagraming.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("Failed to open anagram log file: %s", err))
	}
	anagramlog := log.New(anagramfile, "sheave:", log.Ldate|log.Ltime|log.Lshortfile)

	con.Connect(ircconfig.Server)
	if err != nil {
		panic(fmt.Sprintf("Error connecting to server: %s", ircconfig.Server))
	}
	con.AddCallback("001", func(e *irc.Event) {

		for _, v := range ircconfig.Channels {
			irclog.Printf("Initializing: %v\n", v)

			hist := history.NewHistory(20)
			heap.Init(hist)
			ChannelHistory[v] = hist

			con.Join(v)
		}
	})
	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		GopherHandler(e, con)
	})
	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		AnagramResponder(e, con, anagramlog)
	})
	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		ChannelHistorian(e)
	})
	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		irclog.Printf("%s %s: %s", e.Arguments[0], e.Nick, e.Arguments[1])
	})
	con.Loop()
}
