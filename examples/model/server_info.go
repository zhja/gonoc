package model

import (
	//"fmt"
	"github.com/zhja/gonoc"
	//"reflect"
)

type Server_info struct {
	Id int `PK`
	Sn string
	Hostname string
	Asset_number string
	Idc_id int
	Event_id int
	IdcName string
	IdcCname string
	Event string
	Status_id int
}

func (art *Server_info) GetHeaderMap() map[string]map[string]string {
	hm := make(map[string]map[string]string)
	hm["IdcName"] = map[string]string{"Idc_info" : "name"}
	hm["IdcCname"] = map[string]string{"Idc_info" : "cname"}
	hm["Event"] = map[string]string{"Purchase_event" : "name"}
	return hm

}

func (art *Server_info) GetHas() map[string]map[string]string {
	has := make(map[string]map[string]string)
	has["Idc_info"] = map[string]string{"id" : "idc_id"}
	has["Purchase_event"] = map[string]string{"id" : "event_id"}
	return has
}

func (art *Server_info) BeforeConditions(where *[][]string) {
	if t, keys := gonoc.ExistsSVT(*where, "Rack"); t {
		gonoc.Db.AndWhere("=", "name", (*where)[keys][1])
		gonoc.Db.Field("id")
		gonoc.Db.Select("idc_line")
		rows := gonoc.Db.QueryRow()
		var id string
		rows.Scan(&id);
		*where = append(*where, []string{"Id", id})
	}
}

func (art *Server_info) AfterQuery(val interface{}, header interface{}) {
	vals := val.(map[string]string)
	headers := header.([]string)
	if t := gonoc.ExistsSV(headers, "Status"); t {
		gonoc.Db.AndWhere("=", "id", vals["status_id"])
		gonoc.Db.Field("name")
		gonoc.Db.Select("dict_status")
		rows := gonoc.Db.QueryRow()
		var name string
		rows.Scan(&name);
		vals["Status"] = name
	}
}