package mirg

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestGetGoldSalePriceConfig(t *testing.T) {
	GetGoldSalePriceConfig()
}

func TestUnsafe(t *testing.T) {
	// 任意指针类型都可以转为unsafe.Pointer
	// unsafe.Pointer可以转换为任意指针类型

	//uintptr可以转为unsafe.Pointer
	//unsafe.Pointer 可以转为uintptr

	// 1.把uint32转为int32(类型转换)
	var f uint32 = 100
	p := *(*int32)(unsafe.Pointer(&f))
	t.Log(fmt.Sprintf("0x%x，%v", p, &p))
	t.Log("")

	// 2.把int32转为uintptr(内存地址)
	t.Log(uintptr(unsafe.Pointer(&f)))
	t.Log(*(*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(&f)))))
	t.Log("")

	// 3.获取对象中的属性中的值 （1.获取对象的pointer  2.对象的pointer装换uintptr地址 + 属性的偏移量  3.根据pointer强转类型）
	type D struct {
		val       int64
		timestamp int64
	}
	data := D{timestamp: 1000}
	value := *(*int64)(unsafe.Pointer(uintptr(unsafe.Pointer(&data)) + unsafe.Offsetof(data.timestamp)))
	t.Log(value, unsafe.Pointer(&data), uintptr(unsafe.Pointer(&data)), unsafe.Offsetof(data.timestamp))
	t.Log("")
}
