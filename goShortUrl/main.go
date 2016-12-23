package main

import (
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

}

func GetOriginUrl(w http.ResponseWriter, r *http.Request, uri string) {
	// w.Header().Set("Location", "http://baidu.com")
	if r.Method == "GET" {
		urlS := &urls{}
		ok, err := engine.Alias("t").Where("t.short_url = ?", uri).Get(urlS)
		if err != nil {
			fmt.Println(err.Error())
			http.Redirect(w, r, frontEnd, 302)
		} else {
			if ok {
				urlS.Count++
				_, err := engine.Id(urlS.Id).Cols("count").Update(urlS)
				if err != nil {
					fmt.Println(err.Error())
				}
				http.Redirect(w, r, urlS.LongUrl, 302)
			} else {
				http.Redirect(w, r, "http://baidu.com", 302)
			}
		}
	} else {
		w.WriteHeader(400)
	}
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
