package bot

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ropes/anagrams"
	"github.com/ropes/sheave/parse"
	"github.com/thoj/go-ircevent"
)

//Events contain both hack/talk night Event structs for global access
type Events struct {
	hacknight parse.CalEvent
	talknight parse.CalEvent
}

var events Events

var AM *anagrams.AnagramMap

//LoadCalendar reads in the Events from their JSON definitions and
//applies it to global 'events' variable for access
func LoadCalendar() {
	c := make(chan parse.CalEvent)
	go parse.ParseEvent("hacknights.json", c)
	go parse.ParseEvent("talknights.json", c)

	events.hacknight, events.talknight = <-c, <-c
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
			msgs := EventResponse(events.talknight, e.Nick, "Talk Night")
			SendPrivMsgs(con, channel, msgs)
		case "!nexthack":
			msgs := EventResponse(events.hacknight, e.Nick, "Hack Night")
			SendPrivMsgs(con, channel, msgs)
		case "!sheavehelp":
			con.Privmsg(channel, "Sheavebot Cmds: !nextmeetup !nexttalk !nexthack")
		}
	}
}

//AnagramHandler will record previous messages in a pool of strings
//
func AnagramHandler(e *irc.Event, con *irc.Connection) {
	channel := e.Arguments[0]
	if cmd := strings.Trim(e.Arguments[1], " "); cmd[0] == '!' && len(e.Arguments) >= 2 {
		if cmd == "!anagram" {
			x := []string{"stop", "trust"}
			s := AM.AnagramSentence(x)
			con.Privmsg(channel, strings.Join(s, " "))
		}
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
	log := log.New(file, "sheave:", log.Ldate|log.Ltime|log.Lshortfile)

	con.Connect(ircconfig.Server)
	if err != nil {
		panic(fmt.Sprintf("Error connecting to server: %s", ircconfig.Server))
	}
	con.AddCallback("001", func(e *irc.Event) {
		con.Join("#pdxbots")
		con.Join("#pdxgo")
		con.Join("#pdxtech")
	})
	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		GopherHandler(e, con)
	})
	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		AnagramHandler(e, con)
	})
	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		log.Printf("%s %s: %s", e.Arguments[0], e.Nick, e.Arguments[1])
	})
	con.Loop()
}
