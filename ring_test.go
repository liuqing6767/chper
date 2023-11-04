package chper

import (
	"reflect"
	"testing"
)

func TestRing(t *testing.T) {
	{
		ring := NewRing[int](1)

		got := ring.Elements(nil)
		var want []int
		if !reflect.DeepEqual(got, want) {
			t.Errorf("want: %v, got: %v", want, got)
		}

		ring.Push(1)
		got = ring.Elements(nil)
		want = []int{1}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("want: %v, got: %v", want, got)
		}

		ring.Push(2)
		got = ring.Elements(nil)
		want = []int{2}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("want: %v, got: %v", want, got)
		}

		got = ring.Elements(func(i int) bool { return false })
		want = nil
		if !reflect.DeepEqual(got, want) {
			t.Errorf("want: %v, got: %v", want, got)
		}
	}

	{
		ring := NewRing[int](10)

		_, ok := ring.First()
		if ok {
			want := false
			got := ok
			t.Errorf("want: %v, got: %v", want, got)
		}

		_, ok = ring.Last()
		if ok {
			want := false
			got := ok
			t.Errorf("want: %v, got: %v", want, got)
		}

		got := ring.Elements(nil)
		var want []int
		if !reflect.DeepEqual(got, want) {
			t.Errorf("want: %v, got: %v", want, got)
		}

		for i := 0; i < 5; i++ {
			ring.Push(i)
		}

		first, ok := ring.First()
		if !ok {
			want := true
			got := ok
			t.Errorf("want: %v, got: %v", want, got)
		}
		if first != 0 {
			want := 0
			got := first
			t.Errorf("want: %v, got: %v", want, got)
		}

		last, ok := ring.Last()
		if !ok {
			want := true
			got := ok
			t.Errorf("want: %v, got: %v", want, got)
		}
		if last != 4 {
			want := 4
			got := last
			t.Errorf("want: %v, got: %v", want, got)
		}

		got = ring.Elements(nil)
		want = []int{0, 1, 2, 3, 4}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("want: %v, got: %v", want, got)
		}

		got = ring.Elements(func(i int) bool { return i%2 == 0 })
		want = []int{0, 2, 4}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("want: %v, got: %v", want, got)
		}

		for i := 5; i < 15; i++ {
			ring.Push(i)
		}

		first, ok = ring.First()
		if !ok {
			want := true
			got := ok
			t.Errorf("want: %v, got: %v", want, got)
		}
		if first != 5 {
			want := 5
			got := first
			t.Errorf("want: %v, got: %v", want, got)
		}

		last, ok = ring.Last()
		if !ok {
			want := true
			got := ok
			t.Errorf("want: %v, got: %v", want, got)
		}
		if last != 14 {
			want := 14
			got := last
			t.Errorf("want: %v, got: %v", want, got)
		}

		got = ring.Elements(nil)
		want = []int{5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("want: %v, got: %v", want, got)
		}
	}
}
