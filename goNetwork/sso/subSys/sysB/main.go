package main

import (
	"log"
	"net/http"
)

const (
	server    = "http://my.sso.com"
	selfAddr  = "http://systemB.com/"
	ssoLogout = server + "/logout"
	port      = ":20002"
)

/**
* protected resource
 */
func protect(w http.ResponseWriter, r *http.Request) {
	//check auth can move to middleware layer(filter)
	c, err := r.Cookie("jwtSSO")
	if err != nil { // not login
		http.Redirect(w, r, server+"/login?redirectUrl="+selfAddr+"/protect&remote="+selfAddr, 302)
		return
	}
	if res := verify(c.Value); res { // check ok
		w.Write([]byte("ok! hello! You are authed in " + selfAddr))
	} else {
		http.Redirect(w, r, server+"/login?redirectUrl="+selfAddr+"/protect&remote="+selfAddr, 302)
		return
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, ssoLogout+"?redirectUrl="+selfAddr+"index", 302)
	return
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
		http.Redirect(w, r, server+"/login?redirectUrl="+selfAddr+"/protect&remote="+selfAddr, 302)
		return
	}
}

/**
* verify the jwt token by request sso
 */
func verify(payload string) bool {
	res, err := http.Get(server + "/validate?jwt=" + payload + "&redirectUrl=" + selfAddr + "/protect&remote=" + selfAddr)
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

//not login
func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("you are not login yet"))
}

func main() {
	http.HandleFunc("/protect", protect)
	http.HandleFunc("/attach", attach)
	http.HandleFunc("/index", index)
	http.HandleFunc("/logout", logout)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err.Error())
	}
}
