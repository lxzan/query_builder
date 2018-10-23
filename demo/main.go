package main

import "github.com/lxzan/query_builder"

func main() {
	opts := make([]*db.ConnectOption, 0)
	opts = append(opts, &db.ConnectOption{
		Host:     "192.168.183.253",
		Port:     "3890",
		User:     "deve",
		Password: "Weiphone2017Q!",
		Database: "discuz",
		Charset:  "utf8mb4",
	})
	db.ConnectCluster(opts)

	res, _ := db.FetchAll("select * from pre_notice limit 10")
	println(res)
}
