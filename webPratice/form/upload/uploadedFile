package main

import (
	"encoding/json"
	// "fmt"
	"github.com/bmizerany/pat"
	"log"
	"net/http"
)

var users map[string]*User

type User struct {
	Id   string
	Name string
}

type Res struct {
	Status  int
	Message string
	Datas   []User
}

func main() {
	users = make(map[string]*User)
	mux := pat.New()

	// a := make(map[string]*User)
	// a["a"] = &User{"1", "hah"}
	// fmt.Println(a)
	// a["a"].Name = "ffff"
	// fmt.Println(*a["a"])
	// return

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
	ret := Res{Status: 0, Message: "failed", Datas: []User{}}
	for _, v := range users {
		if v.Name == name {
			ret.Datas = append(ret.Datas, *v)
			ret.Status = 1
			ret.Message = "successed"
			break
		}
	}
	body, _ := json.Marshal(ret)
	w.Write(body)
}

func addUser(w http.ResponseWriter, r *http.Request) {
	res := Res{Status: 0, Message: "failed", Datas: []User{}}
	name := r.PostFormValue("name")
	id := r.PostFormValue("id")
	if name == "" || id == "" {
		res.Message = "invalid param"
		// w.Write([]byte("invalid param!"))
	} else {
		_, ok := users[id]
		if ok {
			res.Message = "user already exists"
		} else {
			users[id] = &User{Id: id, Name: name}
			res.Datas = append(res.Datas, *users[id])
			res.Status = 1
			res.Message = "successed"
		}
	}
	body, _ := json.Marshal(res)
	w.Write(body)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	res := Res{Status: 0, Message: "failed", Datas: []User{}}
	id := r.PostFormValue("id")
	name := r.PostFormValue("name")
	if id == "" || name == "" {
		res.Message = "invalid param!"
	} else {
		_, ok := users[id]
		if ok {
			users[id].Name = name
			res.Status = 1
			res.Message = "successed"
			res.Datas = append(res.Datas, *users[id])
		} else {
			res.Message = "user does not exists"
		}
	}
	body, _ := json.Marshal(res)
	w.Write(body)
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	res := Res{Status: 0, Message: "failed", Datas: []User{}}
	if len(users) > 0 {
		res.Status = 1
		res.Message = "successed"
	}
	for _, v := range users {
		res.Datas = append(res.Datas, *v)
	}
	body, _ := json.Marshal(res)
	w.Write(body)
}
