package main

import (
	"fmt"
	"math/rand"
	"time"
)

func trySend[T any](c chan<- T, v T) bool {
	select {
	case c <- v:
		return true
	default:
		return false
	}
}

func raceSame[T any](cs ...<-chan T) T {
	done := make(chan T)
	defer close(done)

	for _, c := range cs {
		go func(c <-chan T) {
			trySend(done, <-c)
		}(c)
	}

	return <-done
}

type Awaiter interface {
	Await()
}

type awaiterChan[T any] <-chan T

func (rc awaiterChan[T]) Await() { <-rc }

type targetChan[T any] struct {
	c      <-chan T
	target *T
}

func (rc targetChan[T]) Await() { *rc.target = <-rc.c }

func raceAny(rs ...Awaiter) {
	done := make(chan struct{})
	defer close(done)

	for _, r := range rs {
		go func(r Awaiter) {
			r.Await()
			trySend(done, struct{}{})
		}(r)
	}

	<-done
}

func main() {
	c1 := make(chan string)
	c2 := make(chan int)
	c3 := make(chan float64)

	go func() {
		defer close(c1)
		time.Sleep(300 * time.Millisecond)
		c1 <- "1"
	}()
	go func() {
		defer close(c2)
		time.Sleep(400 * time.Millisecond)
		c2 <- 2
	}()
	go func() {
		defer close(c3)
		time.Sleep(200 * time.Millisecond)
		c3 <- 3.0
	}()

	var result2 int
	var result3 float64

	raceAny(
		awaiterChan[string](c1),
		targetChan[int]{c2, &result2},
		targetChan[float64]{c3, &result3},
	)

	fmt.Println(result2, result3)
}
