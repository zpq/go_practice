package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	dir        = "F:\\go\\workspace\\src\\go_practice\\webPratice\\room\\imgs\\"
	url        = "http://www.qq.com"
	urlRegexp  = regexp.MustCompile(`<a.*href="(.*?)"`)
	imgRegexp  = regexp.MustCompile(`<img.*src="(.*?)"`)
	httpRegexp = regexp.MustCompile(`^https?://.*`)
	maxCur     chan bool
)

const maxDepth = 3

func getContent(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	return string(body)
}

func filterURL(urls [][]string) []string {
	var newUrls []string
	for _, v := range urls {
		ok := httpRegexp.MatchString(v[1])
		if ok {
			newUrls = append(newUrls, v[1])
		}
		// else {
		// 	newUrls = append(newUrls, "http:"+v[1])
		// }
	}
	return newUrls
}

func download(url string) bool {
	maxCur <- true
	defer func() {
		fmt.Println(<-maxCur)
	}()
	imgData := []byte(getContent(url))
	pos := strings.LastIndex(url, "/")
	filename := url[pos+1:]
	f, err := os.Create(dir + filename)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	defer f.Close()
	_, err = f.Write(imgData)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	fmt.Println(url, " finished downloaded!")
	return true
}

func handle(url string, depth int) {
	if depth <= maxDepth {
		html := getContent(url)
		urls := urlRegexp.FindAllStringSubmatch(html, -1)
		httpURL := filterURL(urls)
		images := imgRegexp.FindAllStringSubmatch(html, -1)
		imgURL := filterURL(images)
		go func() {
			for _, v := range imgURL {
				go download(v)
			}
		}()
		for _, v := range httpURL {
			go handle(v, depth+1)
		}
	}
}

func main() {
	maxCur = make(chan bool, 10)
	handle(url, 1)
	time.Sleep(time.Second * 60)
}
