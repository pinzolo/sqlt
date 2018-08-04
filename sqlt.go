package sqlt

import (
	"bytes"
	"database/sql"
	"regexp"
	"text/template"
	"time"
)

const (
	// LeftDelim is start delimiter for SQL template.
	LeftDelim = "/*%"
	// RightDelim is end delimiter for SQL template.
	RightDelim = "%*/"
	// Connector is delimiter that is used when `sqlt` makes original argument name.
	Connector = "__"
)

var (
	strRegex = regexp.MustCompile(`%\*/'[^']*'`)
	inRegex  = regexp.MustCompile(`%\*/\([^()]*\)`)
	valRegex = regexp.MustCompile(`%\*/\S*`)
)

// SQLTemplate is template struct.
type SQLTemplate struct {
	dialect Dialect
	// TimeFunc used `time` and `now` function in template.
	// This func should return current time.
	// If this function is not set, used `time.Now()` as default function.
	TimeFunc func() time.Time
	// CustomFuncs are custom functions that are used in template.
	CustomFuncs map[string]interface{}
}

// New template initialized with dialect.
func New(dialect Dialect) SQLTemplate {
	return SQLTemplate{dialect: dialect, CustomFuncs: make(map[string]interface{})}
}

// AddFunc add custom tempalte func.
func (st SQLTemplate) AddFunc(name string, fn interface{}) SQLTemplate {
	st.CustomFuncs[name] = fn
	return st
}

// AddFuncs add custom template functions.
func (st SQLTemplate) AddFuncs(funcs map[string]interface{}) SQLTemplate {
	for k, v := range funcs {
		st.CustomFuncs[k] = v
	}
	return st
}

// Exec executes given template with given map parameters.
// This function replaces to normal placeholder.
func (st SQLTemplate) Exec(text string, m map[string]interface{}) (string, []interface{}, error) {
	c := newContext(false, st.dialect, st.TimeFunc, m)
	s, err := st.exec(c, text, m)
	if err != nil {
		return "", nil, err
	}
	return s, c.Args(), c.err
}

// ExecNamed executes given template with given map parameters.
// This function replaces to named placeholder.
func (st SQLTemplate) ExecNamed(text string, m map[string]interface{}) (string, []sql.NamedArg, error) {
	c := newContext(true, st.dialect, st.TimeFunc, m)
	s, err := st.exec(c, text, m)
	if err != nil {
		return "", nil, err
	}
	return s, c.NamedArgs(), c.err
}

func (st SQLTemplate) exec(c *context, text string, m map[string]interface{}) (string, error) {
	t, err := template.New("").Funcs(c.funcMap(st.CustomFuncs)).Delims(LeftDelim, RightDelim).Parse(dropSample(text))
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	if err = t.Execute(buf, m); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func dropSample(text string) string {
	s := strRegex.ReplaceAllString(text, RightDelim)
	s = inRegex.ReplaceAllString(s, RightDelim)
	return valRegex.ReplaceAllString(s, RightDelim)
}
