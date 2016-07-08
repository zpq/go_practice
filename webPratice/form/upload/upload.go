package main

import (
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func upload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("upload.html")
		h := md5.New()
		cutime := time.Now().Unix()
		io.WriteString(h, strconv.FormatInt(cutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))
		f, _ := os.Create("md5")
		defer f.Close()
		f.Write([]byte(token))
		t.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, handle, err := r.FormFile("file")
		if err != nil {
			log.Fatal(err.Error())
		}
		defer file.Close()

		token, _ := ioutil.ReadFile("md5")
		if r.Form["token"][0] != string(token) {
			w.WriteHeader(403)
			return
		}
		fmt.Fprint(w, handle.Header)
		f, _ := os.Create("./uploadedFile")
		defer f.Close()
		io.Copy(f, file)
	}
}

func main() {
	http.HandleFunc("/upload", upload)
	http.ListenAndServe(":8001", nil)
}
