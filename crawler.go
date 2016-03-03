package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	// "strconv"
	"strings"
	"time"
)

var url string = "http://www.qq.com"
var dir string = "F:\\imgs\\go\\"
var (
	// ptnIndexItem    = regexp.MustCompile(`<a target="_blank" href="(.+\.html)" title=".+" >(.+)</a>`)
	// ptnContentRough = regexp.MustCompile(`(?s).*<div class="artcontent">(.*)<div id="zhanwei">.*`)
	// ptnBrTag        = regexp.MustCompile(`<br>`)
	// ptnHTMLTag      = regexp.MustCompile(`(?s)</?.*?>`)
	// ptnSpace        = regexp.MustCompile(`(^\s+)|( )`)
	src    = regexp.MustCompile(`<img.*src="(.*?)"`)
	regUrl = regexp.MustCompile(`<a.*href="(.*?)"`)
)

func main() {
	content, statusCode := GetHtmlSource(url)
	if statusCode != 200 {
		fmt.Println("error code :", statusCode)
		return
	}
	imgs := src.FindAllStringSubmatch(content, -1)
	urls := regUrl.FindAllStringSubmatch(content, -1)
	lens := len(imgs)
	for i := 0; i < lens; i++ {
		matches := strings.Split(imgs[i][1], "/")
		go downLoadFile(imgs[i][1], matches[len(matches)-1])
	}

	time.Sleep(time.Second * 10)
	for _, v := range urls {
		fmt.Println(v[1])
	}
	time.Sleep(time.Second * 10)
}

// func urlHandler() {

// }

func downLoadFile(url, filename string) {
	resp, err := http.Get(url)
	// fmt.Println(url)
	if err != nil {
		fmt.Println("get error:" + err.Error())
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("io read error:" + err.Error())
		return
	}

	file, err := os.Create(dir + filename)
	defer file.Close()
	if err != nil {
		fmt.Println("file os error:" + err.Error())
		return
	}

	_, err = file.Write(data)
	if err != nil {
		fmt.Println("file write error:" + err.Error())
		return
	}
	fmt.Println("download success ", dir+filename)
}

func GetHtmlSource(url string) (content string, statusCode int) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		statusCode = 0
		return "", statusCode
	}
	defer resp.Body.Close()
	buff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		statusCode = 0
		return "", statusCode
	}
	content = string(buff)
	return content, resp.StatusCode
}
