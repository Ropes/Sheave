package history

import (
	"container/heap"
	"fmt"
	"testing"
)

func TestInit(t *testing.T) {
	x := []string{"x", "y", "z"}
	a := []string{"a", "b", "c"}
	fmt.Println(x, a)

	hh := make(HistoryHeap, 20)
	heap.Push(&hh, x)
	heap.Push(&hh, a)

	fmt.Printf("hh: %#v\n", hh)
	if hh == nil {
		t.Errorf("First history stored is: %#v", hh)
	}
}
