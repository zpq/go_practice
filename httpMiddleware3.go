package main

import (
	"net/http"
	"net/http/httptest"
)

type httpMiddleware struct {
	handler http.Handler
}

func (this *httpMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rec := httptest.NewRecorder()
	this.handler.ServeHTTP(rec, r)

	for k, v := range rec.Header() {
		w.Header()[k] = v
	}

	w.Header().Set("Server", "Goginx")
	w.WriteHeader(418)
	w.Write([]byte("Golang http middleware!\n"))
	w.Write(rec.Body.Bytes())
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/hello" {
		hello(w, r)
	} else {
		w.Write([]byte("Your uri is invalid!"))
	}
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello Method!"))
}

func main() {
	hm := &httpMiddleware{http.HandlerFunc(myHandler)}
	http.ListenAndServe(":3003", hm)
}
