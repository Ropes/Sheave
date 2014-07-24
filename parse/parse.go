package parse

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

func ParseConfig(path string) IRCConfig {
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

type CalEvent struct {
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

func ParseEvent(path string, e chan CalEvent) {
	contents, err := ioutil.ReadFile(path)
	var cal CalEvent
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

func ParseCalendar(path string, e chan interface{}) {
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
