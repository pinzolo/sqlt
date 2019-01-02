package sqlt_test

import (
	"testing"
	"time"

	"github.com/pinzolo/sqlt"
)

func TestMySQLP(t *testing.T) {
	s := `SELECT * FROM users WHERE id = /*%p "id" %*/1`
	query, args, err := sqlt.New(sqlt.MySQL).Exec(s, singleMap("id", 1))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT * FROM users WHERE id = ?`
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

func TestMySQLRepeatedP(t *testing.T) {
	s := `
SELECT *
FROM users
WHERE family_name = /*%p "name" %*/'foo'
OR given_name = /*%p "name" %*/'bar'
OR nick_name = /*%p "name" %*/'baz'`
	query, args, err := sqlt.New(sqlt.MySQL).Exec(s, singleMap("name", "test"))
	if err != nil {
		t.Error(err)
	}

	eSQL := `
SELECT *
FROM users
WHERE family_name = ?
OR given_name = ?
OR nick_name = ?`
	if eSQL != query {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
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
	s := `SELECT * FROM users WHERE id = /*%p "id" %*/1`
	query, args, err := sqlt.New(sqlt.MySQL).ExecNamed(s, singleMap("id", 1))
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

func TestMySQLRepeatedPNamed(t *testing.T) {
	s := `
SELECT *
FROM users
WHERE family_name = /*%p "name" %*/'foo'
OR given_name = /*%p "name" %*/'bar'
OR nick_name = /*%p "name" %*/'baz'`
	query, args, err := sqlt.New(sqlt.MySQL).ExecNamed(s, singleMap("name", "test"))
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

func TestMySQLIn(t *testing.T) {
	s := `SELECT * FROM users WHERE id IN /*%in "ids" %*/(1, 2)`
	query, args, err := sqlt.New(sqlt.MySQL).Exec(s, singleMap("ids", []int{1, 2}))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT * FROM users WHERE id IN (?, ?)`
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

func TestMySQLInNamed(t *testing.T) {
	s := `SELECT * FROM users WHERE id IN /*%in "ids" %*/(1, 2)`
	query, args, err := sqlt.New(sqlt.MySQL).ExecNamed(s, singleMap("ids", []int{1, 2}))
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

func TestMySQLInWithSingleValue(t *testing.T) {
	s := `SELECT * FROM users WHERE id IN /*%in "ids" %*/(1, 2)`
	query, args, err := sqlt.New(sqlt.MySQL).Exec(s, singleMap("ids", 1))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT * FROM users WHERE id IN (?)`
	if eSQL != query {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
	}
	if len(args) != 1 {
		t.Errorf("exec failed: values should have 2 length, but got %v", args)
	}
	if isInvalidInt(args[0], 1) {
		t.Errorf("exec failed: values should have 1, but got %v", args)
	}
}

func TestMySQLInNamedWithSingleValue(t *testing.T) {
	s := `SELECT * FROM users WHERE id IN /*%in "ids" %*/(1, 2)`
	query, args, err := sqlt.New(sqlt.MySQL).ExecNamed(s, singleMap("ids", 1))
	if err != nil {
		t.Error(err)
	}

	eSQL := `SELECT * FROM users WHERE id IN (:ids)`
	if eSQL != query {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
	}
	if len(args) != 1 {
		t.Errorf("exec failed: values should have 2 length, but got %v", args)
	}
	if isInvalidIntArg(args[0], "ids", 1) {
		t.Errorf("exec failed: values should have id = 1, but got %v", args)
	}
}

func TestMySQLOtherTemplateFeature(t *testing.T) {
	s := `
SELECT *
FROM users
WHERE id = /*%p "id" %*/1
/*%- if get "onlyMale" %*/
AND sex = 'MALE'
/*%- end%*/
ORDER BY /*% get "order" %*/id`
	query, args, err := sqlt.New(sqlt.MySQL).Exec(s, map[string]interface{}{
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
WHERE id = ?
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

func TestMySQLTime(t *testing.T) {
	bt := time.Now()
	s := `
INSERT INTO users (
	name
  , created_at
  , updated_at
) VALUES (
	/*% p "name" %*/'John Doe'
  , /*% time %*/'2000-01-01'
  , /*% time %*/'2000-01-01'
)`
	query, args, err := sqlt.New(sqlt.MySQL).Exec(s, map[string]interface{}{
		"name": "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	eSQL := `
INSERT INTO users (
	name
  , created_at
  , updated_at
) VALUES (
	?
  , ?
  , ?
)`
	et := time.Now()
	if eSQL != query {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
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
	s := `
INSERT INTO users (
	name
  , created_at
  , updated_at
) VALUES (
	/*% p "name" %*/'John Doe'
  , /*% now %*/'2000-01-01'
  , /*% now %*/'2000-01-01'
)`
	query, args, err := sqlt.New(sqlt.MySQL).Exec(s, map[string]interface{}{
		"name": "test",
	})
	if err != nil {
		t.Fatal(err)
	}
	eSQL := `
INSERT INTO users (
	name
  , created_at
  , updated_at
) VALUES (
	?
  , ?
  , ?
)`
	et := time.Now()
	if eSQL != query {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
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

func TestMySQLLikeEscape(t *testing.T) {
	s := `
SELECT *
FROM items
WHERE note1 LIKE /*% infix "note" %*/''
OR note2 LIKE /*% prefix "note" %*/''
OR note3 LIKE /*% suffix "note" %*/''
OR note4 = /*% p "note" %*/''`
	query, args, err := sqlt.New(sqlt.MySQL).Exec(s, map[string]interface{}{
		"note": `abc%def_ghi％jkl＿mno[pqr\stu`,
	})
	if err != nil {
		t.Error(err)
	}

	eSQL := `
SELECT *
FROM items
WHERE note1 LIKE '%' || ? || '%' ESCAPE '\'
OR note2 LIKE ? || '%' ESCAPE '\'
OR note3 LIKE '%' || ? ESCAPE '\'
OR note4 = ?`
	if eSQL != query {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
	}
	if len(args) != 4 {
		t.Errorf("exec failed: values should have 1 length, but got %v", args)
	}
	if args[0] != `abc\%def\_ghi％jkl＿mno[pqr\\stu` {
		t.Errorf("exec failed: 1st value %q is invalid", args[0])
	}
	if args[1] != `abc\%def\_ghi％jkl＿mno[pqr\\stu` {
		t.Errorf("exec failed: 2nd value %q is invalid", args[1])
	}
	if args[2] != `abc\%def\_ghi％jkl＿mno[pqr\\stu` {
		t.Errorf("exec failed: 3rd value %q is invalid", args[2])
	}
	if args[3] != `abc%def_ghi％jkl＿mno[pqr\stu` {
		t.Errorf("exec failed: 4th value %q is invalid", args[3])
	}
}

func TestMySQLLikeEscapeNamed(t *testing.T) {
	s := `
SELECT *
FROM items
WHERE note1 LIKE /*% infix "note" %*/''
OR note2 LIKE /*% prefix "note" %*/''
OR note3 LIKE /*% suffix "note" %*/''
OR note4 = /*% p "note" %*/''`
	query, args, err := sqlt.New(sqlt.MySQL).ExecNamed(s, map[string]interface{}{
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
	if isInvalidStringArg(args[0], "note__esc", `abc\%def\_ghi％jkl＿mno[pqr\\stu`) {
		t.Errorf("exec failed: escaped value %v is invalid", args[0])
	}
	if isInvalidStringArg(args[1], "note", `abc%def_ghi％jkl＿mno[pqr\stu`) {
		t.Errorf("exec failed: escaped value %v is invalid", args[1])
	}
}

func TestMySQLLikeEscapeWithoutWildcard(t *testing.T) {
	s := `
SELECT *
FROM items
WHERE note1 LIKE /*% infix "note" %*/''
OR note2 LIKE /*% prefix "note" %*/''
OR note3 LIKE /*% suffix "note" %*/''
OR note4 = /*% p "note" %*/''`
	query, args, err := sqlt.New(sqlt.MySQL).Exec(s, map[string]interface{}{
		"note": `abcde`,
	})
	if err != nil {
		t.Error(err)
	}

	eSQL := `
SELECT *
FROM items
WHERE note1 LIKE '%' || ? || '%' ESCAPE '\'
OR note2 LIKE ? || '%' ESCAPE '\'
OR note3 LIKE '%' || ? ESCAPE '\'
OR note4 = ?`
	if eSQL != query {
		t.Errorf("exec failed: expected %s, but got %s", eSQL, query)
	}
	if len(args) != 4 {
		t.Error("exec failed: when not exist wildcard char should reuse original")
	}
	if args[0] != `abcde` {
		t.Errorf("exec failed: 1st value %q is invalid", args[0])
	}
	if args[1] != `abcde` {
		t.Errorf("exec failed: 2nd value %q is invalid", args[1])
	}
	if args[2] != `abcde` {
		t.Errorf("exec failed: 3rd value %q is invalid", args[2])
	}
	if args[3] != `abcde` {
		t.Errorf("exec failed: 4th value %q is invalid", args[3])
	}
}

func TestMySQLLikeEscapeNamedWithoutWildcard(t *testing.T) {
	s := `
SELECT *
FROM items
WHERE note1 LIKE /*% infix "note" %*/''
OR note2 LIKE /*% prefix "note" %*/''
OR note3 LIKE /*% suffix "note" %*/''
OR note4 = /*% p "note" %*/''`
	query, args, err := sqlt.New(sqlt.MySQL).ExecNamed(s, map[string]interface{}{
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
		t.Errorf("exec failed: escaped value %v is invalid", args[0])
	}
}
