package main

import (
	"net/http"
	"os/exec"
	"time"
	"os"
	"fmt"
)

const (
	articleURL = "http://sheaned.com/articles/16"
	interval   = time.Second * 1800
)

func checkArticle(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintln(os.Stdout, time.Now().String() + " " + err.Error())
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
		fmt.Fprintln(os.Stdout, time.Now().String() + " " +err.Error())
	}
}

func main() {
	t := time.NewTicker(interval)
	for {
		select {
		case <-t.C:
			if !checkArticle(articleURL) {
				restartRedis()
				fmt.Fprintln(os.Stdout, time.Now().String() + " broken down!\n")
			} else {
				fmt.Fprintln(os.Stdout, time.Now().String() + " redis test ok!\n")
			}
		default:
		}
	}
}
