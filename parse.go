package sheave

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

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

type Event struct {
	Time        time.Time
	Localtime   string
	Location    string
	Lat         float64
	Lon         float64
	Link        string
	Tweet       string
	Topic       string
	Description []string
	Notes       string
}

func parseEvent(path string, e chan Event) {
	contents, err := ioutil.ReadFile(path)
	var cal Event
	if err != nil {
		fmt.Println(err)
		e <- cal
	}
	err = json.Unmarshal(contents, &cal)
	if err != nil {
		fmt.Println(err)
		e <- cal
	}
	e <- cal
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
