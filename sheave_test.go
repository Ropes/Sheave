package sheave

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestLocalEventParsing(t *testing.T) {
	LoadCalendar()
	if events.talknight == nil || events.hacknight == nil {
		t.Errorf("Events nil!: %#v", events)
	}
}

type Talk struct {
	Localtime   string
	Location    string
	Lat         float64
	Topic       string
	Description []string
}

func parseCal(path string, e chan Talk) {
	contents, err := ioutil.ReadFile(path)
	var cal Talk
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

func TestTalkTargetEventParsing(t *testing.T) {

	c := make(chan Talk)
	go parseCal("testing/resources/talking.json", c)
	jsn := <-c
	if jsn.Location != "ESRI" {
		t.Errorf("Location incorrect: %#v", jsn)
	}
	if jsn.Lat != 45.521525 {
		t.Errorf("Latitude incorrect(45.521525): %#v", jsn)
	}
}
