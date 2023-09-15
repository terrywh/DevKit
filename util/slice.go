package util

func SlicePtrStringIdx(slice []*string, idx int, def string) string {
	if len(slice) < 1 {
		return def
	}
	return *slice[0]
}

type Integer interface {
	int64 | uint64 | int32 | uint32 | int16 | uint16 | int8 | uint8
}

func SliceRange[T any, I Integer](slice []T, s, e I) []T {
	start := int(s)
	end := int(e)
	if start < 0 {
		start = 0
	}
	if start >= len(slice) {
		start = len(slice)
	}
	if end < 0 {
		end = 0
	}
	if end >= len(slice) {
		end = len(slice)
	}
	return slice[start:end]
}
