package controller

import (
	"fmt"
	"github.com/zhja/gonoc"
	"../model"
	"encoding/json"
	//"strconv"
)

type ApiController struct {
	gonoc.Controller
}

func (this *ApiController) Test() {
	Server_info := model.Server_info{}
	rss, _ := this.ResApi(&Server_info)
	//b,_ := strconv.Atoi(rss[0]["status_id"])
	
	rss_json, _ := json.Marshal(rss)
	fmt.Fprintf(gonoc.Requests.W, string(rss_json))
}