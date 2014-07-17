package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/mikeclarke/go-irclib"
)

type Events struct {
	hacknight Event
	talknight Event
}

var events Events

func LoadCalendar() {
	c := make(chan Event)
	go parseEvent("hacknights.json", c)
	go parseEvent("talknights.json", c)

	events.hacknight, events.talknight = <-c, <-c
}

//EventResponse creates a []string of useful information of an Event(struct)
//which will be sent to inquiring user.
func EventResponse(e Event, user string, etype string) []string {
	resp := make([]string, 0)
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
func SendPrivMsgs(event *irc.Event, channel string, msgs []string) {
	client := event.Client
	for _, msg := range msgs {
		client.Privmsg(channel, msg)
	}
}

//GopherHandler which responds with the next meeting type for the !nextmeetup command
func GopherHandler(event *irc.Event) {
	client := event.Client
	channel := event.Arguments[0]

	if channel == "#pdxgo" || channel == "#pdxgotest" || channel == "#pdxbots" {
		if len(event.Arguments) >= 2 {
			LoadCalendar()
			cmd := strings.Trim(event.Arguments[1], " ")
			log.Printf("Event: %#v\n", event)
			user := PrivMsgUser(event)
			//log.Printf("Message:'%#v'\n", event.Arguments)
			log.Printf("Channel: %+v", channel)
			switch cmd {
			case "!nextmeetup":
				msg := fmt.Sprintf("%s: %s", user, " TODO: meetingtime!")
				client.Privmsg(channel, msg)
			case "!nexttalk":
				msgs := EventResponse(events.talknight, user, "Talk Night")
				SendPrivMsgs(event, channel, msgs)
			case "!nexthack":
				msgs := EventResponse(events.hacknight, user, "Hack Night")
				SendPrivMsgs(event, channel, msgs)
			case "!sheavehelp":
				client.Privmsg(channel, "Sheavebot Cmds: !nextmeetup !nexttalk !nexthack")
			}
		}
	}
}

//PrivMsgUser returns the user's name who sent a message via the Event object
func PrivMsgUser(event *irc.Event) string {
	prefix := event.Prefix
	split := strings.Split(prefix, "!")
	return split[0]
}

//IRCConnect initializes and runs the irc connection and adds the GopherHandler to its event loop for parsing messages
func IRCConnect(ircconfig IRCConfig) {
	ircc := irc.New(ircconfig.UserName, ircconfig.UserName)
	ircc.Server = ircconfig.Server
	ircc.RealName = ircconfig.RealName
	ircc.Password = ircconfig.Passwd

	connErr := ircc.Connect(ircc.Server)
	if connErr != nil {
		fmt.Println("Connection Error: \n", connErr)
	}

	for _, v := range ircconfig.Channels {
		fmt.Println("Joining: ", v)
		ircc.Join(v)
	}
	ircc.AddHandler(GopherHandler)
	ircc.Run()
}

func main() {
	ircconfig := parseConfig("conf.json")
	IRCConnect(ircconfig)
}
