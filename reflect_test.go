package main_test

import (
	"reflect"
	"testing"
)

// A

func apple(value,value2 int)int{
	return value
}

func reflect_build_args_and_use_func_A(fn interface{}) func() int{
	fnt := reflect.TypeOf(fn)
	fnv := reflect.ValueOf(fn)

	param := fnt.In(0)
	param2 := fnt.In(1)
	return func() int {
		arg := reflect.New(param).Elem()
		arg2 := reflect.New(param2).Elem()
		resp := fnv.Call([]reflect.Value{arg, arg2}) // 严重耗时
		return resp[0].Interface().(int)
	}
}

func BenchmarkReflectA(b *testing.B)  {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fn := reflect_build_args_and_use_func_A(apple)
		fn()
	}
}

// B

type banana struct {
	value int
	value2 int
}

type iBanana interface {
	apple() int
}

func (b *banana) apple() int {
	return b.value
}

func reflect_build_args_and_use_func_B(b iBanana) func() int{
	param := reflect.TypeOf(b).Elem()

	return func() int {
		arg := reflect.New(param).Interface().(iBanana)
		return arg.apple() // 无需反射
	}
}

func BenchmarkReflectB(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fn := reflect_build_args_and_use_func_B(&banana{})
		fn()
	}
}
