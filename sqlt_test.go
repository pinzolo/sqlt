package sqlt

import (
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
	WHERE id IN ($1,$2,$3)
	AND name = $4
	AND sex = 'MALE'
	ORDER BY name DESC`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}
	if len(vals) != 4 {
		t.Errorf("exec failed: values should have 2 length, but got %v", vals)
	}
	if v, ok := vals[0].(int); !ok || v != 1 {
		t.Errorf("exec failed: values should have 1, but got %v", vals)
	}
	if v, ok := vals[1].(int); !ok || v != 2 {
		t.Errorf("exec failed: values should have 2, but got %v", vals)
	}
	if v, ok := vals[2].(int); !ok || v != 3 {
		t.Errorf("exec failed: values should have 3, but got %v", vals)
	}
	if v, ok := vals[3].(string); !ok || v != "Alex" {
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
	WHERE id IN (:ids1,:ids2,:ids3)
	AND name = :name
	ORDER BY name DESC`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}
	if len(vals) != 4 {
		t.Errorf("exec failed: values should have 2 length, but got %v", vals)
	}
	v1 := vals[0]
	if v, ok := v1.Value.(int); v1.Name != "ids1" || !ok || v != 1 {
		t.Errorf("exec failed: values should have ids1 = 1, but got %v", vals)
	}
	v2 := vals[1]
	if v, ok := v2.Value.(int); v2.Name != "ids2" || !ok || v != 2 {
		t.Errorf("exec failed: values should have ids2 = 2, but got %v", vals)
	}
	v3 := vals[2]
	if v, ok := v3.Value.(int); v3.Name != "ids3" || !ok || v != 3 {
		t.Errorf("exec failed: values should have ids3 = 3, but got %v", vals)
	}
	v4 := vals[3]
	if v, ok := v4.Value.(string); v4.Name != "name" || !ok || v != "Alex" {
		t.Errorf("exec failed: values should have name = 'Alex', but got %v", vals)
	}
}

func singleMap(k string, v interface{}) map[string]interface{} {
	return map[string]interface{}{
		k: v,
	}
}
