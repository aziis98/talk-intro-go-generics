package main

import (
	"fmt"
	"log"
)

//
// Database
//

type QueryResult interface {
	Scan(v any) error
}

type Database interface {
	Get(query string, values ...any) QueryResult
}

//
// Library
//

// DatabaseTable represents a typed database table with
type DatabaseTable[T any] struct {
	Table    string
	IdKey    string
	GetIdPtr func(entry T) *string
}

func (t DatabaseTable[T]) RefForId(id string) DatabaseRef[T] {
	return DatabaseRef[T]{id}
}

func (t DatabaseTable[T]) RefForValue(v T) DatabaseRef[T] {
	return t.RefForId(*t.GetIdPtr(v))
}

type DatabaseRef[T any] struct{ Id string }

func DatabaseRead[T any](db Database, table DatabaseTable[T], ref DatabaseRef[T]) (*T, error) {
	query := fmt.Sprintf(`select * from %s where %s = ?`, table.Table, table.IdKey)

	result := db.Get(query, ref.Id)

	var value T
	if err := result.Scan(&value); err != nil {
		return nil, err
	}

	return &value, nil
}

// DatabaseWrite creates a new entry in the db using the provided id in the value, this returns a typed ref to the db entry
func DatabaseWrite[T any](db Database, table DatabaseTable[T], value *T) (DatabaseRef[T], error) {
	panic("not implemented")
}

// DatabaseCreate creates a new entry in the db retrieving the new generated id and returning it as a typed ref
func DatabaseCreate[T any](db Database, table DatabaseTable[T], value *T) (DatabaseRef[T], error) {
	panic("not implemented")
}

//
// Client code
//

type User struct {
	Username  string
	FirstName string
	LastName  string
}

var UsersTable = DatabaseTable[User]{
	Table: "users",
	IdKey: "username",
	GetIdPtr: func(u User) *string {
		return &u.Username
	},
}

func _(db Database) error {
	u1 := &User{"j.smith", "John", "Smith"}
	log.Println(u1)

	ref1, err := DatabaseWrite(db, UsersTable, u1)
	if err != nil {
		return err
	}

	u2, err := DatabaseRead(db, UsersTable, ref1)
	if err != nil {
		return err
	}

	log.Println(u2)

	return nil
}
