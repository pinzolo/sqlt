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

// Config is configuration for executing template.
type config struct {
	timeFunc func() time.Time
}

// SQLTemplate is template struct.
type SQLTemplate struct {
	dialect Dialect
	// customFuncs are custom functions that are used in template.
	customFuncs map[string]interface{}
	config      *config
}

// New template initialized with dialect.
func New(dialect Dialect) *SQLTemplate {
	return &SQLTemplate{
		dialect:     dialect,
		customFuncs: make(map[string]interface{}),
		config:      &config{},
	}
}

// AddFunc add custom template func.
func (st *SQLTemplate) AddFunc(name string, fn interface{}) *SQLTemplate {
	st.customFuncs[name] = fn
	return st
}

// AddFuncs add custom template functions.
func (st *SQLTemplate) AddFuncs(funcs map[string]interface{}) *SQLTemplate {
	for k, v := range funcs {
		st.customFuncs[k] = v
	}
	return st
}

// WithOptions apply given options.
func (st *SQLTemplate) WithOptions(opts ...Option) *SQLTemplate {
	for _, opt := range opts {
		opt(st.config)
	}
	return st
}

// Exec executes given template with given map parameters.
// This function replaces to normal placeholder.
func (st *SQLTemplate) Exec(text string, m map[string]interface{}) (string, []interface{}, error) {
	c := newContext(false, st.dialect, m, st.config)
	s, err := st.exec(c, text, m)
	if err != nil {
		return "", nil, err
	}
	return s, c.Args(), c.err
}

// ExecNamed executes given template with given map parameters.
// This function replaces to named placeholder.
func (st *SQLTemplate) ExecNamed(text string, m map[string]interface{}) (string, []sql.NamedArg, error) {
	c := newContext(true, st.dialect, m, st.config)
	s, err := st.exec(c, text, m)
	if err != nil {
		return "", nil, err
	}
	return s, c.NamedArgs(), c.err
}

func (st *SQLTemplate) exec(c *context, text string, m map[string]interface{}) (string, error) {
	t, err := template.New("").Funcs(c.funcMap(st.customFuncs)).Delims(LeftDelim, RightDelim).Parse(dropSample(text))
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	if err = t.Execute(buf, nil); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func dropSample(text string) string {
	s := strRegex.ReplaceAllString(text, RightDelim)
	s = inRegex.ReplaceAllString(s, RightDelim)
	return valRegex.ReplaceAllString(s, RightDelim)
}
