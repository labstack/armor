// https://www.npmjs.com/package/fixed-array

package cube

import (
	"sync"
)

type (
	Window struct {
		size   int
		i      int
		values []int64
		mutex  sync.RWMutex
	}
)

func NewWindow(size int) *Window {
	return &Window{
		size:   size,
		values: make([]int64, size),
	}
}

func (w Window) Push(v int64) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	w.values[w.i] = v
	w.i++
	if w.i == w.size {
		// Reset
		w.i = 0
	}
}

func (w *Window) Sum() int64 {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	var t int64
	for _, v := range w.values {
		t += v
	}
	return t
}

func (w *Window) Mean() int64 {
	if len(w.values) == 0 {
		return 0
	}
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.Sum() / int64(len(w.values))
}
