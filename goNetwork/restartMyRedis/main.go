package main

import (
	"log"
	"net/http"
	"os/exec"
	"time"
)

const (
	articleURL = "http://sheaned.com/articles/16"
	interval   = time.Second * 3600
)

func checkArticle(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	if resp.StatusCode == 500 {
		return false
	}
	return true
}

func restartRedis() {
	cmd := exec.Command("supervisorctl", "restart", "redis")
	if err := cmd.Run(); err != nil {
		log.Println(err.Error())
	}
}

func main() {
	t := time.NewTicker(interval)
	for {
		select {
		case <-t.C:
			if !checkArticle(articleURL) {
				restartRedis()
			}
		default:
		}
	}
}
