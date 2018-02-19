package sqlt

import (
	"database/sql"
	"testing"
)

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

func TestP(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id = /*%p "id" %*/1`
	sql, vals, err := Exec(s, sql.Named("id", 1))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE id = $1`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}

	if len(vals) != 1 {
		t.Errorf("exec failed: values should have 1 length, but got %v", vals)
	}

	if v, ok := vals[0].(int); !ok || v != 1 {
		t.Errorf("exec failed: values should have 1, but got %v", vals)
	}
}

func TestPNamed(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id = /*%p "id" %*/1`
	sql, vals, err := ExecNamed(s, sql.Named("id", 1))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE id = :id`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}
	if len(vals) != 1 {
		t.Errorf("exec failed: values should have 1 length, but got %v", vals)
	}
	if vals[0].Name != "id" {
		t.Errorf("exec failed: named args should have arg named 'id', but got %v", vals)
	}
	v1 := vals[0]
	if v, ok := v1.Value.(int); v1.Name != "id" || !ok || v != 1 {
		t.Errorf("exec failed: values should have id = 1, but got %v", vals)
	}
}

func TestRepeatedPNamed(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE family_name = /*%p "name" %*/'foo'
	OR given_name = /*%p "name" %*/'bar'
	OR nick_name = /*%p "name" %*/'baz'`
	sql, vals, err := ExecNamed(s, sql.Named("name", "test"))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE family_name = :name
	OR given_name = :name
	OR nick_name = :name`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}
	if len(vals) != 1 {
		t.Errorf("exec failed: values should have 1 length, but got %v", vals)
	}
	if vals[0].Name != "name" {
		t.Errorf("exec failed: named args should have arg named 'name', but got %v", vals)
	}
	v1 := vals[0]
	if v, ok := v1.Value.(string); v1.Name != "name" || !ok || v != "test" {
		t.Errorf("exec failed: values should have name = 'test', but got %v", vals)
	}
}

func TestIn(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id IN /*%in "ids" %*/(1, 2)`
	sql, vals, err := Exec(s, sql.Named("ids", []int{1, 2}))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE id IN ($1,$2)`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}

	if len(vals) != 2 {
		t.Errorf("exec failed: values should have 2 length, but got %v", vals)
	}

	if v, ok := vals[0].(int); !ok || v != 1 {
		t.Errorf("exec failed: values should have 1, but got %v", vals)
	}

	if v, ok := vals[1].(int); !ok || v != 2 {
		t.Errorf("exec failed: values should have 2, but got %v", vals)
	}
}

func TestInNamed(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id IN /*%in "ids" %*/(1, 2)`
	sql, vals, err := ExecNamed(s, sql.Named("ids", []int{1, 2}))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE id IN (:ids1,:ids2)`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}

	if len(vals) != 2 {
		t.Errorf("exec failed: values should have 2 length, but got %v", vals)
	}

	v1 := vals[0]
	if v, ok := v1.Value.(int); v1.Name != "ids1" || !ok || v != 1 {
		t.Errorf("exec failed: values should have id = 1, but got %v", vals)
	}

	v2 := vals[1]
	if v, ok := v2.Value.(int); v2.Name != "ids2" || !ok || v != 2 {
		t.Errorf("exec failed: values should have id = 2, but got %v", vals)
	}
}

func TestInWithSingleValue(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id IN /*%in "ids" %*/(1, 2)`
	sql, vals, err := Exec(s, sql.Named("ids", 1))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE id IN ($1)`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}

	if len(vals) != 1 {
		t.Errorf("exec failed: values should have 2 length, but got %v", vals)
	}

	if v, ok := vals[0].(int); !ok || v != 1 {
		t.Errorf("exec failed: values should have 1, but got %v", vals)
	}
}

func TestInNamedWithSingleValue(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id IN /*%in "ids" %*/(1, 2)`
	sql, vals, err := ExecNamed(s, sql.Named("ids", 1))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE id IN (:ids)`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}

	if len(vals) != 1 {
		t.Errorf("exec failed: values should have 2 length, but got %v", vals)
	}

	v1 := vals[0]
	if v, ok := v1.Value.(int); v1.Name != "ids" || !ok || v != 1 {
		t.Errorf("exec failed: values should have id = 1, but got %v", vals)
	}
}

func TestOtherTemplateFeature(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id = /*%p "id" %*/1
	/*%- if .onlyMale %*/
	AND sex = 'MALE'
	/*%- end%*/
	ORDER BY /*% .order %*/id`
	sql, vals, err := Exec(s, sql.Named("id", 1), sql.Named("order", "name DESC"), sql.Named("onlyMale", true))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE id = $1
	AND sex = 'MALE'
	ORDER BY name DESC`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}

	if len(vals) != 1 {
		t.Errorf("exec failed: values should have 1 length, but got %v", vals)
	}

	if v, ok := vals[0].(int); !ok || v != 1 {
		t.Errorf("exec failed: values should have 1, but got %v", vals)
	}
}