/*
* watcher the change of file in the src dir, copy them to the dst dir, something like file backup
 */

package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

var md5Maps map[string][16]byte = make(map[string][16]byte, 128)
var md5MapsMux sync.Mutex

func main() {
	flag.Parse()
	for _, v := range flag.Args() {
		fmt.Println(v)
	}

	srcDir, dstDir := "F:\\gotest\\src\\", "F:\\gotest\\dst\\"
	for {

		dirs, err := ioutil.ReadDir(srcDir)
		if err != nil {
			fmt.Println("read dir error:", err.Error())
			return
		}
		for _, v := range dirs {
			func() {
				if v.IsDir() {
					return
				}
				fs, fd := srcDir+v.Name(), dstDir+v.Name()
				srcData, err := ioutil.ReadFile(fs)
				if err != nil {
					return
				}
				ok := parseFile(srcData, fs)
				if !ok {
					n, err := CopyFile(fs, fd)
					if err != nil {
						return
					} else {
						if n != -1 {
							fmt.Printf("%s copy ==> to %s, total byte %d\n", fs, fd, n)
						}
					}
				}
			}()
		}

		fmt.Println("I am watching...")
		time.Sleep(time.Second * 30)
	}
}

func CopyFile(src, dst string) (n int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return -1, err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return -1, err
	}
	defer dstFile.Close()
	return io.Copy(dstFile, srcFile)
}

func parseFile(srcData []byte, filename string) (ok bool) {
	md5MapsMux.Lock() //sync
	cipherText1 := md5.Sum(srcData)
	if cipherText1 == md5Maps[filename] {
		ok = true
	} else {
		md5Maps[filename] = cipherText1
		ok = false
	}
	md5MapsMux.Unlock()
	return
}
