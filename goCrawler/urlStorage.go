package main

type urlStorage interface {
	getOneURL(k string) string
	addOneURL(k, v string)
}
