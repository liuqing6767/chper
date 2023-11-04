package chper

import (
	"reflect"
)

// MapKeys return map's key
func MapKeys[T comparable, V any](m map[T]V) []T {
	out := make([]T, 0, len(m))
	for t := range m {
		out = append(out, t)
	}

	return out
}

// MapValues return vaules slice
func MapValues[K comparable, V any](m map[K]V) []V {
	res := make([]V, 0, len(m))
	for _, v := range m {
		res = append(res, v)
	}

	return res
}

func boolConvertor[K comparable, V any, V1 bool](k K, v V) (data V1, skip bool, err error) {
	return true, false, nil
}

// MapConvertBool  convert map, replace all values to true
func MapConvertBool[K comparable, V any](m map[K]V) map[K]bool {
	data, _ := MapConvertF(m, boolConvertor[K, V, bool])
	return data
}

// MapConvertF convert map[K]V to map[K]V1, convertor specified K -> K1 logic
func MapConvertF[K comparable, V, V1 any](m map[K]V,
	convertor func(k K, v V) (data V1, skip bool, err error)) (map[K]V1, error) {

	out := make(map[K]V1, len(m))
	for k, v := range m {
		data, skip, err := convertor(k, v)
		if err != nil {
			return nil, err
		}
		if skip {
			continue
		}

		out[k] = data

	}

	return out, nil
}

// MapCompareF compare two map
// justAInculde include elements which just occur in a
// justBInculde include elements which just occur in b
// bothButNotEqual include elements which occur in a and b, but value is diffrent
func MapCompareF[K comparable, T any](a, b map[K]T, isSame func(a, b T) bool) (
	justAInclude, justBInclude, bothButNotEqual map[K]bool) {

	justAInclude, justBInclude, bothButNotEqual = map[K]bool{}, map[K]bool{}, map[K]bool{}

	if a == nil {
		a = map[K]T{}
	}
	if b == nil {
		b = map[K]T{}
	}

	for k, av := range a {
		bv, ok := b[k]
		if !ok {
			justAInclude[k] = true
		} else if !isSame(av, bv) {
			bothButNotEqual[k] = true
		}
	}

	for k := range b {
		_, ok := a[k]
		if !ok {
			justBInclude[k] = true
		}
	}

	return
}

func isSame[T any](a, b T) bool {
	return reflect.DeepEqual(a, b)
}

// MapCompare compare two map, isSame is reflect.DeepEqual
func MapCompare[K comparable, T any](a, b map[K]T) (
	justAInclude, justBInclude, bothButNotEqual map[K]bool) {

	return MapCompareF(a, b, isSame[T])
}

// MapMerge merge two map
// if inPlace is true, baseline will use a, if is false, baseline is a new map
// if the same key has diffrent value in a and b, b's value will be used
func MapMerge[K comparable, V any](a, b map[K]V, inPlace bool) map[K]V {
	dest := a
	if !inPlace {
		dest = make(map[K]V, len(a)+len(b))
		for k, v := range a {
			dest[k] = v
		}
	}

	if dest == nil {
		dest = make(map[K]V, len(b))
	}

	for k, v := range b {
		dest[k] = v
	}

	return dest
}

// MapFileter filter map, just key which occur in keep map's items
func MapFilter[K comparable, V any](origin map[K]V, keepMap map[K]bool) map[K]V {
	return MapFilterF(origin, func(k K, v V) bool {
		return keepMap[k]
	})
}

// MapFilterF is same as MapFilter, keepFunc specify whether keep
func MapFilterF[K comparable, V any](origin map[K]V, keepFunc func(k K, v V) (keep bool)) map[K]V {
	return MapShallowCopy(origin, keepFunc)
}

// MapShallowCopy return a copy, value reuse origin vaule
func MapShallowCopy[K comparable, V any](origin map[K]V,
	keepFuncOpt ...func(k K, v V) (keep bool)) map[K]V {

	if origin == nil {
		return nil
	}

	var keepFunc func(k K, v V) (keep bool)
	if len(keepFuncOpt) == 1 {
		keepFunc = keepFuncOpt[0]
	}

	m := make(map[K]V, len(origin))
	for k, v := range origin {
		if keepFunc != nil && keepFunc(k, v) {
			m[k] = v
		}
	}

	return m
}
