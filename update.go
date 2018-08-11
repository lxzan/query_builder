package db

import (
	"strings"
	"reflect"
	"database/sql"
)

type Updater struct {
	update     string
	set        string
	conditions []Condition
	limit      string
}

func Update(table string) *Updater {
	var obj = new(Updater)
	obj.update = "UPDATE " + table
	return obj
}

func (c *Updater) Set(data Json) *Updater {
	var arr = make([]string, 0)
	for k, v := range data {
		item := Build(`{key} = '{val}'`, Form{
			"key": k,
			"val": Interface2String(v),
		})
		arr = append(arr, item)
	}
	c.set = " SET " + strings.Join(arr, ", ")
	return c
}

// val only supports int, int64, string type
func (c *Updater) Where(key string, val interface{}, op ...string) *Updater {
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

func (c *Updater) Limit(nums ...int64) *Updater {
	var arr = make([]string, 0)
	for _, item := range nums {
		arr = append(arr, ToString(item))
	}
	c.limit = " LIMIT " + strings.Join(arr, ", ")
	return c
}

func (c *Updater) SQL() string {
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
	return c.update + c.set + where + c.limit
}

func (c *Updater) Exec() (sql.Result, error) {
	return Exec(c.SQL())
}
