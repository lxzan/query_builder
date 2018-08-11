# query_builder: golang查询构造器

> 为了脱离框架束缚, 随处使用自己喜欢的方式查询数据, 我封装了这个库
> 支持链式增删改查, 但链式查询器不支持连表, 本人不鼓励使用此方法进行复杂查询; 对于复杂查询, 推荐使用原生SQL
> 使用Any类型储存查询结果, 方便数据转换

### Install
```shell
go get github.com/lxzan/query_builder
```

### Connect
```go
// Single Mode
Connect(&ConnectOption{
	Host:     "127.0.0.1",
	Port:     "3306",
	User:     "root",
	Password: "lxz",
	Database: "test",
	Charset:  "utf8mb4",
})

// Cluster Mode
// one master, multi slave, the first is master
opts := make([]*ConnectOption, 0)
opts = append(opts, &ConnectOption{
	Host:     "127.0.0.1",
	Port:     "3306",
	User:     "root",
	Password: "lxz",
	Database: "test",
	Charset:  "utf8mb4",
})
ConnectCluster(opts)
```

### Query
```go
result, _ := Fetch("select * from pre_notice")

results, _ := FetchAll("select * from pre_notice limit 5")

query := Select().From("pre_notice").Where("uuid", "31856de1ba784169a4e6e3ca7c0249f8")
result,_ := query.Fetch()
println(result["uuid"].String(), query.SQL())
// print
31856de1ba784169a4e6e3ca7c0249f8 
SELECT * FROM pre_notice WHERE uuid = '31856de1ba784169a4e6e3ca7c0249f8' LIMIT 1
```
