package main

import (
	"fmt"
	"reflect"
)

type A struct {
	Name  string
	Value int
}

func main() {
	a := true

	v := reflect.ValueOf(a)
	fmt.Println(v.IsZero())
	fmt.Println(v)

}
