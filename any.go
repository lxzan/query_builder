package db

import (
	"strconv"
)


type Any struct {
	v  interface{}
	s  string
	i  int
	b  bool
	i6 int64
}

func NewAny(v interface{}) *Any {
	var st = new(Any)
	st.v = v
	return st
}

func (u *Any) String() string {
	if u.s == "" {
		u.s = Interface2String(u.v)
	}
	return u.s
}

func (u *Any) Int() int {
	if u.i == 0 {
		u.String()
		i, _ := strconv.ParseInt(u.s, 10, 64)
		u.i = int(i)
	}
	return u.i
}

func (u *Any) Int64() int64 {
	if u.i6 == 0 {
		u.String()
		u.i6, _ = strconv.ParseInt(u.s, 10, 64)
	}
	return u.i6
}

func (u *Any) Bool() bool {
	if u.b == false {
		u.String()
		if u.s == "true" {
			u.b = true
		} else {
			u.b = false
		}
	}
	return u.b
}
