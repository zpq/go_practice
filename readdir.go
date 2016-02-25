package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
)

func main() {
	dir := "F:\\gotest\\src\\"
	dirs, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, v := range dirs {
		// fmt.Println(v.Name(), " ", v.IsDir(), " ", v.ModTime(), " ", v.Mode(), " ", v.Size(), " ")
		if !v.IsDir() {
			srcData, err := ioutil.ReadFile(dir + v.Name())
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			cipherText1 := md5.Sum(srcData)
			fmt.Print(v.Name() + " ")
			fmt.Printf(" md5 encrypto is : %x \n", cipherText1)
		}

		// hash.Write(srcData)

		// cipherText2 := hash.Sum(nil)

		// hexText := make([]byte, 32)

		// hex.Encode(hexText, cipherText2)

		// fmt.Println("md5 encrypto is \"iyannik0215\":", string(hexText))
	}
}
