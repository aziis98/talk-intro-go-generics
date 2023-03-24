package genericmethods_test

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

type Validator interface {
	Validate() error
}

type FooRequest struct {
	A int    `json:"a"`
	B string `json:"b"`
}

func (foo FooRequest) Validate() error {
	if foo.A < 0 {
		return fmt.Errorf(`parameter "a" cannot be lesser than zero`)
	}
	if !strings.HasPrefix(foo.B, "baz-") {
		return fmt.Errorf(`parameter "b" has wrong prefix`)
	}

	return nil
}

func DecodeAndValidateJSON_Generic[T Validator](r *http.Request) (T, error) {
	var value T
	err := json.NewDecoder(r.Body).Decode(&value)
	if err != nil {
		var zero T
		return zero, err
	}

	if err := value.Validate(); err != nil {
		var zero T
		return zero, err
	}

	return value, nil
}

func TestDecodeAndValidateJSON_Generic(t *testing.T) {
	m := http.NewServeMux()
	m.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		foo, err := DecodeAndValidateJSON_Generic[FooRequest](r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(foo)
	})
}

func DecodeAndValidateJSON_Interface(r *http.Request, target *Validator) error {
	err := json.NewDecoder(r.Body).Decode(target)
	if err != nil {
		return err
	}

	if err := (*target).Validate(); err != nil {
		return err
	}

	return nil
}

func TestDecodeAndValidateJSON_Interface(t *testing.T) {
	m := http.NewServeMux()
	m.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		var foo Validator = FooRequest{}
		if err := DecodeAndValidateJSON_Interface(r, &foo); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(foo)
	})
}

type WithPrimaryKey interface {
	PrimaryKeyPtr() *string
}

type Ref[T WithPrimaryKey] string

type Table[T WithPrimaryKey] struct {
	Name     string
	PkColumn string
}

func IdToRef[T WithPrimaryKey](table Table[T], id string) (Ref[T], error) {
	if !strings.HasPrefix(id, table.Name+":") {
		var zero Ref[T]
		return zero, fmt.Errorf(`invalid reference %v for table %v`, id, table)
	}

	return Ref[T](id), nil
}

func Read[T WithPrimaryKey](db *sql.DB, table Table[T], ref Ref[T]) (*T, error) {
	result := db.QueryRow(
		fmt.Sprintf(
			`SELECT * FROM %s WHERE %s = ?`,
			table.Name, table.PkColumn,
		),
		string(ref),
	)

	var value T
	if err := result.Scan(&value); err != nil {
		return nil, err
	}

	return &value, nil
}

type User struct {
	Username  string
	FirstName string
	LastName  string
}

var _ WithPrimaryKey = &User{}

func (u *User) PrimaryKeyPtr() *string {
	return &u.Username
}
