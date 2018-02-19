package sqlt

import (
	"bytes"
	"database/sql"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

const (
	// LeftDelim is start delimiter for SQL template.
	LeftDelim = "/*%"
	// RightDelim is end delimiter for SQL template.
	RightDelim = "%*/"
)

type param struct {
	sql.NamedArg
	Index int
}

type context struct {
	named     bool
	namedArgs []sql.NamedArg
	params    []*param
	values    []interface{}
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

func (c *context) p(name string) string {
	p := c.get(name)
	if p == nil {
		return ""
	}

	if p.Index == 0 {
		c.values = append(c.values, p.Value)
		p.Index = len(c.values)
		c.namedArgs = append(c.namedArgs, p.NamedArg)
	}
	if c.named {
		return ":" + p.Name
	}
	return "$" + strconv.Itoa(p.Index)
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
			placeholder = ":" + name + strconv.Itoa(i+1)
			c.namedArgs = append(c.namedArgs, sql.Named(name+strconv.Itoa(i+1), sv))
		} else {
			c.values = append(c.values, sv)
			placeholder = "$" + strconv.Itoa(len(c.values))
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

func Exec(text string, args ...sql.NamedArg) (string, []interface{}, error) {
	c := newContext(false, args...)
	s, err := exec(c, text)
	if err != nil {
		return "", nil, err
	}
	return s, c.values, nil
}

func ExecNamed(text string, args ...sql.NamedArg) (string, []sql.NamedArg, error) {
	c := newContext(true, args...)
	s, err := exec(c, text)
	if err != nil {
		return "", nil, err
	}
	return s, c.namedArgs, nil
}

func exec(c *context, text string) (string, error) {
	t, err := template.New("").Funcs(c.funcMap()).Delims(LeftDelim, RightDelim).Parse(dropSample(text))
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	if err = t.Execute(buf, c.parameters()); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func dropSample(text string) string {
	endPat := `%\*/`
	str := regexp.MustCompile(endPat + `'[^']*'`)
	in := regexp.MustCompile(endPat + `\([^\(\)]*\)`)
	val := regexp.MustCompile(endPat + `\S*`)
	s := str.ReplaceAllString(text, RightDelim)
	s = in.ReplaceAllString(s, RightDelim)
	return val.ReplaceAllString(s, RightDelim)
}

func newContext(named bool, args ...sql.NamedArg) *context {
	params := make([]*param, len(args))
	for i, arg := range args {
		params[i] = &param{NamedArg: arg}
	}
	return &context{
		named:     named,
		namedArgs: make([]sql.NamedArg, 0),
		params:    params,
		values:    make([]interface{}, 0),
	}
}
