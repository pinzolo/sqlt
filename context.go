package sqlt

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

const (
	escapeChar   = '\\'
	escapeClause = " ESCAPE '" + string(escapeChar) + "'"
)

type param struct {
	name  string
	value interface{}
}

func newParam(name string, value interface{}) *param {
	return &param{
		name:  name,
		value: value,
	}
}

type context struct {
	named   bool
	dialect Dialect
	params  []*param
	args    []*param
	timer   *timer
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
		named:   named,
		dialect: dialect,
		params:  params,
		args:    make([]*param, 0),
		timer:   newTimer(fn),
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

func (c *context) addArg(name string, value interface{}) {
	c.args = append(c.args, newParam(name, value))
}

func (c *context) mergeArg(name string, value interface{}) {
	for _, arg := range c.args {
		if arg.name == name {
			return
		}
	}
	c.addArg(name, value)
}

func (c *context) argIndex(name string) int {
	for i, arg := range c.args {
		if arg.name == name {
			return i + 1
		}
	}
	return 0
}

func (c *context) Args() []interface{} {
	v := make([]interface{}, len(c.args))
	for i, arg := range c.args {
		v[i] = arg.value
	}
	return v
}

func (c *context) NamedArgs() []sql.NamedArg {
	v := make([]sql.NamedArg, len(c.args))
	for i, arg := range c.args {
		v[i] = sql.Named(arg.name, arg.value)
	}
	return v
}

func (c *context) Placeholder(name string) string {
	if c.named {
		return c.dialect.NamedPlaceholderPrefix() + name
	}
	if c.dialect.IsOrdinalPlaceholderSupported() {
		return c.dialect.OrdinalPlaceholderPrefix() + strconv.Itoa(c.argIndex(name))
	}
	return c.dialect.Placeholder()
}

func unknownParamOutput(name string) string {
	return fmt.Sprintf("/*! %s is unknown */", name)
}
