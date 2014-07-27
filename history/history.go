package history

type HistoryHeap struct {
	heap  [][]string
	limit int
}

func (hh HistoryHeap) Len() int { return len(hh.heap) }

//Note: This is broken currently, need to decide how to handle comparison
func (hh HistoryHeap) Less(i, j int) bool { return true }
func (hh HistoryHeap) Swap(i, j int)      { hh.heap[i], hh.heap[j] = hh.heap[j], hh.heap[i] }

func (hh *HistoryHeap) Push(newString interface{}) {
	if len(hh.heap) < hh.limit {
		hh.heap = append(hh.heap, newString.([]string))
	} else {
		hh.Pop()
		hh.heap = append(hh.heap, newString.([]string))
	}
}

func (hh *HistoryHeap) Pop() interface{} {
	old := hh.heap
	n := len(old)
	x := old[n-1]
	hh.heap = old[0 : n-1]
	return x
}

func NewHistory(limit int) *HistoryHeap {
	hh := HistoryHeap{limit: limit}
	hh.heap = make([][]string, limit)
	return &hh
}
