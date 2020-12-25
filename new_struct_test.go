package main

import (
	"testing"
)

type testStruct struct {
	value int
}

// 三种构造结构体的方法

func NewStruct() *testStruct { // 逃逸
	return &testStruct{}
}

func NewStruct2() testStruct { // 不逃逸 但传值时发生拷贝
	return testStruct{}
}

func (s *testStruct) New() *testStruct {
	// 不逃逸，只拷贝了指针
	// 更重要的是收束到了 testStruct 这个命名空间里
	return s
}



func TestNew(t *testing.T){
	_ = (&testStruct{}).New()
	//fmt.Println(s)
}

func BenchmarkNew(b *testing.B){
	for i := 0;i<b.N;i++{
		_ = (&testStruct{}).New()
	}
}

func BenchmarkNew2(b *testing.B){
	for i := 0;i<b.N;i++{
		_ = NewStruct()
	}
}