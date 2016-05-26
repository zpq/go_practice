package main

import (
	// "fmt"
	"net/http"
)

type SingleHost struct {
	handler   http.Handler
	allowHost []string
}

func (this *SingleHost) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	allowed := false
	for _, v := range this.allowHost {
		if r.Host == v {
			allowed = true
			break
		}
	}
	if allowed {
		w.Write([]byte("I am a http middleware!\n"))
		this.handler.ServeHTTP(w, r)
	} else {
		w.WriteHeader(403)
	}
}

//实际的业务逻辑函数
func myHandler(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/hello" {
		hello(w, r)
	} else {
		w.Write([]byte("Your uri is invalid!"))
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello!I am hello method!"))
}

func main() {
	sh := &SingleHost{
		handler:   http.HandlerFunc(myHandler),
		allowHost: []string{"localhost:3001", "127.0.0.1:3001"},
	}
	http.ListenAndServe(":3001", sh)
}
