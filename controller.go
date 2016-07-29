package gonoc

import (
	//"fmt"
	"reflect"
	"encoding/json"
	"strconv"
)

var RouterMap map[string]interface{}
var Requests Request

type Controller struct {
}

func (ci *Controller) ResApi(model interface{}) (rs []map[string]string, err error) {
	var header []string
	json.Unmarshal([]byte(Requests.Get("header")), &header)
	var where [][]string
	json.Unmarshal([]byte(Requests.Get("where")), &where)
	var order [][]string
	json.Unmarshal([]byte(Requests.Get("order")), &order)
	//header := []string{"Id", "Sn", "Hostname", "IdcName", "IdcCname", "Event", "Status", "Status_id"}
	//where := [][]string{{"Sn" , "10", "like", "and"}, {"Event", "2016Q2"}, {"Status", "上线"}}
	//where := [][]string{{"Rack", "NT015F2-Y9-11-6"}}
	//order := [][]string{{"Idc_id", "DESC"},{"Sn", "ASC"}}
	start,_ := strconv.Atoi(Requests.Get("start"))
	offset,_ :=  strconv.Atoi(Requests.Get("offset"))
	limit := []int{start, offset}
	models := reflect.ValueOf(model)

	Before := models.MethodByName("BeforeConditions")
	if Before.IsValid() {
		BeforeInput := make([]reflect.Value, 1)
		BeforeInput[0] = reflect.ValueOf(&where)
		Before.Call(BeforeInput)
	}
	Db.SelectStruct(model, header, where, order, limit)
	rs, err = Db.Query()
	if err == nil {
		//反射models中的函数，对部分字段进行转换
		After := models.MethodByName("AfterQuery")
		if After.IsValid() {
			for _, row := range rs {
				AfterInput := make([]reflect.Value, 2)
				AfterInput[0] = reflect.ValueOf(row)
				AfterInput[1] = reflect.ValueOf(header)
				go After.Call(AfterInput)
			}
		}
	}
	return
}
