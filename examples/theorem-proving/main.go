package main

//
// Encoding the semantics
//

type Bool interface{ isBool() }

type Term interface{ isTerm() }

type Term2Term interface{ isTerm2Term() }

// trick to encode higher-kinded types
type V[H Term2Term, T Term] Term

//
// Naturals
//

type Zero Term
type Succ Term2Term

// Some aliases

type One = V[Succ, Zero]
type Two = V[Succ, V[Succ, Zero]]
type Three = V[Succ, V[Succ, V[Succ, Zero]]]

//
// Equality
//

type Eq[A, B any] Bool

// Eq_Refl ~> forall x : x = x
func Eq_Reflexive[T any]() Eq[T, T] { panic("axiom") }

// Eq_Symmetric ~> forall a, b: a = b => b = a
func Eq_Symmetric[A, B any](_ Eq[A, B]) Eq[B, A] { panic("axiom") }

// Eq_Transitive ~> forall a, b, c: a = b e b = c => a = c
func Eq_Transitive[A, B, C any](_ Eq[A, B], _ Eq[B, C]) Eq[A, C] { panic("axiom") }

// Function_Eq ~> forall f function, forall a, b terms: a = b  => f(a) = f(b)
func Function_Eq[F Term2Term, A, B Term](_ Eq[A, B]) Eq[V[F, A], V[F, B]] {
	panic("axiom")
}

//
// Plus Axioms
//

type Plus[L, R Term] Term

// "a + 0 = a"

// Plus_Zero ~> forall a, b: a + succ(b) = succ(a + b)
func Plus_Zero[N Term]() Eq[Plus[N, Zero], N] { panic("axiom") }

// "a + (b + 1) = (a + b) + 1"

// Plus_Sum ~> forall a, b: a + succ(b) = succ(a + b)
func Plus_Sum[N, M Term]() Eq[Plus[N, V[Succ, M]], V[Succ, Plus[N, M]]] { panic("axiom") }

//
//	1 + 1 = 2
//

// Theorem_OnePlusOneEqTwo
func Theorem_OnePlusOneEqTwo() Eq[Plus[One, One], Two] {
	var en1 Eq[
		Plus[One, Zero],
		One,
	] = Plus_Zero[One]()

	var en2 Eq[
		V[Succ, Plus[One, Zero]],
		Two,
	] = Function_Eq[Succ](en1)

	var en3 Eq[
		Plus[One, One],
		V[Succ, Plus[One, Zero]],
	] = Plus_Sum[One, Zero]()

	return Eq_Transitive(en3, en2)
}
