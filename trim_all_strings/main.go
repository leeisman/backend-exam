package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func TrimAllStrings(a any) {
	visited := make(map[uintptr]bool) // 用來記錄已處理過的指標，避免 cycle
	trimValue(reflect.ValueOf(a), visited)
}

func trimValue(v reflect.Value, visited map[uintptr]bool) {
	if !v.IsValid() {
		return
	}

	switch v.Kind() {

	case reflect.Pointer:
		if v.IsNil() {
			return
		}
		// 檢查是否已拜訪過這個指標（處理 cycle）
		ptr := v.Pointer()
		if visited[ptr] {
			return
		}
		visited[ptr] = true
		trimValue(v.Elem(), visited)

	case reflect.Interface:
		if v.IsNil() {
			return
		}
		trimValue(v.Elem(), visited)

	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			trimValue(field, visited)
		}

	case reflect.String:
		if v.CanSet() {
			v.SetString(strings.TrimSpace(v.String()))
		}

	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			trimValue(v.Index(i), visited)
		}

	case reflect.Map:
		for _, key := range v.MapKeys() {
			val := v.MapIndex(key)
			trimValue(val, visited)
		}

	// 其他型別不處理
	default:
		return
	}
}

func main() {
	type Person struct {
		Name string
		Age  int
		Next *Person
	}

	a := &Person{
		Name: " name ",
		Age:  20,
		Next: &Person{
			Name: " name2 ",
			Age:  21,
			Next: &Person{
				Name: " name3 ",
				Age:  22,
			},
		},
	}

	TrimAllStrings(&a)

	m, _ := json.Marshal(a)

	fmt.Println(string(m))

	a.Next = a

	TrimAllStrings(&a)

	fmt.Println(a.Next.Next.Name == "name")
}
