package utils

import "reflect"

//=>Chunk size array.
func ChunkSize[T any](slice T, chunkSize int) []T {
	sliceType := reflect.TypeOf(slice)
	sliceVal := reflect.ValueOf(slice)
	length := sliceVal.Len()
	if sliceType.Kind() != reflect.Slice {
		panic("parameter must be []T")
	}
	n := 0
	if length%chunkSize > 0 {
		n = 1
	}
	SST := reflect.MakeSlice(reflect.SliceOf(sliceType), 0, length/chunkSize+n)
	st, ed := 0, 0
	for st < length {
		ed = st + chunkSize
		if ed > length {
			ed = length
		}
		SST = reflect.Append(SST, sliceVal.Slice(st, ed))
		st = ed
	}
	return SST.Interface().([]T)
}
