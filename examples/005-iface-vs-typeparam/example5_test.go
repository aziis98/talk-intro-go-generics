package example5_test

import (
	"bytes"
	"io"
	"testing"
)

func WriteSingleWithInterface(w io.Writer, data byte) {
	w.Write([]byte{data})
}

//go:noinline
func WriteSingleWithInterfaceNoInline(w io.Writer, data byte) {
	w.Write([]byte{data})
}

func WriteSingleWithGeneric[T io.Writer](w T, data byte) {
	w.Write([]byte{data})
}

//go:noinline
func WriteSingleWithGenericNoInline[T io.Writer](w T, data byte) {
	w.Write([]byte{data})
}

func Test1(t *testing.T) {
	t.Log(`Ok!`)
}

func BenchmarkInterface(b *testing.B) {
	d := &bytes.Buffer{}

	for i := 0; i < b.N; i++ {
		WriteSingleWithInterface(d, 42)
	}
}

func BenchmarkInterfaceNoInline(b *testing.B) {
	d := &bytes.Buffer{}

	for i := 0; i < b.N; i++ {
		WriteSingleWithInterfaceNoInline(d, 42)
	}
}

func BenchmarkGeneric(b *testing.B) {
	d := &bytes.Buffer{}

	for i := 0; i < b.N; i++ {
		WriteSingleWithGeneric(d, 42)
	}
}

func BenchmarkGenericNoInline(b *testing.B) {
	d := &bytes.Buffer{}

	for i := 0; i < b.N; i++ {
		WriteSingleWithGenericNoInline(d, 42)
	}
}
