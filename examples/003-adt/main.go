package main

// type Option[T] = Some(T) | None

type Option[T any] interface {
	Match(
		caseSome func(value T),
		caseNone func(),
	)
}

type Some[T any] struct{ Value T }

func (v Some[T]) Match(caseSome func(value T), caseNone func()) {
	caseSome(v.Value)
}

type None[T any] struct{}

func (v None[T]) Match(caseSome func(value T), caseNone func()) {
	caseNone()
}

// type Either[A, B] = Left A | Right B

type Either[A, B any] interface {
	Match(
		caseLeft func(value A),
		caseRight func(value B),
	)
}

type Left[A, B any] struct{ Value A }

func (v Left[A, B]) Match(
	caseLeft func(value A),
	caseRight func(value B),
) {
	caseLeft(v.Value)
}

type Right[A, B any] struct{ Value B }

func (v Right[A, B]) Match(
	caseLeft func(value A),
	caseRight func(value B),
) {
	caseRight(v.Value)
}
