package db

import (
	"strconv"
	"strings"
	"math/rand"
	"time"
	"fmt"
	"regexp"
)

func Reverse(arr []string) {
	var n int
	var length = len(arr)
	n = length / 2
	for i := 0; i < n; i++ {
		arr[i], arr[length-1-i] = Swap(arr[i], arr[length-1-i])
	}
}

func ToString(num int64) string {
	return strconv.Itoa(int(num))
}

func Interface2String(v interface{}) string {
	var s = ""
	if f, ok := v.(string); ok {
		s = f
	} else if f, ok := v.(int64); ok {
		s = ToString(f)
	} else if f, ok := v.(float64); ok {
		s = fmt.Sprintf("%.6f", f)
		re, _ := regexp.Compile(`\.*[0]+$`)
		s = re.ReplaceAllString(s, "")
	} else {
		panic("val only supports int, int64, float, string type")
	}
	return s
}

func Swap(a string, b string) (string, string) {
	return b, a
}

// 转义引号
func AddSlashes(str string) string {
	str = strings.Replace(str, "'", "\\'", -1)
	str = strings.Replace(str, "\"", "\\\"", -1)
	str = strings.Replace(str, "`", "\\`", -1)
	return str
}

// 字符串模板
func Build(message string, bind Form) string {
	for k, v := range bind {
		var re = "{" + k + "}"
		message = strings.Replace(message, re, v, -1)
	}
	return message
}

// @return min <= x <=10
func Rand(min int, max int) int {
	var seed = rand.New(rand.NewSource(time.Now().UnixNano()))
	var delta = max - min + 1
	var val = min + seed.Intn(delta)
	return val
}
