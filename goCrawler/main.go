package main

import (
	"errors"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

const (
	maxThreadPerWebSite = 20
	interval            = 2 // seconds
	maxPage             = 1
	dsn                 = "root:123456@/goCrawler?charset=utf8"
	redisIP             = "192.168.2.131"
	redisPORT           = "6379"
)

type openOpenCrawler struct {
	crawler
}

func (c *openOpenCrawler) initRules() {
	c.rule = &rule{
		articleList: regexp.MustCompile(`<h3><a href="(.*?)"`),
		title:       regexp.MustCompile(`<h1 id="articleTitle" >(.*?)</h1>`),
		time:        regexp.MustCompile(`<div class="meta"[\s\S]+?<span class=item>(.*?)</span>`),
		content:     regexp.MustCompile(`<article>([\s\S]+)</article>`),
		tagURL:      regexp.MustCompile(`<li class="nav__item"><a href="(.*?)"`),
		tagName:     regexp.MustCompile(`<li class="nav__item"><a href=".*?>(.*?)</a>`),
	}
}

func (c *openOpenCrawler) getTags(html string) {
	tagURL := c.rule.tagURL.FindAllStringSubmatch(html, -1)
	for _, v := range tagURL {
		tg := &tag{url: v[1]}
		tagName := strings.SplitAfter(v[1], "/")
		c.tags[tagName[len(tagName)-1]] = tg
	}
}

func (c *openOpenCrawler) getURLList(html, tag string) {
	articleList := c.rule.articleList.FindAllStringSubmatch(html, -1)
	for _, v := range articleList {
		c.urlProvider.addOneURL(tag, v[1])
	}
}

func (c *openOpenCrawler) getArticle(html, tag string) error {
	title := c.rule.title.FindAllStringSubmatch(html, -1)
	article := c.rule.content.FindAllStringSubmatch(html, -1)
	timeString := c.rule.time.FindAllStringSubmatch(html, -1)

	if title == nil || len(title[0]) < 2 ||
		timeString == nil || len(timeString[0]) < 2 ||
		article == nil || len(article[0]) < 2 {
		return errors.New("catch error")
	}

	res := &result{
		Tag:     tag,
		Title:   title[0][1],
		Time:    timeString[0][1],
		Content: strings.TrimSpace(article[0][1]),
	}

	// c.syn.Lock()
	// c.results = append(c.results, res)
	// c.syn.Unlock()
	_, err := c.db.Insert(res)
	if err != nil {
		return err
	}
	log.Println("insert ok!" + res.Title + " " + strconv.Itoa(len(res.Content)))
	return nil
}

func (c *openOpenCrawler) run() {
	c.initRules()
	html, err := c.getHTML(c.initURL + "/lib")
	if err != nil {
		log.Println("getHTML error: " + err.Error())
		return
	}
	log.Println("runing....")
	c.getTags(html)
	for k := range c.tags {
		go func(k string) {
			// c.maxCo <- true
			// defer log.Println(<-c.maxCo)
			// synG.Add(1)
			i := 0
			for i <= maxPage {
				// log.Printf("i=%d\n", i)
				body, err := c.getHTML(c.initURL + "/lib/tag/" + k + "?pn=" + strconv.Itoa(i))
				if err != nil {
					log.Println(err.Error())
					continue
				}
				c.getURLList(body, k)
				i++
				time.Sleep(time.Second * interval * 2)
			}
			// synG.Done()
		}(k)

		go func(k string) {
			// c.maxCo <- true
			// defer log.Println(<-c.maxCo)
			for {
				url := c.urlProvider.getOneURL(k)
				if url == "" {
					time.Sleep(time.Second * interval)
				}
				body, err := c.getHTML(c.initURL + url)
				if err != nil {
					log.Println(err.Error())
					continue
				}
				err = c.getArticle(body, k)
				if err != nil {
					log.Println(err.Error())
					continue
				}
				time.Sleep(time.Second * interval)
			}
		}(k)
	}
}

var synG sync.WaitGroup

func main() {
	var err error
	redisPool, err := newRedisPool("")
	if err != nil {
		log.Fatal("redis connection error : ", err.Error())
	}
	oc := &openOpenCrawler{
		crawler: crawler{
			initURL:     "http://www.open-open.com",
			maxCo:       make(chan bool, maxThreadPerWebSite),
			tags:        make(map[string]*tag),
			urlProvider: redisPool,
		},
	}
	engine, err := xorm.NewEngine("mysql", dsn)
	if err != nil {
		log.Fatalln(err.Error())
	}
	oc.db = engine
	oc.run()
	time.Sleep(time.Second * 30)
}
