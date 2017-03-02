## gonoc

仿照noc中resapi实现的go web框架。

### get

```php
go get github.com/zhja/gonoc
```

### server

```php
router.Init()
gonoc.Open("mysql","xxx:xxxxxx@tcp(xxx.xxx.xxx.xxx:8000)/xxx?charset=utf8")
gonoc.Run("/usr/local/mygo/static", "/static/", ":80")
```

router.Init() 注册路由

gonoc.Open() 链接数据库

gonoc.Run() 启动web服务器，静态路由路径、根目录、监听域名和端口

### router

./router/router.go 中配置。

举例:
1./index/test
```php
gonoc.RegisterRouter("Index", &controller.IndexController{})
```
2./api/list
```php
gonoc.RegisterRouter("Api", &controller.ApiController{})
```
### controller
./controller/ 中，想要正常访问controller需要再router中注册。
每个controller为一个结构体。同时gonoc.controller提供了类似于resapi的功能，后面会详细讲解，结构体只需要继承就可以正常使用。
```php
type ApiController struct {
	gonoc.Controller
}
```
### db
链接数据库后会生成一个实例，可以通过gonoc.Db访问。
#### 函数介绍
AndWhere 		设置where条件
OrderBy 		设置排序
GroupBy 		设置group
Limit 			设置 start， offset
Field 			设置获取字段
Select 			生成sql语句
Count 			获取总条数
Query 			获取数据（多条）
QueryRow 		获取单条
Save 			保存和修改数据
Delete 			删除数据
SetSql 			直接执行sql
SelectStruct 		支持结构体获取（主要用于resapi）
ShowSql 		查看生成的sql语句
ShowCountSql 	查看生成的count sql语句
#### 添加数据
```php
type idc_line struct {
	Id int PK
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
```
#### 修改数据
```php
var rack idc_line
rack.Id = 36319
rack.Name = "test-rack-edit"
gonoc.Db.Save(&rack)
fmt.Println(rack.Id)
```
#### 查询数据
```php
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
```
##### 结果
```php
[map[id:36312 name:SHPBS012FD05-G-16-19 ip:112.65.119.95] map[id:36311 name:SHPBS012FD05-G-15-20 ip:112.65.119.94] map[id:36310 name:SHPBS012FD05-G-16-20 ip:112.65.119.94] map[id:28455 name:SHPBS012FD05-F-13-13 ip:10.111.0.60] map[id:28454 name:SHPBS012FD05-F-13-12 ip:10.111.0.59]]
```
### model
./model/ 中，每个model是一个结构体。

主要有五部分组成：
结构体字段 ：获取数据时指定的字段是否有效。
GetHeaderMap ：设置非本表字段的对应关系。

```php
hm["IdcName"] = map[string]string{"Idc_info" : "name"}
```
转为：
```php
"idc_info.name as idcname"
```
GetHas ：设置非本表字段的关联关系。
```php
has["Idc_info"] = map[string]string{"id" : "idc_id"}
```
转为：
```php
"LEFT JOIN idc_info ON server_info.idc_id= idc_info.id"
```
BeforeConditions ：query前执行，主要用于处理search条件
AfterQuery ： query后执行，处理获取的数据
 ### resapi
举例：获取服务器数据
./model/server_info.go中配置
```php
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
```
发送请求：
```php
<script language="javascript" type="text/javascript">
var header = ["Id", "Sn", "Hostname", "IdcName", "IdcCname", "Event", "Status", "Status_id"];
var where = [["Sn" , "10", "like", "and"],["Event", "2016Q2"],["Status", "上线"]];
var order = [["Idc_id", "DESC"],["Sn", "ASC"]];
header = JSON.stringify(header);
where = JSON.stringify(where);
order = JSON.stringify(order);
var data = {
	header : header,
	where : where,
	order : order,
	start : 10,
	offset : 0
};
$.ajax({
   type: "POST",
   url: "http://192.168.210.130/api/test",
   data: data,
   success: function(msg){
     //alert(1);
   }
});
</script>
```
结果：
```php
[
    {
        "Event": "2016Q2",
        "IdcCname": "娄底联通",
        "IdcName": "LDUN01",
        "hostname": "TMP-6100450500671961",
        "id": "31437",
        "sn": "6100450500671961",
        "status_id": "13"
    },
    {
        "Event": "2016Q2",
        "IdcCname": "娄底联通",
        "IdcName": "LDUN01",
        "hostname": "TMP-6100450500672335",
        "id": "31438",
        "sn": "6100450500672335",
        "status_id": "13"
    }
]
```
### template
暂时使用 "html/template"
### request
http参数获取，支持get和post请求。
```php
gonoc.Requests.Get("header")
```
### 问题及后续工作
1.  性能优化
2.  支持rule
3.  日志系统，自动保存和生成日志文件
4.  错误处理
5.  支持事务
6.  request和resapi所有数据转为string
7.  nosql支持
8.  noc中model配置转为gonoc版本工具
9.  类型转换功能
