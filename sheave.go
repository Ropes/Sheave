package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/mikeclarke/go-irclib"
)

func GopherHandler(event *irc.Event) {
	client := event.Client

	fmt.Println(event)
	channel := event.Arguments[0]
	fmt.Printf("args: %+v", channel)
	//Array: [#pdxgo ryana: I'll have to check pretty out]

	if channel == "#pdxgo" || channel == "#pdxgotest" {
		if len(event.Arguments) >= 2 {
			cmd := strings.Trim(event.Arguments[1], " ")
			fmt.Printf("CMD:'%s'\n", cmd)
			msg := "Hack Night 2014/7/9"
			switch cmd {
			case "!nextmeetup":
				client.Privmsg("#pdxgotest", msg)
			}
		}
	}
}

type Config struct {
	Server   string
	UserName string
	RealName string
	Passwd   string
	Channels []string
}

func parseConfig(path string) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	var conf Config
	err = json.Unmarshal(contents, &conf)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%#v\n", conf)
}

func main() {
	parseConfig("conf.json")

	/*
		ircc := irc.New("sheave", "sheave")
		ircc.Server = "chat.freenode.net:6667"
		ircc.RealName = "PDX Gopher"

		fmt.Println(ircc)
		connErr := ircc.Connect("chat.freenode.net:6667")
		if connErr != nil {
			fmt.Println("Connection Error: \n", connErr)
		}

		fmt.Printf("Nick: %+v\n", ircc.GetNick())
		ircc.Join("#pdxgotest")
		ircc.AddHandler(GopherHandler)
		ircc.Run()

	*/
}
