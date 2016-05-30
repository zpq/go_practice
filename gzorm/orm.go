package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"reflect"
	"strconv"
	"strings"
)

/*
* o := NewOrm()
* o.RegisterModel(Book{})
* where := make(map[string]string)
* where["id"] = 1
* result, err := o.Select("*").Where("id > :id", where).FindAll()
* if err != nil { log.Fatal(err.Error()) }
* log.PrintLn(result)
 */

type gzorm struct {
	db          *sql.DB
	tablePrefix string
	tableName   string
	sqlRaw      string
	fields      string
	params      []interface{}
	from        string
	where       string
	limit       int
	offset      int
	groupBy     string
	orderBy     string
}

type Book struct {
	Id            int
	Bookname      string
	Booknum       int
	Status        int8
	BookIntroduce string
}

func main() {
	o := NewOrm()
	where := []interface{}{0, "1"}
	o.RegisterModel(Book{}).SetTablePrefix("db_")
	result, err := o.Fields("*").Where("id > ? and status = ?", where).OrderBy("id desc").Limit(10, 2).FindAll()
	if err != nil {
		log.Fatal(err.Error())
	}
	for _, v := range result {
		log.Println(v)
	}
}

func connect() *sql.DB {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
	if err != nil {
		log.Fatal(err.Error())
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}
	return db
}

func NewOrm() *gzorm {
	return &gzorm{db: connect()}
}

func (this *gzorm) SetTablePrefix(s string) *gzorm {
	this.tablePrefix = s
	return this
}

func (this *gzorm) RegisterModel(m interface{}) *gzorm {
	f := reflect.TypeOf(m)
	this.tableName = strings.ToLower(f.Name())
	return this
}

func (this *gzorm) Fields(s string) *gzorm {
	this.fields = s
	this.from = this.tablePrefix + this.tableName
	return this
}

func (this *gzorm) Where(s string, params []interface{}) *gzorm {
	this.params = params
	this.where = s
	return this
}

func (this *gzorm) Limit(offset, limit int) *gzorm {
	this.limit = limit
	this.offset = offset
	return this
}

func (this *gzorm) GroupBy(s string) *gzorm {
	this.groupBy = s
	return this
}

func (this *gzorm) OrderBy(s string) *gzorm {
	this.orderBy = s
	return this
}

func (this *gzorm) FindOne() (map[string]interface{}, error) {
	data, err := this.query(1)
	if err != nil {
		return nil, err
	}
	return data[0], nil
}

func (this *gzorm) FindAll() ([]map[string]interface{}, error) {
	data, err := this.query(0)
	if err != nil {
		return nil, err
	}
	return data, nil
}

//raw sql
// func (this *gzorm) RunSql(s string, map[string]interface{}) (interface{}, error) {
// this.setSql(s)
// data,err := this.query(0)
// if err != nil {
// 	return nil, err
// }
// return data, nil
// }

func (this *gzorm) setSql(rawSql string) string {
	var sqlText string
	if rawSql == "" {
		sqlText = "select " + this.fields + " from " + this.from
		if this.where != "" {
			sqlText += " where " + this.where
		}
		if this.groupBy != "" {
			sqlText += " group by " + this.groupBy
		}
		if this.orderBy != "" {
			sqlText += " order by " + this.orderBy
		}
		if this.limit != 0 {
			sqlText += " limit " + strconv.Itoa(this.offset) + "," + strconv.Itoa(this.limit)
		}
	} else {
		sqlText = rawSql
	}
	return sqlText
}

func (this *gzorm) query(num int) ([]map[string]interface{}, error) {
	sqlText := this.setSql("")
	rows, err := this.db.Query(sqlText, this.params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	values := make([]sql.RawBytes, len(columns))
	scans := make([]interface{}, len(columns))
	for i := range values {
		scans[i] = &values[i]
	}

	var result []map[string]interface{}
	for rows.Next() {
		rows.Scan(scans...)
		each := make(map[string]interface{})
		for i, col := range values {
			each[columns[i]] = string(col)
		}
		result = append(result, each)
		if num == 1 {
			break
		}
	}
	return result, nil
}

func (this *gzorm) Insert(m interface{}) (int, error) {

}

/*




func (this *gzorm) Update(m interface{}) (int, error) {

}



func (this *gzorm) BeginTran() *gzorm {

}

func (this *gzorm) Commit() *gzorm {

}

func (this *gzorm) Rollback() *gzorm {

}
*/
