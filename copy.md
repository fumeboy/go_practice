# 一个快速拷贝被类型转换为 interface 的 struct 的方法

这个方法生产环境下大概是不能用的

方法是 从 interface 获取 struct 的地址, 然后将 struct 转成 []byte， 拷贝后再将 []byte 转成 struct

详见 copy_test.go 注释