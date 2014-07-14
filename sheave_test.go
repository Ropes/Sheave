package main

import "testing"

func TestLocalLoadCalendarParsing(t *testing.T) {
	LoadCalendar()
	//fmt.Println(events.hacknight.Time.Unix())
	if events.talknight.Time.Unix() == -62135596800 || events.hacknight.Time.Unix() == -62135596800 {
		t.Errorf("Events nil!: %#v", events)
	}
}

func TestTalkTargetEventParsing(t *testing.T) {
	c := make(chan Event)
	go parseEvent("testing/resources/talking.json", c)
	jsn := <-c
	//fmt.Printf("%#v\n", jsn)
	if jsn.Location != "ESRI" {
		t.Errorf("Location incorrect: %#v", jsn)
	}
	if jsn.Lat != 45.521525 {
		t.Errorf("Latitude incorrect(45.521525): %#v", jsn)
	}
	if jsn.Localtime != "Tuesday" {
		t.Errorf("Localtime incorrect(Tuesday): %#v", jsn)
	}
}
