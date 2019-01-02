package sqlt_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/pinzolo/sqlt"
)

func BenchmarkExec(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := `
SELECT *
FROM users
WHERE id = /*%p "id" %*/1
AND name = /*% p "name" %*/'John Doe'
/*%- if get "onlyMale" %*/
AND sex = 'MALE'
/*%- end %*/
ORDER BY /*% out "order"" %*/id`
		m := map[string]interface{}{
			"id":       1,
			"order":    "name DESC",
			"onlyMale": true,
			"name":     "Alex",
		}
		_, _, err := sqlt.New(sqlt.Postgres).Exec(s, m)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkExecNamed(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := `
SELECT *
FROM users
WHERE id = /*%p "id" %*/1
AND name = /*% p "name" %*/'John Doe'
/*%- if get "onlyMale" %*/
AND sex = 'MALE'
/*%- end %*/
ORDER BY /*% out "order"" %*/id`
		m := map[string]interface{}{
			"id":       1,
			"order":    "name DESC",
			"onlyMale": true,
			"name":     "Alex",
		}
		_, _, err := sqlt.New(sqlt.Postgres).ExecNamed(s, m)
		if err != nil {
			b.Error(err)
		}
	}
}

func TestExec(t *testing.T) {
	s := `
SELECT *
FROM users
WHERE id IN /*% in "ids" %*/(1, 2)
AND name = /*% p "name" %*/'John Doe'
/*%- if get "onlyMale" %*/
AND sex = 'MALE'
/*%- end %*/
ORDER BY /*% out "order" %*/id`
	query, args, err := sqlt.New(sqlt.Postgres).Exec(s, map[string]interface{}{
		"ids":      []int{1, 2, 3},
		"order":    "name DESC",
		"onlyMale": true,
		"name":     "Alex",
	})
	if err != nil {
		t.Error(err)
	}

	eSQL := `
SELECT *
FROM users
WHERE id IN ($1, $2, $3)
AND name = $4
AND sex = 'MALE'
ORDER BY name DESC`
	if eSQL != query {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
	}
	if len(args) != 4 {
		t.Errorf("exec failed: values should have 4 length, but got %v", args)
	}
	if isInvalidInt(args[0], 1) {
		t.Errorf("exec failed: values should have 1, but got %v", args)
	}
	if isInvalidInt(args[1], 2) {
		t.Errorf("exec failed: values should have 2, but got %v", args)
	}
	if isInvalidInt(args[2], 3) {
		t.Errorf("exec failed: values should have 3, but got %v", args)
	}
	if isInvalidString(args[3], "Alex") {
		t.Errorf("exec failed: values should have 'Alex', but got %v", args)
	}
}

func TestExecNamed(t *testing.T) {
	s := `
SELECT *
FROM users
WHERE id IN /*% in "ids" %*/(1, 2)
AND name = /*% p "name" %*/'John Doe'
/*%- if get "onlyMale" %*/
AND sex = 'MALE'
/*%- end%*/
ORDER BY /*% out "order" %*/id`
	query, args, err := sqlt.New(sqlt.Postgres).ExecNamed(s, map[string]interface{}{
		"ids":      []int{1, 2, 3},
		"order":    "name DESC",
		"onlyMale": false,
		"name":     "Alex",
	})
	if err != nil {
		t.Error(err)
	}

	eSQL := `
SELECT *
FROM users
WHERE id IN (:ids__1, :ids__2, :ids__3)
AND name = :name
ORDER BY name DESC`
	if eSQL != query {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
	}
	if len(args) != 4 {
		t.Errorf("exec failed: values should have 4 length, but got %v", args)
	}
	if isInvalidIntArg(args[0], "ids__1", 1) {
		t.Errorf("exec failed: values should have ids__1 = 1, but got %v", args)
	}
	if isInvalidIntArg(args[1], "ids__2", 2) {
		t.Errorf("exec failed: values should have ids__2 = 2, but got %v", args)
	}
	if isInvalidIntArg(args[2], "ids__3", 3) {
		t.Errorf("exec failed: values should have ids__3 = 3, but got %v", args)
	}
	if isInvalidStringArg(args[3], "name", "Alex") {
		t.Errorf("exec failed: values should have name = 'Alex', but got %v", args)
	}
}

func TestExecWithNilParams(t *testing.T) {
	s := `SELECT * FROM users`
	query, args, err := sqlt.New(sqlt.Postgres).Exec(s, nil)
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT * FROM users`
	if eSQL != query {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
	}
	if len(args) != 0 {
		t.Errorf("exec failed: values should have 0 length, but got %v", args)
	}
}

func TestExecNamedWithNilParams(t *testing.T) {
	s := `SELECT * FROM users`
	query, args, err := sqlt.New(sqlt.Postgres).ExecNamed(s, nil)
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT * FROM users`
	if eSQL != query {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
	}
	if len(args) != 0 {
		t.Errorf("exec failed: values should have 0 length, but got %v", args)
	}
}

func TestExecWithInvalidTemplate(t *testing.T) {
	s := `
SELECT *
FROM users
WHERE id IN /*% in "ids" %*/(1, 2)
AND name = /*% pp "name" %*/'John Doe'
/*%- if get "onlyMale" %*/
AND sex = 'MALE'
/*%- end%*/
ORDER BY /*% out "order" %*/id`
	_, _, err := sqlt.New(sqlt.Postgres).Exec(s, map[string]interface{}{
		"ids":      []int{1, 2, 3},
		"order":    "name DESC",
		"onlyMale": true,
		"name":     "Alex",
	})
	if err == nil {
		t.Error("Exec with invalid template should raise error")
	}
}

func TestExecNamedWithInvalidTemplate(t *testing.T) {
	s := `
SELECT *
FROM users
WHERE id IN /*% in "ids" %*/(1, 2)
AND name = /*% pp "name" %*/'John Doe'
/*%- if get "onlyMale" %*/
AND sex = 'MALE'
/*%- end%*/
ORDER BY /*% out "order" %*/id`
	_, _, err := sqlt.New(sqlt.Postgres).ExecNamed(s, map[string]interface{}{
		"ids":      []int{1, 2, 3},
		"order":    "name DESC",
		"onlyMale": true,
		"name":     "Alex",
	})
	if err == nil {
		t.Error("ExecNamed with invalid template should raise error")
	}
}

func TestWithInvalidParamNameOnParamFunc(t *testing.T) {
	s := `
SELECT *
FROM users
WHERE id IN /*% in "ids" %*/(1, 2)
AND name = /*% p "userName" %*/'John Doe'
/*%- if get "onlyMale" %*/
AND sex = 'MALE'
/*%- end%*/
ORDER BY /*% out "order" %*/id`
	_, _, err := sqlt.New(sqlt.Postgres).Exec(s, map[string]interface{}{
		"ids":      []int{1, 2, 3},
		"order":    "name DESC",
		"onlyMale": true,
		"name":     "Alex",
	})
	if err == nil {
		t.Errorf("exec failed: should raise error when unknown param")
	}
}

func TestWithInvalidParamNameOnInFunc(t *testing.T) {
	s := `
SELECT *
FROM users
WHERE id IN /*% in "idList" %*/(1, 2)
AND name = /*% p "name" %*/'John Doe'
/*%- if get "onlyMale" %*/
AND sex = 'MALE'
/*%- end%*/
ORDER BY /*% out "order" %*/id`
	_, _, err := sqlt.New(sqlt.Postgres).Exec(s, map[string]interface{}{
		"ids":      []int{1, 2, 3},
		"order":    "name DESC",
		"onlyMale": true,
		"name":     "Alex",
	})
	if err == nil {
		t.Errorf("exec failed: should raise error when unknown param")
	}
}

func TestContinuousIn(t *testing.T) {
	s := `
SELECT *
FROM users
WHERE name IN /*% in "keywords" %*/('')
OR email IN /*% in "keywords" %*/('')`
	_, _, err := sqlt.New(sqlt.Postgres).Exec(s, map[string]interface{}{
		"keywords": []string{"foo", "bar"},
	})
	if err != nil {
		t.Error(err)
	}
}

func TestContinuousInNamed(t *testing.T) {
	s := `
SELECT *
FROM users
WHERE name IN /*% in "keywords" %*/('')
OR email IN /*% in "keywords" %*/('')`
	query, args, err := sqlt.New(sqlt.Postgres).ExecNamed(s, map[string]interface{}{
		"keywords": []string{"foo", "bar"},
	})
	if err != nil {
		t.Error(err)
	}

	eSQL := `
SELECT *
FROM users
WHERE name IN (:keywords__1, :keywords__2)
OR email IN (:keywords__1, :keywords__2)`
	if eSQL != query {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
	}
	if len(args) != 2 {
		t.Errorf("exec failed: values should have 2 length, but got %v", args)
	}
	if isInvalidStringArg(args[0], "keywords__1", "foo") {
		t.Errorf("exec failed: values should have 1, but got %v", args)
	}
	if isInvalidStringArg(args[1], "keywords__2", "bar") {
		t.Errorf("exec failed: values should have 1, but got %v", args)
	}
}

func TestTime(t *testing.T) {
	bt := time.Now()
	s := `
INSERT INTO users (
	name
  , created_at
  , updated_at
) VALUES (
	/*% p "name" %*/'John Doe'
  , /*% time %*/'2000-01-01'
  , /*% time %*/'2000-01-01'
)`
	query, args, err := sqlt.New(sqlt.Postgres).Exec(s, map[string]interface{}{
		"name": "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	eSQL := `
INSERT INTO users (
	name
  , created_at
  , updated_at
) VALUES (
	$1
  , $2
  , $2
)`
	et := time.Now()
	if eSQL != query {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
	}
	tm, ok := args[1].(time.Time)
	if !ok {
		t.Errorf("exec failed: 2nd arg expected time, but got %t", args[1])
	}
	if !(bt.Unix() <= tm.Unix() && tm.Unix() <= et.Unix()) {
		t.Errorf("time should be current time, but got %v", tm)
	}
}

func TestTimeWithTimeFunc(t *testing.T) {
	bt := time.Now()
	s := `
INSERT INTO users (
	name
  , created_at
  , updated_at
) VALUES (
	/*% p "name" %*/'John Doe'
  , /*% time %*/'2000-01-01'
  , /*% time %*/'2000-01-01'
)`
	st := sqlt.New(sqlt.Postgres).WithOptions(sqlt.TimeFunc(func() time.Time {
		return bt.AddDate(0, 0, 1)
	}))
	query, args, err := st.Exec(s, map[string]interface{}{
		"name": "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	eSQL := `
INSERT INTO users (
	name
  , created_at
  , updated_at
) VALUES (
	$1
  , $2
  , $2
)`
	et := time.Now()
	if eSQL != query {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
	}
	tm, ok := args[1].(time.Time)
	if !ok {
		t.Errorf("exec failed: 2nd arg expected time, but got %t", args[1])
	}
	if tm.Unix() != bt.AddDate(0, 0, 1).Unix() && bt.Unix() <= tm.Unix() && tm.Unix() <= et.Unix() {
		t.Errorf("time should return by TimeFunc, but got %v", tm)
	}
}

func TestTimeNamed(t *testing.T) {
	bt := time.Now()
	s := `
INSERT INTO users (
	name
  , created_at
  , updated_at
) VALUES (
	/*% p "name" %*/'John Doe'
  , /*% time %*/'2000-01-01'
  , /*% time %*/'2000-01-01'
)`
	query, args, err := sqlt.New(sqlt.Postgres).ExecNamed(s, map[string]interface{}{
		"name": "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	eSQL := `
INSERT INTO users (
	name
  , created_at
  , updated_at
) VALUES (
	:name
  , :time__
  , :time__
)`
	et := time.Now()
	if eSQL != query {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
	}
	arg := args[1]
	if arg.Name != "time__" {
		t.Errorf("default time arg name shoud be %q, but got %q", "time__", arg.Name)
	}
	tm, ok := arg.Value.(time.Time)
	if !ok {
		t.Errorf("exec failed: 2nd arg value expected time, but got %t", arg.Value)
	}
	if !(bt.Unix() <= tm.Unix() && tm.Unix() <= et.Unix()) {
		t.Errorf("time should be current time, but got %v", tm)
	}
}

func TestTimeNamedWithNameParam(t *testing.T) {
	bt := time.Now()
	s := `
INSERT INTO users (
	name
  , created_at
  , updated_at
) VALUES (
	/*% p "name" %*/'John Doe'
  , /*% time %*/'2000-01-01'
  , /*% time %*/'2000-01-01'
)`
	query, args, err := sqlt.New(sqlt.Postgres).ExecNamed(s, map[string]interface{}{
		"name": "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	eSQL := `
INSERT INTO users (
	name
  , created_at
  , updated_at
) VALUES (
	:name
  , :time__
  , :time__
)`
	et := time.Now()
	if eSQL != query {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
	}
	arg := args[1]
	if arg.Name != "time__" {
		t.Errorf("default time arg name shoud be %q, but got %q", "time__", arg.Name)
	}
	tm, ok := arg.Value.(time.Time)
	if !ok {
		t.Errorf("exec failed: 2nd arg value expected time, but got %t", arg.Value)
	}
	if !(bt.Unix() <= tm.Unix() && tm.Unix() <= et.Unix()) {
		t.Errorf("time should be current time, but got %v", tm)
	}
}

func TestNow(t *testing.T) {
	bt := time.Now()
	s := `
INSERT INTO users (
	name
  , created_at
  , updated_at
) VALUES (
	/*% p "name" %*/'John Doe'
  , /*% now %*/'2000-01-01'
  , /*% now %*/'2000-01-01'
)`
	query, args, err := sqlt.New(sqlt.Postgres).Exec(s, map[string]interface{}{
		"name": "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	eSQL := `
INSERT INTO users (
	name
  , created_at
  , updated_at
) VALUES (
	$1
  , $2
  , $3
)`
	et := time.Now()
	if eSQL != query {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
	}
	tm1, ok := args[1].(time.Time)
	if !ok {
		t.Errorf("exec failed: 2nd arg expected time, but got %t", args[1])
	}
	if !isBetweenTime(bt, tm1, et) {
		t.Errorf("time should be current time, but got %v", tm1)
	}
	tm2, ok := args[2].(time.Time)
	if !ok {
		t.Errorf("exec failed: 3rd arg expected time, but got %t", args[2])
	}
	if !isBetweenTime(bt, tm2, et) {
		t.Errorf("time should be current time, but got %v", tm2)
	}
}

func TestNowNamed(t *testing.T) {
	bt := time.Now()
	s := `
INSERT INTO users (
	name
  , created_at
  , updated_at
) VALUES (
	/*% p "name" %*/'John Doe'
  , /*% now %*/'2000-01-01'
  , /*% now %*/'2000-01-01'
)`
	query, args, err := sqlt.New(sqlt.Postgres).ExecNamed(s, map[string]interface{}{
		"name": "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	eSQL := `
INSERT INTO users (
	name
  , created_at
  , updated_at
) VALUES (
	:name
  , :now__1
  , :now__2
)`
	et := time.Now()
	if eSQL != query {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
	}
	arg := args[1]
	if arg.Name != "now__1" {
		t.Errorf("default time arg name shoud be %q, but got %q", "now__1", arg.Name)
	}
	tm1, ok := arg.Value.(time.Time)
	if !ok {
		t.Errorf("exec failed: 2nd arg value expected time, but got %t", arg.Value)
	}
	if !isBetweenTime(bt, tm1, et) {
		t.Errorf("time should be current time, but got %v", tm1)
	}
	arg = args[2]
	if arg.Name != "now__2" {
		t.Errorf("default time arg name shoud be %q, but got %q", "now__2", arg.Name)
	}
	tm2, ok := arg.Value.(time.Time)
	if !ok {
		t.Errorf("exec failed: 3rd arg value expected time, but got %t", arg.Value)
	}
	if !isBetweenTime(bt, tm2, et) {
		t.Errorf("time should be current time, but got %v", tm2)
	}
}

func TestNowNamedWithNameParam(t *testing.T) {
	bt := time.Now()
	s := `
INSERT INTO users (
	name
  , created_at
  , updated_at
) VALUES (
	/*% p "name" %*/'John Doe'
  , /*% now %*/'2000-01-01'
  , /*% now %*/'2000-01-01'
)`
	query, args, err := sqlt.New(sqlt.Postgres).ExecNamed(s, map[string]interface{}{
		"name": "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	eSQL := `
INSERT INTO users (
	name
  , created_at
  , updated_at
) VALUES (
	:name
  , :now__1
  , :now__2
)`
	et := time.Now()
	if eSQL != query {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
	}
	arg := args[1]
	if arg.Name != "now__1" {
		t.Errorf("default time arg name shoud be %q, but got %q", "now__1", arg.Name)
	}
	tm1, ok := arg.Value.(time.Time)
	if !ok {
		t.Errorf("exec failed: 2nd args value expected time, but got %t", arg.Value)
	}
	if !isBetweenTime(bt, tm1, et) {
		t.Errorf("time should be current time, but got %v", tm1)
	}
	arg = args[2]
	if arg.Name != "now__2" {
		t.Errorf("default time arg name shoud be %q, but got %q", "now__2", arg.Name)
	}
	tm2, ok := arg.Value.(time.Time)
	if !ok {
		t.Errorf("exec failed: 3rd args value expected time, but got %t", arg.Value)
	}
	if !isBetweenTime(bt, tm2, et) {
		t.Errorf("time should be current time, but got %v", tm2)
	}
}

func TestCustomFuncs(t *testing.T) {
	s := `
SELECT *
FROM users
WHERE name LIKE /*% infix "name" %*/''
/*% paging 3 50 %*/`
	query, args, err := sqlt.New(sqlt.Postgres).AddFuncs(map[string]interface{}{
		"paging": func(offset, limit int) string {
			return fmt.Sprintf("OFFSET %d LIMIT %d", offset, limit)
		},
		"infix": func() {
			panic("should not called")
		},
	}).Exec(s, map[string]interface{}{
		"name": "Alex",
	})
	if err != nil {
		t.Error(err)
	}

	eSQL := `
SELECT *
FROM users
WHERE name LIKE '%' || $1 || '%' ESCAPE '\'
OFFSET 3 LIMIT 50`
	if eSQL != query {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
	}
	if len(args) != 1 {
		t.Errorf("exec failed: values should have 1 length, but got %v", args)
	}
	if args[0] != `Alex` {
		t.Error("exec failed: embedded function should not be overwritten")
	}
}

func TestCustomFuncsContinuous(t *testing.T) {
	s := `
SELECT *
FROM users
WHERE name LIKE /*% infix "name" %*/''
/*% paging 3 50 %*/`
	query, args, err := sqlt.New(sqlt.Postgres).AddFunc("paging", func(offset, limit int) string {
		return fmt.Sprintf("OFFSET %d LIMIT %d", offset, limit)
	}).AddFunc("infix", func() {
		panic("should not called")
	}).Exec(s, map[string]interface{}{
		"name": "Alex",
	})
	if err != nil {
		t.Error(err)
	}

	eSQL := `
SELECT *
FROM users
WHERE name LIKE '%' || $1 || '%' ESCAPE '\'
OFFSET 3 LIMIT 50`
	if eSQL != query {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
	}
	if len(args) != 1 {
		t.Errorf("exec failed: values should have 1 length, but got %v", args)
	}
	if args[0] != `Alex` {
		t.Error("exec failed: embedded function should not be overwritten")
	}
}

func TestGetFunc(t *testing.T) {
	s := `
SELECT *
FROM users
ORDER BY /*% get "col" %*/id`
	data := []struct {
		col   string
		valid bool
		tag   string
	}{
		{"name", true, "valid"},
		{"name'", false, "single quotation"},
		{"name;", false, "semi colon"},
		{"name--", false, "line comment"},
		{"name/*", false, "block comment begin"},
		{"name*/", false, "block comment end"},
	}
	for _, d := range data {
		t.Run(d.tag, func(t *testing.T) {
			_, _, err := sqlt.New(sqlt.Postgres).Exec(s, map[string]interface{}{"col": d.col})
			if d.valid {
				if err != nil {
					t.Error(err)
				}
			} else {
				if err == nil {
					t.Error("should raise error on value contains prohibited character")
				}
			}
		})
	}
}

func TestOutFunc(t *testing.T) {
	s := `
SELECT *
FROM users
ORDER BY /*% out "col" %*/id`
	data := []struct {
		col   string
		valid bool
		tag   string
	}{
		{"name", true, "valid"},
		{"name'", false, "single quotation"},
		{"name;", false, "semi colon"},
		{"name--", false, "line comment"},
		{"name/*", false, "block comment begin"},
		{"name*/", false, "block comment end"},
	}
	for _, d := range data {
		t.Run(d.tag, func(t *testing.T) {
			_, _, err := sqlt.New(sqlt.Postgres).Exec(s, map[string]interface{}{"col": d.col})
			if d.valid {
				if err != nil {
					t.Error(err)
				}
			} else {
				if err == nil {
					t.Error("should raise error on value contains prohibited character")
				}
			}
		})
	}
}

func TestAnnotation(t *testing.T) {
	s := `
SELECT *
FROM users
WHERE id IN /*% in "ids" %*/(0) 
/*%- if get "sex" %*/
AND sex = /*% p "sex" %*/
/*%- end %*/
AND name LIKE /*% infix "name" %*/''
ORDER BY /*% out "col" %*/id`
	query, _, err := sqlt.New(sqlt.Postgres).WithOptions(sqlt.Annotation()).Exec(s, map[string]interface{}{
		"ids":  []int{1, 2, 3},
		"sex":  "MALE",
		"name": "foo",
		"col":  "name",
	})
	if err != nil {
		t.Fatal(err)
	}
	eSQL := `
SELECT *
FROM users
WHERE id IN ($1, $2, $3)/*# ids */
AND sex = $4/*# sex */
AND name LIKE '%' || $5/*# name */ || '%' ESCAPE '\'
ORDER BY name/*# col */`
	if eSQL != query {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
	}
}

func TestAnnotationError(t *testing.T) {
	data := []struct {
		tmpl string
		val  string
		want string
		tag  string
	}{
		{`/*% p "bar" %*/`, "test", `/*# error: "bar" is unknown param */`, "name not found"},
		{`/*% out "foo" %*/`, "test;", `/*# error: "foo" contains prohibited character(semi colon) */`, "prohibited character"},
		{`/*% get "foo" %*/`, "test;", "<no value>", "get does not annotate error"},
	}
	for _, d := range data {
		t.Run(d.tag, func(t *testing.T) {
			query, _, err := sqlt.New(sqlt.Postgres).WithOptions(sqlt.Annotation()).Exec(d.tmpl, map[string]interface{}{"foo": d.val})
			if err == nil {
				t.Error("should raise error")
			}
			if d.want != query {
				t.Errorf("exec failed: expected %v, but got %s", d.want, query)
			}
		})
	}
}

type Foo struct {
	Value string
}

type Bar struct {
	Foo    Foo
	FooPtr *Foo
	Baz    Baz
}

func (b Bar) Value() string {
	return b.Foo.Value
}

func (b Bar) Prop() Foo {
	return b.Foo
}

func (b Bar) PropPtr() *Foo {
	return b.FooPtr
}

func (b Bar) FnIn(s string) string {
	return s + b.Foo.Value
}

func (b Bar) FnOut2() (string, error) {
	return b.Foo.Value, nil
}

type Baz interface {
	Value() string
}

type baz struct{}

func (z *baz) Value() string {
	return "Alex"
}

func TestExecStruct(t *testing.T) {
	data := []struct {
		name  string
		value interface{}
		pArg  string
		tag   string
	}{
		{"foo", Foo{"Alex"}, "foo.Value", "simple struct"},
		{"foo", &Foo{"Alex"}, "foo.Value", "simple struct ptr"},
		{"bar", Bar{Foo: Foo{"Alex"}}, "bar.Foo.Value", "nested struct"},
		{"bar", Bar{FooPtr: &Foo{"Alex"}}, "bar.FooPtr.Value", "nested struct ptr"},
		{"bar", Bar{Foo: Foo{"Alex"}}, "bar.Value", "getter"},
		{"bar", Bar{Foo: Foo{"Alex"}}, "bar.Prop.Value", "nested getter"},
		{"bar", Bar{FooPtr: &Foo{"Alex"}}, "bar.PropPtr.Value", "nested getter ptr"},
		{"bar", Bar{Baz: &baz{}}, "bar.Baz.Value", "interface"},
	}
	for _, d := range data {
		t.Run(d.tag, func(t *testing.T) {
			p := `/*% p "` + d.pArg + `" %*/''`
			s := fmt.Sprintf(`SELECT * FROM users WHERE first_name = %s OR last_name = %s`, p, p)
			query, args, err := sqlt.New(sqlt.Postgres).Exec(s, map[string]interface{}{
				d.name: d.value,
			})
			if err != nil {
				t.Error(err)
			}
			eSQL := `SELECT * FROM users WHERE first_name = $1 OR last_name = $1`
			if eSQL != query {
				t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
			}
			if len(args) != 1 {
				t.Errorf("exec failed: values should have 1 length, but got %v", args)
			}
			if args[0] != "Alex" {
				t.Errorf("exec failed: expected value %q, but got %q", "Alex", args[0])
			}
		})
	}
}

func TestExecNamedStruct(t *testing.T) {
	data := []struct {
		name    string
		value   interface{}
		pArg    string
		argName string
		tag     string
	}{
		{"foo", Foo{"Alex"}, "foo.Value", "foo__Value", "simple struct"},
		{"foo", &Foo{"Alex"}, "foo.Value", "foo__Value", "simple struct ptr"},
		{"bar", Bar{Foo: Foo{"Alex"}}, "bar.Foo.Value", "bar__Foo__Value", "nested struct"},
		{"bar", Bar{FooPtr: &Foo{"Alex"}}, "bar.FooPtr.Value", "bar__FooPtr__Value", "nested struct ptr"},
		{"bar", Bar{Foo: Foo{"Alex"}}, "bar.Value", "bar__Value", "getter"},
		{"bar", Bar{Foo: Foo{"Alex"}}, "bar.Prop.Value", "bar__Prop__Value", "nested getter"},
		{"bar", Bar{FooPtr: &Foo{"Alex"}}, "bar.PropPtr.Value", "bar__PropPtr__Value", "nested getter ptr"},
		{"bar", Bar{Baz: &baz{}}, "bar.Baz.Value", "bar__Baz__Value", "interface"},
	}
	for _, d := range data {
		t.Run(d.tag, func(t *testing.T) {
			p := `/*% p "` + d.pArg + `" %*/''`
			s := fmt.Sprintf(`SELECT * FROM users WHERE first_name = %s OR last_name = %s`, p, p)
			query, args, err := sqlt.New(sqlt.Postgres).ExecNamed(s, map[string]interface{}{
				d.name: d.value,
			})
			if err != nil {
				t.Error(err)
			}
			eSQL := fmt.Sprintf(`SELECT * FROM users WHERE first_name = :%s OR last_name = :%s`, d.argName, d.argName)
			if eSQL != query {
				t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
			}
			if len(args) != 1 {
				t.Errorf("exec failed: values should have 1 length, but got %v", args)
			}
			if args[0].Name != d.argName {
				t.Errorf("exec failed: expected args name %q, but got %q", "foo__Value", args[0].Name)
			}
			if s, ok := args[0].Value.(string); !ok || s != "Alex" {
				t.Errorf("exec failed: expected value %q, but got %q", "Alex", args[0].Value)
			}
		})
	}
}

func TestExecStructError(t *testing.T) {
	data := []struct {
		value  interface{}
		pArg   string
		errMsg string
		tag    string
	}{
		{Bar{Foo: Foo{"Alex"}}, "bar.Qux.Value", `error: "bar.Qux" is unknown param`, "unknown"},
		{Bar{Foo: Foo{"Alex"}}, "bar.Value.Length", `error: "bar.Value" is not struct`, "not struct"},
		{Bar{}, "bar.FooPtr.Value", `error: "bar.FooPtr" is nil value`, "nil field"},
		{Bar{}, "bar.PropPtr.Value", `error: "bar.PropPtr" is nil value`, "nil getter"},
		{Bar{}, "bar.Baz.Value", `error: "bar.Baz" is nil value`, "nil interface"},
		{Bar{Foo: Foo{"Alex"}}, "bar.FnIn.Value", `error: "bar.FnIn" is invalid method`, "invalid in arg num method"},
		{Bar{Foo: Foo{"Alex"}}, "bar.FnOut2.Value", `error: "bar.FnOut2" is invalid method`, "invalid out get num method"},
		{nil, "bar.FooPtr.Value", `error: "bar" is nil value`, "nil root"},
		{Bar{Foo: Foo{"Alex"}}, "baz.Foo.Value", `error: "baz" is unknown param`, "unknown root"},
	}
	for _, d := range data {
		t.Run(d.tag, func(t *testing.T) {
			s := `SELECT * FROM users WHERE name = /*% p "` + d.pArg + `" %*/''`
			query, _, err := sqlt.New(sqlt.Postgres).WithOptions(sqlt.Annotation()).Exec(s, map[string]interface{}{
				"bar": d.value,
			})
			if err == nil {
				t.Error("should raise error")
			}
			eSQL := fmt.Sprintf(`SELECT * FROM users WHERE name = /*# %s */`, d.errMsg)
			if eSQL != query {
				t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
			}
		})
	}
}

func singleMap(k string, v interface{}) map[string]interface{} {
	return map[string]interface{}{
		k: v,
	}
}

func isInvalidInt(i interface{}, n int) bool {
	v, ok := i.(int)
	if !ok {
		return true
	}
	return v != n
}

func isInvalidIntArg(arg sql.NamedArg, name string, n int) bool {
	v, ok := arg.Value.(int)
	if !ok {
		return true
	}
	if arg.Name != name {
		return true
	}
	return v != n
}

func isInvalidString(i interface{}, s string) bool {
	v, ok := i.(string)
	if !ok {
		return true
	}
	return v != s
}

func isInvalidStringArg(arg sql.NamedArg, name string, s string) bool {
	v, ok := arg.Value.(string)
	if !ok {
		return true
	}
	if arg.Name != name {
		return true
	}
	return v != s
}

func isBetweenTime(bt, tm, et time.Time) bool {
	return bt.Unix() <= tm.Unix() && tm.Unix() <= et.Unix()
}
