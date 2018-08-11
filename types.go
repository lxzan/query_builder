package db

type ConnectOption struct {
	User         string
	Password     string
	Host         string
	Port         string
	Database     string
	Charset      string
	MaxOpenConns int // default 200
	MaxIdleConns int // default 100
}

type Form map[string]string

type Json map[string]interface{}

type QueryResult map[string]*Any

