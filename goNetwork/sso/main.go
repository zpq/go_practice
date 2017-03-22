package main

import (
	"crypto/md5"
	"encoding/hex"
	"html/template"
	"log"
	"net/http"

	"fmt"

	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	USERNAME = "admin"
	PASSWORD = "admin"
)

var (
	mySignedKey = []byte("dgsasagsa")
	myRsaKey    = "rsakey_1233"
	rp          *redisPool
)

type user struct {
	name string
}

type clientInfo struct {
	RedirectURL string
	Host        string
}

type myClaimJwt struct {
	Username  string `json:"username"`
	SessionID string `json:"sessionId"`
	jwt.StandardClaims
}

func login(w http.ResponseWriter, r *http.Request) {
	ci := clientInfo{
		r.FormValue("redirectUrl"),
		r.FormValue("remote"),
	}
	sid, err := r.Cookie("sid")
	if err == nil { // has sid cookie, then check it
		if ok, _ := checkAuth(sid.Value); ok { // already login
			v, err := rp.pool.Get().Do("HGET", sid.Value, "username")
			if err != nil {
				return
			}
			vv := v.([]byte)
			jwtSign := makeJwtSign(string(vv), sid.Value)
			if jwtSign == "" {
				log.Println("jwt token gen error: ", err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			if redirectURL := r.FormValue("redirectUrl"); redirectURL != "" {
				http.Redirect(w, r, r.FormValue("remote")+"attach?jwt="+jwtSign+"&redirectUrl="+redirectURL, 302)
			} else {
				http.Redirect(w, r, r.Referer(), 302)
			}
			return
		}
	}
	//above code can move to middleware layer

	t, err := template.ParseFiles("./static/login.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	t.Execute(w, ci)
}

func doLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		username := r.PostFormValue("username")
		password := r.PostFormValue("password")
		if username != USERNAME || password != PASSWORD {
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		// generate sid, create global session
		// mr := key.NewRsaGen(2048)
		// err := mr.GenKey()
		// if err != nil {
		// 	log.Println("genkey error: ", err.Error())
		// 	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		// 	return
		// }
		// v, err := mr.Encrypt([]byte(fmt.Sprintf("%s%s%s%s", username, password, time.Now().String(), myRsaKey)))
		// if err != nil {
		// 	log.Println("encrypt error: ", err.Error())
		// 	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		// 	return
		// }

		v := myMd5(fmt.Sprintf("%s%s%s%s", username, password, time.Now().String(), myRsaKey))

		// store session in redis
		_, err := rp.pool.Get().Do("HMSET", v, "sessionid", v, "username", username)
		if err != nil {
			log.Println("redis hash set error: ", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		//set redis expired
		_, err = rp.pool.Get().Do("EXPIRE", v, "14400")
		if err != nil {
			log.Println("redis hash set expire error: ", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// set cookie for client
		cookie := http.Cookie{
			Name:     "sid",
			Value:    v,
			HttpOnly: true,
			Domain:   "http://127.0.0.1:20000/",
			Path:     "/",
			Expires:  time.Now().Add(time.Second * 7200),
			// MaxAge:   int(time.Second * 7200),
		}
		http.SetCookie(w, &cookie)

		// generate jwt token for subSystem
		jwtSign := makeJwtSign(username, string(v))
		if jwtSign == "" {
			log.Println("jwt token gen error: ", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		// log.Println("jwt = ", jwtSign)
		http.Redirect(w, r, r.PostFormValue("host")+"attach?redirectUrl="+r.PostFormValue("redirectUrl")+"&jwt="+jwtSign, 302)
	} else {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func myMd5(data string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(data))
	cipherStr := md5Ctx.Sum(nil)
	// fmt.Print(cipherStr)
	// fmt.Print("\n")
	return hex.EncodeToString(cipherStr)
}

func checkAuth(auth string) (bool, *user) {
	sid := auth
	if sid == "" {
		return false, nil
	}
	v, err := rp.pool.Get().Do("HGET", sid, "sessionid")
	if err != nil {
		return false, nil
	}

	if v == nil {
		return false, nil
	}

	vv := v.([]byte)
	// fmt.Printf("%t %v\n", vv, string(vv))
	if vv == nil || string(vv) == "" {
		return false, nil
	}
	return true, &user{}
}

func makeJwtSign(data ...string) string {
	claims := myClaimJwt{
		data[0],
		data[1],
		jwt.StandardClaims{
			ExpiresAt: int64(time.Second * 10),
			Issuer:    "sso",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySignedKey)
	if err != nil {
		log.Println("JWT SIGN ERROR: ", err.Error())
		return ""
	}
	return ss
}

func validate(w http.ResponseWriter, r *http.Request) {
	log.Println("subsystem call validate ", r.Referer())
	j := r.FormValue("jwt")
	token, err := jwt.ParseWithClaims(j, &myClaimJwt{}, func(token *jwt.Token) (interface{}, error) {
		return mySignedKey, nil
	})
	if err != nil {
		log.Println("jwt parse error:", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if claims, ok := token.Claims.(*myClaimJwt); ok && token.Valid {
		// fmt.Printf("%v %v %v\n", claims.SessionID, claims.Username, claims.StandardClaims.ExpiresAt)
		if ok, _ := checkAuth(claims.SessionID); ok {
			w.WriteHeader(200)
			fmt.Fprintln(w, "success")
		} else {
			log.Println("jwt get sessionID error")
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else {
		log.Println("jwt parse error 2")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func main() {
	var err error
	rp, err = newRedisPool("")
	if err != nil {
		log.Fatal("redis connect error: ", err.Error())
	}

	http.HandleFunc("/validate", validate)
	http.HandleFunc("/login", login)
	http.HandleFunc("/doLogin", doLogin)
	if err := http.ListenAndServe(":20000", nil); err != nil {
		log.Fatal(err.Error())
	}
}
