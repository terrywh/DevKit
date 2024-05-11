package util

func FirstIndexOf[T comparable](s []T, a T) int {
	for i, x := range s {
		if x == a {
			return i
		}
	}
	return -1
}

func Contains[T comparable](s []T, a T) bool {
	return FirstIndexOf(s, a) > -1
}
