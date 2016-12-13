package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"time"
)

var material = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

func makeNumber() (res []string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	length := len(material)
	i := 0
	for i < 4 {
		t := material[r.Intn(length)]
		if !isInArray(t, res) {
			res = append(res, t)
			i++
		}
	}
	return res
}

func isInArray(search string, haystack []string) bool {
	l := len(haystack)
	for i := 0; i < l; i++ {
		if haystack[i] == search {
			return true
		}
	}
	return false
}

func checkNumber(s string, gn []string) (bool, string) {
	a, b, ok := 0, 0, false
	for i := 0; i < len(gn); i++ {
		if s[i:i+1] == gn[i] {
			a++
		} else if isInArray(s[i:i+1], gn) {
			b++
		}
	}
	if a == 4 {
		ok = true
	}
	return ok, fmt.Sprintf("%dA%dB", a, b)
}

func validate(s []byte) bool {
	return regexp.MustCompile("^[0-9]{4}$").Match(s)
}

func main() {
	total, win := 0, 0
	for {
		gn := makeNumber()
		chance := 7
		for chance > 0 {
			fmt.Printf("\nU left %d chance!\nPlease input your answer:\n", chance)
			input, _, err := bufio.NewReader(os.Stdin).ReadLine()
			if err != nil {
				log.Fatal(err.Error())
			}

			if !validate(input) {
				fmt.Println("you input is invalid! Please input 4 numbers without space!")
				continue
			}

			right, res := checkNumber(string(input), gn)
			if right {
				fmt.Printf("%s!You are right!The number is %s\n", res, gn)
				win++
				break
			} else {
				fmt.Printf("You are wrong!Tip is %s\n", res)
			}
			chance--
			if chance == 0 {
				fmt.Println("the number is ", gn)
			}
		}
		total++
		fmt.Println("1: continue  2: exit")
		goOn, _, _ := bufio.NewReader(os.Stdin).ReadLine()
		if string(goOn) != "1" {
			fmt.Printf("total %d, won %d\n", total, win)
			break
		}
	}

}
