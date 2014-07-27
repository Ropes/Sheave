package history

import (
	"container/heap"
	"fmt"
	"testing"
)

var x []string
var a []string

func TestInit(t *testing.T) {
	x = []string{"x", "y", "z"}
	a = []string{"a", "b", "c"}
	//fmt.Println(x, a)

	hh := make(HistoryHeap, 20)
	heap.Push(&hh, x)
	if hh[0][0] != "x" {
		t.Errorf("First element incorrect: %#v", hh[0][0])
	}
	heap.Push(&hh, a)

	fmt.Printf("hh: %#v\n", hh)
	if hh == nil {
		t.Errorf("First history stored is: %#v", hh)
	}
	if hh[0][0] != "a" {
		t.Errorf("First element incorrect: %#v", hh[0][0])
	}
}
