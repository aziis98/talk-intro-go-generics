package main

import (
	"log"

	"golang.org/x/exp/constraints"
)

//
// The old way
//

func MinInt(x, y int) int {
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

//
// Generics
//

func Min[T constraints.Ordered](x, y T) T {
	if x < y {
		return x
	}
	return y
}

func main() {
	shortMin := Min[int16] // func(int16, int16) int16
	log.Printf(`min(2, 3) = %v`, shortMin(2, 3))
}
