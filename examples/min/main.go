package main

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

// before go generics

func MinInt8(x, y int8) int8 {
	if x < y {
		return x
	}

	return y
}

func MinInt16(x, y int16) int16 {
	if x < y {
		return x
	}

	return y
}

func MinInt32(x, y int32) int32 {
	if x < y {
		return x
	}

	return y
}

func MinInt64(x, y int64) int64 {
	if x < y {
		return x
	}

	return y
}

func MinFloat32(x, y float32) float32 {
	if x < y {
		return x
	}

	return y
}

func MinFloat64(x, y float64) float64 {
	if x < y {
		return x
	}

	return y
}

// with go generics

func Min[T constraints.Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}

type Liter int

func main() {
	{
		var a, b int = 1, 2
		fmt.Println(Min(a, b))
	}
	{
		var a, b float64 = 3.14, 2.71
		fmt.Println(Min(a, b))
	}
	{
		var a, b Liter = 1, 2
		fmt.Println(Min(a, b))
	}
}
