package util

func StrPtr(str string) (s *string) {
	s = new(string)
	*s = str
	return
}

func IntPtr(i int) (n *int) {
	n = new(int)
	*n = i
	return
}

func Int64Ptr(i int64) (n *int64) {
	n = new(int64)
	*n = i
	return
}

func Uint64Ptr(i uint64) (n *uint64) {
	n = new(uint64)
	*n = i
	return
}
