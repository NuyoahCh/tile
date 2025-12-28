package slice

import "github.com/NuyoahCh/tile/internal/errs"

// Add 添加元素
func Add[T any](src []T, element T, index int) ([]T, error) {
	length := len(src)
	// 数组下标越界
	if index < 0 || index > length {
		return nil, errs.NewErrIndexOutOfRange(length, index)
	}

	// 先将 src 扩展为一个元素
	var zeroValue T
	src = append(src, zeroValue)
	// 先位移
	for i := len(src) - 1; i > index; i-- {
		if i-1 >= 0 {
			src[i] = src[i-1]
		}
	}
	// 插入元素
	src[index] = element
	return src, nil
}
