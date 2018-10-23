package db

import (
	q "database/sql"
	_ "github.com/go-sql-driver/mysql"
	"regexp"
	"sort"
	"strings"
)

var conns = make([]*q.DB, 0)

func Connect(opt *ConnectOption) *q.DB {
	dsn := Build("{user}:{password}@tcp({host}:{port})/{database}?charset={charset}", Form{
		"user":     opt.User,
		"password": opt.Password,
		"host":     opt.Host,
		"port":     opt.Port,
		"database": opt.Database,
		"charset":  opt.Charset,
	})

	if opt.MaxOpenConns == 0 {
		opt.MaxOpenConns = 200
	}
	if opt.MaxIdleConns == 0 {
		opt.MaxIdleConns = 100
	}
	conn, _ := q.Open("mysql", dsn)
	conn.SetMaxOpenConns(opt.MaxOpenConns)
	conn.SetMaxIdleConns(opt.MaxIdleConns)
	conns = append(conns, conn)
	return conn
}

func ConnectCluster(opts []*ConnectOption) {
	var count = len(opts)
	conns = make([]*q.DB, 0)
	for i := 0; i < count; i++ {
		Connect(opts[i])
	}
}

func Master() *q.DB {
	return conns[0]
}

func Slave() *q.DB {
	var count = len(conns)
	if count == 1 {
		return conns[0]
	}
	var index = Rand(1, count-1)
	return conns[index]
}

func FetchAll(sql string, bind ...Json) ([]QueryResult, error) {
	sql = BuildSQL(sql, bind...)
	var conn = Slave()
	rows, _ := conn.Query(sql)
	columns, err := rows.Columns()
	if err != nil {
		return []QueryResult{}, nil
	}

	var length = len(columns)
	var res = make([]QueryResult, 0)
	for rows.Next() {
		var fields = make([]interface{}, length)
		for i := 0; i < length; i++ {
			var tmp = ""
			fields[i] = &tmp
		}
		rows.Scan(fields...)

		doc := QueryResult{}
		for i, item := range fields {
			s, _ := item.(*string)
			doc[columns[i]] = NewAny(*s)
		}
		res = append(res, doc)
	}
	defer rows.Close()
	return res, nil
}

// if result == nil, query result is empty
func Fetch(sql string, bind ...Json) (result QueryResult, err error) {
	sql = BuildSQL(sql, bind...)
	re, _ := regexp.Compile(`(?im:limit [0-9]+$)`)
	if !re.MatchString(sql) {
		sql += " LIMIT 1"
	}

	res, err := FetchAll(sql)
	if err != nil {
		return nil, err
	}

	if len(res) > 0 {
		return res[0], nil
	}
	return nil, nil
}

func Exec(sql string, bind ...Json) (q.Result, error) {
	var client = Master()
	sql = BuildSQL(sql, bind...)
	return client.Exec(sql)
}

func BuildSQL(sql string, bind ...Json) string {
	var params = Json{}
	if len(bind) > 0 {
		params = bind[0]
	}

	var keys []string
	for key, _ := range params {
		keys = append(keys, key)
	}
	if len(keys) == 0 {
		return sql
	}
	sort.Strings(keys)
	Reverse(keys)

	for _, key := range keys {
		v := Interface2String(params[key])
		val := AddSlashes(v)
		sql = strings.Replace(sql, ":"+key, "'"+val+"'", -1)
	}

	re, _ := regexp.Compile(`(?imU:limit.*$)`)
	matches := re.FindAllString(sql, -1)
	if len(matches) > 0 {
		match := matches[len(matches)-1]
		tmp := strings.Replace(match, "'", "", -1)
		sql = strings.Replace(sql, match, tmp, 1)
	}
	return sql
}

func In(arr []string) string {
	return "'" + strings.Join(arr, "', '") + "'"
}

func Insert(tableName string, data Json) (q.Result, error) {
	var keys []string
	var values []string
	for key, _ := range data {
		keys = append(keys, key)
		values = append(values, ":"+key)
	}
	var keyString = strings.Join(keys, ", ")
	var valueString = strings.Join(values, ", ")
	var sql = "INSERT INTO {tableName} ({keys}) VALUES ({values})"
	sql = Build(sql, Form{
		"tableName": tableName,
		"keys":      keyString,
		"values":    valueString,
	})
	return Exec(sql, data)
}

func InsertAll(tableName string, data []Json) (q.Result, error) {
	var keys []string
	for key, _ := range data[0] {
		keys = append(keys, key)
	}
	var keyString = strings.Join(keys, ", ")
	var vals []string
	for _, item := range data {
		var str = "('" + strings.Join(ArrayValues(item, keys), "', '") + "')"
		vals = append(vals, str)
	}
	var values = strings.Join(vals, ", ")
	var sql = "INSERT INTO {tableName} ({keys}) VALUES {values}"
	sql = Build(sql, Form{
		"tableName": tableName,
		"keys":      keyString,
		"values":    values,
	})
	return Exec(sql)
}

func ArrayValues(form Json, keys []string) []string {
	var result []string
	for _, key := range keys {
		s := Interface2String(form[key])
		result = append(result, AddSlashes(s))
	}
	return result
}

// 插入更新
func Replace(tableName string, data Json) (q.Result, error) {
	var keys []string
	var values []string
	for key, _ := range data {
		keys = append(keys, key)
		values = append(values, ":"+key)
	}
	var keyString = strings.Join(keys, ", ")
	var valueString = strings.Join(values, ", ")
	var sql = "REPLACE INTO {tableName} ({keys}) VALUES ({values})"
	sql = Build(sql, Form{
		"tableName": tableName,
		"keys":      keyString,
		"values":    valueString,
	})
	return Exec(sql, data)
}

// 批量插入更新
func ReplaceAll(tableName string, data []Json) (q.Result, error) {
	var keys []string
	for key, _ := range data[0] {
		keys = append(keys, key)
	}
	var keyString = strings.Join(keys, ", ")
	var vals []string
	for _, item := range data {
		var str = "('" + strings.Join(ArrayValues(item, keys), "', '") + "')"
		vals = append(vals, str)
	}
	var values = strings.Join(vals, ", ")
	var sql = "REPLACE INTO {tableName} ({keys}) VALUES {values}"
	sql = Build(sql, Form{
		"tableName": tableName,
		"keys":      keyString,
		"values":    values,
	})
	return Exec(sql)
}
