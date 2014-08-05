package bot

import (
	"container/heap"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/ropes/sheave/history"
	"github.com/ropes/sheave/parse"
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

/*
//Events contain both hack/talk night Event structs for global access
type Events struct {
	hacknight parse.CalEvent
	talknight parse.CalEvent
}

//var CalEvents Events
*/

func loadCalendar(hackPath, talkPath string) {
	c := make(chan parse.CalEvent)
	go parse.ParseEvent(hackPath, c)
	go parse.ParseEvent(talkPath, c)

	CalEvents.hacknight, CalEvents.talknight = <-c, <-c
}

func TestLocalLoadCalendarParsing(t *testing.T) {

	//bot.LoadCalendar()
	loadCalendar("../hacknights.json", "../talknights.json")
	if CalEvents.talknight.Time.Unix() == -62135596800 || CalEvents.hacknight.Time.Unix() == -62135596800 {
		t.Errorf("Events nil!: %#v", CalEvents)
	}
}

func TestEventResponse(t *testing.T) {
	c := make(chan parse.CalEvent)
	go parse.ParseEvent("testing/resources/talking.json", c)
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
	c := make(chan parse.CalEvent)
	go parse.ParseEvent("testing/resources/talking.json", c)
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

func TestHistoryHeapUsage(t *testing.T) {
	hh1 := history.NewHistory(20)
	x := []string{"x", "y", "z"}
	a := []string{"a", "b", "c"}
	c := []string{"d", "e", "f"}
	g := []string{"j", "k", "l"}
	ChannelHistory := make(map[string]*history.HistoryHeap)

	heap.Init(hh1)

	fail := hh1.Hist(2)
	if len(fail) != 0 {
		t.Errorf("Asserting that empty string list returned from requesting bad history")
	}

	heap.Push(hh1, x)
	heap.Push(hh1, a)
	ChannelHistory["a"] = hh1

	fail = hh1.Hist(2)
	if len(fail) != 0 {
		t.Errorf("Asserting that empty string list returned from requesting bad history")
	}

	if hh1.Len() != 2 {
		t.Errorf("Messages not being pushed correctly onto heap: %#v\n", hh1)
	}

	hh2 := history.NewHistory(20)
	heap.Init(hh2)

	ChannelHistory["b"] = hh2

	d := ChannelHistory["b"]
	heap.Push(d, c)

	if ChannelHistory["b"].Len() != 1 {
		t.Errorf("item was not put onto the heap")
	}

	heap.Push(hh2, g)
	if ChannelHistory["b"].Len() != 2 {
		t.Errorf("Second item not pushed into heap")
	}

}

func TestAnagramParsing(t *testing.T) {
	x := "2!anagram"
	ret := AnagramCmdParse(x)
	if ret[1] != "2" {
		t.Errorf("Number not being parsed correctly")
	}

	x = "!anagram"
	ret = AnagramCmdParse(x)

	x = "!anagram ropes tis flat"
	ret = AnagramCmdParse(x)
	if len(ret) > 2 && ret[1] != "!anagram" {
		t.Errorf("Base command failed to be parsed")
	}
	fmt.Println(ret)

}
