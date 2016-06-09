package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

func myHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Println(r.URL.Path)
		t, _ := template.ParseFiles("index.html")
		cutime := time.Now().Unix()
		fmt.Println(cutime)
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(cutime, 10))
		f, _ := os.Create("./md5.txt")
		defer f.Close()
		token := fmt.Sprintf("%x", h.Sum(nil))
		fmt.Println(token)
		f.Write([]byte(token))
		t.Execute(w, token)
	} else {
		r.ParseForm()
		token := r.PostFormValue("token")
		data, err := ioutil.ReadFile("./md5.txt")
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if token != string(data) {
			w.WriteHeader(403)
			return
		}
		username := r.PostFormValue("username")
		for k, v := range r.Form {
			fmt.Println(k, " : ", v)
		}
		w.Write([]byte("Hello " + username + "! You are login"))
	}
}

func main() {
	http.HandleFunc("/", myHandler)
	http.ListenAndServe(":8000", nil)
}
