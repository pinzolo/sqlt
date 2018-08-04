package sqlt

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/pinzolo/tagt"
)

func BenchmarkExec(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := `SELECT *
		FROM users
		WHERE id = /*%p "id" %*/1
		AND name = /*% p "name" %*/'John Doe'
		/*%- if .onlyMale %*/
		AND sex = 'MALE'
		/*%- end%*/
		ORDER BY /*% .order %*/id`
		m := map[string]interface{}{
			"id":       1,
			"order":    "name DESC",
			"onlyMale": true,
			"name":     "Alex",
		}
		_, _, err := New(Postgres).Exec(s, m)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkExecNamed(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s := `SELECT *
		FROM users
		WHERE id = /*%p "id" %*/1
		AND name = /*% p "name" %*/'John Doe'
		/*%- if .onlyMale %*/
		AND sex = 'MALE'
		/*%- end%*/
		ORDER BY /*% .order %*/id`
		m := map[string]interface{}{
			"id":       1,
			"order":    "name DESC",
			"onlyMale": true,
			"name":     "Alex",
		}
		_, _, err := New(Postgres).ExecNamed(s, m)
		if err != nil {
			b.Error(err)
		}
	}
}

func TestDropSample(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id IN /*%in "ids" %*/(1, 2, 3)
	AND created_at = /*%p "time"%*/'2000-01-01 12:34:56'
	AND name = /*%p "name"%*/'foo'
	AND age = /*%p "age"%*/18`
	expected := `SELECT *
	FROM users
	WHERE id IN /*%in "ids" %*/
	AND created_at = /*%p "time"%*/
	AND name = /*%p "name"%*/
	AND age = /*%p "age"%*/`
	if actual := dropSample(s); actual != expected {
		t.Errorf("dropSample faild: expected %s, but got %s", expected, actual)
	}
}
func TestExec(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id IN /*% in "ids" %*/(1, 2)
	AND name = /*% p "name" %*/'John Doe'
	/*%- if .onlyMale %*/
	AND sex = 'MALE'
	/*%- end%*/
	ORDER BY /*% .order %*/id`
	sql, args, err := New(Postgres).Exec(s, map[string]interface{}{
		"ids":      []int{1, 2, 3},
		"order":    "name DESC",
		"onlyMale": true,
		"name":     "Alex",
	})
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE id IN ($1, $2, $3)
	AND name = $4
	AND sex = 'MALE'
	ORDER BY name DESC`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
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
	s := `SELECT *
	FROM users
	WHERE id IN /*% in "ids" %*/(1, 2)
	AND name = /*% p "name" %*/'John Doe'
	/*%- if .onlyMale %*/
	AND sex = 'MALE'
	/*%- end%*/
	ORDER BY /*% .order %*/id`
	sql, args, err := New(Postgres).ExecNamed(s, map[string]interface{}{
		"ids":      []int{1, 2, 3},
		"order":    "name DESC",
		"onlyMale": false,
		"name":     "Alex",
	})
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE id IN (:ids__1, :ids__2, :ids__3)
	AND name = :name
	ORDER BY name DESC`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
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
	s := `SELECT *
	FROM users`
	sql, args, err := New(Postgres).Exec(s, nil)
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}
	if len(args) != 0 {
		t.Errorf("exec failed: values should have 0 length, but got %v", args)
	}
}

func TestExecNamedWithNilParams(t *testing.T) {
	s := `SELECT *
	FROM users`
	sql, args, err := New(Postgres).ExecNamed(s, nil)
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}
	if len(args) != 0 {
		t.Errorf("exec failed: values should have 0 length, but got %v", args)
	}
}

func TestExecWithInvalidTemplate(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id IN /*% in "ids" %*/(1, 2)
	AND name = /*% pp "name" %*/'John Doe'
	/*%- if .onlyMale %*/
	AND sex = 'MALE'
	/*%- end%*/
	ORDER BY /*% .order %*/id`
	_, _, err := New(Postgres).Exec(s, map[string]interface{}{
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
	s := `SELECT *
	FROM users
	WHERE id IN /*% in "ids" %*/(1, 2)
	AND name = /*% pp "name" %*/'John Doe'
	/*%- if .onlyMale %*/
	AND sex = 'MALE'
	/*%- end%*/
	ORDER BY /*% .order %*/id`
	_, _, err := New(Postgres).ExecNamed(s, map[string]interface{}{
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
	s := `SELECT *
	FROM users
	WHERE id IN /*% in "ids" %*/(1, 2)
	AND name = /*% p "userName" %*/'John Doe'
	/*%- if .onlyMale %*/
	AND sex = 'MALE'
	/*%- end%*/
	ORDER BY /*% .order %*/id`
	sql, _, err := New(Postgres).Exec(s, map[string]interface{}{
		"ids":      []int{1, 2, 3},
		"order":    "name DESC",
		"onlyMale": true,
		"name":     "Alex",
	})
	if err == nil {
		t.Errorf("exec failed: should raise error when unknown param")
	}

	eSQL := `SELECT *
	FROM users
	WHERE id IN ($1, $2, $3)
	AND name = /*! unknown param: userName */
	AND sex = 'MALE'
	ORDER BY name DESC`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}
}

func TestWithInvalidParamNameOnInFunc(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id IN /*% in "idList" %*/(1, 2)
	AND name = /*% p "name" %*/'John Doe'
	/*%- if .onlyMale %*/
	AND sex = 'MALE'
	/*%- end%*/
	ORDER BY /*% .order %*/id`
	sql, _, err := New(Postgres).Exec(s, map[string]interface{}{
		"ids":      []int{1, 2, 3},
		"order":    "name DESC",
		"onlyMale": true,
		"name":     "Alex",
	})
	if err == nil {
		t.Errorf("exec failed: should raise error when unknown param")
	}

	eSQL := `SELECT *
	FROM users
	WHERE id IN /*! unknown param: idList */
	AND name = $1
	AND sex = 'MALE'
	ORDER BY name DESC`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}
}

func TestContinuousIn(t *testing.T) {
	s := `SELECT *
FROM users
WHERE name IN /*% in "keywords" %*/('')
OR email IN /*% in "keywords" %*/('')`
	sql, args, err := New(Postgres).Exec(s, map[string]interface{}{
		"keywords": []string{"foo", "bar"},
	})
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
FROM users
WHERE name IN ($1, $2)
OR email IN ($1, $2)`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}
	if len(args) != 2 {
		t.Errorf("exec failed: values should have 2 length, but got %v", args)
	}
	if isInvalidString(args[0], "foo") {
		t.Errorf("exec failed: values should have 1, but got %v", args)
	}
	if isInvalidString(args[1], "bar") {
		t.Errorf("exec failed: values should have 1, but got %v", args)
	}
}

func TestContinuousInNamed(t *testing.T) {
	s := `SELECT *
FROM users
WHERE name IN /*% in "keywords" %*/('')
OR email IN /*% in "keywords" %*/('')`
	sql, args, err := New(Postgres).ExecNamed(s, map[string]interface{}{
		"keywords": []string{"foo", "bar"},
	})
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
FROM users
WHERE name IN (:keywords__1, :keywords__2)
OR email IN (:keywords__1, :keywords__2)`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
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
	s := `INSERT INTO users (
	    name
	  , created_at
	  , updated_at
	) VALUES (
	    /*% p "name" %*/'John Doe'
	  , /*% time %*/'2000-01-01'
	  , /*% time %*/'2000-01-01'
	)`
	sql, args, err := New(Postgres).Exec(s, map[string]interface{}{
		"name": "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	eSQL := `INSERT INTO users (
	    name
	  , created_at
	  , updated_at
	) VALUES (
	    $1
	  , $2
	  , $2
	)`
	et := time.Now()
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
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
	s := `INSERT INTO users (
	    name
	  , created_at
	  , updated_at
	) VALUES (
	    /*% p "name" %*/'John Doe'
	  , /*% time %*/'2000-01-01'
	  , /*% time %*/'2000-01-01'
	)`
	st := New(Postgres)
	st.TimeFunc = func() time.Time {
		return bt.AddDate(0, 0, 1)
	}
	sql, args, err := st.Exec(s, map[string]interface{}{
		"name": "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	eSQL := `INSERT INTO users (
	    name
	  , created_at
	  , updated_at
	) VALUES (
	    $1
	  , $2
	  , $2
	)`
	et := time.Now()
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
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
	s := `INSERT INTO users (
	    name
	  , created_at
	  , updated_at
	) VALUES (
	    /*% p "name" %*/'John Doe'
	  , /*% time %*/'2000-01-01'
	  , /*% time %*/'2000-01-01'
	)`
	sql, args, err := New(Postgres).ExecNamed(s, map[string]interface{}{
		"name": "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	eSQL := `INSERT INTO users (
	    name
	  , created_at
	  , updated_at
	) VALUES (
	    :name
	  , :time__
	  , :time__
	)`
	et := time.Now()
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
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
	s := `INSERT INTO users (
	    name
	  , created_at
	  , updated_at
	) VALUES (
	    /*% p "name" %*/'John Doe'
	  , /*% time %*/'2000-01-01'
	  , /*% time %*/'2000-01-01'
	)`
	sql, args, err := New(Postgres).ExecNamed(s, map[string]interface{}{
		"name": "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	eSQL := `INSERT INTO users (
	    name
	  , created_at
	  , updated_at
	) VALUES (
	    :name
	  , :time__
	  , :time__
	)`
	et := time.Now()
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
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
	s := `INSERT INTO users (
	    name
	  , created_at
	  , updated_at
	) VALUES (
	    /*% p "name" %*/'John Doe'
	  , /*% now %*/'2000-01-01'
	  , /*% now %*/'2000-01-01'
	)`
	sql, args, err := New(Postgres).Exec(s, map[string]interface{}{
		"name": "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	eSQL := `INSERT INTO users (
	    name
	  , created_at
	  , updated_at
	) VALUES (
	    $1
	  , $2
	  , $3
	)`
	et := time.Now()
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}
	tm1, ok := args[1].(time.Time)
	if !ok {
		t.Errorf("exec failed: 2nd arg expected time, but got %t", args[1])
	}
	if !(bt.Unix() <= tm1.Unix() && tm1.Unix() <= et.Unix()) {
		t.Errorf("time should be current time, but got %v", tm1)
	}
	tm2, ok := args[2].(time.Time)
	if !ok {
		t.Errorf("exec failed: 3rd arg expected time, but got %t", args[2])
	}
	if !(bt.Unix() <= tm2.Unix() && tm2.Unix() <= et.Unix()) {
		t.Errorf("time should be current time, but got %v", tm2)
	}
}

func TestNowNamed(t *testing.T) {
	bt := time.Now()
	s := `INSERT INTO users (
	    name
	  , created_at
	  , updated_at
	) VALUES (
	    /*% p "name" %*/'John Doe'
	  , /*% now %*/'2000-01-01'
	  , /*% now %*/'2000-01-01'
	)`
	sql, args, err := New(Postgres).ExecNamed(s, map[string]interface{}{
		"name": "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	eSQL := `INSERT INTO users (
	    name
	  , created_at
	  , updated_at
	) VALUES (
	    :name
	  , :now__1
	  , :now__2
	)`
	et := time.Now()
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}
	arg := args[1]
	if arg.Name != "now__1" {
		t.Errorf("default time arg name shoud be %q, but got %q", "now__1", arg.Name)
	}
	tm1, ok := arg.Value.(time.Time)
	if !ok {
		t.Errorf("exec failed: 2nd arg value expected time, but got %t", arg.Value)
	}
	if !(bt.Unix() <= tm1.Unix() && tm1.Unix() <= et.Unix()) {
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
	if !(bt.Unix() <= tm2.Unix() && tm2.Unix() <= et.Unix()) {
		t.Errorf("time should be current time, but got %v", tm2)
	}
}

func TestNowNamedWithNameParam(t *testing.T) {
	bt := time.Now()
	s := `INSERT INTO users (
	    name
	  , created_at
	  , updated_at
	) VALUES (
	    /*% p "name" %*/'John Doe'
	  , /*% now %*/'2000-01-01'
	  , /*% now %*/'2000-01-01'
	)`
	sql, args, err := New(Postgres).ExecNamed(s, map[string]interface{}{
		"name": "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	eSQL := `INSERT INTO users (
	    name
	  , created_at
	  , updated_at
	) VALUES (
	    :name
	  , :now__1
	  , :now__2
	)`
	et := time.Now()
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}
	arg := args[1]
	if arg.Name != "now__1" {
		t.Errorf("default time arg name shoud be %q, but got %q", "now__1", arg.Name)
	}
	tm1, ok := arg.Value.(time.Time)
	if !ok {
		t.Errorf("exec failed: 2nd args value expected time, but got %t", arg.Value)
	}
	if !(bt.Unix() <= tm1.Unix() && tm1.Unix() <= et.Unix()) {
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
	if !(bt.Unix() <= tm2.Unix() && tm2.Unix() <= et.Unix()) {
		t.Errorf("time should be current time, but got %v", tm2)
	}
}

func TestCustomFuncs(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE name LIKE /*% infix "name" %*/''
	/*% paging 3 50 %*/`
	sql, args, err := New(Postgres).AddFuncs(map[string]interface{}{
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

	eSQL := `SELECT *
	FROM users
	WHERE name LIKE '%' || $1 || '%' ESCAPE '\'
	OFFSET 3 LIMIT 50`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}
	if len(args) != 1 {
		t.Errorf("exec failed: values should have 1 length, but got %v", args)
	}
	if args[0] != `Alex` {
		t.Error("exec failed: embedded function should not be overwritten")
	}
}

func TestCustomFuncsContinuous(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE name LIKE /*% infix "name" %*/''
	/*% paging 3 50 %*/`
	sql, args, err := New(Postgres).AddFunc("paging", func(offset, limit int) string {
		return fmt.Sprintf("OFFSET %d LIMIT %d", offset, limit)
	}).AddFunc("infix", func() {
		panic("should not called")
	}).Exec(s, map[string]interface{}{
		"name": "Alex",
	})
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE name LIKE '%' || $1 || '%' ESCAPE '\'
	OFFSET 3 LIMIT 50`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}
	if len(args) != 1 {
		t.Errorf("exec failed: values should have 1 length, but got %v", args)
	}
	if args[0] != `Alex` {
		t.Error("exec failed: embedded function should not be overwritten")
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
		tt := tagt.New(t, d.tag)
		p := `/*% p "` + d.pArg + `" %*/''`
		s := fmt.Sprintf(`SELECT *
FROM users
WHERE first_name = %s OR last_name = %s`, p, p)
		sql, args, err := New(Postgres).Exec(s, map[string]interface{}{
			d.name: d.value,
		})
		if err != nil {
			t.Error(err)
		}
		eSQL := `SELECT *
FROM users
WHERE first_name = $1 OR last_name = $1`
		if eSQL != sql {
			tt.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
		}
		if len(args) != 1 {
			tt.Errorf("exec failed: values should have 1 length, but got %v", args)
		}
		if args[0] != "Alex" {
			tt.Errorf("exec failed: expected value %q, but got %q", "Alex", args[0])
		}
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
		tt := tagt.New(t, d.tag)
		p := `/*% p "` + d.pArg + `" %*/''`
		s := fmt.Sprintf(`SELECT *
FROM users
WHERE first_name = %s OR last_name = %s`, p, p)
		sql, args, err := New(Postgres).ExecNamed(s, map[string]interface{}{
			d.name: d.value,
		})
		if err != nil {
			t.Error(err)
		}
		eSQL := fmt.Sprintf(`SELECT *
FROM users
WHERE first_name = :%s OR last_name = :%s`, d.argName, d.argName)
		if eSQL != sql {
			tt.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
		}
		if len(args) != 1 {
			tt.Errorf("exec failed: values should have 1 length, but got %v", args)
		}
		if args[0].Name != d.argName {
			tt.Errorf("exec failed: expected args name %q, but got %q", "foo__Value", args[0].Name)
		}
		if s, ok := args[0].Value.(string); !ok || s != "Alex" {
			tt.Errorf("exec failed: expected value %q, but got %q", "Alex", args[0].Value)
		}
	}
}

func TestExecStructError(t *testing.T) {
	data := []struct {
		value  interface{}
		pArg   string
		errMsg string
		tag    string
	}{
		{Bar{Foo: Foo{"Alex"}}, "bar.Qux.Value", "unknown param: bar.Qux", "unknown"},
		{Bar{Foo: Foo{"Alex"}}, "bar.Value.Length", "not struct: bar.Value", "not struct"},
		{Bar{}, "bar.FooPtr.Value", "nil value: bar.FooPtr", "nil field"},
		{Bar{}, "bar.PropPtr.Value", "nil value: bar.PropPtr", "nil getter"},
		{Bar{}, "bar.Baz.Value", "nil value: bar.Baz", "nil interface"},
		{Bar{Foo: Foo{"Alex"}}, "bar.FnIn.Value", "invalid method: bar.FnIn", "invalid in arg num method"},
		{Bar{Foo: Foo{"Alex"}}, "bar.FnOut2.Value", "invalid method: bar.FnOut2", "invalid out val num method"},
		{nil, "bar.FooPtr.Value", "nil value: bar", "nil root"},
		{Bar{Foo: Foo{"Alex"}}, "baz.Foo.Value", "unknown param: baz", "unknown root"},
	}
	for _, d := range data {
		tt := tagt.New(t, d.tag)
		s := `SELECT *
FROM users
WHERE name = /*% p "` + d.pArg + `" %*/''`
		sql, _, err := New(Postgres).Exec(s, map[string]interface{}{
			"bar": d.value,
		})
		if err == nil {
			tt.Error("should raise error")
		}
		eSQL := fmt.Sprintf(`SELECT *
FROM users
WHERE name = /*! %s */`, d.errMsg)
		if eSQL != sql {
			tt.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
		}
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
