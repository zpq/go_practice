package main

import (
	"crypto/md5"
	"fmt"
	"github.com/gorilla/sessions"
	"html/template"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

var store = sessions.NewCookieStore([]byte("mychat"))

func MyHttpWeb(handler http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "Goginx")
		session, _ := store.Get(r, "session-name3")
		if session.IsNew {
			session.Values["sid"] = generateSessionId()
			//need login
		}
		session.Save(r, w)
		// fmt.Println(session.IsNew, session.Name(), session.Values["sid"])
		handler.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func generateSessionId() string {
	nano := time.Now().UnixNano()
	rand.Seed(nano)
	rndNum := rand.Int63()
	sessionId := Md5(Md5(strconv.FormatInt(nano, 10)) + Md5(strconv.FormatInt(rndNum, 10)))
	return sessionId
}

func Md5(text string) string {
	hashMd5 := md5.New()
	io.WriteString(hashMd5, text)
	return fmt.Sprintf("%x", hashMd5.Sum(nil))
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		t, err := template.ParseFiles("../frontend/index.html")
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(500)
		} else {
			t.Execute(w, nil)
		}
	} else if r.URL.Path == "/chat" {

	} else {
		w.WriteHeader(404)
	}
}

func main() {
	myweb := MyHttpWeb(http.HandlerFunc(myHandler))
	err := http.ListenAndServe(":3005", myweb)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
