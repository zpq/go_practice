package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// var sy sync.WaitGroup

// func handle() {
// 	sy.Add(1)
// 	cmd := exec.Command("F:\\nodejs\\phantomjs\\bin\\phantomjs.exe", "F:\\nodejs\\phantomjs\\workspace\\login\\index.js")
// 	// err := cmd.Start()
// 	out, err := cmd.Output()
// 	if err != nil {
// 		log.Println(err.Error())
// 	}
// 	log.Println(string(out))
// 	sy.Done()
// }

func main() {

	for index := 0; index < 50000; index++ {
		go func() {
			client := &http.Client{}
			req, err := http.NewRequest("GET", "http://systemA.com/protect", strings.NewReader("name=cjb"))
			if err != nil {
				log.Println(err.Error())
				return
			}
			cookie := http.Cookie{
				Name:     "sid",
				Value:    "35f2cde457a3c2c0d4f796816dbdafe6",
				Domain:   ".my.sso.com",
				HttpOnly: true,
				Path:     "/",
				Expires:  time.Now().Add(time.Second * 7200),
			}
			req.AddCookie(&cookie)
			cookie2 := http.Cookie{
				Name:     "jwtSSO",
				Value:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwic2Vzc2lvbklkIjoiMzVmMmNkZTQ1N2EzYzJjMGQ0Zjc5NjgxNmRiZGFmZTYiLCJpc3MiOiJzc28ifQ.b8jLDYwUxqusur0Ph3qK7Ou-vg175cIdm8j3nVB92OU",
				Domain:   "systemA.com",
				HttpOnly: true,
				Path:     "/",
				Expires:  time.Now().Add(time.Second * 7200),
			}
			req.AddCookie(&cookie2)

			resp, err := client.Do(req)
			if err != nil {
				log.Println(err.Error())
				return
			}

			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println(err.Error())
				return
			}
			log.Println(string(body))
		}()
		time.Sleep(time.Millisecond * 10) // 100qps
	}
}
