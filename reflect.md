# 一个针对 golang reflect 对 func 反射过慢的优化技巧

golang 的 func 的反射使用，目前还没有很合适的优化手段

但是可以很巧妙地绕过对 func 的反射


## 低效例子
gobench 测试大约 300 ns/op

```go
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

func apple(value,value2 int)int{
	return value
}

fn := reflect_build_args_and_use_func_A(apple)
fn()
```

## 高效例子
gobench 测试大约 50 ns/op
```go
func reflect_build_args_and_use_func_B(b iBanana) func() int{
	param := reflect.TypeOf(b).Elem()

	return func() int {
		arg := reflect.New(param).Interface().(iBanana)
		return arg.apple() // 无需反射
	}
}

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

fn := reflect_build_args_and_use_func_B(&banana{})
fn()
```

