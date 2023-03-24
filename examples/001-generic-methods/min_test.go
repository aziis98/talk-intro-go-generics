package genericmethods_test

import (
	"testing"

	"golang.org/x/exp/constraints"
)

func Min[T constraints.Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}

type Liter int

func TestInt(t *testing.T) {
	{
		var a, b int = 1, 2
		Min(a, b)
	}
	{
		var a, b float64 = 3.14, 2.71
		Min(a, b)
	}
	{
		var a, b Liter = 1, 2
		Min(a, b)
	}
}
