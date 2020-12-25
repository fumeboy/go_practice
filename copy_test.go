package main_test

import (
	"fmt"
	"testing"
	"unsafe"
)

type beCopy struct {
	value int
}

type emptyInterface struct {
	typ  *struct{size       uintptr}
	word unsafe.Pointer
}
type slice struct {
	array unsafe.Pointer
	len   int
	cap   int
}
func Copy(v interface{}) interface{} {
	// 传进来的是一个 struct， 但是类型转换为了 interface
	// 目的是拷贝这个 struct
	// 因为 interface 本质是一个携带了原来类型信息的指针
	// 所以直接 值传递 拷贝是不行的， 值传递拷贝只能再得到一个这样的指针

	var v_p = (*emptyInterface)(unsafe.Pointer(&v))
	var length = int(v_p.typ.size)



	var vslice = []byte{}
	var vslice_p = (*slice)(unsafe.Pointer(&vslice))
	// 所以这里将 struct 的地址替换到 slice 的地址位， 使 vslice 指向的一串内存就是 struct 的内存
	vslice_p.array = v_p.word
	vslice_p.len = length
	vslice_p.cap = length

	vvslice := make([]byte, length) // 再创建一个 slice
	copy(vvslice, vslice) // 将 struct 的内存拷贝到新的 slice

	vv := v // 拷贝一个 interface 指针
	((*emptyInterface)(unsafe.Pointer(&vv))).word = (*slice)(unsafe.Pointer(&vvslice)).array // 将 新 slice 的内存地址 替换为结构体指针指向的地址
	return vv // 返回深拷贝后的 interface 指针
	// 大致意思就是，将 struct 转成 []byte， 拷贝后再将 []byte 转成 struct
}

func TestCopy(t *testing.T){
	b := beCopy{value: 3}
	d := Copy(b)

	b.value++
	e := Copy(b)
	fmt.Println(b,d,e)
}

func BenchmarkCopy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c := beCopy{value: 3}
		Copy(c)
	}
}