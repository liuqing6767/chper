package chper

import (
	"math"
	"math/rand"
	"sort"
)

type Numeric interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

// Range create an slice containing a range of element
// step is optional, default is 1, cant be 0
// return an slice of elements from start to end
func Range[V Numeric](start, stop V, stepOpt ...V) []V {
	var step V = 1
	if len(stepOpt) == 1 {
		step = stepOpt[0]
	}

	if step == 0 {
		panic("step must not be zero")
	}

	if start == stop || (start < stop && step < 0) || (start > stop && step > 0) {
		return nil
	}

	size := int(math.Ceil(float64(stop-start) / float64(step)))
	if V(size)*step == stop-start {
		size++
	}
	res := make([]V, size)
	for i := 0; i < size; i++ {
		res[i] = start + step*V(i)
	}

	return res
}

// SliceWalk applices call to each element of the slice
func SliceWalk[V any](list []V, action func(int, V)) {
	for index, item := range list {
		action(index, item) //每个元素执行的操作
	}
}

// https://www.php.net/manual/en/function.arsort.php

// SliceMap apply the callback to the lementes of the given list
// return a new list with the results of the callback
func SliceMap[V, V1 any](list []V, callback func(i int, element V) V1) []V1 {
	result := make([]V1, len(list))
	for i, element := range list {
		result[i] = callback(i, element)
	}

	return result
}

// SliceReduce applies iteratively the callback function to the elements of the slice,
// so as to reduce the slice to a single value.
func SliceReduce[V any](list []V, callback func(carry, item V) V, initValueOpt ...V) V {
	var rst V
	if len(initValueOpt) > 0 {
		rst = initValueOpt[0]
	}

	for _, one := range list {
		rst = callback(rst, one)
	}

	return rst
}

// SliceCombine creates an map by using keys slice for keys and vaules slice for its values
// if len(keys) != len(values), result just include the shorten keys and values
func SliceCombine[K comparable, V any](keys []K, values []V) map[K]V {
	size := len(keys)
	if s := len(values); s < size {
		size = s
	}

	result := make(map[K]V, size)

	for i := 0; i < size; i++ {
		result[keys[i]] = values[i]
	}

	return result
}

// SliceCountValues counts the occurrences of each distinct value in slice
// return a map using the values of slice as keys  and their frequency in slice as vaules
func SliceCountValues[V comparable](list []V) map[V]int {
	result := make(map[V]int)
	for _, v := range list {
		result[v]++
	}

	return result
}

// SliceAsKey return a map, key is slice element, vaule is true
func SliceAsKey[V comparable](ss []V) map[V]bool {
	return Slice2Map(ss, func(v V) (key V, val bool, skip bool) {
		return v, true, false
	})
}

// SliceDiff computes the difference of slices
// return an slice containing all the entries from slice as
// that are not present in any of the other slices
func SliceDiff[V comparable](as []V, bs ...[]V) []V {
	bm := map[V]struct{}{}
	SliceWalk(as, func(i int, v V) {
		bm[v] = struct{}{}
	})

	for _, s := range bs {
		for _, v := range s {
			delete(bm, v)

			if len(bm) == 0 {
				return nil
			}
		}
	}

	return MapKeys(bm)
}

// SliceIntersect computes the intersection of slices
func SliceIntersect[V comparable](sa []V, sb ...[]V) []V {
	val2frequency := map[V]int8{}
	SliceWalk(sa, func(i int, v V) {
		val2frequency[v] = 1
	})

	for _, s := range sb {
		for _, v := range s {
			val2frequency[v]++
		}
	}

	result := make([]V, 0, len(sa))
	for k, frequency := range val2frequency {
		if frequency > 1 {
			result = append(result, k)
		}
	}

	return result
}

// SliceFill fill a slice with vaulue
// size must be greater than 0
func SliceFill[V any](start, size int, v V) []V {
	if size < 0 {
		panic("SliceFill: size < 0")
	}
	result := make([]V, start+size)
	for i := start; i < start+size; i++ {
		result[i] = v
	}

	return result
}

// SliceFilter filters elements of a slice using keepMap
func SliceFilter[V comparable](list []V, keepMap map[V]bool) []V {
	if keepMap == nil {
		return nil
	}

	return SliceFilterF(list, func(i int, v V) (keep bool) {
		return keepMap[v]
	})
}

// SliceFilterF filters elements of a slice using callback function
func SliceFilterF[V any](list []V, keepFunc func(i int, V V) (keep bool)) []V {
	res := make([]V, 0, len(list))
	for i, s := range list {
		if keepFunc(i, s) {
			res = append(res, s)
		}
	}

	return res
}

// SliceFlip exchanges all keys with thieir associated values in slice
// if a value has several occurrences, the lastest index will be used as its vaule
func SliceFlip[V comparable](list []V) map[V]int {
	result := make(map[V]int, len(list))
	SliceWalk(list, func(i int, v V) {
		result[v] = i
	})

	return result
}

// SliceMerge mergee one or more slices
func SliceMerge[V any](ss ...[]V) []V {
	size := 0
	for _, s := range ss {
		size += len(s)

	}
	result := make([]V, 0, size)
	for _, s := range ss {
		result = append(result, s...)
	}

	return result
}

// SlicePad returns a copy of the slice padded to size specified by length with value value.
// If length is positive then the slice is padded on the right, if it's negative then on the left.
// If the absolute value of length is less than or equal to the length of the slice then no padding takes place.
func SlicePad[V any](list []V, wantLength int, valOpt ...V) []V {
	var v V
	if len(valOpt) == 1 {
		v = valOpt[0]
	}

	padLeftSize, padRightSize := 0, 0
	if wantLength < 0 {
		if -wantLength > len(list) {
			padLeftSize = -wantLength - len(list)
		}
	} else {
		if wantLength > len(list) {
			padRightSize = wantLength - len(list)
		}
	}

	data := make([]V, 0, len(list)+padLeftSize+padRightSize)
	for i := 0; i < padLeftSize; i++ {
		data = append(data, v)
	}
	data = append(data, list...)
	for i := 0; i < padRightSize; i++ {
		data = append(data, v)
	}

	return data
}

// SliceReverse return a slice with elements in reverse order
func SliceReverse[V any](list []V) []V {
	if list == nil {
		return nil
	}

	rst := make([]V, 0, len(list))
	for i := len(list) - 1; i >= 0; i-- {
		rst = append(rst, list[i])
	}

	return rst
}

func SliceUnique[V comparable](list []V) []V {
	if list == nil {
		return nil
	}

	tmp := map[V]struct{}{}
	rst := []V{}

	for _, one := range list {
		if _, ok := tmp[one]; !ok {
			tmp[one] = struct{}{}
			rst = append(rst, one)
		}
	}

	return rst
}

// Slice2Map convert slice to map
// keyFactory return key and  whether skip
func Slice2Map[V, V1 any, K comparable](ss []V, keyFactory func(V) (key K, val V1, skip bool)) map[K]V1 {
	rst := make(map[K]V1)
	SliceWalk(ss, func(i int, s V) {
		key, val, skip := keyFactory(s)
		if !skip {
			rst[key] = val
		}
	})

	return rst
}

// SlliceConvert convert element by convertor
func SliceConvert[V, V1 any](ss []V, convertor func(V) (data V1, skip bool, err error)) (converted []V1, err error) {
	converted = make([]V1, 0, len(ss))
	for _, s := range ss {
		data, skip, err := convertor(s)
		if err != nil {
			return nil, err
		}
		if skip {
			continue
		}

		converted = append(converted, data)
	}

	return converted, nil
}

type bigAble interface {
	Numeric | string | rune | byte
}

func SliceSort[V bigAble](ss []V) {
	SliceSortF(ss, func(v1, v2 V) bool { return v1 < v2 })
}

// SliceSortF Sort an array by values using a user-defined comparison function
func SliceSortF[V any](ss []V, cmp func(v1, v2 V) bool) {
	if cmp == nil {
		return
	}

	sb := &sortBy[V]{
		vs:  ss,
		cmp: cmp,
	}

	sort.Sort(sb)
}

type sortBy[V any] struct {
	vs  []V
	cmp func(v1, v2 V) bool
}

func (sb *sortBy[V]) Len() int           { return len(sb.vs) }
func (sb *sortBy[V]) Swap(i, j int)      { sb.vs[i], sb.vs[j] = sb.vs[j], sb.vs[i] }
func (sb *sortBy[V]) Less(i, j int) bool { return sb.cmp(sb.vs[i], sb.vs[j]) }

// SliceExistF Checks if a value exists in an array, return 0, false or index, true
func SliceExistF[V any](ss []V, compareFunc func(v V) (match bool)) (int, bool) {
	for i, v := range ss {
		if compareFunc(v) {
			return i, true
		}
	}

	return 0, false
}

// SliceExist same with SliceExistF
func SliceExist[V comparable](ss []V, data V) (int, bool) {
	return SliceExistF(ss, func(v V) (match bool) { return v == data })
}

// SliceShuffle shuffle slice
func SliceShuffle[V comparable](ss []V) {
	rand.Shuffle(len(ss), func(i, j int) { ss[i], ss[j] = ss[j], ss[i] })
}
