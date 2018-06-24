package sqlt

import (
	"testing"
	"time"
)

func TestMySQLP(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id = /*%p "id" %*/1`
	sql, args, err := New(MySQL).Exec(s, singleMap("id", 1))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE id = ?`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}
	if len(args) != 1 {
		t.Errorf("exec failed: values should have 1 length, but got %v", args)
	}
	if isInvalidInt(args[0], 1) {
		t.Errorf("exec failed: values should have 1, but got %v", args)
	}
}

func TestMySQLRepeatedP(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE family_name = /*%p "name" %*/'foo'
	OR given_name = /*%p "name" %*/'bar'
	OR nick_name = /*%p "name" %*/'baz'`
	sql, args, err := New(MySQL).Exec(s, singleMap("name", "test"))
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
	if len(args) != 3 {
		t.Errorf("exec failed: values should have 3 length, but got %v", args)
	}
	for _, v := range args {
		if isInvalidString(v, "test") {
			t.Errorf("exec failed: values should have 'test', but got %v", args)
		}
	}
}

func TestMySQLPNamed(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id = /*%p "id" %*/1`
	sql, args, err := New(MySQL).ExecNamed(s, singleMap("id", 1))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE id = :id`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}
	if len(args) != 1 {
		t.Errorf("exec failed: values should have 1 length, but got %v", args)
	}
	if isInvalidIntArg(args[0], "id", 1) {
		t.Errorf("exec failed: values should have id = 1, but got %v", args)
	}
}

func TestMySQLRepeatedPNamed(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE family_name = /*%p "name" %*/'foo'
	OR given_name = /*%p "name" %*/'bar'
	OR nick_name = /*%p "name" %*/'baz'`
	sql, args, err := New(MySQL).ExecNamed(s, singleMap("name", "test"))
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
	if len(args) != 1 {
		t.Errorf("exec failed: values should have 1 length, but got %v", args)
	}
	if isInvalidStringArg(args[0], "name", "test") {
		t.Errorf("exec failed: values should have name = 'test', but got %v", args)
	}
}

func TestMySQLIn(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id IN /*%in "ids" %*/(1, 2)`
	sql, args, err := New(MySQL).Exec(s, singleMap("ids", []int{1, 2}))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE id IN (?, ?)`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}
	if len(args) != 2 {
		t.Errorf("exec failed: values should have 2 length, but got %v", args)
	}
	if isInvalidInt(args[0], 1) {
		t.Errorf("exec failed: values should have 1, but got %v", args)
	}
	if isInvalidInt(args[1], 2) {
		t.Errorf("exec failed: values should have 2, but got %v", args)
	}
}

func TestMySQLInNamed(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id IN /*%in "ids" %*/(1, 2)`
	sql, args, err := New(MySQL).ExecNamed(s, singleMap("ids", []int{1, 2}))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE id IN (:ids1, :ids2)`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}
	if len(args) != 2 {
		t.Errorf("exec failed: values should have 2 length, but got %v", args)
	}
	if isInvalidIntArg(args[0], "ids1", 1) {
		t.Errorf("exec failed: values should have id = 1, but got %v", args)
	}
	if isInvalidIntArg(args[1], "ids2", 2) {
		t.Errorf("exec failed: values should have id = 2, but got %v", args)
	}
}

func TestMySQLInWithSingleValue(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id IN /*%in "ids" %*/(1, 2)`
	sql, args, err := New(MySQL).Exec(s, singleMap("ids", 1))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE id IN (?)`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}
	if len(args) != 1 {
		t.Errorf("exec failed: values should have 2 length, but got %v", args)
	}
	if isInvalidInt(args[0], 1) {
		t.Errorf("exec failed: values should have 1, but got %v", args)
	}
}

func TestMySQLInNamedWithSingleValue(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id IN /*%in "ids" %*/(1, 2)`
	sql, args, err := New(MySQL).ExecNamed(s, singleMap("ids", 1))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE id IN (:ids)`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}
	if len(args) != 1 {
		t.Errorf("exec failed: values should have 2 length, but got %v", args)
	}
	if isInvalidIntArg(args[0], "ids", 1) {
		t.Errorf("exec failed: values should have id = 1, but got %v", args)
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
	sql, args, err := New(MySQL).Exec(s, map[string]interface{}{
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
	if len(args) != 1 {
		t.Errorf("exec failed: values should have 1 length, but got %v", args)
	}
	if isInvalidInt(args[0], 1) {
		t.Errorf("exec failed: values should have 1, but got %v", args)
	}
}

func TestMySQLTime(t *testing.T) {
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
	sql, args, err := New(MySQL).Exec(s, map[string]interface{}{
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
	    ?
	  , ?
	  , ?
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
	if tm1 != tm2 {
		t.Errorf("time should return same time on each calling, but tm1 is %v and tm2 is %v", tm1, tm2)
	}
}

func TestMySQLNow(t *testing.T) {
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
	sql, args, err := New(MySQL).Exec(s, map[string]interface{}{
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
	    ?
	  , ?
	  , ?
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
	if tm1 == tm2 {
		t.Errorf("now should not return same time on each calling")
	}
}
