package main

import (
	"fmt"
)

func letterSystem(num int) string {
	letter := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	llen := len(letter)
	shang := num / llen
	yushu := num % llen
	if shang == 0 {
		return letter[yushu]
	} else {
		return letterSystem(shang) + letter[yushu]
	}
}

func main() {
	for i := 0; i < 1000; i++ {
		fmt.Println(i, " = ", letterSystem(i))
		if (i+1)%26 == 0 {
			fmt.Println("进位")
		}
	}
}
