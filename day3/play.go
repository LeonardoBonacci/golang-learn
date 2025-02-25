package main

import (
	"fmt"
	"reflect"
	"strconv"
)

var i, j = 1, 2

func swap(x, y, z string) (string, string, int, error) {
	num, err := strconv.Atoi(x)
	return y, z, num, err
}

func main() {
	// a, b, c, err := swap("10", "world", "20")
	// fmt.Println(a, b, c, err)

	var c, python, java = true, false, "foo"
	fmt.Println(i, j, c, python, reflect.TypeOf(java))

	v := "42" // change me!
	v = "foo"
	fmt.Printf("v is of type %T\n", v)

	const Truth = true
	fmt.Println("Go rules?", Truth)
}
