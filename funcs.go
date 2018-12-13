package sqlt

import (
	"fmt"
	"reflect"
	"strings"
	"text/template"
)

func (c *context) paramWithFunc(name string, fn func(string, interface{}) (string, interface{})) string {
	p, err := c.Get(name)
	if err != nil {
		return c.errorOutput(err)
	}

	nm := p.name
	v := p.value
	if fn != nil {
		nm, v = fn(nm, v)
	}

	if c.named || c.dialect.IsOrdinalPlaceholderSupported() {
		c.MergeArg(nm, v)
	} else {
		c.AddArg(nm, v)
	}
	return c.Placeholder(nm)
}

func (c *context) val(name string) interface{} {
	p, err := c.Get(name)
	if err != nil {
		c.err = err
		return nil
	}
	if s, ok := p.value.(string); ok {
		if err = safe(s); err != nil {
			c.err = fmt.Errorf("%q contains prohibited character(%s)", name, err.Error())
			return nil
		}
	}
	return p.value
}

func (c *context) param(name string) string {
	return c.paramWithFunc(name, nil)
}

func (c *context) in(name string) string {
	p, err := c.Get(name)
	if err != nil {
		return c.errorOutput(err)
	}

	v := reflect.ValueOf(p.value)
	if v.Kind() != reflect.Slice {
		return "(" + c.param(name) + ")"
	}

	placeholders := make([]string, v.Len())
	for i := 0; i < v.Len(); i++ {
		sv := v.Index(i).Interface()
		argName := fmt.Sprintf("%s%s%d", name, Connector, i+1)
		if c.named || c.dialect.IsOrdinalPlaceholderSupported() {
			c.MergeArg(argName, sv)
		} else {
			c.AddArg(argName, sv)
		}
		placeholders[i] = c.Placeholder(argName)
	}
	return "(" + strings.Join(placeholders, ", ") + ")"
}

func (c *context) time() string {
	name := "time" + Connector
	tm := c.timer.time()
	if c.named || c.dialect.IsOrdinalPlaceholderSupported() {
		c.MergeArg(name, tm)
	} else {
		c.AddArg(name, tm)
	}

	return c.Placeholder(name)
}

func (c *context) now() string {
	name := fmt.Sprintf("now%s%d", Connector, c.timer.nowCnt)
	c.AddArg(name, c.timer.now())
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
	return c.paramWithFunc(name, func(name string, v interface{}) (string, interface{}) {
		nv := c.escapeLike(v)
		if nv == v {
			return name, v
		}
		return name + Connector + "esc", nv
	})
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
	fm["value"] = c.val
	fm["val"] = c.val
	fm["v"] = c.val
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

func safe(s string) error {
	if strings.Contains(s, "'") {
		return fmt.Errorf("single quotation")
	}
	if strings.Contains(s, ";") {
		return fmt.Errorf("semi colon")
	}
	if strings.Contains(s, "--") || strings.Contains(s, "/*") || strings.Contains(s, "*/") {
		return fmt.Errorf("comment")
	}
	return nil
}
