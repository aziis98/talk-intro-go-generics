package genericmethods_test

type Cons[L, R any] struct {
	Left  L
	Right L
}

// interface cast type check
func _[T any]() LiftLower[BoxT, T] { return &Box[T]{} }

type BoxT struct{}
type Box[T any] struct {
	Content T
}

func (c Box[T]) Lift() V[BoxT, T] {
	return c
}
func (c *Box[T]) Lower(v V[BoxT, T]) {
	*c = v.(Box[T])
}

// interface cast type check
func _[T any]() LiftLower[ChestT, T] { return &Chest[T]{} }

type ChestT struct{}
type Chest[T any] struct {
	Treasure T
}

func (c Chest[T]) Lift() V[ChestT, T] {
	return c
}
func (c *Chest[T]) Lower(v V[ChestT, T]) {
	*c = v.(Chest[T])
}

type LiftLower[F ~struct{}, T any] interface {
	Lift() V[F, T]
	Lower(v V[F, T])
}

func BoxToChest[T any](b Box[T]) Chest[T] {
	return Chest[T]{b.Content}
}

func BoxToChestLifted[T any](b V[BoxT, T]) V[ChestT, T] {
	var bb Box[T]
	bb.Lower(b)
	return Chest[T]{bb.Content}.Lift()
}

type V[F ~struct{}, T any] any

type ConsF[F, T, R any] struct{}

func _() {
	// l1 := Cons[Box[int], Box[string]]{}

}
