package sqlt

import (
	"testing"
)

func TestOracleP(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id = /*%p "id" %*/1`
	sql, args, err := New(Oracle).Exec(s, singleMap("id", 1))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE id = :1`
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

func TestOracleRepeatedP(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE family_name = /*%p "name" %*/'foo'
	OR given_name = /*%p "name" %*/'bar'
	OR nick_name = /*%p "name" %*/'baz'`
	sql, args, err := New(Oracle).Exec(s, singleMap("name", "test"))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE family_name = :1
	OR given_name = :1
	OR nick_name = :1`
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

func TestOraclePNamed(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id = /*%p "id" %*/1`
	sql, args, err := New(Oracle).ExecNamed(s, singleMap("id", 1))
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

func TestOracleRepeatedPNamed(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE family_name = /*%p "name" %*/'foo'
	OR given_name = /*%p "name" %*/'bar'
	OR nick_name = /*%p "name" %*/'baz'`
	sql, args, err := New(Oracle).ExecNamed(s, singleMap("name", "test"))
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

func TestOracleIn(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id IN /*%in "ids" %*/(1, 2)`
	sql, args, err := New(Oracle).Exec(s, singleMap("ids", []int{1, 2}))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE id IN (:1, :2)`
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

func TestOracleInNamed(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id IN /*%in "ids" %*/(1, 2)`
	sql, args, err := New(Oracle).ExecNamed(s, singleMap("ids", []int{1, 2}))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE id IN (:ids__1, :ids__2)`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}
	if len(args) != 2 {
		t.Errorf("exec failed: values should have 2 length, but got %v", args)
	}
	if isInvalidIntArg(args[0], "ids__1", 1) {
		t.Errorf("exec failed: values should have id = 1, but got %v", args)
	}
	if isInvalidIntArg(args[1], "ids__2", 2) {
		t.Errorf("exec failed: values should have id = 2, but got %v", args)
	}
}

func TestOracleInWithSingleValue(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id IN /*%in "ids" %*/(1, 2)`
	sql, args, err := New(Oracle).Exec(s, singleMap("ids", 1))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE id IN (:1)`
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

func TestOracleInNamedWithSingleValue(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id IN /*%in "ids" %*/(1, 2)`
	sql, args, err := New(Oracle).ExecNamed(s, singleMap("ids", 1))
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
		t.Errorf("exec failed: values should have 1 length, but got %v", args)
	}
	if isInvalidIntArg(args[0], "ids", 1) {
		t.Errorf("exec failed: values should have id = 1, but got %v", args)
	}
}

func TestOracleOtherTemplateFeature(t *testing.T) {
	s := `SELECT *
	FROM users
	WHERE id = /*%p "id" %*/1
	/*%- if .onlyMale %*/
	AND sex = 'MALE'
	/*%- end%*/
	ORDER BY /*% .order %*/id`
	sql, args, err := New(Oracle).Exec(s, map[string]interface{}{
		"id":       1,
		"order":    "name DESC",
		"onlyMale": true,
	})
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM users
	WHERE id = :1
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

func TestOracleLikeEscape(t *testing.T) {
	s := `SELECT *
	FROM items
	WHERE note1 LIKE /*% infix "note" %*/''
	OR note2 LIKE /*% prefix "note" %*/''
	OR note3 LIKE /*% suffix "note" %*/''`
	sql, args, err := New(Oracle).Exec(s, map[string]interface{}{
		"note": `abc%def_ghi％jkl＿mno[pqr\stu`,
	})
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT *
	FROM items
	WHERE note1 LIKE '%' || :1 || '%' ESCAPE '\'
	OR note2 LIKE :1 || '%' ESCAPE '\'
	OR note3 LIKE '%' || :1 ESCAPE '\'`
	if eSQL != sql {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, sql)
	}
	if len(args) != 1 {
		t.Errorf("exec failed: values should have 1 length, but got %v", args)
	}
	if args[0] != `abc\%def\_ghi\％jkl\＿mno[pqr\\stu` {
		t.Errorf("exec failed: escaped value %q is invalid", args[0])
	}
}
