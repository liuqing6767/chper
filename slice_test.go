package chper

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

func TestSliceConvert(t *testing.T) {
	type T struct {
		Name string
	}

	tests := []struct {
		name string

		ss        []T
		convertor func(T) (data string, skip bool, err error)

		wantConverted []string
		wantErr       bool
	}{
		{
			name: "case1",
			ss:   []T{{"1"}, {"2"}, {"3"}, {"4"}, {"5"}, {"6"}, {"7"}},
			convertor: func(t T) (data string, skip bool, err error) {
				return t.Name, false, nil
			},

			wantConverted: []string{"1", "2", "3", "4", "5", "6", "7"},
			wantErr:       false,
		},
		{
			name: "case2",
			ss:   []T{{"1"}, {"2"}, {"3"}, {"4"}, {"5"}, {"6"}, {"7"}},
			convertor: func(t T) (data string, skip bool, err error) {
				return t.Name, true, nil
			},

			wantConverted: []string{},
			wantErr:       false,
		},
		{
			name: "case3",
			ss:   []T{{"1"}, {"2"}, {"3"}, {"4"}, {"5"}, {"6"}, {"7"}},
			convertor: func(t T) (data string, skip bool, err error) {
				return t.Name, true, fmt.Errorf("")
			},

			wantConverted: []string{},
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotConverted, err := SliceConvert(tt.ss, tt.convertor)
			if (err != nil) != tt.wantErr {
				t.Errorf("SliceConvert() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if !reflect.DeepEqual(gotConverted, tt.wantConverted) {
				t.Errorf("SliceConvert() = %v, want %v", gotConverted, tt.wantConverted)
			}
		})
	}
}

func TestForeach(t *testing.T) {
	total := 0

	SliceWalk([]int{1, 2, 3}, func(i, j int) {
		total += j
	})

	if total != 6 {
		t.Errorf("want 6, got %v", total)
	}
}

func TestSliceCombine(t *testing.T) {
	type args struct {
		keys   []int
		values []int
	}
	tests := []struct {
		name    string
		args    args
		want    map[int]int
		wantErr bool
	}{
		{
			name: "case0",
			args: args{
				keys:   []int{5, 6, 7, 7},
				values: []int{5, 6, 7, 8},
			},
			want: map[int]int{5: 5, 6: 6, 7: 8},
		},
		{
			name: "case1",
			args: args{
				keys:   []int{5, 6},
				values: []int{5, 6, 7, 8},
			},
			want: map[int]int{5: 5, 6: 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SliceCombine(tt.args.keys, tt.args.values)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SliceCombine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceDiff(t *testing.T) {
	type args struct {
		as []int
		bs [][]int
	}
	tests := []struct {
		name    string
		args    args
		want    []int
		wantErr bool
	}{
		{
			name: "case0",
			args: args{
				as: []int{5, 6, 7, 7},
				bs: [][]int{{5, 6, 7, 8}},
			},
			want: []int{},
		},
		{
			name: "case1",
			args: args{
				as: []int{5, 6, 7, 8},
				bs: [][]int{{5, 6}},
			},
			want: []int{7, 8},
		},
		{
			name: "case2",
			args: args{
				as: nil,
				bs: [][]int{{5, 6}},
			},
			want: []int{},
		},
		{
			name: "case3",
			args: args{
				as: []int{5, 6},
				bs: [][]int{{5, 6}},
			},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SliceDiff(tt.args.as, tt.args.bs...)
			gots, wants := fmt.Sprint(got), fmt.Sprint(tt.want)
			if !reflect.DeepEqual(gots, wants) {
				t.Errorf("SliceDiff() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestSliceIntersect(t *testing.T) {
	type args struct {
		as []int
		bs []int
	}
	tests := []struct {
		name    string
		args    args
		want    []int
		wantErr bool
	}{
		{
			name: "case0",
			args: args{
				as: []int{5, 6, 7, 7},
				bs: []int{5, 6, 7, 8},
			},
			want: []int{5, 6, 7},
		},
		{
			name: "case1",
			args: args{
				as: []int{5, 6},
				bs: []int{5, 6, 7, 8},
			},
			want: []int{5, 6},
		},
		{
			name: "case2",
			args: args{
				as: []int{5, 6, 7, 8},
				bs: []int{5, 6},
			},
			want: []int{5, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SliceIntersect(tt.args.as, tt.args.bs)
			SliceSort(got)
			SliceSort(tt.want)
			gots, wants := fmt.Sprint(got), fmt.Sprint(tt.want)
			if !reflect.DeepEqual(gots, wants) {
				t.Errorf("SliceIntersect() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestSliceFill(t *testing.T) {
	{
		got := SliceFill(0, 2, 1)
		want := []int{1, 1}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %v, want: %v", got, want)
		}
	}

	{
		got := SliceFill(1, 2, 1)
		want := []int{0, 1, 1}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %v, want: %v", got, want)
		}
	}

	{
		got := SliceFill(1, 0, 1)
		want := []int{0}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %v, want: %v", got, want)
		}
	}
}

func TestSliceFilter(t *testing.T) {
	{
		got := SliceFilter(nil, map[int]bool{4: true, 5: true})
		want := []int{}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %v, want: %v", got, want)
		}
	}
	{
		got := SliceFilter([]int{0, 4, 1}, map[int]bool{4: true, 5: true})
		want := []int{4}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %v, want: %v", got, want)
		}
	}
	{
		got := SliceFilter([]int{0, 4, 1}, nil)
		want := []int{}
		gots, wants := fmt.Sprint(got), fmt.Sprint(want)
		if !reflect.DeepEqual(gots, wants) {
			t.Errorf("got: %v, want: %v", got, want)
		}
	}
}

func TestSliceFlip(t *testing.T) {
	got := SliceFlip([]int{1, 2, 3, 3})
	want := map[int]int{1: 0, 2: 1, 3: 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func TestSliceMerge(t *testing.T) {
	got := SliceMerge([]int{1, 2, 3, 3}, []int{2, 2}, nil, []int{3, 3})
	want := []int{1, 2, 3, 3, 2, 2, 3, 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func TestSliceMap(t *testing.T) {
	got := SliceMap([]int{1, 2, 3}, func(i int, element int) string { return strconv.Itoa(element) })
	want := []string{"1", "2", "3"}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func TestSliceReverse(t *testing.T) {
	{
		got := SliceReverse([]int{1, 2, 3})
		want := []int{3, 2, 1}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %v, want: %v", got, want)
		}
	}
	{

		var data []int
		got := SliceReverse(data)
		var want []int

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %v, want: %v", got, want)
		}
	}
}

func TestSliceUnique(t *testing.T) {
	{
		got := SliceUnique([]int{1, 2, 3, 2, 3})
		want := []int{1, 2, 3}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %v, want: %v", got, want)
		}
	}
	{

		var data []int
		got := SliceUnique(data)
		var want []int

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got: %v, want: %v", got, want)
		}
	}
}

func TestSliceAsKey(t *testing.T) {
	got := SliceAsKey([]int{1, 2, 3, 2, 3})
	want := map[int]bool{1: true, 2: true, 3: true}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func TestSliceCountValues(t *testing.T) {
	type args struct {
		ss []int
	}
	tests := []struct {
		name    string
		args    args
		want    map[int]int
		wantErr bool
	}{
		{
			name: "case0",
			args: args{
				ss: []int{5, 6, 7, 7},
			},
			want: map[int]int{5: 1, 6: 1, 7: 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SliceCountValues(tt.args.ss)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SliceCountValues() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRangeInt(t *testing.T) {
	tests := []struct {
		name              string
		start, stop, step int
		want              []int
	}{
		{
			name:  "int0",
			start: 0,
			stop:  2,
			step:  3,
			want:  []int{0},
		},
		{
			name:  "int1",
			start: 0,
			stop:  3,
			step:  3,
			want:  []int{0, 3},
		},
		{
			name:  "int2",
			start: 0,
			stop:  4,
			step:  3,
			want:  []int{0, 3},
		},

		{
			name:  "int3",
			start: 0,
			stop:  2,
			step:  -3,
			want:  nil,
		},
		{
			name:  "int4",
			start: 10,
			stop:  2,
			step:  -3,
			want:  []int{10, 7, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Range(tt.start, tt.stop, tt.step); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Range() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRangeFloat(t *testing.T) {
	tests := []struct {
		name              string
		start, stop, step float64
		want              []float64
	}{
		{
			name:  "float640",
			start: 0,
			stop:  1.4,
			step:  0.5,
			want:  []float64{0, 0.5, 1.0},
		},
		{
			name:  "float641",
			start: 0,
			stop:  1.5,
			step:  0.5,
			want:  []float64{0, 0.5, 1.0, 1.5},
		},
		{
			name:  "float642",
			start: 0,
			stop:  1.6,
			step:  0.5,
			want:  []float64{0, 0.5, 1.0, 1.5},
		},

		{
			name:  "float643",
			start: 0,
			stop:  2,
			step:  -0.3,
			want:  nil,
		},
		{
			name:  "float644",
			start: 1.5,
			stop:  0,
			step:  -0.5,
			want:  []float64{1.5, 1.0, 0.5, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Range(tt.start, tt.stop, tt.step); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Range() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlicePad(t *testing.T) {
	type args struct {
		list       []int
		wantLength int
		valOpt     []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "case0",
			args: args{
				list:       []int{5, 6, 7},
				wantLength: 2,
				valOpt:     nil,
			},
			want: []int{5, 6, 7},
		},
		{
			name: "case1",
			args: args{
				list:       []int{5, 6, 7},
				wantLength: 3,
				valOpt:     nil,
			},
			want: []int{5, 6, 7},
		},
		{
			name: "case2",
			args: args{
				list:       []int{5, 6, 7},
				wantLength: 4,
				valOpt:     nil,
			},
			want: []int{5, 6, 7, 0},
		},
		{
			name: "case3",
			args: args{
				list:       []int{5, 6, 7},
				wantLength: -4,
				valOpt:     nil,
			},
			want: []int{0, 5, 6, 7},
		},
		{
			name: "case4",
			args: args{
				list:       []int{5, 6, 7},
				wantLength: 4,
				valOpt:     []int{1},
			},
			want: []int{5, 6, 7, 1},
		},
		{
			name: "case5",
			args: args{
				list:       []int{5, 6, 7},
				wantLength: -4,
				valOpt:     []int{1},
			},
			want: []int{1, 5, 6, 7},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SlicePad(tt.args.list, tt.args.wantLength, tt.args.valOpt...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SlicePad() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceReduce(t *testing.T) {
	type args struct {
		list         []int
		callback     func(carry, item int) int
		initValueOpt []int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "case1",
			args: args{
				list: []int{1, 2, 3},
				callback: func(carry, item int) int {
					return carry + item
				},
				initValueOpt: nil,
			},
			want: 6,
		},
		{
			name: "case2",
			args: args{
				list: []int{1, 2, 3},
				callback: func(carry, item int) int {
					return carry + item
				},
				initValueOpt: []int{5},
			},
			want: 11,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SliceReduce(tt.args.list, tt.args.callback, tt.args.initValueOpt...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SliceReduce() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceSort(t *testing.T) {
	type args struct {
		ss  []int
		cmp func(v1, v2 int) bool
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "case1",
			args: args{
				ss: []int{1, 2, 3},
				cmp: func(v1, v2 int) bool {
					return v1 < v2
				},
			},
			want: []int{1, 2, 3},
		},
		{
			name: "case2",
			args: args{
				ss: []int{1, 3, 2},
				cmp: func(v1, v2 int) bool {
					return v1 < v2
				},
			},
			want: []int{1, 2, 3},
		},
		{
			name: "case3",
			args: args{
				ss:  []int{1, 3, 2},
				cmp: nil,
			},
			want: []int{1, 3, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SliceSortF(tt.args.ss, tt.args.cmp)
			if !reflect.DeepEqual(tt.args.ss, tt.want) {
				t.Errorf("TestSliceSort() = %v, want %v", tt.args.ss, tt.want)
			}
		})
	}
}

func TestSliceExist(t *testing.T) {
	{
		index, exist := SliceExist([]int{1, 2, 3}, 2)
		wantIndex := 1
		wantExist := true
		if wantExist != exist || index != wantIndex {
			t.Errorf("want exist: %v, index: %v, got exist: %v, index: %v", wantIndex, wantExist, index, exist)
		}
	}

	{
		index, exist := SliceExist([]int{1, 2, 3}, 20)
		wantIndex := 0
		wantExist := false
		if wantExist != exist || index != wantIndex {
			t.Errorf("want exist: %v, index: %v, got exist: %v, index: %v", wantIndex, wantExist, index, exist)
		}
	}
}

func TestSliceShuffle(t *testing.T) {
	var ss []int
	SliceShuffle(ss)

	SliceShuffle([]int{1})
	// just no panic
}
