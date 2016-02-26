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
	"strings"
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
	if flag.NArg() != 2 {
		fmt.Println("args error")
		return
	}
	srcDir, dstDir := flag.Arg(0), flag.Arg(1)
	// srcDir, dstDir := "F:\\gotest\\src", "F:\\gotest\\dst"
	for {
		go core(srcDir, dstDir)
		fmt.Println("I am watching...")
		time.Sleep(time.Second * 30)
	}
}

func core(srcDir, dstDir string) {
	dirs, err := ioutil.ReadDir(srcDir)
	if err != nil {
		fmt.Println("read dir error:", err.Error())
		return
	}
	for _, v := range dirs {
		if v.IsDir() {
			subSrcDir := srcDir + v.Name() + "/"
			subDstDir := dstDir + v.Name() + "/"
			if !fileExists(subDstDir) {
				if err := os.Mkdir(dstDir+v.Name(), 0777); err != nil {
					fmt.Println("mkdir error:", err.Error())
					return
				}
			}
			go core(subSrcDir, subDstDir)
		} else {
			fi, err := os.Lstat(srcDir + v.Name())
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			if fi.Mode().IsDir() { //is a link file
				link, err := os.Readlink(srcDir + v.Name())
				if err != nil {
					fmt.Println(err.Error())
					return
				}
				links := strings.Split("/", link)
				for _, value := range links[len(links)-1:] {
					dstDir += "/" + value
				}

				if !fileExists(dstDir) {
					if err := os.Mkdir(dstDir, 0777); err != nil {
						fmt.Println("mkdir error:", err.Error())
						return
					}
				}
				// fmt.Println(link)
				go core(link, dstDir)
			} else {
				fs, fd := srcDir+v.Name(), dstDir+v.Name()
				srcData, err := ioutil.ReadFile(fs)
				if err != nil {
					fmt.Println("read file error:", err.Error())
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
			}
		}
	}
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
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
