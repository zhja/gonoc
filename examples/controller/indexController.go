package controller

import (
	"fmt"
	//"reflect"
	"github.com/zhja/gonoc"
	"html/template"
	//"net/http"
)

type IndexController struct {
	gonoc.Controller
}

func (this *IndexController) Test() {
	t, _ := template.ParseFiles("view/test.html")
	err := t.Execute(gonoc.Requests.W, nil)
	if err != nil {
        fmt.Println("There was an error:", err)
    }
}