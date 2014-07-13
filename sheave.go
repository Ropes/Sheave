package sheave

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/mikeclarke/go-irclib"
)

type Events struct {
	hacknight interface{}
	talknight interface{}
}

var events Events

func LoadCalendar() {
	c := make(chan interface{})
	go parseCalendar("hacknights.json", c)
	go parseCalendar("talknights.json", c)

	events.hacknight, events.talknight = <-c, <-c
}

//GopherHandler which responds with the next meeting type for the !nextmeetup command
func GopherHandler(event *irc.Event) {
	client := event.Client
	channel := event.Arguments[0]

	if channel == "#pdxgo" || channel == "#pdxgotest" {
		if len(event.Arguments) >= 2 {
			cmd := strings.Trim(event.Arguments[1], " ")
			log.Printf("Event: %#v\n", event)
			user := PrivMsgUser(event)
			//log.Printf("Message:'%#v'\n", event.Arguments)
			switch cmd {
			case "!nextmeetup":
				log.Printf("Channel: %+v", channel)
				msg := fmt.Sprintf("%s: %s", user, "meetingtime!")
				client.Privmsg(channel, msg)
			case "!nexttalk":
				log.Printf("Channel: %+v", channel)
				msg := fmt.Sprintf("%s: Next Talk night: %s", user)
				client.Privmsg(channel, msg)
			case "!nexthack":
				log.Printf("Channel: %+v", channel)
				msg := fmt.Sprintf("%s: %s", user, "meetingtime!")
				client.Privmsg(channel, msg)
			case "!sheavehelp":
				client.Privmsg(channel, "Help msg")

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

type IRCConfig struct {
	Server   string
	UserName string
	RealName string
	Passwd   string
	Channels []string
}

func parseConfig(path string) IRCConfig {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	var conf IRCConfig
	err = json.Unmarshal(contents, &conf)
	if err != nil {
		fmt.Println(err)
	}
	return conf
}

type Hacknights struct {
	data map[string]interface{}
}

func parseCalendar(path string, e chan interface{}) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		e <- nil
	}
	var cal interface{}
	err = json.Unmarshal(contents, &cal)
	if err != nil {
		fmt.Println(err)
		e <- nil
	}
	e <- cal
}

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
	/*
		hacknight, err := parseCalendar("hacknights.json")
		if err != nil {
			panic(fmt.Sprintf("Error parsing calendar: %v", err))
		}
		talknight, err := parseCalendar("talknights.json")
		if err != nil {
			panic(fmt.Sprintf("Error parsing calendar: %v", err))
		}
	*/
	//ircconfig := parseConfig("conf.json")
	LoadCalendar()
	fmt.Printf("events:\n%#v\n", events.talknight)
	//IRCConnect(ircconfig)
}
