package main

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

func TestLocalLoadCalendarParsing(t *testing.T) {
	LoadCalendar()
	//fmt.Println(events.hacknight.Time.Unix())
	if events.talknight.Time.Unix() == -62135596800 || events.hacknight.Time.Unix() == -62135596800 {
		t.Errorf("Events nil!: %#v", events)
	}
}

func TestEventResponse(t *testing.T) {
	c := make(chan CalEvent)
	go parseEvent("testing/resources/talking.json", c)
	jsn := <-c

	out := EventResponse(jsn, "ropes", "Talk Night")
	fmt.Printf("%#v\n", out)
	if len(out) != 3 {
		t.Errorf("Three strings not returned!")
	}
	for i, v := range out {
		if v == "" {
			t.Errorf("String %d empty!: '%s'", i, v)
		}
	}
}

func TestTalkTargetEventParsing(t *testing.T) {
	c := make(chan CalEvent)
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
