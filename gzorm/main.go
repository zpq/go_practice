package mains

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type Book struct {
	Id            int
	Bookname      string
	Booknum       int
	Status        int8
	BookIntroduce string
}

func Find(db *sql.DB) {
	rows, err := db.Query("select * from db_book where id > ?", 0)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer rows.Close()

	var book []*Book
	for rows.Next() {
		tbook := Book{}
		err := rows.Scan(&tbook.Id, &tbook.Bookname, &tbook.Booknum, &tbook.Status, &tbook.BookIntroduce)
		if err != nil {
			log.Fatal(err.Error())
		}
		book = append(book, &tbook)
	}

	for _, v := range book {
		log.Println(v)
	}
}

func Insert(db *sql.DB) {
	stmt, err := db.Prepare("insert into db_book (bookname, booknumber, status, Book_introduce) values (?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer stmt.Close()

	book := Book{0, "golang", 200, 1, "golang book for you"}
	lastId, err := stmt.Exec(book.Bookname, book.Booknum, book.Status, book.BookIntroduce)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Insert Success!", lastId)
}

func main() {
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()
	// Insert(db)
	Find(db)
}
