package sheave

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
