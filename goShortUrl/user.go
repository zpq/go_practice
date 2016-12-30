package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const (
	tokenSecret = "shortUrl_secret_zpq"
	myMd5Secret = "shortUrl_md5_secret_zpq"
)

type User struct {
	Id       int `xorm:"pk"`
	Username string
	Password string
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username, password := strings.TrimSpace(r.PostFormValue("username")), strings.TrimSpace(r.PostFormValue("password"))
	res := &Res{
		Status:  0,
		Message: "fail to login",
	}
	if username != "" && password != "" {
		user := &User{}
		ok, err := engine.Alias("t").Where("t.username = ? and t.password = ?", username, myMd5(password)).Get(user)
		if ok && err == nil {
			jToken, err := createToken(user.Id, username)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				res.Status = 1
				res.Message = "success to login"
				res.Datas = append(res.Datas, []string{jToken})
			}
		}
	}
	body, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	w.Write([]byte(body))
}

func UserRegister(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username, password := strings.TrimSpace(r.PostFormValue("username")), strings.TrimSpace(r.PostFormValue("password"))
	res := &Res{
		Status:  0,
		Message: "fail to register",
	}
	if username != "" && password != "" {
		user := &User{}
		ok, err := engine.Alias("t").Where("t.username = ?", username).Get(user)
		if !ok && err == nil { // can register
			user.Username = username
			user.Password = myMd5(password)
			lastID, err := engine.Insert(user)
			if err == nil && lastID > 0 { // success to insert
				res.Status = 1
				res.Message = "success to register"
			} else {
				fmt.Println(err.Error())
			}
		}
	}
	body, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	w.Write([]byte(body))
}

func myMd5(s string) string {
	m := md5.New()
	m.Write([]byte(s))
	return hex.EncodeToString(m.Sum([]byte(myMd5Secret)))
}

func createToken(userid int, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"userid":   userid,
		"exp":      time.Now().Add(time.Second * 3600),
		"nbf":      time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})
	return token.SignedString([]byte(tokenSecret))
}

func checkToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
