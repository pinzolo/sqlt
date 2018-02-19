package sqlt

import (
	"bytes"
	"database/sql"
	"regexp"
	"text/template"
)

const (
	// LeftDelim is start delimiter for SQL template.
	LeftDelim = "/*%"
	// RightDelim is end delimiter for SQL template.
	RightDelim = "%*/"
)

// SQLTemplate is template struct.
type SQLTemplate struct {
	dialect Dialect
}

// New template initialized with dialect.
func New(dialect Dialect) SQLTemplate {
	return SQLTemplate{dialect: dialect}
}

// Exec executes given template with given arguments.
// This function replaces to normal placeholder.
func (st SQLTemplate) Exec(text string, args ...sql.NamedArg) (string, []interface{}, error) {
	c := newContextWithArgs(false, st.dialect, args...)
	s, err := st.exec(c, text)
	if err != nil {
		return "", nil, err
	}
	return s, c.values, nil
}

// ExecWithMap executes given template with given map parameters.
// This function replaces to normal placeholder.
func (st SQLTemplate) ExecWithMap(text string, m map[string]interface{}) (string, []interface{}, error) {
	c := newContextWithMap(false, st.dialect, m)
	s, err := st.exec(c, text)
	if err != nil {
		return "", nil, err
	}
	return s, c.values, nil
}

// ExecNamed executes given template with given arguments.
// This function replaces to named placeholder.
func (st SQLTemplate) ExecNamed(text string, args ...sql.NamedArg) (string, []sql.NamedArg, error) {
	c := newContextWithArgs(true, st.dialect, args...)
	s, err := st.exec(c, text)
	if err != nil {
		return "", nil, err
	}
	return s, c.namedArgs, nil
}

// ExecNamedWithMap executes given template with given map parameters.
// This function replaces to named placeholder.
func (st SQLTemplate) ExecNamedWithMap(text string, m map[string]interface{}) (string, []sql.NamedArg, error) {
	c := newContextWithMap(true, st.dialect, m)
	s, err := st.exec(c, text)
	if err != nil {
		return "", nil, err
	}
	return s, c.namedArgs, nil
}

func (st SQLTemplate) exec(c *context, text string) (string, error) {
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
