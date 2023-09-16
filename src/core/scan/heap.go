package scan

import (
	"container/heap"
	"encoding/json"
	"os"
	"sort"
	"time"
)

type RecentFileItem struct {
	Index
	ModTime string    `json:"-"`
	modTime time.Time `json:"-"`
}

func NewRecentFile(size int) *RecentFile {
	return &RecentFile{
		MaxSize: size,
		items:   []RecentFileItem{},
	}
}

type RecentFile struct {
	MaxSize int
	items   []RecentFileItem
	init    bool
}

func (rf *RecentFile) Len() int {
	return len(rf.items)
}

func (rf *RecentFile) Less(i, j int) bool {
	return rf.items[i].modTime.Before(rf.items[j].modTime)
}

func (rf *RecentFile) Swap(i, j int) {
	rf.items[i], rf.items[j] = rf.items[j], rf.items[i]
}

func (rf *RecentFile) Push(x interface{}) {
	rf.items = append(rf.items, x.(RecentFileItem))
}

func (rf *RecentFile) Pop() interface{} {
	old := rf.items
	n := len(old)
	x := old[n-1]
	rf.items = old[0 : n-1]
	return x
}

func (rf *RecentFile) PushItem(item RecentFileItem) {
	item.modTime, _ = time.Parse(time.DateTime, item.ModTime)
	if rf.init {
		heap.Push(rf, item)
		heap.Pop(rf)
		return
	}

	rf.items = append(rf.items, item)
	if len(rf.items) < rf.MaxSize {
		return
	}

	heap.Init(rf)
	rf.init = true
}

func (rf *RecentFile) WriteTo(filename string) error {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}

	defer f.Close()

	sort.Slice(rf.items, func(i, j int) bool {
		return rf.items[i].modTime.After(rf.items[j].modTime)
	})

	for _, v := range rf.items {
		b, err := json.Marshal(v)
		if err != nil {
			return err
		}
		if _, err := f.Write(b); err != nil {
			return err
		}
		if _, err := f.Write([]byte{'\n'}); err != nil {
			return err
		}
	}

	return nil
}
