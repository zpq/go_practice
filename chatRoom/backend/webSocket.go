package main

import (
	"fmt"
	"github.com/gorilla/sessions"
	"html/template"
	"net/http"
	"os"
)

var store = sessions.NewCookieStore([]byte("mychat"))

func MyHttpWeb(handler http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "Goginx")

		session, _ := store.Get(r, "session-name")
		session.Values["foo"] = "bar"
		session.Values[42] = 43
		session.Save(r, w)
		fmt.Println(session)

		handler.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
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
