package history

import "fmt"

//HistoryHeap stuct represents indexed queue of message strings of limited number.
//New strings replace old ones once limit is reached.
type HistoryHeap struct {
	heap  [][]string
	limit int
}

//Len returns length of heap
func (hh HistoryHeap) Len() int {
	for i := 0; i < len(hh.heap); i++ {
		if hh.heap[i] == nil {
			return i
		}
	}
	return 0 // len(hh.heap)
}

//Less isn't functional, not used since no need to sort
func (hh HistoryHeap) Less(i, j int) bool { return true }

//Swap swaps elements at the indexes given...
func (hh HistoryHeap) Swap(i, j int) { hh.heap[i], hh.heap[j] = hh.heap[j], hh.heap[i] }

//Push adds string to the heap and removes an element if the limit has been reached.
func (hh *HistoryHeap) Push(newString interface{}) {
	fmt.Printf("\nPushing: %#v\n", newString)
	/*
		if len(hh.heap) < hh.limit {
			hh.heap = append(hh.heap, newString.([]string))
		} else {
			hh.Pop()
			hh.heap = append(hh.heap, newString.([]string))
		}
	*/

	for i := len(hh.heap) - 1; i > 0; i-- {
		hh.heap[i] = hh.heap[i-1]
	}
	hh.heap[0] = newString.([]string)
}

func (hh *HistoryHeap) lastIndex() int {
	if endi := len(hh.heap) - 1; hh.heap[endi] != nil {
		return endi
	}
	if hh.heap[0] == nil {
		return 0
	}
	for i := 1; i < len(hh.heap); i++ {
		if hh.heap[i] == nil {
			return i - 1
		}
	}
	return len(hh.heap)
}

//Add string to the heap and removes an element if the limit has been reached.
func (hh *HistoryHeap) Add(newString interface{}) {
	for i := len(hh.heap) - 1; i > 0; i-- {
		hh.heap[i] = hh.heap[i-1]
	}
	hh.heap[0] = newString.([]string)
	//hh.PrintDump()
}

//Pop removes and returns the oldest element in the heap
func (hh *HistoryHeap) Pop() interface{} {
	//n := len(old)
	n := hh.lastIndex()
	x := hh.heap[n]
	hh.heap[n] = nil
	return x
}

func (hh *HistoryHeap) Hist(i int) []string {
	return hh.heap[i]
}

func (hh *HistoryHeap) PrintDump() {
	fmt.Println("HistoryHeap, limit: ", hh.limit)
	for i, v := range hh.heap {
		fmt.Println(i, v)
	}
}

//NewHistory initializes a new HistoryHeap given the limit variable
func NewHistory(limit int) *HistoryHeap {
	hh := HistoryHeap{limit: limit}
	hh.heap = make([][]string, limit)
	return &hh
}
