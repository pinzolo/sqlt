package sqlt

import (
	"testing"
)

func TestMySQLP(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id = /*%p "id" %*/1`
	sql, vals, err := New(MySQL).Exec(s, singleMap("id", 1))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE id = ?`
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

func TestMySQLRepeatedP(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE family_name = /*%p "name" %*/'foo'
	OR given_name = /*%p "name" %*/'bar'
	OR nick_name = /*%p "name" %*/'baz'`
	sql, vals, err := New(MySQL).Exec(s, singleMap("name", "test"))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE family_name = ?
	OR given_name = ?
	OR nick_name = ?`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}
	if len(vals) != 3 {
		t.Errorf("exec failed: values should have 3 length, but got %v", vals)
	}
	for _, v := range vals {
		if v, ok := v.(string); !ok || v != "test" {
			t.Errorf("exec failed: values should have 'test', but got %v", vals)
		}
	}
}

func TestMySQLPNamed(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id = /*%p "id" %*/1`
	sql, vals, err := New(MySQL).ExecNamed(s, singleMap("id", 1))
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

func TestMySQLRepeatedPNamed(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE family_name = /*%p "name" %*/'foo'
	OR given_name = /*%p "name" %*/'bar'
	OR nick_name = /*%p "name" %*/'baz'`
	sql, vals, err := New(MySQL).ExecNamed(s, singleMap("name", "test"))
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

func TestMySQLIn(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id IN /*%in "ids" %*/(1, 2)`
	sql, vals, err := New(MySQL).Exec(s, singleMap("ids", []int{1, 2}))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE id IN (?,?)`
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

func TestMySQLInNamed(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id IN /*%in "ids" %*/(1, 2)`
	sql, vals, err := New(MySQL).ExecNamed(s, singleMap("ids", []int{1, 2}))
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

func TestMySQLInWithSingleValue(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id IN /*%in "ids" %*/(1, 2)`
	sql, vals, err := New(MySQL).Exec(s, singleMap("ids", 1))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE id IN (?)`
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

func TestMySQLInNamedWithSingleValue(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id IN /*%in "ids" %*/(1, 2)`
	sql, vals, err := New(MySQL).ExecNamed(s, singleMap("ids", 1))
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

func TestMySQLOtherTemplateFeature(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id = /*%p "id" %*/1
	/*%- if .onlyMale %*/
	AND sex = 'MALE'
	/*%- end%*/
	ORDER BY /*% .order %*/id`
	sql, vals, err := New(MySQL).Exec(s, map[string]interface{}{
		"id":       1,
		"order":    "name DESC",
		"onlyMale": true,
	})
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE id = ?
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
