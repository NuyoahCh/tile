package tile

// Comparator 用于比较两个对象的大小：src < dst，返回 -1，src = dst，返回 0，scr > dst，返回 1
type Comparator[T any] func(src T, dst T) int

// ComparatorRealNumber 比较值方法
func ComparatorRealNumber[T RealNumber](src T, dst T) int {
	if src < dst {
		return -1
	} else if src == dst {
		return 0
	} else {
		return 1
	}
}
