package history

import (
	"container/heap"
	"testing"
)

var x []string
var a []string

func TestInit(t *testing.T) {
	x = []string{"x", "y", "z"}
	a = []string{"a", "b", "c"}
	//fmt.Println(x, a)

	hh := NewHistory(20)
	heap.Push(hh, x)
	if hh.heap[0][0] != "x" {
		t.Errorf("First element incorrect: %#v", hh.heap[0][0])
	}
	heap.Push(hh, a)

	//fmt.Printf("hh: %#v\n", hh)
	if hh == nil {
		t.Errorf("First history stored is: %#v", hh)
	}
	if hh.heap[0][0] != "a" {
		t.Errorf("First element incorrect: %#v", hh.heap[0][0])
	}
}

func TestPop(t *testing.T) {
	x = []string{"x", "y", "z"}
	a = []string{"a", "b", "c"}
	//fmt.Println(x, a)

	hh := NewHistory(2)
	heap.Push(hh, x)
	heap.Push(hh, a)

	poped := heap.Pop(hh).([]string)
	if poped[0] != "a" {
		t.Errorf("Wrong element poped from stack: %#v\n", poped)
	}
	if len(hh.heap) != 1 {
		t.Errorf("Too many items in heap? %#v\n", hh)
	}
}

func TestRound(t *testing.T) {
	x = []string{"x", "y", "z"}
	a = []string{"a", "b", "c"}
	c := []string{"d", "e", "f"}
	g := []string{"j", "k", "l"}

	hh := NewHistory(2)
	heap.Push(hh, x)
	heap.Push(hh, a)
	heap.Push(hh, c)
	heap.Push(hh, g)

	if len(hh.heap) > 2 {
		t.Errorf("Heap grew beyond limit? \n%#v\n", hh)
	}
	//fmt.Println(hh) :D
}

func TestInsertion(t *testing.T) {
	c := []string{"d", "e", "f"}
	//g := []string{"j", "k", "l"}

	hh := NewHistory(20)
	heap.Push(hh, c)
	heap.Push(hh, x)
	heap.Push(hh, a)

	if hh.Hist(0) == nil {
		t.Errorf("Heap Push not adding to correct end of list")
	}
	hh.PrintDump()
}
