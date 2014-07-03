package main

import (
	"fmt"

	"github.com/mikeclarke/go-irclib"
)

//var ircc = &client.IRCClient{Nickname: "sheave"}

func main() {
	ircc := irc.New("sheave", "sheave")
	ircc.Server = "chat.freenode.net:6667"

	/*
		var ircc = &irc.IRCClient{
			Nickname: "sheave",
			Server:   "chat.freenode.net:6667",
			Username: "sheave",
			//SSL:    true,
			//heartbeatInterval: 1,
		}
	*/
	connErr := ircc.Connect("chat.freenode.net:6667")
	ircc.Join("pdxgo")
	if connErr != nil {
		fmt.Println("Connection Error: \n", connErr)
	}
	ircc.Run()

}
