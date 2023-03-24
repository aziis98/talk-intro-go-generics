package main

import (
	"database/sql"
	"fmt"
)

type IdTable interface {
	PrimaryKey() *string
	Columns() []any
}

type Ref[T IdTable] string

type Table[T IdTable] struct {
	Name     string
	PkColumn string
	// Columns  []string
}

type DB = *sql.DB

func Create[T IdTable](d DB, t Table[T], row T) (Ref[T], error) {
	var zero Ref[T]
	return zero, nil
}

func Insert[T IdTable](d DB, t Table[T], row T) (Ref[T], error) {
	var zero Ref[T]
	return zero, nil
}

func Read[T IdTable](d DB, t Table[T], ref Ref[T]) (*T, error) {
	result := d.QueryRow(fmt.Sprintf(`SELECT * FROM %s WHERE %s = ?`, t.Name, t.PkColumn), string(ref))

	var value T
	if err := result.Scan(value.Columns()...); err != nil {
		return nil, err
	}

	return &value, nil
}

func Update[T IdTable](d DB, t Table[T], row T) error {
	return nil
}

func Delete[T IdTable](d DB, t Table[T], id string) error {
	return nil
}
