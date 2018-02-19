package sqlt

import (
	"database/sql"
	"reflect"
	"strconv"
	"strings"
	"text/template"
)

type param struct {
	sql.NamedArg
	Index int
}

type context struct {
	named     bool
	dialect   Dialect
	namedArgs []sql.NamedArg
	params    []*param
	values    []interface{}
}

func newContext(named bool, dialect Dialect, args ...sql.NamedArg) *context {
	params := make([]*param, len(args))
	for i, arg := range args {
		params[i] = &param{NamedArg: arg}
	}
	return &context{
		named:     named,
		dialect:   dialect,
		namedArgs: make([]sql.NamedArg, 0),
		params:    params,
		values:    make([]interface{}, 0),
	}
}

func (c *context) parameters() map[string]interface{} {
	paramMap := make(map[string]interface{})
	for _, p := range c.params {
		paramMap[p.Name] = p.Value
	}
	return paramMap
}

func (c *context) get(name string) *param {
	for _, p := range c.params {
		if p.Name == name {
			return p
		}
	}
	return nil
}

func (c *context) addNamedArgs(name string, value interface{}) {
	for _, arg := range c.namedArgs {
		if arg.Name == name {
			return
		}
	}
	c.namedArgs = append(c.namedArgs, sql.Named(name, value))
}

func (c *context) p(name string) string {
	p := c.get(name)
	if p == nil {
		return ""
	}

	if c.named {
		c.addNamedArgs(p.Name, p.Value)
		return c.dialect.NamedPlaceholderPrefix() + p.Name
	}

	if c.dialect.IsOrdinalPlaceholderSupported() {
		if p.Index == 0 {
			c.values = append(c.values, p.Value)
			p.Index = len(c.values)
		}
		return c.dialect.OrdinalPlaceHolderPrefix() + strconv.Itoa(p.Index)
	}

	c.values = append(c.values, p.Value)
	return c.dialect.Placeholder()
}

func (c *context) in(name string) string {
	p := c.get(name)
	if p == nil {
		return ""
	}

	v := reflect.ValueOf(p.Value)
	if v.Kind() != reflect.Slice {
		return "(" + c.p(name) + ")"
	}

	placeholders := make([]string, v.Len())
	for i := 0; i < v.Len(); i++ {
		sv := v.Index(i).Interface()
		var placeholder string
		if c.named {
			placeholder = c.dialect.NamedPlaceholderPrefix() + name + strconv.Itoa(i+1)
			c.namedArgs = append(c.namedArgs, sql.Named(name+strconv.Itoa(i+1), sv))
		} else {
			c.values = append(c.values, sv)
			if c.dialect.IsOrdinalPlaceholderSupported() {
				placeholder = c.dialect.OrdinalPlaceHolderPrefix() + strconv.Itoa(len(c.values))
			} else {
				placeholder = c.dialect.Placeholder()
			}
		}
		placeholders[i] = placeholder
	}
	return "(" + strings.Join(placeholders, ",") + ")"
}

func (c *context) funcMap() template.FuncMap {
	return template.FuncMap{
		"p":  c.p,
		"in": c.in,
	}
}
