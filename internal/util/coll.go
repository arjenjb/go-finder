package util

import "errors"

// Detect an element in slice `coll` using matching function `f` or return an error
func Detect[N any](coll []N, f func(n N) bool) (*N, error) {
	for _, each := range coll {
		if f(each) {
			return &each, nil
		}
	}
	return nil, errors.New("could not find item")
}

// MustDetect is the error ignoring version of Detect
func MustDetect[N any](children []N, f func(n N) bool) *N {
	return Must(Detect(children, f))
}

// Reduce a list to a single value, starting with an accumulator
func Reduce[T any, R any](coll []T, f func(acc R, e T) R, initial R) R {
	acc := initial
	for _, each := range coll {
		acc = f(acc, each)
	}
	return acc
}

// Fold reduces a list to a single value starting with the first two elements
func Fold[T any](coll []T, f func(a, b T) T) (T, error) {
	if len(coll) < 1 {
		var t T
		return t, errors.New("cannot fold an empty collection")
	} else if len(coll) == 1 {
		return coll[0], nil
	} else {
		return Reduce(coll[1:], f, coll[0]), nil
	}
}

// MustFold panics if Fold fails
func MustFold[T any](coll []T, f func(a, b T) T) T {
	return Must(Fold(coll, f))
}

// Map the values of a slice using function f. This function differs from the `lo` package in that the mapping function
// does not receive the index of the current item. This makes it easier to use with existing functions.
func Map[T any, R any](coll []T, f func(e T) R) []R {
	var result []R
	for _, e := range coll {
		result = append(result, f(e))
	}
	return result
}

// Map the values of a slice first by function f, then f2
func Map2[A any, B any, R any](coll []A, f func(e A) B, f2 func(e B) R) []R {
	var result []R
	for _, e := range coll {
		result = append(result, f2(f(e)))
	}
	return result
}

// MustMap maps the elements of slice `coll` by function `f` and wraps the output with `Must` to ignore error handling
func MustMap[T any, R any](coll []T, f func(T) (i R, err error)) []R {
	var result []R
	for _, e := range coll {
		v := Must(f(e))
		result = append(result, v)
	}
	return result
}

func Contains[T comparable](common []T, item T) bool {
	for _, e := range common {
		if e == item {
			return true
		}
	}
	return false
}

func Filter[T any](coll []T, f func(p T) bool) []T {
	var result []T
	for _, e := range coll {
		if f(e) {
			result = append(result, e)
		}
	}
	return result
}

// Chunked chunks slice `coll` into chunks of length `i`
func Chunked[T any](coll []T, i int) [][]T {
	var chunk []T
	var result [][]T

	for _, e := range coll {
		chunk = append(chunk, e)
		if len(chunk) == i {
			result = append(result, chunk)
			chunk = nil
		}
	}

	if len(chunk) > 0 {
		result = append(result, chunk)
	}

	return result
}

// AnySatisfy returns true if at least one element of `coll` satisfies function `f`
func AnySatisfy[T any](coll []T, f func(n T) bool) bool {
	for _, each := range coll {
		if f(each) {
			return true
		}
	}
	return false
}

// AllSatisfy returns true if each element of `coll` satisfies function `f`
func AllSatisfy[T any](coll []T, f func(n T) bool) bool {
	for _, each := range coll {
		if !f(each) {
			return false
		}
	}
	return true
}

// WithoutLast returns a copy of slice `coll` without the last N element
func WithoutLast[T any](coll []T, n int) []T {
	return coll[:len(coll)-n]
}

// Last returns the last element of slice `coll`
func Last[T any](coll []T) T {
	return coll[len(coll)-1]
}

// CopyUpUntil copies the elements from slice coll until it encounters an element that evaluates to true with func `f`
func CopyUpUntil[T any](coll []T, f func(T) bool) []T {
	var result []T
	for _, e := range coll {
		result = append(result, e)
		if f(e) {
			break
		}
	}
	return result
}

// FindIndex returns the first index for which function `f` returns true, if not found returns -1
func FindIndex[T any](coll []T, f func(T) bool) int {
	for idx, each := range coll {
		if f(each) {
			return idx
		}
	}

	return -1
}
