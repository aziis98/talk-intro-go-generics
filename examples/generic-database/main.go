package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
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

//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//

type Transporter interface {
	Transport(peopleCount int) error
}

type Car struct {
	Owner string
}

func (c *Car) Transport(count int) error {
	if count > 5 {
		return fmt.Errorf(`troppe persone`)
	}

	return nil
}

type Bus struct{}

func (c *Bus) Transport(count int) error {
	if count > 40 {
		return fmt.Errorf(`troppe persone`)
	}

	return nil
}

type BrokenCar struct{}

func (c *BrokenCar) Transport(count int) error {
	return fmt.Errorf(`rotta!`)
}

func OurTrip(part1 Transporter, part2 Transporter) error {
	if err := part1.Transport(3); err != nil {
		return err
	}
	if err := part2.Transport(3); err != nil {
		return err
	}

	return nil
}
