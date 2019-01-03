# sqlt

[![Build Status](https://travis-ci.org/pinzolo/sqlt.png)](http://travis-ci.org/pinzolo/sqlt)
[![Coverage Status](https://coveralls.io/repos/github/pinzolo/sqlt/badge.svg?branch=master)](https://coveralls.io/github/pinzolo/sqlt?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/pinzolo/sqlt)](https://goreportcard.com/report/github.com/pinzolo/sqlt)
[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/pinzolo/sqlt)
[![license](http://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/pinzolo/sqlt/master/LICENSE)

## Description

Simple SQL template for 2 way SQL.

## Sample (PostgreSQL)

### Template

Template is compatible with SQL.  
Use `/*%` and `%*/` as delimiter instead of `{{` and `}}` for processing template.

```sql
SELECT *
FROM users
WHERE id IN /*% in "ids" %*/(1, 2)
AND name = /*% p "name" %*/'John Doe'
/*%- if get "onlyMale" %*/
AND sex = 'MALE'
/*%- end %*/
ORDER BY /*% out "order" %*/id
```

### Go code

* func `param` or `p` replace to placeholder by name.
* func `in` deploy slice values to parentheses and placeholders.
* func `time` returns current time and cache it, this func always same time in same template.
* func `now` returns current time each calling.
* func `escape`, `prefix`, `inffix`, `suffix` replace to placeholder with escape for `LIKE` keyword.  
* If you want to use value for building SQL only or embedding value to SQL directly, you must use `get` or `out` func. This func check that value contains prohibited character(s) for avoiding SQL injection.`out` is annotative, but `get` is not annotative.  
  Prohibited characters are:
 	* Single quotation
	* Semi colon
	* Line comment (--)
	* Block comment (/* or */)
* If database driver that you use supports `sql.NamedArg`, you should call `ExecNamed` func.

```go
// sql is generated SQL from template.
// args are arguments for generated SQL.
query, args, err := sqlt.New(sqlt.Postgres).Exec(s, map[string]interface{}{
	"ids":      []int{1, 2, 3},
	"order":    "name DESC",
	"onlyMale": false,
	"name":     "Alex",
})
rows, err := db.Query(query, args...)
```

#### options

* `TimeFunc`: For using customized time in template.
* `Annotation`: Output meta data for debugging to rendered SQL.

### Generated SQL

#### call `Exec`

```sql
SELECT *
FROM users
WHERE id IN ($1, $2, $3)
AND name = $4
ORDER BY name DESC
```

#### call `ExecNamed`

Currently there are also many drivers who do not support `sql.NamedArg`.  
In future, driver support `sql.NamedArg`, you only need to change `Exec` to `ExecNamed`.

```sql
SELECT *
FROM users
WHERE id IN (:ids__1, :ids__2, :ids__3)
AND name = :name
ORDER BY name DESC
```

## Install

```bash
$ go get github.com/pinzolo/sqlt
```

## Suppor

### Go version

Go 1.9 or later

### Databases

* PostgreSQL
* MySQL
* Oracle
* SQL Server

## Contribution

1. Fork ([https://github.com/pinzolo/sqlt/fork](https://github.com/pinzolo/sqlt/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[pinzolo](https://github.com/pinzolo)
