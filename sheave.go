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
	confPath := flag.String("confPath",
		"/home/ropes/.config/sheave.conf",
		"Path to config file for Sheave bot.")
	dictPath := flag.String("dictPath",
		"/usr/share/dict/words",
		"Path to dictionary word list file")
	flag.Parse()

	words, err := anagrams.ReadCustomWords(*dictPath)
	if err != nil {
		fmt.Println("No error reading word list")
	}
	anagrammap := anagrams.AnagramList(words)
	bot.AM = &anagrams.AnagramMap{Mapping: anagrammap}

	bot.ChannelHistory = make(map[string]*history.HistoryHeap)

	ircconfig := parse.ParseConfig(*confPath)
	if ircconfig.UserName == "" || ircconfig.Passwd == "" {
		log.Fatal("IRC Config failed to parse important information.")
	}
	bot.IRCConnect(ircconfig)
}

func main() {
	bootBot()
}
