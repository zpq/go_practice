package main

import (
	"./core/model"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type SecKillServer struct {
	handler http.Handler
}

func serverError(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (this *SecKillServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Goginx")
	this.handler.ServeHTTP(w, r)
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	paths := strings.Split(path, "/")

	if paths[1] == "" || paths[1] == "index" {
		indexHandler(w, r, paths)
	} else if paths[1] == "detail" && len(paths) > 2 {
		detailHandler(w, r, paths)
	} else if paths[1] == "execute" && len(paths) > 2 {
		executeHandler(w, r, paths)
	} else {
		w.WriteHeader(404)
		loadTemplateHtml(w, "404.html", nil)
	}
}

//首页控制器
func indexHandler(w http.ResponseWriter, r *http.Request, paths []string) {
	if r.Method == "GET" {
		loadTemplateHtml(w, "index.html", nil)
	} else {
		w.WriteHeader(404)
		loadTemplateHtml(w, "404.html", nil)
	}
}

//商品详情控制器
func detailHandler(w http.ResponseWriter, r *http.Request, paths []string) {
	if r.Method == "GET" {
		loadTemplateHtml(w, "detail.html", nil)
	} else {
		w.WriteHeader(404)
		loadTemplateHtml(w, "404.html", nil)
	}
}

//秒杀控制器
func executeHandler(w http.ResponseWriter, r *http.Request, paths []string) {
	if r.Method == "GET" && paths[2] == "time" { //获取系统时间  /execute/time
		w.WriteHeader(200)
	} else if r.Method == "GET" && len(paths) > 3 && paths[2] != "" && paths[3] == "exposer" { //获取秒杀地址 /execute/{id}/exposer
		w.WriteHeader(200)
	} else if r.Method == "POST" && len(paths) > 3 && paths[2] != "" && paths[3] == "doExecute" { //执行秒杀 /execute/{id}/doExecute //POST VALUE MD5
		w.WriteHeader(200)
	} else {
		w.WriteHeader(403)
	}
}

//加载html模版
func loadTemplateHtml(w http.ResponseWriter, tempName string, data interface{}) {
	w.Header().Set("content-type", "text/html;charset=utf-8")
	t, err := template.ParseFiles(tempName)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	t.Execute(w, data)
}

func main() {
	db := model.NewDb()
	log.Println(db)
	now := time.Now().Format("2006-01-02 15:04:05")
	product := model.Product{
		Name:       "iphone6",
		Stock:      1000,
		End_time:   now,
		Start_time: now,
	}
	log.Println(product)

	model.Insert(db, product)
	s := &SecKillServer{http.HandlerFunc(myHandler)}
	http.ListenAndServe(":8008", s)
}
