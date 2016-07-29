package controller

import (
	"fmt"
	//"reflect"
	"github.com/zhja/gonoc"
	//"net/http"
)

type RackController struct {
	gonoc.Controller
}

type idc_line struct {
	Id int `PK`
	Name string
	Idc_id int
	Col_id int
	Row_id int
	Room_id int
	Ip string
	Status_id int
	Type_id int
	Modify_time string
}

func (this *RackController) Add() {
	var rack idc_line
	rack.Name = "test-rack"
	rack.Idc_id = 8
	rack.Col_id = 8
	rack.Row_id = 8
	rack.Room_id = 8
	rack.Ip = "127.0.0.1"
	rack.Status_id = 8
	rack.Type_id = 8

	rack.Modify_time = "2016-07-29"

	gonoc.Db.Save(&rack)
	fmt.Println(rack.Id)
}

func (this *RackController) Edit() {
	var rack idc_line
	rack.Id = 36319
	rack.Name = "test-rack-edit"
	gonoc.Db.Save(&rack)
	fmt.Println(rack.Id)
}

func (this *RackController) List() {
	gonoc.Db.AndWhere("=", "status_id", "1")
	gonoc.Db.AndWhere("=", "idc_id", "26")
	gonoc.Db.OrderBy("id", "desc")
	gonoc.Db.OrderBy("idc_id", "desc")
	gonoc.Db.GroupBy("id")
	gonoc.Db.Limit(5, 0)
	gonoc.Db.Field("id, name, ip")
	gonoc.Db.Select("idc_line")
	//fmt.Println(gonoc.Db.ShowSql())
	//fmt.Println(gonoc.Db.ShowCountSql())
	count, _ := gonoc.Db.Count()
	fmt.Println(count)
	rows, _ := gonoc.Db.Query()
	fmt.Println(rows)
	
}