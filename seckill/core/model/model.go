package model

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"reflect"
	"time"
)

// type BaseModel struct {
// 	Id          uint `gorm:"primary_key"`
// 	Create_time time.Time
// 	time.Time
// 	DeletedAt *time.Time
// }

type Model struct {
	db          *sql.DB
	tableName   string
	valueMap    map[string]string
	pk          string
	selectField string
	where       string
	orderBy     string
	join        string
	limit       string
}

type Product struct {
	Id          int
	Name        string `orm:"rel(fk)"`
	Stock       int
	End_time    string
	Start_time  string
	Create_time string
}

type Order struct {
	Prod_id     int
	User_phone  string
	State       int8
	Create_time time.Time
}

func NewDb() *sql.DB {
	db, err := sql.Open("mysql", "root:123456@/seckill")
	if err != nil {
		log.Fatalf("open mysql error: %s\n", err.Error())
	}
	defer db.Close()
	return db
}

//raw sql
//func (this *Model) {
//
//}


//初始化model对象
func NewModel(m interface{}) {

}

func Insert(db *sql.DB, datas interface{}) {
	t := reflect.TypeOf(datas)
	fmt.Println("Type : ", t.Name())
	v := reflect.ValueOf(datas)
	fmt.Println("Fields : ")
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		val := v.Field(i).Interface()
		tag := f.Tag.Get("orm")
		fmt.Printf("%6s : %v = %v Tag : %v\n", f.Name, f.Type, val, tag)
	}

	stmt, err := db.Prepare("INSERT INTO product(name, stock, start_time, end_time) VALUES(?, ?, ?, ?)")
	defer stmt.Close()
	if err != nil {
		log.Fatalln(err.Error())
	}
	stmt.Exec(datas)
}
