package main

import (
	"fmt"
)

type ints []int

func (i ints) Iterator() *Iterator {
	return &Iterator{
		data:  i,
		index: 0,
	}
}

type Iterator struct {
	data  ints
	index int
}

func (i *Iterator) hasNext() bool {
	return i.index < len(i.data)
}

func (i *Iterator) next() int {
	v := i.data[i.index]
	i.index++
	return v
}

func main() {
	d := ints{1, 2, 3}
	for it := d.Iterator(); it.hasNext(); {
		fmt.Println(it.next())
	}
}
