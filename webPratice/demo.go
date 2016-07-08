package main

import (
	"fmt"
)

type A struct {
	name string
}

type AB interface {
	Say() string
}

func (a A) Say() string {
	return a.name
}

func Talk(ab AB) {
	fmt.Println(ab.Say())
}

func main() {
	a := A{"jack"}
	Talk(a)
	var x interface{} = "asd"
	fmt.Println(x.(string))
}
