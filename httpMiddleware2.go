package main

import (
	"net/http"
)

func SingleHost(handler http.Handler, allowHost string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if allowHost == r.Host {
			handler.ServeHTTP(w, r)
			w.Write([]byte("I am a http middleware!\n"))
		} else {
			w.WriteHeader(403)
		}
	}
	return http.HandlerFunc(fn)
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello!"))
}

func main() {
	single := SingleHost(http.HandlerFunc(myHandler), "127.0.0.1:3002")
	http.ListenAndServe(":3002", single)
}
