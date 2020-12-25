package main_test

import (
	_ "runtime/pprof"
	"sync"
	"testing"
)

var bytePool = sync.Pool{
	New: func() interface{} {
		b := make([]byte, 1024)
		return &b
	},
}

var customPool = make(chan []byte, 1024)

func BenchmarkPA1(b *testing.B){
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			obj := (bytePool.Get().(*[]byte))
			_ = len(*obj)
			bytePool.Put(obj)
		}
	})
}

func BenchmarkPA2(b *testing.B){
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			obj := make([]byte, 1024)
			l := len(obj)
			_ = l
		}
	})
}

func BenchmarkPA3(b *testing.B){
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			obj := make([]byte, 1024)
			_ = len(obj)
			_ = obj
		}
	})
}
