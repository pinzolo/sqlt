package sqlt

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"text/template"
)

func (c *context) paramWithFunc(name string, fn func(interface{}) interface{}) string {
	p := c.get(name)
	if p == nil {
		return unknownParamOutput(name)
	}

	v := p.value
	if fn != nil {
		v = fn(v)
	}

	if c.named || c.dialect.IsOrdinalPlaceholderSupported() {
		c.mergeArg(p.name, v)
	} else {
		c.addArg(name, v)
	}
	return c.Placeholder(p.name)
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
		argName := fmt.Sprintf("%s__%d", name, i+1)
		if c.named || c.dialect.IsOrdinalPlaceholderSupported() {
			c.mergeArg(argName, sv)
		} else {
			c.addArg(argName, sv)
		}
		placeholders[i] = c.Placeholder(argName)
	}
	return "(" + strings.Join(placeholders, ", ") + ")"
}

func (c *context) time() string {
	name := "time__"
	tm := c.timer.time()
	if c.named || c.dialect.IsOrdinalPlaceholderSupported() {
		c.mergeArg(name, tm)
	} else {
		c.addArg(name, tm)
	}

	return c.Placeholder(name)
}

func (c *context) now() string {
	name := "now__" + strconv.Itoa(c.timer.nowCnt)
	c.addArg(name, c.timer.now())
	return c.Placeholder(name)
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
		if c.isEscapee(r) {
			v = append(v, escapeChar)
		}
		v = append(v, r)
	}
	return string(v)
}

func (c *context) isEscapee(r rune) bool {
	for _, w := range c.dialect.WildcardRunes() {
		if r == w || r == escapeChar {
			return true
		}
	}
	return false
}

func (c *context) funcMap(funcs map[string]interface{}) template.FuncMap {
	fm := template.FuncMap(funcs)
	fm["param"] = c.param
	fm["p"] = c.param
	fm["in"] = c.in
	fm["time"] = c.time
	fm["now"] = c.now
	fm["prefix"] = c.prefix
	fm["infix"] = c.infix
	fm["suffix"] = c.suffix
	fm["escape"] = c.escapeLike
	return fm
}
