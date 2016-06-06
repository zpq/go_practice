package main

import (
	"fmt"
)

func letterSystem(num int) string {
	letter := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	shang := num / 26
	yushu := num % 26
	if shang <= 26 {
		if shang == 0 {
			return letter[yushu-1]
		}
		if yushu == 0 && shang == 1 {
			return "Z"
		}
		if yushu == 0 && shang > 1 {
			return letter[shang-2] + "Z"
		}
		return letter[shang-1] + letter[yushu-1]
	} else {
		return letter[25] + letterSystem((shang-26)*26+yushu)
	}
}

func main() {
	for i := 1; i < 1000; i++ {
		fmt.Println(i, " = ", letterSystem(i))
		if i%26 == 0 {
			fmt.Println("进位")
		}
	}
}
