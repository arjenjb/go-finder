package util

func Must[T any](val T, e error) T {
	if e != nil {
		panic(e)
	}
	return val
}
