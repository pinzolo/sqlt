package sqlt

import (
	"database/sql"
	"testing"
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
	sql, vals, err := New(Postgres).Exec(s, map[string]interface{}{
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
	if len(vals) != 4 {
		t.Errorf("exec failed: values should have 4 length, but got %v", vals)
	}
	if isInvalidInt(vals[0], 1) {
		t.Errorf("exec failed: values should have 1, but got %v", vals)
	}
	if isInvalidInt(vals[1], 2) {
		t.Errorf("exec failed: values should have 2, but got %v", vals)
	}
	if isInvalidInt(vals[2], 3) {
		t.Errorf("exec failed: values should have 3, but got %v", vals)
	}
	if isInvalidString(vals[3], "Alex") {
		t.Errorf("exec failed: values should have 'Alex', but got %v", vals)
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
	sql, vals, err := New(Postgres).ExecNamed(s, map[string]interface{}{
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
	WHERE id IN (:ids1, :ids2, :ids3)
	AND name = :name
	ORDER BY name DESC`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}
	if len(vals) != 4 {
		t.Errorf("exec failed: values should have 4 length, but got %v", vals)
	}
	if isInvalidIntArg(vals[0], "ids1", 1) {
		t.Errorf("exec failed: values should have ids1 = 1, but got %v", vals)
	}
	if isInvalidIntArg(vals[1], "ids2", 2) {
		t.Errorf("exec failed: values should have ids2 = 2, but got %v", vals)
	}
	if isInvalidIntArg(vals[2], "ids3", 3) {
		t.Errorf("exec failed: values should have ids3 = 3, but got %v", vals)
	}
	if isInvalidStringArg(vals[3], "name", "Alex") {
		t.Errorf("exec failed: values should have name = 'Alex', but got %v", vals)
	}
}

func TestExecWithNilParams(t *testing.T) {
	s := `SELECT *
	FROM users`
	sql, vals, err := New(Postgres).Exec(s, nil)
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}
	if len(vals) != 0 {
		t.Errorf("exec failed: values should have 0 length, but got %v", vals)
	}
}

func TestExecNamedWithNilParams(t *testing.T) {
	s := `SELECT *
	FROM users`
	sql, vals, err := New(Postgres).ExecNamed(s, nil)
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}
	if len(vals) != 0 {
		t.Errorf("exec failed: values should have 0 length, but got %v", vals)
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
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE id IN ($1, $2, $3)
	AND name = /*! userName is unknown */
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
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE id IN /*! idList is unknown */
	AND name = $1
	AND sex = 'MALE'
	ORDER BY name DESC`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}
}

func TestQueryArgs(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id IN /*% in "ids" %*/(1, 2)
	AND name = /*% p "name" %*/'John Doe'
	/*%- if .onlyMale %*/
	AND sex = 'MALE'
	/*%- end%*/
	ORDER BY /*% .order %*/id`
	sql, vals, err := New(Postgres).Exec(s, map[string]interface{}{
		"ids":      []int{1, 2, 3},
		"order":    "name DESC",
		"onlyMale": true,
		"name":     "Alex",
	})
	if err != nil {
		t.Error(err)
	}

	args := QueryArgs(sql, vals)
	if len(args) != 5 {
		t.Errorf("QueryArgs failed: arguments should have 5 length, but got %d", len(args))
	}
	if arg, ok := args[0].(string); ok {
		if arg != sql {
			t.Errorf("1st argument should be sql, but got %s", arg)
		}
	} else {
		t.Errorf("1st argument should be string.")
	}
	if isInvalidInt(args[1], 1) {
		t.Errorf("QueryArgs failed: arguments should have 1, but got %v", args)
	}
	if isInvalidInt(args[2], 2) {
		t.Errorf("QueryArgs failed: arguments should have 2, but got %v", args)
	}
	if isInvalidInt(args[3], 3) {
		t.Errorf("QueryArgs failed: arguments should have 3, but got %v", args)
	}
	if isInvalidString(args[4], "Alex") {
		t.Errorf("QueryArgs failed: arguments should have 'Alex', but got %v", args)
	}
}

func TestQueryArgsNamed(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id IN /*% in "ids" %*/(1, 2)
	AND name = /*% p "name" %*/'John Doe'
	/*%- if .onlyMale %*/
	AND sex = 'MALE'
	/*%- end%*/
	ORDER BY /*% .order %*/id`
	query, vals, err := New(Postgres).ExecNamed(s, map[string]interface{}{
		"ids":      []int{1, 2, 3},
		"order":    "name DESC",
		"onlyMale": false,
		"name":     "Alex",
	})
	if err != nil {
		t.Error(err)
	}

	args := QueryArgsNamed(query, vals)
	if len(args) != 5 {
		t.Errorf("QueryArgsNamed failed: arguments should have 5 length, but got %d", len(args))
	}
	if arg, ok := args[0].(string); ok {
		if arg != query {
			t.Errorf("1st argument should be sql, but got %s", arg)
		}
	} else {
		t.Errorf("1st argument should be string.")
	}
	if arg, ok := args[1].(sql.NamedArg); ok {
		if isInvalidIntArg(arg, "ids1", 1) {
			t.Errorf("QueryArgsNamed failed: arguments should have ids1 = 1, but got %v", arg)
		}
	} else {
		t.Error("QueryArgsNamed failed: 2nd argument should be named args")
	}
	if arg, ok := args[2].(sql.NamedArg); ok {
		if isInvalidIntArg(arg, "ids2", 2) {
			t.Errorf("QueryArgsNamed failed: arguments should have ids2 = 2, but got %v", arg)
		}
	} else {
		t.Error("QueryArgsNamed failed: 3rd argument should be named args")
	}
	if arg, ok := args[3].(sql.NamedArg); ok {
		if isInvalidIntArg(arg, "ids3", 3) {
			t.Errorf("QueryArgsNamed failed: arguments should have ids3 = 3, but got %v", arg)
		}
	} else {
		t.Error("QueryArgsNamed failed: 4th argument should be named args")
	}
	if arg, ok := args[4].(sql.NamedArg); ok {
		if isInvalidStringArg(arg, "name", "Alex") {
			t.Errorf("QueryArgsNamed failed: arguments should have name = 'Alex', but got %v", arg)
		}
	} else {
		t.Error("QueryArgsNamed failed: 5th argument should be named args")
	}
}

func TestQueryArgsWithoutValues(t *testing.T) {
	s := "SELECT * FROM users"
	params := [][]interface{}{
		nil,
		[]interface{}{},
	}
	for _, p := range params {
		args := QueryArgs(s, p)
		if len(args) != 1 {
			t.Error("QueryArgs without values should returns only sql.")
		}
		if arg, ok := args[0].(string); ok {
			if arg != s {
				t.Errorf("1st argument should be sql, but got %s", arg)
			}
		} else {
			t.Errorf("1st argument should be string.")
		}
	}
}

func TestQueryArgsNamedWithoutValues(t *testing.T) {
	s := "SELECT * FROM users"
	params := [][]sql.NamedArg{
		nil,
		[]sql.NamedArg{},
	}
	for _, p := range params {
		args := QueryArgsNamed(s, p)
		if len(args) != 1 {
			t.Error("QueryArgsNamed without values should returns only sql.")
		}
		if arg, ok := args[0].(string); ok {
			if arg != s {
				t.Errorf("1st argument should be sql, but got %s", arg)
			}
		} else {
			t.Errorf("1st argument should be string.")
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
