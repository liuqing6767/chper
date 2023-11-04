package chper

import "sync"

// Ring is a sorted set with fixed capcity
// all element sorted by Push index
// early elelements will be removed if Ring is filled
type Ring[V any] struct {
	capacity int
	begin    int
	size     int

	lock sync.Mutex

	elements []V
}

func NewRing[V any](capacity int) *Ring[V] {
	if capacity < 1 {
		panic("bad capacity")
	}

	return &Ring[V]{
		capacity: capacity,
		begin:    0,
		size:     0,

		elements: make([]V, capacity),
	}
}

// Push append one element to Ring
// early elelements will be removed if Ring is filled
/*
 one example: NewRing(3), [ring header]
	action 		elements
	Push(1)		[1]
	Push(2)		[1] 2
	Push(3)		[1] 2 	3
	Push(4)		4   [2] 3
	Push(5)		4   5 	[3]
*/
func (r *Ring[V]) Push(data V) {
	r.lock.Lock()

	if r.capacity == r.size { // 满了
		r.elements[r.begin] = data // 覆盖当前的值

		r.begin++
		r.begin %= r.capacity
	} else {
		r.elements[r.begin+r.size] = data
		r.size++
	}

	r.lock.Unlock()
}

// Elements return matched elements, sorted by push index
func (r *Ring[V]) Elements(filter func(V) bool) (elements []V) {
	r.lock.Lock()
	defer r.lock.Unlock()

	for i := r.begin; i < r.begin+r.size; i++ {
		ele := r.elements[i%r.capacity]
		if filter == nil || filter(ele) {
			elements = append(elements, ele)
		}
	}

	return
}

// First return first element and exist
func (r *Ring[V]) First() (V, bool) {
	r.lock.Lock()
	defer r.lock.Unlock()

	return r.elements[(r.begin)], r.size != 0
}

// Last return last element and exist
func (r *Ring[V]) Last() (v V, ok bool) {
	r.lock.Lock()
	defer r.lock.Unlock()

	if r.size == 0 {
		return
	}

	return r.elements[(r.begin+r.size-1)%r.capacity], true
}

// Size return ring's elements count
func (r *Ring[V]) Size() int {
	r.lock.Lock()
	defer r.lock.Unlock()

	return r.size
}
