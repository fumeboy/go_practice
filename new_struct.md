# 一个构造结构体的方法

假设要构造 testStruct 这个结构体, 常规方法是写一个 `func NewStruct() testStruct` 函数

这里的方法是 将这个函数 改写成结构体方法 `func (s *testStruct) New() *testStruct`

比较一下使用时的不同:

```go
_ = NewStruct()

_ = (&testStruct{}).New()

```

可以发现 New 这个函数被约束到了 testStruct 命名空间下, 这也是我的主要目的

## 完整的测试用例:
基准测试中, 两种方式的执行效率几乎一致,当 struct 很大时, `(&testStruct{}).New()` 略有优势

```go
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
```
