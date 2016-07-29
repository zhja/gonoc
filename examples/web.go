package main

import (
	"github.com/zhja/gonoc"
	"./router"
)

func main() {
	router.Init()
	gonoc.Open("mysql", "xxx:xxxxxx@tcp(xxx.xxx.xxx.xxx)/xxxxxx?charset=utf8")
    gonoc.Run("/usr/local/mygo/static", "/static/", ":80")
}