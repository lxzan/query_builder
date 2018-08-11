package db

import (
	"strings"
	"reflect"
	"database/sql"
)

type Deleter struct {
	delete     string
	conditions []Condition
	limit      string
}

func Delete(table string) *Deleter {
	var obj = new(Deleter)
	obj.delete = "DELETE FROM " + table
	return obj
}

// val only supports int, int64, string type
func (c *Deleter) Where(key string, val interface{}, op ...string) *Deleter {
	var operator = "="
	if len(op) > 0 {
		operator = op[0]
	}
	var v = ""
	var valType = reflect.TypeOf(val).Name()
	if valType == "string" {
		v = reflect.ValueOf(val).String()
	} else if valType == "int" || valType == "int64" {
		tmp := reflect.ValueOf(val).Int()
		v = Interface2String(tmp)
	} else {
		panic("val only supports int, int64, string type")
	}
	c.conditions = append(c.conditions, Condition{
		Key:   key,
		Value: v,
		Op:    operator,
	})
	return c
}

func (c *Deleter) Limit(nums ...int64) *Deleter {
	var arr = make([]string, 0)
	for _, item := range nums {
		arr = append(arr, ToString(item))
	}
	c.limit = " LIMIT " + strings.Join(arr, ", ")
	return c
}

func (c *Deleter) SQL() string {
	wheres := make([]string, 0)
	for _, item := range c.conditions {
		var ele = Build("{key} {op} '{val}'", Form{
			"key": item.Key,
			"op":  item.Op,
			"val": AddSlashes(item.Value),
		})
		wheres = append(wheres, ele)
	}
	var where = ""
	if len(c.conditions) > 0 {
		where = " WHERE " + strings.Join(wheres, " AND ")
	}
	return c.delete + where + c.limit
}

func (c *Deleter) Exec() (sql.Result, error) {
	return Exec(c.SQL())
}
