package tile

// IfThenElse 根据条件返回对应的泛型结果，注意避免空指针问题
func IfThenElse[T any](condition bool, trueValue, falseValue T) T {
	if condition {
		return trueValue
	}
	return falseValue
}

// IfThenElseFunc 根据条件执行对应的函数并返回泛型结果
func IfThenElseFunc[T any](condition bool, trueFun, falseFunc func() (T, error)) (T, error) {
	if condition {
		return trueFun()
	}
	return falseFunc()
}
