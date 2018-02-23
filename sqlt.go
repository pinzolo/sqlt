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

var (
	strRegex = regexp.MustCompile(`%\*/'[^']*'`)
	inRegex  = regexp.MustCompile(`%\*/\([^\(\)]*\)`)
	valRegex = regexp.MustCompile(`%\*/\S*`)
)

// SQLTemplate is template struct.
type SQLTemplate struct {
	dialect Dialect
}

// New template initialized with dialect.
func New(dialect Dialect) SQLTemplate {
	return SQLTemplate{dialect: dialect}
}

// Exec executes given template with given map parameters.
// This function replaces to normal placeholder.
func (st SQLTemplate) Exec(text string, m map[string]interface{}) (string, []interface{}, error) {
	c := newContext(false, st.dialect, m)
	s, err := st.exec(c, text)
	if err != nil {
		return "", nil, err
	}
	return s, c.values, nil
}

// ExecNamed executes given template with given map parameters.
// This function replaces to named placeholder.
func (st SQLTemplate) ExecNamed(text string, m map[string]interface{}) (string, []sql.NamedArg, error) {
	c := newContext(true, st.dialect, m)
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
	s := strRegex.ReplaceAllString(text, RightDelim)
	s = inRegex.ReplaceAllString(s, RightDelim)
	return valRegex.ReplaceAllString(s, RightDelim)
}
