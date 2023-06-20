package query_builder

import (
	"strings"
)

type Filter struct {
	Expressions []string
	Args        []interface{}
}

func NewFilter() *Filter {
	return new(Filter)
}

func (c *Filter) append(key string, val interface{}, op string) {
	builder := strings.Builder{}
	builder.WriteString(key)
	builder.WriteString(" ")
	builder.WriteString(op)
	switch op {
	case "IN", "NOT IN":
		builder.WriteString(" (?)")
	default:
		builder.WriteString(" ?")
	}
	c.Expressions = append(c.Expressions, builder.String())
	c.Args = append(c.Args, val)
}

func (c *Filter) Equal(key string, val interface{}) *Filter {
	c.append(key, val, "=")
	return c
}

func (c *Filter) NotEqual(key string, val interface{}) *Filter {
	c.append(key, val, "!=")
	return c
}

func (c *Filter) Gt(key string, val interface{}) *Filter {
	c.append(key, val, ">")
	return c
}

func (c *Filter) Lt(key string, val interface{}) *Filter {
	c.append(key, val, "<")
	return c
}

func (c *Filter) Gte(key string, val interface{}) *Filter {
	c.append(key, val, ">=")
	return c
}

func (c *Filter) Lte(key string, val interface{}) *Filter {
	c.append(key, val, "<=")
	return c
}

func (c *Filter) addPercent(str string) string {
	if str == "" {
		return str
	}
	var n = len(str)
	if str[0] == '%' || str[n-1] == '%' {
		return str
	}
	return "%" + str + "%"
}

func (c *Filter) Like(key string, val string) *Filter {
	c.append(key, c.addPercent(val), "LIKE")
	return c
}

func (c *Filter) NotLike(key string, val string) *Filter {
	c.append(key, c.addPercent(val), "NOT LIKE")
	return c
}

func (c *Filter) In(key string, val ...interface{}) *Filter {
	c.append(key, val, "IN")
	return c
}

func (c *Filter) NotIn(key string, val ...interface{}) *Filter {
	c.append(key, val, "NOT IN")
	return c
}

func (c *Filter) IsNull(key string) *Filter {
	c.Expressions = append(c.Expressions, key+" IS NULL")
	return c
}

func (c *Filter) With(key string, val ...interface{}) *Filter {
	c.Expressions = append(c.Expressions, key)
	c.Args = append(c.Args, val...)
	return c
}

func (c *Filter) GetExpression() string {
	return strings.Join(c.Expressions, " AND ")
}
