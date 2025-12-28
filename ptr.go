package tile

// ToPtr 转化为指针类型
func ToPtr[T any](t T) *T {
	return &t
}
