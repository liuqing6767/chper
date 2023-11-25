chper means `container helper`, a golang library, make use of Go's Generics feature to simplify programming, inspired by [PHP Array functions](https:-www.php.net/manual/en/ref.array.php)

## Slice Function
- Range: create an slice containing a range of element
- SliceWalk: applices call to each element of the slice
- SliceMap: apply the callback to the lementes of the given list
- SliceReduce: applies iteratively the callback function to the elements of the slice, so as to reduce the slice to a single value
- SliceCombine: creates an map by using keys slice for keys and vaules slice for its values
- SliceCountValues: counts the occurrences of each distinct value in slice
- SliceAsKey: return a map, key is slice element, vaule is true
- SliceDiff: computes the difference of slices
- SliceIntersect: computes the intersection of slices
- SliceFill: fill a slice with vaulue
- SliceFilter: filters elements of a slice using keepMap
- SliceFilterF: filters elements of a slice using callback function
- SliceFlip: exchanges all keys with thieir associated values in slice
- SliceMerge: mergee one or more slices
- SlicePad: returns a copy of the slice padded to size specified by length with value value.
- SliceReverse: return a slice with elements in reverse order
- Slice2Map: convert slice to map
- SliceSortF: Sort an array by values using a user-defined comparison function
- SliceExistF: Checks if a value exists in an array, return 0, false or index, true
- SliceExist: same with SliceExistF
- SliceShuffle: shuffle slice


## Map Function
- MapKeys: return map's key
- MapValues: return vaules slice
- MapConvertBool:  convert map, replace all values to true
- MapConvertF: convert map[K]V to map[K]V1, convertor specified K -> K1 logic
- MapCompareF: compare two map
- MapCompare: compare two map, isSame is reflect.DeepEqual
- MapMerge: merge two map
- MapFileter: filter map, just key which occur in keep map's items
- MapFilterF: is same as MapFilter, keepFunc specify whether keep
- MapShallowCopy: return a copy, value reuse origin vaule

## Ring
Ring is a sorted set with fixed capcity
- Push: append one element
- First: get first element
- Last: get last element
- Element: get all element
- Size: get Ring size

## Consistent Hashing
- NewCHash[Node any](nodes []Node, options ...chashOptionFunc[Node]) (*CHash[Node], error): new consistent hash
- (ch *CHash[Node]) Hash(data []byte) (node Node, err error): get node by data
- (ch *CHash[Node]) HashIDer(ider IDer) (node Node, err error): get node by IDer
- (ch *CHash[Node]) AddNode(node Node) error: add one node
- (ch *CHash[Node]) RemoveNode(node Node) error: remove one node