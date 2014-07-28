package history

//HistoryHeap stuct represents indexed queue of message strings of limited number.
//New strings replace old ones once limit is reached.
type HistoryHeap struct {
	heap  [][]string
	limit int
}

//Len returns length of heap
func (hh HistoryHeap) Len() int { return len(hh.heap) }

//Less isn't functional, not used since no need to sort
func (hh HistoryHeap) Less(i, j int) bool { return true }

//Swap swaps elements at the indexes given...
func (hh HistoryHeap) Swap(i, j int) { hh.heap[i], hh.heap[j] = hh.heap[j], hh.heap[i] }

//Push adds string to the heap and removes an element if the limit has been reached.
func (hh *HistoryHeap) Push(newString interface{}) {
	if len(hh.heap) < hh.limit {
		hh.heap = append(hh.heap, newString.([]string))
	} else {
		hh.Pop()
		hh.heap = append(hh.heap, newString.([]string))
	}
}

//Pop removes and returns the oldest element in the heap
func (hh *HistoryHeap) Pop() interface{} {
	old := hh.heap
	n := len(old)
	x := old[n-1]
	hh.heap = old[0 : n-1]
	return x
}

//NewHistory initializes a new HistoryHeap given the limit variable
func NewHistory(limit int) *HistoryHeap {
	hh := HistoryHeap{limit: limit}
	hh.heap = make([][]string, limit)
	return &hh
}
