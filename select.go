package db

import (
	"reflect"
	"strings"
)

type Condition struct {
	Key   string
	Value string
	Op    string
}

type Selecter struct {
	selects    string
	from       string
	conditions []Condition
	orderBy    string
	groupBy    string
	limit      string
}

func Select(fields ...string) *Selecter {
	var obj = new(Selecter)
	if len(fields) == 0 {
		obj.selects = "SELECT *"
	} else {
		obj.selects = "SELECT " + strings.Join(fields, ", ")
	}
	return obj
}

func (c *Selecter) From(table string) *Selecter {
	c.from = " FROM " + table
	return c
}

// val only supports int, int64, string type
func (c *Selecter) Where(key string, val interface{}, op ...string) *Selecter {
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
		v = ToString(tmp)
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

func (c *Selecter) OrderBy(col string, rank string) *Selecter {
	c.orderBy = " ORDER BY " + col + " " + strings.ToUpper(rank)
	return c
}

func (c *Selecter) GroupBy(col string) *Selecter {
	c.groupBy = " GROUP BY " + col
	return c
}

func (c *Selecter) Limit(nums ...int64) *Selecter {
	var arr = make([]string, 0)
	for _, item := range nums {
		arr = append(arr, ToString(item))
	}
	c.limit = " LIMIT " + strings.Join(arr, ", ")
	return c
}

func (c *Selecter) SQL() string {
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
	return c.selects + c.from + where + c.orderBy + c.groupBy + c.limit
}

/**
@param result: 查询结果
@param exist: 是否存在
 */
func (c *Selecter) Fetch() (result QueryResult, err error) {
	c.Limit(1)
	sql := c.SQL()
	return Fetch(sql)
}

func (c *Selecter) FetchAll() (results []QueryResult, err error) {
	sql := c.SQL()
	return FetchAll(sql)
}
