package chper

import (
	"reflect"
	"testing"
)

func TestMapKey(t *testing.T) {
	got := MapKeys(map[string]int{"1": 1, "2": 2})
	want := []string{"1", "2"}
	SliceSort(got)
	SliceSort(want)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want: %v, got: %v", want, got)
	}
}

func TestMapValue(t *testing.T) {
	got := MapValues(map[string]int{"1": 1, "2": 2})
	want := []int{1, 2}
	SliceSort(got)
	SliceSort(want)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want: %v, got: %v", want, got)
	}
}

func TestMapConvertBool(t *testing.T) {
	got := MapConvertBool(map[string]int{"1": 1, "2": 2})
	want := map[string]bool{"1": true, "2": true}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want: %v, got: %v", want, got)
	}
}

func TestMapCompare(t *testing.T) {
	tests := []struct {
		name string

		a      map[string]int
		b      map[string]int
		isSame func(a, b int) bool

		wantJustAInclude   map[string]bool
		wantJustBInclude   map[string]bool
		wantBothButNotSame map[string]bool
	}{
		{
			name:   "nil",
			a:      nil,
			b:      nil,
			isSame: func(a int, b int) bool { return a == b },

			wantJustAInclude:   map[string]bool{},
			wantJustBInclude:   map[string]bool{},
			wantBothButNotSame: map[string]bool{},
		},
		{
			name:   "case1",
			a:      map[string]int{"a": 1},
			b:      map[string]int{"a": 1},
			isSame: func(a int, b int) bool { return a == b },

			wantJustAInclude:   map[string]bool{},
			wantJustBInclude:   map[string]bool{},
			wantBothButNotSame: map[string]bool{},
		},
		{
			name:   "case2",
			a:      map[string]int{"a": 1},
			b:      map[string]int{"a": 2},
			isSame: func(a int, b int) bool { return a == b },

			wantJustAInclude:   map[string]bool{},
			wantJustBInclude:   map[string]bool{},
			wantBothButNotSame: map[string]bool{"a": true},
		},
		{
			name:   "case3",
			a:      map[string]int{"a": 1},
			b:      map[string]int{"b": 2},
			isSame: func(a int, b int) bool { return a == b },

			wantJustAInclude:   map[string]bool{"a": true},
			wantJustBInclude:   map[string]bool{"b": true},
			wantBothButNotSame: map[string]bool{},
		},
		{
			name:   "case4",
			a:      map[string]int{"a": 1, "b": 2, "c": 3},
			b:      map[string]int{"b": 2, "c": 30, "d": 4},
			isSame: func(a int, b int) bool { return a == b },

			wantJustAInclude:   map[string]bool{"a": true},
			wantJustBInclude:   map[string]bool{"d": true},
			wantBothButNotSame: map[string]bool{"c": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotJustAInclude, gotJustBInclude, gotBothButNotSame := MapCompareF(tt.a, tt.b, tt.isSame)
			if !reflect.DeepEqual(gotJustAInclude, tt.wantJustAInclude) {
				t.Errorf("MapCompare() gotJustAInclude = %v, want %v", gotJustAInclude, tt.wantJustAInclude)
			}
			if !reflect.DeepEqual(gotJustBInclude, tt.wantJustBInclude) {
				t.Errorf("MapCompare() gotJustBInclude = %v, want %v", gotJustBInclude, tt.wantJustBInclude)
			}
			if !reflect.DeepEqual(gotBothButNotSame, tt.wantBothButNotSame) {
				t.Errorf("MapCompare() gotBothButNotSame = %v, want %v", gotBothButNotSame, tt.wantBothButNotSame)
			}
		})
	}
}

func TestMapMerge(t *testing.T) {
	{
		a, b := map[string]int{"1": 1, "2": 2}, map[string]int{"1": 11, "3": 3}
		got := MapMerge(a, b, false)
		want := map[string]int{"1": 11, "2": 2, "3": 3}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("want: %v, got: %v", want, got)
		}

		if reflect.DeepEqual(a, got) {
			t.Error("want not equal")
		}
	}

	{
		a, b := map[string]int{"1": 1, "2": 2}, map[string]int{"1": 11, "3": 3}
		got := MapMerge(a, b, true)
		want := map[string]int{"1": 11, "2": 2, "3": 3}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("want: %v, got: %v", want, got)
		}

		if !reflect.DeepEqual(a, got) {
			t.Error("want equal")
		}
	}
}

func TestMapFilter(t *testing.T) {
	a := map[string]int{"1": 1, "2": 2}
	got := MapFilter(a, map[string]bool{"1": true})
	want := map[string]int{"1": 1}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want: %v, got: %v", want, got)
	}
}
