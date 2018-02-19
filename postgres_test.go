package sqlt

import (
	"database/sql"
	"testing"
)

func TestPostgresP(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id = /*%p "id" %*/1`
	sql, vals, err := New(Postgres).Exec(s, sql.Named("id", 1))
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

func TestPostgresRepeatedP(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE family_name = /*%p "name" %*/'foo'
	OR given_name = /*%p "name" %*/'bar'
	OR nick_name = /*%p "name" %*/'baz'`
	sql, vals, err := New(Postgres).Exec(s, sql.Named("name", "test"))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE family_name = $1
	OR given_name = $1
	OR nick_name = $1`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}
	if len(vals) != 1 {
		t.Errorf("exec failed: values should have 1 length, but got %v", vals)
	}
	if v, ok := vals[0].(string); !ok || v != "test" {
		t.Errorf("exec failed: values should have 'test', but got %v", vals)
	}
}

func TestPostgresPNamed(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id = /*%p "id" %*/1`
	sql, vals, err := New(Postgres).ExecNamed(s, sql.Named("id", 1))
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

func TestPostgresRepeatedPNamed(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE family_name = /*%p "name" %*/'foo'
	OR given_name = /*%p "name" %*/'bar'
	OR nick_name = /*%p "name" %*/'baz'`
	sql, vals, err := New(Postgres).ExecNamed(s, sql.Named("name", "test"))
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

func TestPostgresIn(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id IN /*%in "ids" %*/(1, 2)`
	sql, vals, err := New(Postgres).Exec(s, sql.Named("ids", []int{1, 2}))
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

func TestPostgresInNamed(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id IN /*%in "ids" %*/(1, 2)`
	sql, vals, err := New(Postgres).ExecNamed(s, sql.Named("ids", []int{1, 2}))
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

func TestPostgresInWithSingleValue(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id IN /*%in "ids" %*/(1, 2)`
	sql, vals, err := New(Postgres).Exec(s, sql.Named("ids", 1))
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

func TestPostgresInNamedWithSingleValue(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id IN /*%in "ids" %*/(1, 2)`
	sql, vals, err := New(Postgres).ExecNamed(s, sql.Named("ids", 1))
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

func TestPostgresOtherTemplateFeature(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id = /*%p "id" %*/1
	/*%- if .onlyMale %*/
	AND sex = 'MALE'
	/*%- end%*/
	ORDER BY /*% .order %*/id`
	sql, vals, err := New(Postgres).Exec(s, sql.Named("id", 1), sql.Named("order", "name DESC"), sql.Named("onlyMale", true))
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