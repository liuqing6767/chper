chper 是 `cotainer helper` 的简称，一个golang类库，通过使用泛型特性简化编程，函数签名参考了 [PHP Array functions](https://www.php.net/manual/en/ref.array.php)

## Slice Function
- Range: 创建一个范围内的slice
- SliceWalk: 给每个元素应用函数
- SliceMap: 给每个元素应用函数，返回一个新的slice
- SliceReduce: 给每个元素应用回调，最后得到一个单个的值
- SliceCombine: 将keys和values一一对应得到一个map
- SliceCountValues: 返回Slice中每个元素出现的次数
- SliceDiff: 计数只在第一个参数中促销案的元素
- SliceIntersect: 计算slice的交集
- SliceFill: 使用元素填充得到一个slice
- SliceFilter: 过滤slice
- SliceFilterF: 过滤slice
- SliceFlip: 将元素的值作为key，返回一个value为最后出现的index的map
- SliceAsKey: 将元素的值作为key，返回一个value为true的map
- SliceMerge: 合并slice
- SlicePad: 填充slice
- SliceReverse: 反转slice
- Slice2Map: 将slice转换为map
- SliceSortF: 排序slice
- SliceExistF: 检查指定元素是否存在
- SliceExist: 检查指定元素是否存在
- SliceShuffle: 打乱slice

## Map Function
- MapKeys: 返回map中的key
- MapValues: 返回map中的值
- MapConvertBool: 转换map，将值都修改为true
- MapConvertF: 转换map
- MapCompareF: 比较map
- MapCompare: 比较map
- MapMerge: 合并map
- MapFileter: 过滤map
- MapFilterF: 过滤map
- MapShallowCopy: 浅拷贝map

## Ring
Ring 是一个有固定容量的排序集合
- Push: 最加一个元素
- First: 得到第一个元素
- Last: 得到最后一个元素
- Element: 得到所有的元素
- Size: 得到环的大小