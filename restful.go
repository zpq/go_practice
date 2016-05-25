package main

import (
	"encoding/json"
	"fmt"
	"github.com/bmizerany/pat"
	"log"
	"net/http"
)

var users map[string]*User

type User struct {
	Id   string
	Name string
}

func main() {
	users = make(map[string]*User)
	mux := pat.New()

	mux.Get("/user/:name/profile", http.HandlerFunc(profile))
	mux.Post("/user", http.HandlerFunc(addUser))
	mux.Put("/user", http.HandlerFunc(updateUser))
	mux.Get("/users", http.HandlerFunc(getAllUsers))

	http.Handle("/", mux)
	http.ListenAndServe(":3001", nil)
	log.Println("listen at port 3001....")
}

func profile(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	name := params.Get(":name")
	ret := &User{}
	for _, v := range users {
		if v.Name == name {
			ret = v
			break
		}
	}
	body, _ := json.Marshal(ret)
	w.Write(body)
}

func addUser(w http.ResponseWriter, r *http.Request) {
	name := r.PostFormValue("name")
	id := r.PostFormValue("id")
	if name == "" || id == "" {
		w.Write([]byte("invalid param!"))
	} else {
		_, ok := users[id]
		if ok {
			w.Write([]byte("user already exists"))
		} else {
			users[id] = &User{Id: id, Name: name}

			body, err := json.Marshal(*users[id])
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Print(users[id])
			}
			w.Write(body)
		}
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	id := r.PostFormValue("id")
	name := r.PostFormValue("name")
	if id == "" || name == "" {
		w.Write([]byte("invalid param!"))
	} else {
		_, ok := users[id]
		if ok {
			users[id].Name = name
			body, _ := json.Marshal(*users[id])
			w.Write(body)
		} else {
			w.Write([]byte("user does not exists"))
		}
	}
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	ret, _ := json.Marshal(users)
	w.Write(ret)
}
