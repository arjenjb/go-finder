package util

import "golang.org/x/exp/constraints"

type Number interface {
	constraints.Integer | constraints.Float
}

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	} else {
		return b
	}
}
func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	} else {
		return b
	}
}

func Add[T Number](a, b T) T {
	return a + b
}

func Sum[T Number](items []T) T {
	if len(items) == 0 {
		return 0
	}
	return MustFold(items, Add[T])
}

func Sqrt[T Number](i T) T {
	return i * i
}

func Mul[T Number](a T, b T) T {
	return a * b
}

func Abs[T Number](i T) T {
	if i < 0 {
		return -i
	} else {
		return i
	}
}
