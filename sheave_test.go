package sheave

import "testing"

func TestLocalEventParsing(t *testing.T) {
	LoadCalendar()
	if events.talknight == nil || events.hacknight == nil {
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
}
