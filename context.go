package sqlt

import (
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
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
	config  *config
	err     error
}

func newContext(named bool, dialect Dialect, m map[string]interface{}, conf *config) *context {
	params := make([]*param, len(m))
	i := 0
	for k, v := range m {
		params[i] = newParam(k, v)
		i++
	}
	return &context{
		named:   named,
		dialect: dialect,
		params:  params,
		args:    make([]*param, 0),
		timer:   newTimer(conf.timeFunc),
		config:  conf,
	}
}

func (c *context) Get(name string) (*param, error) {
	if strings.Contains(name, ".") {
		return c.Dig(strings.Split(name, "."))
	}
	for _, p := range c.params {
		if p.name == name {
			return p, nil
		}
	}
	return nil, fmt.Errorf("unknown param: %s", name)
}

func (c *context) Dig(names []string) (*param, error) {
	p, err := c.Get(names[0])
	if err != nil {
		return nil, err
	}

	qname := names[0]
	if p.value == nil {
		return nil, fmt.Errorf("nil value: %s", qname)
	}
	v := reflect.ValueOf(p.value)
	for _, name := range names[1:] {
		v, err = findValue(v, name, qname)
		if err != nil {
			return nil, err
		}
		qname = qname + "." + name
		if v.Kind() == reflect.Ptr {
			if v.IsNil() {
				return nil, fmt.Errorf("nil value: %s", qname)
			}
		}
	}
	return newParam(strings.Join(names, Connector), v.Interface()), nil
}

func findValue(val reflect.Value, name string, prefix string) (reflect.Value, error) {
	// cannot find field from pointer or interface, but can find method from those.
	// so at first search method.
	v, err := findMethodValue(val, name, prefix)
	if err != nil {
		return v, err
	}
	if v.IsValid() {
		return v, nil
	}

	return findFieldValue(val, name, prefix)
}

func findMethodValue(val reflect.Value, name string, prefix string) (reflect.Value, error) {
	// Finding method raises panic when Interface and nil.
	if val.Kind() == reflect.Interface && val.IsNil() {
		return val, fmt.Errorf("nil value: %s", prefix)
	}

	v := val.MethodByName(name)
	qname := prefix + "." + name
	if !v.IsValid() {
		return v, nil
	}
	t := v.Type()
	if t.NumIn() != 0 || t.NumOut() != 1 {
		return v, fmt.Errorf("invalid method: %s", qname)
	}
	v = v.Call([]reflect.Value{})[0]
	return v, nil
}

func findFieldValue(val reflect.Value, name string, prefix string) (reflect.Value, error) {
	// When value is pointer or interface, must call `Elem()` recursively.
	if val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		return findFieldValue(val.Elem(), name, prefix)
	}

	// Finding field raises panic when not struct
	if val.Kind() != reflect.Struct {
		return val, fmt.Errorf("not struct: %s", prefix)
	}

	v := val.FieldByName(name)
	qname := prefix + "." + name
	if !v.IsValid() {
		return v, fmt.Errorf("unknown param: %s", qname)
	}
	return v, nil
}

func (c *context) AddArg(name string, value interface{}) {
	c.args = append(c.args, newParam(name, value))
}

func (c *context) MergeArg(name string, value interface{}) {
	for _, arg := range c.args {
		if arg.name == name {
			return
		}
	}
	c.AddArg(name, value)
}

func (c *context) ArgIndex(name string) int {
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
		return c.dialect.OrdinalPlaceholderPrefix() + strconv.Itoa(c.ArgIndex(name))
	}
	return c.dialect.Placeholder()
}

func (c *context) errorOutput(err error) string {
	c.err = err
	return fmt.Sprintf("/*! %s */", err.Error())
}
