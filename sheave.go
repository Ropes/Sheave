package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ropes/anagrams"
	"github.com/thoj/go-ircevent"
)

//Events contain both hack/talk night Event structs for global access
type Events struct {
	hacknight CalEvent
	talknight CalEvent
}

var events Events

//LoadCalendar reads in the Events from their JSON definitions and
//applies it to global 'events' variable for access
func LoadCalendar() {
	c := make(chan CalEvent)
	go parseEvent("hacknights.json", c)
	go parseEvent("talknights.json", c)

	events.hacknight, events.talknight = <-c, <-c
}

//EventResponse creates a []string of useful information of an Event(struct)
//which will be sent to inquiring user.
func EventResponse(e CalEvent, user string, etype string) []string {
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

	LoadCalendar()
	cmd := strings.Trim(e.Arguments[1], " ")
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

/*
//PrivMsgUser returns the user's name who sent a message via the Event object
func PrivMsgUser(event *irc.Event) string {
	prefix := event.Prefix
	split := strings.Split(prefix, "!")
	return split[0]
}
*/

//IRCConnect initializes and runs the irc connection and adds the GopherHandler to its event loop for parsing messages
func IRCConnect(ircconfig IRCConfig) {
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
	})
	//con.AddCallback("JOIN", func(e *irc.Event) { con.Privmsg("#pdxbots", "hihi!") })
	con.AddCallback("PRIVMSG", func(e *irc.Event) {
		log.Printf("%s %s: %s", e.Arguments[0], e.Nick, e.Arguments[1])
	})
	con.Loop()
}

func main() {
	words, err := anagrams.ReadSystemWords()
	if err != nil {
		fmt.Println("No error reading word list")
	}
	anagrammap := anagrams.AnagramList(words)
	AM := &anagrams.AnagramMap{Mapping: anagrammap}

	word := "god"
	ana := AM.AnagramOfWord(word)
	fmt.Println(ana)

	ircconfig := parseConfig("conf.json")
	IRCConnect(ircconfig)
}
