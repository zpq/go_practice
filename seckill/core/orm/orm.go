package orm

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"reflect"
	"strconv"
	"strings"
)

type gzorm struct {
	db          *sql.DB
	tx          *sql.Tx
	txIsBegin   bool
	tablePrefix string
	tableName   string
	sqlRaw      string
	fields      string
	params      []interface{}
	from        string
	where       string
	limit       int
	offset      int
	groupBy     string
	orderBy     string
}

func connect(dsn string) *sql.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}
	return db
}

func NewOrm(dsn string) *gzorm {
	return &gzorm{db: connect(dsn)}
}

func (this *gzorm) SetTablePrefix(s string) *gzorm {
	this.tablePrefix = s
	this.from = this.tablePrefix + this.tableName
	return this
}

func (this *gzorm) RegisterModel(m interface{}) *gzorm {
	f := reflect.TypeOf(m)
	this.tableName = strings.ToLower(f.Name())
	this.from = this.tableName
	return this
}

func (this *gzorm) Fields(s string) *gzorm {
	this.fields = s
	return this
}

func (this *gzorm) Where(s string, params []interface{}) *gzorm {
	this.params = params
	this.where = s
	return this
}

func (this *gzorm) Limit(offset, limit int) *gzorm {
	this.limit = limit
	this.offset = offset
	return this
}

func (this *gzorm) GroupBy(s string) *gzorm {
	this.groupBy = s
	return this
}

func (this *gzorm) OrderBy(s string) *gzorm {
	this.orderBy = s
	return this
}

func (this *gzorm) FindOne() (map[string]interface{}, error) {
	data, err := this.query(1)
	if err != nil {
		return nil, err
	}
	return data[0], nil
}

func (this *gzorm) FindAll() ([]map[string]interface{}, error) {
	data, err := this.query(0)
	if err != nil {
		return nil, err
	}
	return data, nil
}

//raw sql
// func (this *gzorm) RunSql(s string, map[string]interface{}) (interface{}, error) {
// this.setSql(s)
// data,err := this.query(0)
// if err != nil {
// 	return nil, err
// }
// return data, nil
// }

func (this *gzorm) setSql(rawSql string) string {
	var sqlText string
	if rawSql == "" {
		if this.fields == "" {
			sqlText = "select * from " + this.from
		} else {
			sqlText = "select " + this.fields + " from " + this.from
		}
		if this.where != "" {
			sqlText += " where " + this.where
		}
		if this.groupBy != "" {
			sqlText += " group by " + this.groupBy
		}
		if this.orderBy != "" {
			sqlText += " order by " + this.orderBy
		}
		if this.limit != 0 {
			sqlText += " limit " + strconv.Itoa(this.offset) + "," + strconv.Itoa(this.limit)
		}
	} else {
		sqlText = rawSql
	}
	return sqlText
}

func (this *gzorm) query(num int) ([]map[string]interface{}, error) {
	sqlText := this.setSql("")
	rows, err := this.db.Query(sqlText, this.params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	values := make([]sql.RawBytes, len(columns))
	scans := make([]interface{}, len(columns))
	for i := range values {
		scans[i] = &values[i]
	}

	var result []map[string]interface{}
	for rows.Next() {
		rows.Scan(scans...)
		each := make(map[string]interface{})
		for i, col := range values {
			each[columns[i]] = string(col)
		}
		result = append(result, each)
		if num == 1 {
			break
		}
	}
	return result, nil
}

//checks a field whether needed or not in insert sql
func selectFieldsCheck(v reflect.Value) bool {
	switch vv := v.Interface().(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		if vv != 0 {
			return true
		}
	case string:
		if vv != "" {
			return true
		}
	case float32, float64:
		if vv != 0.0 {
			return true
		}
	default:
		return false
	}
	return false
}

func (this *gzorm) Insert(m interface{}) (int64, error) {
	insertSql := "insert into " + this.tablePrefix + this.tableName + " ("
	fieldSql := ""
	paramsSql := ""
	var insertParams []interface{}
	t := reflect.TypeOf(m)
	v := reflect.ValueOf(m)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		tagField := field.Tag.Get("field") //get customer field name if it exists
		if selectFieldsCheck(value) {
			if tagField != "" {
				fieldSql += tagField + ","
			} else {
				fieldSql += strings.ToLower(field.Name) + ","
			}
			paramsSql += "?,"
			insertParams = append(insertParams, value.Interface())
		}
	}
	fieldSql = strings.TrimRight(fieldSql, ",")
	paramsSql = strings.TrimRight(paramsSql, ",")
	insertSql += fieldSql + ") values (" + paramsSql + ")"
	log.Println(insertSql)
	if this.txIsBegin {
		stmt, err := this.tx.Prepare(insertSql)
		if err != nil {
			return 0, err
		}
		result, err := stmt.Exec(insertParams...)
		if err != nil {
			return 0, err
		}
		lastInsertId, _ := result.LastInsertId()
		return lastInsertId, nil
	} else {
		stmt, err := this.db.Prepare(insertSql)
		if err != nil {
			return 0, err
		}
		result, err := stmt.Exec(insertParams...)
		if err != nil {
			return 0, err
		}
		lastInsertId, _ := result.LastInsertId()
		return lastInsertId, nil
	}

}

//update tableName set bookname = ?, booknumber = ? where id = ?
func (this *gzorm) Update(m interface{}) (int64, error) {
	updateSql := "update " + this.tablePrefix + this.tableName + " set "
	fieldSql := ""
	t := reflect.TypeOf(m)
	v := reflect.ValueOf(m)
	var updateParams []interface{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)
		tagField := field.Tag.Get("field") //get customer field name if it exists
		if selectFieldsCheck(value) {
			if tagField != "" {
				fieldSql += tagField + " = ?,"
			} else {
				fieldSql += strings.ToLower(field.Name) + " = ?,"
			}
			updateParams = append(updateParams, value.Interface())
		}
	}
	fieldSql = strings.TrimRight(fieldSql, ",")
	updateSql += fieldSql

	// add conditions
	if this.where != "" {
		updateSql += " where " + this.where
	}
	// for _, pv := range this.params {
	// 	updateParams = append(updateParams, pv)
	// }
	updateParams = append(updateParams, this.params...)
	if this.txIsBegin {
		db := this.tx
		stmt, err := db.Prepare(updateSql)
		if err != nil {
			return 0, err
		}
		result, err := stmt.Exec(updateParams...)
		if err != nil {
			return 0, err
		}
		num, err := result.RowsAffected()
		if err != nil {
			return 0, err
		}
		return num, nil
	} else {
		db := this.db
		stmt, err := db.Prepare(updateSql)
		if err != nil {
			return 0, err
		}
		result, err := stmt.Exec(updateParams...)
		if err != nil {
			return 0, err
		}
		num, err := result.RowsAffected()
		if err != nil {
			return 0, err
		}
		return num, nil
	}

}

func (this *gzorm) Delete() (int64, error) {
	deleteSql := "DELETE FROM " + this.tablePrefix + this.tableName
	if this.where != "" {
		deleteSql += " where " + this.where
	}
	if this.txIsBegin {
		db := this.tx
		stmt, err := db.Prepare(deleteSql)
		if err != nil {
			return 0, err
		}
		res, err := stmt.Exec(this.params...)
		if err != nil {
			return 0, err
		}
		num, err := res.RowsAffected()
		if err != nil {
			return 0, err
		}
		return num, nil
	} else {
		db := this.db
		stmt, err := db.Prepare(deleteSql)
		if err != nil {
			return 0, err
		}
		res, err := stmt.Exec(this.params...)
		if err != nil {
			return 0, err
		}
		num, err := res.RowsAffected()
		if err != nil {
			return 0, err
		}
		return num, nil
	}
}

func (this *gzorm) BeginTran() error {
	tx, err := this.db.Begin()
	if err != nil {
		return err
	}
	this.tx = tx
	this.txIsBegin = true
	return nil
}

func (this *gzorm) Commit() {
	this.tx.Commit()
	// this.tx = nil
	this.txIsBegin = false
}

func (this *gzorm) Rollback() {
	this.tx.Rollback()
	// this.tx = nil
	this.txIsBegin = false
}
