package main

type Expr[T any] interface{}

func Match() {

}

type numExpr struct {
	value int
}

func Num(value int) Expr[int] {
	return &numExpr{value}
}

type sumExpr struct {
	lhs, rhs Expr[int]
}

func Sum(lhs, rhs Expr[int]) Expr[int] {
	return &sumExpr{lhs, rhs}
}

type leqExpr struct {
	lhs, rhs Expr[int]
}

func Leq(lhs, rhs Expr[int]) Expr[bool] {
	return &leqExpr{lhs, rhs}
}

type ifExpr[R any] struct {
	cond    Expr[bool]
	ifTrue  Expr[R]
	ifFalse Expr[R]
}

func If[R any](cond Expr[bool], ifTrue Expr[R], ifFalse Expr[R]) Expr[R] {
	return &ifExpr[R]{cond, ifTrue, ifFalse}
}

func main() {

}
