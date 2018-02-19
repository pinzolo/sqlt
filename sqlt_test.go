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
		_, _, err := New(Postgres).Exec(s, sql.Named("id", 1), sql.Named("order", "name DESC"), sql.Named("onlyMale", true), sql.Named("name", "Alex"))
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
		_, _, err := New(Postgres).ExecNamed(s, sql.Named("id", 1), sql.Named("order", "name DESC"), sql.Named("onlyMale", true), sql.Named("name", "Alex"))
		if err != nil {
			b.Error(err)
		}
	}
}
func BenchmarkExecWithMap(b *testing.B) {
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
		_, _, err := New(Postgres).ExecWithMap(s, m)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkExecNamedWithMap(b *testing.B) {
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
		_, _, err := New(Postgres).ExecNamedWithMap(s, m)
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

func TestExecWithMap(t *testing.T) {
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
	sql, vals, err := New(Postgres).ExecWithMap(s, m)
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE id = $1
	AND name = $2
	AND sex = 'MALE'
	ORDER BY name DESC`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}
	if len(vals) != 2 {
		t.Errorf("exec failed: values should have 2 length, but got %v", vals)
	}
	if v, ok := vals[0].(int); !ok || v != 1 {
		t.Errorf("exec failed: values should have 1, but got %v", vals)
	}
	if v, ok := vals[1].(string); !ok || v != "Alex" {
		t.Errorf("exec failed: values should have 'Alex', but got %v", vals)
	}
}

func TestExecNamedWithMap(t *testing.T) {
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
		"onlyMale": false,
		"name":     "Alex",
	}
	sql, vals, err := New(Postgres).ExecNamedWithMap(s, m)
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE id = :id
	AND name = :name
	ORDER BY name DESC`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}
	if len(vals) != 2 {
		t.Errorf("exec failed: values should have 2 length, but got %v", vals)
	}
	v1 := vals[0]
	if v, ok := v1.Value.(int); v1.Name != "id" || !ok || v != 1 {
		t.Errorf("exec failed: values should have id = 1, but got %v", vals)
	}
	v2 := vals[1]
	if v, ok := v2.Value.(string); v2.Name != "name" || !ok || v != "Alex" {
		t.Errorf("exec failed: values should have name = 'Alex', but got %v", vals)
	}
}
