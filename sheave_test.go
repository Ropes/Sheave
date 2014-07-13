package sheave

import (
	"fmt"
	"testing"
)

func TestLocalEventParsing(t *testing.T) {
	LoadCalendar()
	if events.talknight == nil || events.hacknight == nil {
		t.Errorf("Events nil!: %#v", events)
	}
}

type Talk struct {
	location    string
	lat         float64
	Description []string
}

func TestTargetEventParsing(t *testing.T) {

	c := make(chan interface{})
	go parseCalendar("testing/resources/talking.json", c)
	jsn := <-c
	fmt.Println(jsn)
	/*
		if jsn.location != "ESRI" {
			t.Errorf("Location incorrect: %#v", jsn)
		}
	*/
}
