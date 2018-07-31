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

const (
	escapeChar   = '\\'
	escapeClause = " ESCAPE '" + string(escapeChar) + "'"
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

func (c *context) paramWithFunc(name string, fn func(interface{}) interface{}) string {
	p := c.get(name)
	if p == nil {
		return unknownParamOutput(name)
	}

	v := p.value
	if fn != nil {
		v = fn(v)
	}

	if c.named {
		c.addNamed(p.name, v)
		return c.dialect.NamedPlaceholderPrefix() + p.name
	}

	if c.dialect.IsOrdinalPlaceholderSupported() {
		if p.index == 0 {
			c.values = append(c.values, v)
			p.index = len(c.values)
		}
		return c.dialect.OrdinalPlaceHolderPrefix() + strconv.Itoa(p.index)
	}
	c.values = append(c.values, v)
	return c.dialect.Placeholder()
}

func (c *context) param(name string) string {
	return c.paramWithFunc(name, nil)
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

func (c *context) prefix(name string) string {
	return c.paramWithEscapeLike(name) + " || '%'" + escapeClause
}

func (c *context) infix(name string) string {
	return "'%' || " + c.paramWithEscapeLike(name) + " || '%'" + escapeClause
}

func (c *context) suffix(name string) string {
	return "'%' || " + c.paramWithEscapeLike(name) + escapeClause

}

func (c *context) paramWithEscapeLike(name string) string {
	return c.paramWithFunc(name, c.escapeLike)
}

func (c *context) escapeLike(i interface{}) interface{} {
	s, ok := i.(string)
	if !ok {
		return i
	}

	rs := []rune(s)
	v := make([]rune, 0)
	for _, r := range rs {
		match := false
		for _, w := range c.dialect.WildcardRunes() {
			if r == w || r == escapeChar {
				match = true
				break
			}
		}
		if match {
			v = append(v, escapeChar)
		}
		v = append(v, r)
	}
	return string(v)
}

func (c *context) funcMap() template.FuncMap {
	return template.FuncMap{
		"param":  c.param,
		"p":      c.param,
		"in":     c.in,
		"time":   c.time,
		"now":    c.now,
		"prefix": c.prefix,
		"infix":  c.infix,
		"suffix": c.suffix,
	}
}

func unknownParamOutput(name string) string {
	return fmt.Sprintf("/*! %s is unknown */", name)
}
