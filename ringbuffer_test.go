package main_test

import (
	"fmt"
	"sync"
	"testing"
)
type beBuffer struct {
	value int
}

type RingBuffer struct {
	buf    []*beBuffer
	size   int
	r      int // next position to read
	w      int // next position to write
	isFull bool
	mu     sync.Mutex
}

func (r *RingBuffer) New(size int) *RingBuffer {
	r.buf = make([]*beBuffer, size)
	//for i := 0;i<size;i++{
	//	r.buf[i] = &beBuffer{}
	//}
	//r.w = size
	r.size = size
	//r.isFull = true
	return r
}

func (r *RingBuffer) Read() (buffer *beBuffer) {
	r.mu.Lock()
	buffer = r.read()
	r.mu.Unlock()
	return
}

func (r *RingBuffer) read() *beBuffer {
	if r.w == r.r && !r.isFull {
		return r.add()
	}

	if r.w > r.r {
		bb := r.buf[r.r]
		r.r = (r.r + 1) % r.size
		return bb
	}

	var bb *beBuffer
	if r.r < r.size {
		bb = r.buf[r.r]
	} else {
		bb = r.buf[0]
	}
	r.r = (r.r + 1) % r.size

	r.isFull = false
	return bb
}

func (r *RingBuffer) add() (buffer *beBuffer) {
	r.write(&beBuffer{})
	buffer = r.read()
	return
}


func (r *RingBuffer) Write(p *beBuffer) {
	r.mu.Lock()
	r.write(p)
	r.mu.Unlock()
}

func (r *RingBuffer) write(p *beBuffer)  {
	if r.isFull {
		return
		// 扩容
	}

	if r.w >= r.r {
		c1 := r.size - r.w
		if c1 >= 1 {
			r.buf[r.w] = p
			r.w += 1
		} else {
			r.buf[0] = p
			r.w = 1
		}
	} else {
		r.buf[r.w] = p
		r.w += 1
	}

	if r.w == r.size {
		r.w = 0
	}
	if r.w == r.r {
		r.isFull = true
	}
}


func TestRingBuffer(t *testing.T){
	ring := (&RingBuffer{}).New(3)
	fmt.Println(ring)
	fmt.Println(ring.Read())
	fmt.Println(ring.Read())
	fmt.Println(ring)
	fmt.Println(ring.Read())
	fmt.Println(ring)
	fmt.Println(ring.Read())
	fmt.Println(ring)

}