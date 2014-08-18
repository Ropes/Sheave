package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/ropes/anagrams"
	"github.com/ropes/sheave/bot"
	"github.com/ropes/sheave/history"
	"github.com/ropes/sheave/parse"
)

func bootBot() {
	words, err := anagrams.ReadSystemWords()
	if err != nil {
		fmt.Println("No error reading word list")
	}
	anagrammap := anagrams.AnagramList(words)
	bot.AM = &anagrams.AnagramMap{Mapping: anagrammap}

	bot.ChannelHistory = make(map[string]*history.HistoryHeap)

	confPath := flag.String("confPath",
		"/home/ropes/.config/sheave.conf",
		"Path to config file for Sheave bot.")
	flag.Parse()

	ircconfig := parse.ParseConfig(*confPath)
	if ircconfig.UserName == "" || ircconfig.Passwd == "" {
		log.Fatal("IRC Config failed to parse important information.")
	}
	fmt.Printf("%#v\n", ircconfig)
	bot.IRCConnect(ircconfig)
}

func main() {
	bootBot()
}
