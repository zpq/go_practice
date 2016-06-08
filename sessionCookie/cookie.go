package main

import (
	"net/http"
	"time"
)

func myHandler(w http.ResponseWriter, r *http.Request) {
	expiration := time.Now()
	expiration = expiration.AddDate(1, 0, 0)
	cookie := http.Cookie{Name: "gocookie", Value: "Cookie set by Golang", Expires: expiration}
	http.SetCookie(w, &cookie)
	w.Write([]byte("Hello World!"))
}

func main() {
	http.HandleFunc("/", myHandler)
	http.ListenAndServe(":8001", nil)
}
