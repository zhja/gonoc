package router

import (
	//"fmt"
	"github.com/zhja/gonoc"
	"../controller"
)

func Init(){
	gonoc.RegisterRouter("Api", &controller.ApiController{})
	gonoc.RegisterRouter("Index", &controller.IndexController{})
	gonoc.RegisterRouter("Rack", &controller.RackController{})
}