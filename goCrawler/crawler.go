package main

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/go-xorm/xorm"
)

type icrawler interface {
	getHTML(url string) (string, error)
	getURLList(html, tag string)
	getArticle(html string)
	getTags(html string)
	run()
}

type tag struct {
	url string
}

type crawler struct {
	initURL     string
	urlProvider urlStorage
	rule        *rule
	tags        map[string]*tag
	maxCo       chan bool
	results     []*result
	syn         sync.Mutex
	db          *xorm.Engine
}

type rule struct {
	articleList *regexp.Regexp
	title       *regexp.Regexp
	time        *regexp.Regexp
	content     *regexp.Regexp
	tagURL      *regexp.Regexp
	tagName     *regexp.Regexp
}

type result struct {
	Id      int `xorm:"pk"`
	Tag     string
	Title   string
	Time    string
	Content string
}

func (c *crawler) getHTML(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (c *crawler) getURLList(html, tag string) {
	articleList := c.rule.articleList.FindAllStringSubmatch(html, -1)
	for _, v := range articleList {
		c.urlProvider.addOneURL(tag, v[1])
	}
}

func (c *crawler) getTags(html string) {
	tagURL := c.rule.tagURL.FindAllStringSubmatch(html, -1)
	for _, v := range tagURL {
		tg := &tag{url: v[1]}
		tagName := strings.SplitAfter(v[1], "/")
		c.tags[tagName[len(tagName)-1]] = tg
	}
}

func (c *crawler) getArticle(html string) {

}

func (c *crawler) run() {

}
