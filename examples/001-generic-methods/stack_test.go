package genericmethods_test

type Stack[T any] []T

func (s *Stack[T]) Push(value T) {
	*s = append(*s, value)
}

func (s Stack[T]) Peek() T {
	return s[len(s)-1]
}

func (s Stack[T]) Len() int {
	return len(s)
}

func (s *Stack[T]) Pop() (T, bool) {
	items := *s

	if len(items) == 0 {
		var zero T
		return zero, false
	}

	newStack, poppedValue := items[:len(items)-1], items[len(items)-1]
	*s = newStack

	return poppedValue, true
}
