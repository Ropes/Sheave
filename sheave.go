package main

import (
	"fmt"

	"github.com/ropes/anagrams"
	"github.com/ropes/sheave/bot"
	"github.com/ropes/sheave/parse"
)

func main() {
	words, err := anagrams.ReadSystemWords()
	if err != nil {
		fmt.Println("No error reading word list")
	}
	anagrammap := anagrams.AnagramList(words)
	bot.AM = &anagrams.AnagramMap{Mapping: anagrammap}

	ircconfig := parse.ParseConfig("~/.config/sheave.json")
	bot.IRCConnect(ircconfig)
}
