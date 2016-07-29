package model

import (
	//"fmt"
	//"reflect"
)

type Article struct {
	Id int `PK`
	Name string
	User_id int
	Class_id int
	User_name string
	Class_name string
	Tag string
}

func (art *Article) GetHeaderMap() map[string]map[string]string {
	hm := make(map[string]map[string]string)
	Tag := make(map[string]string)
	UserName := make(map[string]string)
	Tag["Class"] = "tag"
	UserName["User"] = "name"
	hm["Tag"] = Tag
	hm["UserName"] = UserName
	return hm
}

func (art *Article) GetHas() map[string]map[string]string {
	has := make(map[string]map[string]string)
	User := make(map[string]string)
	Class := make(map[string]string)
	User["id"] = "user_id"
	Class["id"] = "class_id"
	has["User"] = User
	has["Class"] = Class
	return has
}

// func (art *Article) BeforeConditions() map[string]map[string]string {
// 	has := make(map[string]map[string]string)
// 	User := make(map[string]string)
// 	Class := make(map[string]string)
// 	User["id"] = "user_id"
// 	Class["id"] = "class_id"
// 	has["User"] = User
// 	has["Class"] = Class
// 	return has
// }
func (art *Article) AfterQuery() int {
	return 111;
}