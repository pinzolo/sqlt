package sqlt

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type param struct {
	name  string
	value interface{}
	index int
}

func newParam(name string, value interface{}) *param {
	return &param{
		name:  name,
		value: value,
	}
}

type context struct {
	named     bool
	dialect   Dialect
	params    []*param
	namedArgs []sql.NamedArg
	values    []interface{}
	paramMap  map[string]interface{}
	timer     *timer
}

func newContext(named bool, dialect Dialect, timeFn func() time.Time, m map[string]interface{}) *context {
	params := make([]*param, len(m))
	i := 0
	for k, v := range m {
		params[i] = newParam(k, v)
		i++
	}
	fn := time.Now
	if timeFn != nil {
		fn = timeFn
	}
	return &context{
		named:     named,
		dialect:   dialect,
		params:    params,
		namedArgs: []sql.NamedArg{},
		values:    []interface{}{},
		paramMap:  m,
		timer:     newTimer(fn),
	}
}

func (c *context) get(name string) *param {
	for _, p := range c.params {
		if p.name == name {
			return p
		}
	}
	return nil
}

func (c *context) addNamed(name string, value interface{}) {
	for _, arg := range c.namedArgs {
		if arg.Name == name {
			return
		}
	}
	c.namedArgs = append(c.namedArgs, sql.Named(name, value))
}

func (c *context) param(name string) string {
	p := c.get(name)
	if p == nil {
		return unknownParamOutput(name)
	}

	if c.named {
		c.addNamed(p.name, p.value)
		return c.dialect.NamedPlaceholderPrefix() + p.name
	}

	if c.dialect.IsOrdinalPlaceholderSupported() {
		if p.index == 0 {
			c.values = append(c.values, p.value)
			p.index = len(c.values)
		}
		return c.dialect.OrdinalPlaceHolderPrefix() + strconv.Itoa(p.index)
	}

	c.values = append(c.values, p.value)
	return c.dialect.Placeholder()
}

func (c *context) in(name string) string {
	p := c.get(name)
	if p == nil {
		return unknownParamOutput(name)
	}

	v := reflect.ValueOf(p.value)
	if v.Kind() != reflect.Slice {
		return "(" + c.param(name) + ")"
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
	return "(" + strings.Join(placeholders, ", ") + ")"
}

func (c *context) time() string {
	name := "time__"
	if c.named {
		c.addNamed(name, c.timer.time())
		return c.dialect.NamedPlaceholderPrefix() + name
	}

	if c.dialect.IsOrdinalPlaceholderSupported() {
		if c.timer.cacheIndex == 0 {
			c.values = append(c.values, c.timer.time())
			c.timer.cacheIndex = len(c.values)
		}
		return c.dialect.OrdinalPlaceHolderPrefix() + strconv.Itoa(c.timer.cacheIndex)
	}

	c.values = append(c.values, c.timer.time())
	return c.dialect.Placeholder()
}

func (c *context) now() string {
	name := "now__" + strconv.Itoa(c.timer.nowCnt)
	if c.named {
		c.addNamed(name, c.timer.now())
		return c.dialect.NamedPlaceholderPrefix() + name
	}

	if c.dialect.IsOrdinalPlaceholderSupported() {
		c.values = append(c.values, c.timer.now())
		return c.dialect.OrdinalPlaceHolderPrefix() + strconv.Itoa(len(c.values))
	}

	c.values = append(c.values, c.timer.now())
	return c.dialect.Placeholder()
}

func (c *context) funcMap() template.FuncMap {
	return template.FuncMap{
		"param": c.param,
		"p":     c.param,
		"in":    c.in,
		"time":  c.time,
		"now":   c.now,
	}
}

func unknownParamOutput(name string) string {
	return fmt.Sprintf("/*! %s is unknown */", name)
}
