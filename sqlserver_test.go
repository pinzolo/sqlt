package sqlt

import (
	"testing"
)

func TestSQLServerP(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id = /*%p "id" %*/1`
	sql, args, err := New(SQLServer).Exec(s, singleMap("id", 1))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE id = @p1`
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

func TestSQLServerRepeatedP(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE family_name = /*%p "name" %*/'foo'
	OR given_name = /*%p "name" %*/'bar'
	OR nick_name = /*%p "name" %*/'baz'`
	sql, args, err := New(SQLServer).Exec(s, singleMap("name", "test"))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE family_name = @p1
	OR given_name = @p1
	OR nick_name = @p1`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}
	if len(args) != 1 {
		t.Errorf("exec failed: values should have 1 length, but got %v", args)
	}
	if isInvalidString(args[0], "test") {
		t.Errorf("exec failed: values should have 'test', but got %v", args)
	}
}

func TestSQLServerPNamed(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id = /*%p "id" %*/1`
	sql, args, err := New(SQLServer).ExecNamed(s, singleMap("id", 1))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE id = @id`
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

func TestSQLServerRepeatedPNamed(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE family_name = /*%p "name" %*/'foo'
	OR given_name = /*%p "name" %*/'bar'
	OR nick_name = /*%p "name" %*/'baz'`
	sql, args, err := New(SQLServer).ExecNamed(s, singleMap("name", "test"))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE family_name = @name
	OR given_name = @name
	OR nick_name = @name`
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

func TestSQLServerIn(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id IN /*%in "ids" %*/(1, 2)`
	sql, args, err := New(SQLServer).Exec(s, singleMap("ids", []int{1, 2}))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE id IN (@p1, @p2)`
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

func TestSQLServerInNamed(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id IN /*%in "ids" %*/(1, 2)`
	sql, args, err := New(SQLServer).ExecNamed(s, singleMap("ids", []int{1, 2}))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE id IN (@ids1, @ids2)`
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

func TestSQLServerInWithSingleValue(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id IN /*%in "ids" %*/(1, 2)`
	sql, args, err := New(SQLServer).Exec(s, singleMap("ids", 1))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE id IN (@p1)`
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

func TestSQLServerInNamedWithSingleValue(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id IN /*%in "ids" %*/(1, 2)`
	sql, args, err := New(SQLServer).ExecNamed(s, singleMap("ids", 1))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE id IN (@ids)`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}
	if len(args) != 1 {
		t.Errorf("exec failed: values should have 1 length, but got %v", args)
	}
	if isInvalidIntArg(args[0], "ids", 1) {
		t.Errorf("exec failed: values should have id = 1, but got %v", args)
	}
}

func TestSQLServerOtherTemplateFeature(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id = /*%p "id" %*/1
	/*%- if .onlyMale %*/
	AND sex = 'MALE'
	/*%- end%*/
	ORDER BY /*% .order %*/id`
	sql, args, err := New(SQLServer).Exec(s, map[string]interface{}{
		"id":       1,
		"order":    "name DESC",
		"onlyMale": true,
	})
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE id = @p1
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
