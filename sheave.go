package main

import (
	"fmt"

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

	bot.ChannelHistory = make(map[string]history.HistoryHeap)

	ircconfig := parse.ParseConfig("/home/ropes/.config/sheave.conf")
	bot.IRCConnect(ircconfig)
}

func main() {
	bootBot()
}
