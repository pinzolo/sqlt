package sqlt_test

import (
	"testing"

	"github.com/pinzolo/sqlt"
)

func TestOracleP(t *testing.T) {
	s := `SELECT * FROM users WHERE id = /*%p "id" %*/1`
	query, args, err := sqlt.New(sqlt.Oracle).Exec(s, singleMap("id", 1))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT * FROM users WHERE id = :1`
	if eSQL != query {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
	}
	if len(args) != 1 {
		t.Errorf("exec failed: values should have 1 length, but got %v", args)
	}
	if isInvalidInt(args[0], 1) {
		t.Errorf("exec failed: values should have 1, but got %v", args)
	}
}

func TestOracleRepeatedP(t *testing.T) {
	s := `
SELECT *
FROM users
WHERE family_name = /*%p "name" %*/'foo'
OR given_name = /*%p "name" %*/'bar'
OR nick_name = /*%p "name" %*/'baz'`
	query, args, err := sqlt.New(sqlt.Oracle).Exec(s, singleMap("name", "test"))
	if err != nil {
		t.Error(err)
	}

	eSQL := `
SELECT *
FROM users
WHERE family_name = :1
OR given_name = :1
OR nick_name = :1`
	if eSQL != query {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
	}
	if len(args) != 1 {
		t.Errorf("exec failed: values should have 1 length, but got %v", args)
	}
	if isInvalidString(args[0], "test") {
		t.Errorf("exec failed: values should have 'test', but got %v", args)
	}
}

func TestOraclePNamed(t *testing.T) {
	s := `SELECT * FROM users WHERE id = /*%p "id" %*/1`
	query, args, err := sqlt.New(sqlt.Oracle).ExecNamed(s, singleMap("id", 1))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT * FROM users WHERE id = :id`
	if eSQL != query {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
	}
	if len(args) != 1 {
		t.Errorf("exec failed: values should have 1 length, but got %v", args)
	}
	if isInvalidIntArg(args[0], "id", 1) {
		t.Errorf("exec failed: values should have id = 1, but got %v", args)
	}
}

func TestOracleRepeatedPNamed(t *testing.T) {
	s := `
SELECT *
FROM users
WHERE family_name = /*%p "name" %*/'foo'
OR given_name = /*%p "name" %*/'bar'
OR nick_name = /*%p "name" %*/'baz'`
	query, args, err := sqlt.New(sqlt.Oracle).ExecNamed(s, singleMap("name", "test"))
	if err != nil {
		t.Error(err)
	}

	eSQL := `
SELECT *
FROM users
WHERE family_name = :name
OR given_name = :name
OR nick_name = :name`
	if eSQL != query {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
	}
	if len(args) != 1 {
		t.Errorf("exec failed: values should have 1 length, but got %v", args)
	}
	if isInvalidStringArg(args[0], "name", "test") {
		t.Errorf("exec failed: values should have name = 'test', but got %v", args)
	}
}

func TestOracleIn(t *testing.T) {
	s := `SELECT * FROM users WHERE id IN /*%in "ids" %*/(1, 2)`
	query, args, err := sqlt.New(sqlt.Oracle).Exec(s, singleMap("ids", []int{1, 2}))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT * FROM users WHERE id IN (:1, :2)`
	if eSQL != query {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
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
	s := `SELECT * FROM users WHERE id IN /*%in "ids" %*/(1, 2)`
	query, args, err := sqlt.New(sqlt.Oracle).ExecNamed(s, singleMap("ids", []int{1, 2}))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT * FROM users WHERE id IN (:ids__1, :ids__2)`
	if eSQL != query {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
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
	s := `SELECT * FROM users WHERE id IN /*%in "ids" %*/(1, 2)`
	query, args, err := sqlt.New(sqlt.Oracle).Exec(s, singleMap("ids", 1))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT * FROM users WHERE id IN (:1)`
	if eSQL != query {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
	}
	if len(args) != 1 {
		t.Errorf("exec failed: values should have 1 length, but got %v", args)
	}
	if isInvalidInt(args[0], 1) {
		t.Errorf("exec failed: values should have 1, but got %v", args)
	}
}

func TestOracleInNamedWithSingleValue(t *testing.T) {
	s := `SELECT * FROM users WHERE id IN /*%in "ids" %*/(1, 2)`
	query, args, err := sqlt.New(sqlt.Oracle).ExecNamed(s, singleMap("ids", 1))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT * FROM users WHERE id IN (:ids)`
	if eSQL != query {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
	}
	if len(args) != 1 {
		t.Errorf("exec failed: values should have 1 length, but got %v", args)
	}
	if isInvalidIntArg(args[0], "ids", 1) {
		t.Errorf("exec failed: values should have id = 1, but got %v", args)
	}
}

func TestOracleOtherTemplateFeature(t *testing.T) {
	s := `
SELECT *
FROM users
WHERE id = /*%p "id" %*/1
/*%- if get "onlyMale" %*/
AND sex = 'MALE'
/*%- end%*/
ORDER BY /*% get "order" %*/id`
	query, args, err := sqlt.New(sqlt.Oracle).Exec(s, map[string]interface{}{
		"id":       1,
		"order":    "name DESC",
		"onlyMale": true,
	})
	if err != nil {
		t.Error(err)
	}

	eSQL := `
SELECT *
FROM users
WHERE id = :1
AND sex = 'MALE'
ORDER BY name DESC`
	if eSQL != query {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
	}
	if len(args) != 1 {
		t.Errorf("exec failed: values should have 1 length, but got %v", args)
	}
	if isInvalidInt(args[0], 1) {
		t.Errorf("exec failed: values should have 1, but got %v", args)
	}
}

func TestOracleLikeEscape(t *testing.T) {
	s := `
SELECT *
FROM items
WHERE note1 LIKE /*% infix "note" %*/''
OR note2 LIKE /*% prefix "note" %*/''
OR note3 LIKE /*% suffix "note" %*/''
OR note4 = /*% p "note" %*/''`
	query, args, err := sqlt.New(sqlt.Oracle).Exec(s, map[string]interface{}{
		"note": `abc%def_ghi％jkl＿mno[pqr\stu`,
	})
	if err != nil {
		t.Error(err)
	}

	eSQL := `
SELECT *
FROM items
WHERE note1 LIKE '%' || :1 || '%' ESCAPE '\'
OR note2 LIKE :1 || '%' ESCAPE '\'
OR note3 LIKE '%' || :1 ESCAPE '\'
OR note4 = :2`
	if eSQL != query {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
	}
	if len(args) != 2 {
		t.Errorf("exec failed: values should have 1 length, but got %v", args)
	}
	if args[0] != `abc\%def\_ghi\％jkl\＿mno[pqr\\stu` {
		t.Errorf("exec failed: 1st value %q is invalid", args[0])
	}
	if args[1] != `abc%def_ghi％jkl＿mno[pqr\stu` {
		t.Errorf("exec failed: 2nd value %q is invalid", args[1])
	}
}

func TestOracleLikeEscapeNamed(t *testing.T) {
	s := `
SELECT *
FROM items
WHERE note1 LIKE /*% infix "note" %*/''
OR note2 LIKE /*% prefix "note" %*/''
OR note3 LIKE /*% suffix "note" %*/''
OR note4 = /*% p "note" %*/''`
	query, args, err := sqlt.New(sqlt.Oracle).ExecNamed(s, map[string]interface{}{
		"note": `abc%def_ghi％jkl＿mno[pqr\stu`,
	})
	if err != nil {
		t.Error(err)
	}

	eSQL := `
SELECT *
FROM items
WHERE note1 LIKE '%' || :note__esc || '%' ESCAPE '\'
OR note2 LIKE :note__esc || '%' ESCAPE '\'
OR note3 LIKE '%' || :note__esc ESCAPE '\'
OR note4 = :note`
	if eSQL != query {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
	}
	if len(args) != 2 {
		t.Errorf("exec failed: values should have 1 length, but got %v", args)
	}
	if isInvalidStringArg(args[0], "note__esc", `abc\%def\_ghi\％jkl\＿mno[pqr\\stu`) {
		t.Errorf("exec failed: 1st value %v is invalid", args[0])
	}
	if isInvalidStringArg(args[1], "note", `abc%def_ghi％jkl＿mno[pqr\stu`) {
		t.Errorf("exec failed: 2nd value %v is invalid", args[1])
	}
}

func TestOracleLikeEscapeWithoutWildcard(t *testing.T) {
	s := `
SELECT *
FROM items
WHERE note1 LIKE /*% infix "note" %*/''
OR note2 LIKE /*% prefix "note" %*/''
OR note3 LIKE /*% suffix "note" %*/''
OR note4 = /*% p "note" %*/''`
	query, args, err := sqlt.New(sqlt.Oracle).Exec(s, map[string]interface{}{
		"note": `abcde`,
	})
	if err != nil {
		t.Error(err)
	}

	eSQL := `
SELECT *
FROM items
WHERE note1 LIKE '%' || :1 || '%' ESCAPE '\'
OR note2 LIKE :1 || '%' ESCAPE '\'
OR note3 LIKE '%' || :1 ESCAPE '\'
OR note4 = :1`
	if eSQL != query {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
	}
	if len(args) != 1 {
		t.Error("exec failed: when not exist wildcard char should reuse original")
	}
	if args[0] != `abcde` {
		t.Errorf("exec failed: 1st value %q is invalid", args[0])
	}
}

func TestOracleLikeEscapeNamedWithoutWildcard(t *testing.T) {
	s := `
SELECT *
FROM items
WHERE note1 LIKE /*% infix "note" %*/''
OR note2 LIKE /*% prefix "note" %*/''
OR note3 LIKE /*% suffix "note" %*/''
OR note4 = /*% p "note" %*/''`
	query, args, err := sqlt.New(sqlt.Oracle).ExecNamed(s, map[string]interface{}{
		"note": `abcde`,
	})
	if err != nil {
		t.Error(err)
	}

	eSQL := `
SELECT *
FROM items
WHERE note1 LIKE '%' || :note || '%' ESCAPE '\'
OR note2 LIKE :note || '%' ESCAPE '\'
OR note3 LIKE '%' || :note ESCAPE '\'
OR note4 = :note`
	if eSQL != query {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
	}
	if len(args) != 1 {
		t.Error("exec failed: when not exist wildcard char should reuse original")
	}
	if isInvalidStringArg(args[0], "note", `abcde`) {
		t.Errorf("exec failed: 1st value %v is invalid", args[0])
	}
}
