package gonoc

import (
	"database/sql"
	"fmt"
	_ "github.com/Go-SQL-Driver/MySQL"
	"reflect"
	"strings"
)

var Db *MysqlDB

type DB struct {
	sql.DB
	Debug      bool
	sqlCountInfo    string
	sqlInfo    string
	sqlWhere   string
	sqlOrderBy string
	sqlGroupBy string
	sqlLimit   string
	sqlField   string
}

func (db *DB) errorInfo(err string) {
	fmt.Println(err)
}

func (db *DB) ShowSql() string {
	return db.sqlInfo
}

func (db *DB) cleanSql() {
	db.sqlCountInfo = ""
	db.sqlInfo = ""
	db.sqlWhere = ""
	db.sqlOrderBy = ""
	db.sqlGroupBy = ""
	db.sqlLimit = ""
	db.sqlField = ""
}

func (db *DB) CloseSql() {
	db.DB.Close()
}

func (db *DB) SetSql(sqlInfo string) {
	db.sqlInfo = sqlInfo
}

//field
func (db *DB) Field(field string) {
	if field == "" {
		db.sqlField = "*"
	} else {
		db.sqlField = field
	}
}

//group by
func (db *DB) GroupBy(field string) {
	if field == "" {
		db.errorInfo("format the error")
	} else {
		var group string
		if !strings.Contains(db.sqlGroupBy, "GROUP BY") {
			group = "GROUP BY"
		}
		db.sqlGroupBy = fmt.Sprintf("%s %s %s",
			db.sqlGroupBy,
			group,
			field,
		)
	}
}

//order by
func (db *DB) OrderBy(field, value string) {
	if field == "" || value == "" {
		db.errorInfo("format the error")
	} else {
		var order string
		if !strings.Contains(db.sqlOrderBy, "ORDER BY") {
			order = "ORDER BY"
		}
		db.sqlOrderBy = fmt.Sprintf("%s, %s `%s` %s",
			db.sqlOrderBy,
			order,
			field,
			value,
		)
	}
}
func (db *DB) MultiOrderBy(order [][]string) {
	for _, value := range order {
		db.OrderBy(value[0], value[1])
	}
	db.sqlOrderBy = Substr(db.sqlOrderBy, 1, len(db.sqlOrderBy))
}

//limit
func (db *DB) Limit(start, offset int) {
	if start == 0 && offset == 0 {
		db.sqlLimit = ""
	} else {
		if offset == 0 {
			db.sqlLimit = fmt.Sprintf("LIMIT %d",
				start,
			)
		} else {
			db.sqlLimit = fmt.Sprintf("LIMIT %d, %d",
				start,
				offset,
			)
		}
	}
}

//where
func (db *DB) where(operate, in, field, value string) {
	if field == "" || operate == "" || value == "" {
		db.errorInfo("format the error")
	} else {
		var where string
		if !strings.Contains(db.sqlWhere, "WHERE") {
			where = "WHERE"
			operate = ""
		}
		db.sqlWhere = fmt.Sprintf("%s %s %s %s %s \"%s\"",
			db.sqlWhere,
			where,
			operate,
			field,
			in,
			value,
		)
	}
}

func (db *DB) AndWhere(in, field, value string) {
	db.where("AND", in, field, value)
}

func (db *DB) OrWhere(in, field, value string) {
	db.where("OR", in, field, value)
}

//select
func (db *DB) Select(tableName string) {
	if db.sqlField == "" {
		db.sqlField = "*"
	}
	db.sqlCountInfo = fmt.Sprintf("SELECT count(*) from %s %s %s",
		tableName,
		db.sqlWhere,
		db.sqlGroupBy,
	)
	db.sqlInfo = fmt.Sprintf("SELECT %s from %s %s %s %s %s",
		db.sqlField,
		tableName,
		db.sqlWhere,
		db.sqlGroupBy,
		db.sqlOrderBy,
		db.sqlLimit,
	)
}

//count
func (db *DB) Count() (count int64, err error) {
    row := db.DB.QueryRow(db.sqlCountInfo)
    err = row.Scan(&count)
    return
}

//query
func (db *DB) Query() ([]map[string]string, error) {
	rows, err := db.DB.Query(db.sqlInfo)
	rs := db.GetRowsScan(rows)
	db.cleanSql()
	return rs, err
}

func (db *DB) QueryRow() *sql.Row {
	rows := db.DB.QueryRow(db.sqlInfo)
	db.cleanSql()
	return rows
}

//save
func (db *DB) Save(info interface{}) {
	v := reflect.ValueOf(info).Elem()
	typeOfV := v.Type()
	var tableName string = typeOfV.Name()

	var primaryKey string
	var primaryValue interface{}
	var primaryI int
	var insertFields, insertValues, updateFields string
	for i := 0; i < v.NumField(); i++ {
		tag := reflect.ValueOf(typeOfV.Field(i).Tag).String()
		if tag == "PK" {
			primaryKey = typeOfV.Field(i).Name
			primaryValue = v.Field(i).Interface()
			primaryI = i
		} else {
			insertFields = fmt.Sprintf("%s, %s",
				insertFields,
				typeOfV.Field(i).Name,
			)
			insertValues = fmt.Sprintf("%s, '%v'",
				insertValues,
				v.Field(i).Interface(),
			)
			updateFields = fmt.Sprintf("%s, %s = '%v'",
				updateFields,
				typeOfV.Field(i).Name,
				v.Field(i).Interface(),
			)
		}
	}
	if primaryValue == 0 {
		//insert
		db.sqlInfo = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
			tableName,
			Substr(insertFields, 1, len(insertFields)),
			Substr(insertValues, 1, len(insertValues)),
		)
		result, _ := db.DB.Exec(db.sqlInfo)
		ids, _ := result.LastInsertId()
		v.Field(primaryI).SetInt(ids)
	} else {
		db.sqlInfo = fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%v'",
			tableName,
			Substr(updateFields, 1, len(updateFields)),
			primaryKey,
			primaryValue,
		)
		db.DB.Exec(db.sqlInfo)
	}
}

//delete
func (db *DB) Delete(info interface{}) {

}

func (db *DB) SelectStruct(model interface{}, header []string, where [][]string, order [][]string, limit []int) {
	models := reflect.ValueOf(model)
	modelElem := models.Elem()
	typeOfV := modelElem.Type()
	var tableName string = db.FieldProcessing(typeOfV.Name())

	var field, leftJoin string
	headerMap := models.MethodByName("GetHeaderMap").Call(nil)[0].Interface().(map[string]map[string]string)
	has := models.MethodByName("GetHas").Call(nil)[0].Interface().(map[string]map[string]string)
	if len(header) == 0 {
		field = " " + tableName + ".*"
	} else {
		for _, mapKey := range header {
			//header是否在struct中设置，没有设置不生成field和left join字段
			fieldByName := modelElem.FieldByName(mapKey)
			if fieldByName.IsValid() {
				_, mapBool := headerMap[mapKey]
				if mapBool {
					for hmTable, hmField := range headerMap[mapKey] {
						//生成显示字段
						field = fmt.Sprintf("%s, `%s`.`%s` AS `%s`",
							field,
							db.FieldProcessing(hmTable),
							db.FieldProcessing(hmField),
							mapKey,
						)
						//生成left join
						hasValue, hasBool := has[hmTable]
						if hasBool {
							for hasId, hasIds := range hasValue {
								if strings.Count(leftJoin, db.FieldProcessing(hmTable)) == 0 {
									leftJoin = fmt.Sprintf("%s LEFT JOIN `%s` ON `%s`.`%s` = `%s`.`%s`",
										leftJoin,
										db.FieldProcessing(hmTable),
										tableName,
										db.FieldProcessing(hasIds),
										db.FieldProcessing(hmTable),
										db.FieldProcessing(hasId),
									)
								}
							}	
						}	
					}
				} else {
					//生成显示字段
					field = fmt.Sprintf("%s, `%s`.`%s` AS `%s`",
						field,
						tableName,
						db.FieldProcessing(mapKey),
						db.FieldProcessing(mapKey),
					)
				}
			}	
		}
		field = Substr(field, 1, len(field))
	}
	leftJoin = Substr(leftJoin, 1, len(leftJoin))
	//获取where条件
	db.getWhereSql(model, where, tableName, header, headerMap)
	db.MultiOrderBy(order)
	db.Limit(limit[0], limit[1])
	db.sqlCountInfo = fmt.Sprintf("SELECT count(*) from `%s` %s %s",
		tableName,
		leftJoin,
		db.sqlWhere,
	)
	db.sqlInfo = fmt.Sprintf("SELECT%s FROM `%s` %s %s %s %s",
		field,
		tableName,
		leftJoin,
		db.sqlWhere,
		db.sqlOrderBy,
		db.sqlLimit,
	)
}

func (db *DB) getWhereSql(model interface{}, where [][]string, tableName string, header []string, headerMap map[string]map[string]string) {
	models := reflect.ValueOf(model)
	modelElem := models.Elem()
	//获取where条件
	var whereSql string
	for _, whereValue := range where {
		if len(whereValue) < 2 {
			return
		}
		//header是否在struct中设置，没有设置不生成field和left join字段
		fieldByName := modelElem.FieldByName(whereValue[0])
		if fieldByName.IsValid() {
			eq := "="
			and := "AND"
			fieldValue := "'" + whereValue[1] + "'"
			if len(whereValue) >= 3 {
				eq = whereValue[2]
				switch whereValue[2]{
					case "like":
						fieldValue = "'%" + whereValue[1] + "%'"
					case "in":
						fieldValue = "(" + whereValue[1] + ")"
					default:
						fieldValue = fieldValue
				}
			}
			if len(whereValue) == 4 {
				and = strings.ToUpper(whereValue[3])
			}
			//判断是否本表字段
			whereTableName := tableName
			whereField := whereValue[0]
			_, whereMapBool := headerMap[whereValue[0]]
			if whereMapBool {
				for whereHmTable, whereHmField := range headerMap[whereValue[0]] {
					whereTableName = whereHmTable
					whereField = whereHmField
				}
			}

			whereSql = fmt.Sprintf("%s %s (`%s`.`%s` %s %s)",
				whereSql,
				and,
				db.FieldProcessing(whereTableName),
				db.FieldProcessing(whereField),
				eq,
				fieldValue,
			)
		}	
	}
	whereSql = Substr(whereSql, 5, len(whereSql))
	if whereSql != "" {
		whereSql = "WHERE " + whereSql
	}
	db.sqlWhere = whereSql
}

func (db *DB) FieldProcessing(field string) string {
	return strings.ToLower(field)
}

func (db *DB) GetRowsScan(rows *sql.Rows) []map[string]string {
	columns, err := rows.Columns()
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
    for i := range values {
        scanArgs[i] = &values[i]
    } 
    var array = make([]map[string]string, 0)
    for rows.Next() {
    	var rs = make(map[string]string)
		err = rows.Scan(scanArgs...)
        if err != nil {
            panic(err.Error())
        }

        for i, col := range values {
        	var value string
        	if col == nil {
                value = ""
            } else {
                value = string(col)
            }
            rs[columns[i]] = value
        }
        array = append(array, rs)
	}
	return array
}

func (db *DB) GetOneScan(rows *sql.Rows) map[string]interface{} {
	columns, err := rows.Columns()
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
    for i := range values {
        scanArgs[i] = &values[i]
    }
    var rs = make(map[string]interface{})
    for rows.Next() {
		err = rows.Scan(scanArgs...)
        if err != nil {
            panic(err.Error())
        }

        for i, col := range values {
        	var value string
        	if col == nil {
                value = "NULL"
            } else {
                value = string(col)
            }
            rs[columns[i]] = value
        }
	}
	return rs
}

func (db *DB) GetValue(item *map[string]string, name string) string {
	items := *item
	return items[name]
}

type MysqlDB struct {
	DB
}

func Open(databaseType, connConfig string) {
	dbs, err := sql.Open(databaseType, connConfig)
	if err != nil {
		return
	}
	db := &MysqlDB{}
	db.DB.DB = *dbs
	Db = db
}
