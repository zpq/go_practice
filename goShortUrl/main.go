package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"io/ioutil"

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
	if r.RequestURI == "/favicon.ico" {
		return
	}
	token := r.Header.Get("authorization")
	if token != "" {
		var flag bool
		res := &Res{400, "invalid authorization", nil}
		if tokens := strings.Split(token, " "); len(tokens) == 2 {
			cliam, err := checkToken(tokens[1])
			if err != nil {
				fmt.Println(err.Error())
				flag = true
			} else {
				expired := cliam["exp"].(time.Time)
				if time.Now().Sub(expired) > 0 { //expired
					flag = true
				}
			}
		} else {
			flag = true
		}
		if flag {
			body, err := json.Marshal(res)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			w.Write(body)
			return
		}
	}
	this.handler.ServeHTTP(w, r)
	w.Header().Set("Server", "Goginx1.0")
}

func Router(w http.ResponseWriter, r *http.Request) {
	uri := r.RequestURI
	if uri == "/" {
		Index(w, r)
	} else if uri == "/user/login" {

	} else if uri == "/user/register" {

	} else if uri == "/api/url/shorten" {
		MakeShortUrlApi(w, r)
	} else if uri == "/api/user/doLogin" {
		UserLogin(w, r)
	} else if uri == "/api/user/doRegister" {
		UserRegister(w, r)
	} else {
		GetOriginUrl(w, r, uri)
	}
}

func MakeShortUrlApi(w http.ResponseWriter, r *http.Request) {
	// r.Header.Set("Access-Control-Allow-Origin", "*")
	res := Res{
		Status:  0,
		Message: "fail to make a short url",
	}
	r.ParseForm()
	longUrl := r.PostFormValue("longUrl")
	if longUrl != "" {
		longUrl = strings.Trim(longUrl, "/")
		var prefix string
		if !strings.HasPrefix(longUrl, "http://") && !strings.HasPrefix(longUrl, "https://") {
			prefix = "http://"
		}
		if prefix != "" {
			longUrl = prefix + longUrl
		}

		//todo test longUrl (get request)
		if checkUrl(longUrl) {

			shortUrl := ShortenURL(longUrl)[rand.Intn(4)]
			newUrl := &urls{
				ShortUrl: "/" + shortUrl,
				LongUrl:  longUrl,
				Active:   1,
			}
			token := r.Header.Get("authorization")
			if token != "" {
				tokens := strings.Split(token, " ")[1]
				claims, _ := checkToken(tokens)
				username := claims["username"].(string)
				tUser := new(User)
				ok, err := engine.Alias("t").Where("t.username = ?", username).Get(tUser)
				if ok && err == nil {
					newUrl.UserId = tUser.Id
				}
			}

			affectNum, err := engine.Insert(newUrl)
			if err != nil || affectNum == 0 {

			} else {
				res.Status = 1
				res.Message = "success to make a short url"
				res.Datas = append(res.Datas, frontEnd+"/"+shortUrl)
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

func GetOriginUrlAPi(w http.ResponseWriter, r *http.Request, uri string) {
	// w.Header().Set("Location", "http://baidu.com")
	r.Header.Set("Access-Control-Allow-Origin", "*")
	res := Res{
		Status:  0,
		Message: "fail to get url",
	}
	if r.Method == "GET" {
		urlS := &urls{}
		ok, err := engine.Alias("t").Where("t.short_url = ?", uri).Get(urlS)
		if err == nil && ok && urlS.Active == 1 {
			urlS.Count++
			_, err := engine.Id(urlS.Id).Cols("count").Update(urlS)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				res.Status = 1
				res.Message = "success to get url"
				res.Datas = append(res.Datas, urlS.LongUrl)
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

func GetOriginUrl(w http.ResponseWriter, r *http.Request, uri string) {
	urlS := &urls{}
	ok, err := engine.Alias("t").Where("t.short_url = ?", uri).Get(urlS)
	if err == nil && ok && urlS.Active == 1 {
		urlS.Count++
		_, err := engine.Id(urlS.Id).Cols("count").Update(urlS)
		if err != nil {
			fmt.Println(err.Error())
			http.Redirect(w, r, frontEnd, 302)
		} else {
			http.Redirect(w, r, urlS.LongUrl, 302)
		}
	} else {
		http.Redirect(w, r, frontEnd, 302)
	}
}

func DeactiveUrl(w http.ResponseWriter, r *http.Request, uri string) {
	r.Header.Set("Access-Control-Allow-Origin", "*")
	r.ParseForm()
	urlV := strings.TrimSpace(r.PostFormValue("url"))
	res := &Res{
		Status:  0,
		Message: "fail to deactive url",
	}
	if urlV != "" {
		token := r.Header.Get("authorization")
		if token != "" {
			tokens := strings.Split(token, " ")[1]
			claims, _ := checkToken(tokens)
			userid := claims["userid"].(int)
			rUrl := &urls{}
			ok, err := engine.Alias("t").Where("t.userid = ? and short_url = ?", userid, urlV).Get(rUrl)
			if ok && err == nil {
				rUrl.Active = 0
				_, err = engine.Id(rUrl.Id).Cols("active").Update(rUrl)
				if err != nil {
					fmt.Println(err.Error())
				} else {
					res.Status = 1
					res.Message = "success to deactive url"
				}
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

func checkUrl(s string) bool {
	resp, err := http.Get(s)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	_, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	// fmt.Println(string(html))
	return true
}

/**
* note: dangerous
* todo ...
 */
func DeleteUrl(w http.ResponseWriter, r *http.Request, uri string) {

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
