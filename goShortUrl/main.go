package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

const (
	dsn      = "root:123456@/shorturl?charset=utf8"
	serv     = ":9999"
	frontEnd = "http://t.sheaned.com"
)

type urls struct {
	Id        int `xorm:"pk"`
	UserId    int
	ShortUrl  string
	LongUrl   string
	Active    int
	Count     int64
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
}

type Res struct {
	Status  int           `json:"status"`
	Message string        `json:"message"`
	Datas   []interface{} `json:"datas"`
}

type myAppHandler struct {
	handler http.Handler
}

func (this *myAppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	this.handler.ServeHTTP(w, r)
	w.Header().Set("Server", "Goginx-1.0")
	if r.RequestURI == "/favicon.ico" {
		return
	}
}

func Router(w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI
	if uri == "/" {
		// http.Redirect(w, r, frontEnd, 302)
	} else if uri == "/makeShortUrl" {
		MakeShortUrl(w, r)
	} else {
		GetOriginUrl(w, r, uri)
	}
}

func MakeShortUrl(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	longUrl := r.PostFormValue("longUrl")
	shortUrl := ShortenURL(longUrl)
	newUrl := &urls{
		ShortUrl: shortUrl[0],
		LongUrl:  longUrl,
		Active:   1,
	}
	affectNum, err := engine.Insert(newUrl)
	res := Res{
		Status:  0,
		Message: "fail to make a short url",
	}
	if err != nil || affectNum == 0 {

	} else {
		res.Status = 1
		res.Message = "success to make a short url"
		res.Datas = append(res.Datas, shortUrl)
	}
	body, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	w.Write(body)
}

func GetOriginUrl(w http.ResponseWriter, r *http.Request, uri string) {
	// w.Header().Set("Location", "http://baidu.com")
	res := Res{
		Status:  0,
		Message: "fail to get url",
	}
	if r.Method == "GET" {
		urlS := &urls{}
		ok, err := engine.Alias("t").Where("t.short_url = ?", uri).Get(urlS)
		if err == nil && ok {
			res.Status = 1
			res.Message = "success to get url"
			res.Datas = append(res.Datas, urlS.LongUrl)

			urlS.Count++
			_, err := engine.Id(urlS.Id).Cols("count").Update(urlS)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
	body, err := json.Marshal(res)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	w.Write(body)
}

var engine *xorm.Engine

func main() {
	var err error
	engine, err = xorm.NewEngine("mysql", dsn)
	checkError(err)
	myh := &myAppHandler{http.HandlerFunc(Router)}
	checkError(http.ListenAndServe(serv, myh))
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error: ", err.Error())
		os.Exit(1)
	}
}
