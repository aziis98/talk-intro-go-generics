package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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

func DecodeAndValidateJSON_Interface(r *http.Request, target Validator) error {
	err := json.NewDecoder(r.Body).Decode(target)
	if err != nil {
		return err
	}

	if err := target.Validate(); err != nil {
		return err
	}

	return nil
}

func main() {
	m := http.NewServeMux()
	m.HandleFunc("/with-generic", func(w http.ResponseWriter, r *http.Request) {
		foo, err := DecodeAndValidateJSON_Generic[FooRequest](r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(foo); err != nil {
			log.Fatal(err)
			return
		}
	})
	m.HandleFunc("/with-interface", func(w http.ResponseWriter, r *http.Request) {
		var foo FooRequest
		if err := DecodeAndValidateJSON_Interface(r, &foo); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(foo); err != nil {
			log.Fatal(err)
			return
		}
	})

	log.Printf(`Starting server on port :4000...`)
	log.Fatal(http.ListenAndServe(":4000", m))
}
