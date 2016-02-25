package main

import (
	// "crypto/md5"
	"fmt"
	"io/ioutil"
)

func main() {
	dir := "F:\\go\\workspace\\"
	level := "|-"
	readdir(dir, level)
}

func readdir(dir, level string) {
	dirs, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, v := range dirs {
		if v.IsDir() {
			fmt.Println(level + v.Name())
			subdir := dir + v.Name() + "\\"
			// fmt.Println(subdir)
			level += "-"
			readdir(subdir, level)
		} else {
			// fmt.Println(v.Name(), " ", v.IsDir(), " ", v.ModTime(), " ", v.Mode(), " ", v.Size(), " ")
			fmt.Println(level + v.Name())
		}
		// if !v.IsDir() {
		// 	srcData, err := ioutil.ReadFile(dir + v.Name())
		// 	if err != nil {
		// 		fmt.Println(err.Error())
		// 		continue
		// 	}
		// }
	}
}
