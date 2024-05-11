package util

import "strconv"

func ToInteger[T uint64 | int64 | uint32 | int32 | int | uint16 | int16 | int8 | uint8](s string) T {
	x, _ := strconv.ParseUint(s, 10, 64)
	return T(x)
}
