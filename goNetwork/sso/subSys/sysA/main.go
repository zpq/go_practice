package main

import (
	"log"
	"net/http"
)

/**
* protected resource
 */
func protect(w http.ResponseWriter, r *http.Request) {
	//check auth can move to middleware layer(filter)
	c, err := r.Cookie("jwtSSO")
	if err != nil { // not login
		http.Redirect(w, r, "http://127.0.0.1:20000/login?redirectUrl=http://127.0.0.1:20001/protect&remote=http://127.0.0.1:20001/", 302)
		return
	}
	if res := verify(c.Value); res { // check ok
		w.Write([]byte("ok! hello! You are authed in 20001"))
	} else {
		http.Redirect(w, r, "http://127.0.0.1:20000/login?redirectUrl=http://127.0.0.1:20001/protect&remote=http://127.0.0.1:20001/", 302)
		return
	}
}

/**
* accecpt cas request with jwt message in the url
 */
func attach(w http.ResponseWriter, r *http.Request) {
	jwt := r.FormValue("jwt")
	redirectURL := r.FormValue("redirectUrl")
	if ok := verify(jwt); ok {
		cookie := http.Cookie{Name: "jwtSSO", Value: jwt, HttpOnly: true}
		http.SetCookie(w, &cookie)
		http.Redirect(w, r, redirectURL, 302)
	} else {
		http.Redirect(w, r, "http://127.0.0.1:20000/login?redirectUrl=http://127.0.0.1:20001/protect&remote=http://127.0.0.1:20001/", 302)
		return
	}
}

/**
* verify the jwt token by request sso
 */
func verify(payload string) bool {
	res, err := http.Get("http://127.0.0.1:20000/validate?jwt=" + payload + "&redirectUrl=http://127.0.0.1:20001/protect&remote=http://127.0.0.1:20001/")
	if err != nil {
		log.Println("validate http error: ", err.Error())
		return false
	}
	// here can parse userinfo from the sso message
	if res.StatusCode == 200 {
		return true
	}
	return false
}

func main() {
	http.HandleFunc("/protect", protect)
	http.HandleFunc("/attach", attach)
	if err := http.ListenAndServe(":20001", nil); err != nil {
		log.Fatal(err.Error())
	}
}
