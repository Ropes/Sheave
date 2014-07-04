package main

import (
	"fmt"
	"strings"

	"github.com/mikeclarke/go-irclib"
)

//var ircc = &client.IRCClient{Nickname: "sheave"}

/*
&{sheave  PDX Gopher sheave       0  chat.freenode.net:6667 false 0xc2100000a0 false  0 0xc21006e000 0xc210038120 0xc210038180 0xc2100381e0 0xc210038240 0xc21005f100 0xc21001d510 {63540025767 787490388 0x7be3a0}  false <nil>}
&{:ropes!~ropes@107.170.244.234 PRIVMSG #pdxgo :schmichael: is there a Go equivalent to pprint.pformat from Python? ropes!~ropes@107.170.244.234 PRIVMSG [#pdxgo schmichael: is there a Go equivalent to pprint.pformat from Python?] 0xc21005b000 false}
*/

func GopherHandler(event *irc.Event) {
	client := event.Client

	fmt.Println(event)
	fmt.Println("Event:Arguments:\n", event.Arguments)
	channel := event.Arguments[0]
	fmt.Println("args:\n", channel)
	//Array: [#pdxgo ryana: I'll have to check pretty out]
	fmt.Println("\n")

	if channel == "#pdxgo" || channel == "#pdxgotest" {
		if len(event.Arguments) >= 2 {
			cmd := strings.Trim(event.Arguments[1], " ")
			fmt.Printf("CMD:'%s'\n", cmd)
			msg := "Hack Night 2014/7/9"
			switch cmd {
			case "!nextmeetup":
				/*client.SendRawf("%s  Hack Night 2014/7/9",
				event.Prefix,
				event.Arguments[len(event.Arguments)-1])*/
				client.Privmsg("#pdxgotest", msg)
			}
		}
	}
}

func main() {
	ircc := irc.New("sheave", "sheave")
	ircc.Server = "chat.freenode.net:6667"
	ircc.RealName = "PDX Gopher"

	fmt.Println(ircc)
	connErr := ircc.Connect("chat.freenode.net:6667")
	if connErr != nil {
		fmt.Println("Connection Error: \n", connErr)
	}
	fmt.Println(ircc.GetNick())
	ircc.Join("#pdxgotest")
	ircc.AddHandler(GopherHandler)
	ircc.Run()

}
