package main

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

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
