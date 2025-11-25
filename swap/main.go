package main

import (
	"fmt"
	"reflect"
)

func swap[T any](a, b T) {
	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)

	if va.Kind() != reflect.Ptr || vb.Kind() != reflect.Ptr {
		panic("swap requires pointer")
	}

	tmp := reflect.New(va.Elem().Type())
	tmp.Elem().Set(va.Elem())
	va.Elem().Set(vb.Elem())
	vb.Elem().Set(tmp.Elem())
}

func main() {
	a := 10
	b := 20

	fmt.Printf("a = %d, &a = %p\n", a, &a)
	fmt.Printf("b = %d, &b = %p\n", b, &b)

	swap(&a, &b)

	fmt.Printf("a = %d, &a = %p\n", a, &a)
	fmt.Printf("b = %d, &b = %p\n", b, &b)
}
