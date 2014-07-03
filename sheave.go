package main

import (
	"fmt"

	"github.com/mikeclarke/go-irclib"
)

//var ircc = &client.IRCClient{Nickname: "sheave"}

func GopherHandler(event *irc.Event) {
	client := event.Client

	switch event.Command {
	case "!nextmeetup":
		fmt.Println("EventAgs: ", event.Arguments)
		client.SendRawf("Hack Night 2014/7/9",
			event.Arguments[len(event.Arguments)-1])
	}
}

func main() {
	ircc := irc.New("sheave", "sheave")
	ircc.Server = "chat.freenode.net:6667"
	ircc.RealName = "PDX Gopher"

	/*
		var ircc = &irc.IRCClient{
			Nickname: "sheave",
			Server:   "chat.freenode.net:6667",
			Username: "sheave",
			//SSL:    true,
			//heartbeatInterval: 1,
		}
	*/
	fmt.Println(ircc)
	connErr := ircc.Connect("chat.freenode.net:6667")
	if connErr != nil {
		fmt.Println("Connection Error: \n", connErr)
	}

	ircc.Join("#pdxgo")
	ircc.AddHandler(GopherHandler)
	ircc.Run()

}
