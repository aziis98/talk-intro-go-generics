package main

import (
	"context"
	"fmt"
	"time"
)

func Aggregate[T any](cs ...<-chan T) <-chan T {
	agg := make(chan T)

	go func() {
		defer close(agg)
		for {
			closed := 0
			for _, c := range cs {
				select {
				case value, more := <-c:
					if more {
						agg <- value
					} else {
						closed++
					}
				default:
				}
			}

			if closed == len(cs) {
				break
			}
		}
	}()

	return agg
}

type TargetChan[T any] struct {
	c      <-chan T
	target *T
}

type TryReceiver interface {
	TryReceive() bool
}

func (tc TargetChan[T]) TryReceive() bool {
	select {
	case value := <-tc.c:
		*tc.target = value
		return true
	default:
		return false
	}
}

func AwaitFirst(cancel context.CancelFunc, ws ...TryReceiver) {
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			for _, w := range ws {
				if w.TryReceive() {
					done <- struct{}{}
					return
				}
			}
		}
	}()

	<-done
	cancel()
	return
}

func main() {
	c1 := make(chan string)
	c2 := make(chan int)
	c3 := make(chan float64)

	c := context.Background()
	c, cancel := context.WithCancel(c)

	go func() {
		defer close(c1)
		for {
			select {
			case <-c.Done():
				return
			case <-time.After(300 * time.Millisecond):
			}

			c1 <- "1"
		}
	}()
	go func() {
		defer close(c2)
		for {
			select {
			case <-c.Done():
				return
			case <-time.After(400 * time.Millisecond):
			}

			c2 <- 2
		}
	}()
	go func() {
		defer close(c3)
		for {
			select {
			case <-c.Done():
				return
			case <-time.After(200 * time.Millisecond):
			}

			c3 <- 3.0
		}
	}()

	var result1 string
	var result2 int
	var result3 float64

	AwaitFirst(
		cancel,
		TargetChan[string]{c1, &result1},
		TargetChan[int]{c2, &result2},
		TargetChan[float64]{c3, &result3},
	)

	fmt.Println(result1, result2, result3)
}
