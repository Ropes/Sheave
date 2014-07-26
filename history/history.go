package history

type HistoryHeap [][]string

func (hh HistoryHeap) Len() int { return len(hh) }

//Note: This is broken currently, need to decide how to handle comparison
func (hh HistoryHeap) Less(i, j int) bool { return true }
func (hh HistoryHeap) Swap(i, j int)      { hh[i], hh[j] = hh[j], hh[i] }

func (hh *HistoryHeap) Push(newString interface{}) {
	*hh = append(*hh, newString.([]string))
}

func (hh *HistoryHeap) Pop() interface{} {
	old := *hh
	n := len(old)
	x := old[n-1]
	*hh = old[0 : n-1]
	return x
}
