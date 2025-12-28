package list

import (
	"errors"
	"github.com/NuyoahCh/tile"
	"github.com/NuyoahCh/tile/internal/errs"
	"math/rand"
)

const (
	FactorP  = float32(0.25) // level i 上的结点 有FactorP的比例出现在level i + 1上
	MaxLevel = 32            // 跳表的最高层数
)

// skipListNode 跳表节点
type skipListNode[T any] struct {
	Val     T
	Forward []*skipListNode[T]
}

// SkipList 跳表数据结构定义
type SkipList[T any] struct {
	header  *skipListNode[T]
	level   int
	compare tile.Comparator[T]
	size    int
}

// newSkipListNode 初始化跳表节点
func newSkipListNode[T any](Val T, level int) *skipListNode[T] {
	return &skipListNode[T]{Val, make([]*skipListNode[T], level)}
}

// AsSlice 转化为 slice 切片
func (sl *SkipList[T]) AsSlice() []T {
	curr := sl.header
	slice := make([]T, 0, sl.size)
	for curr.Forward[0] != nil {
		slice = append(slice, curr.Forward[0].Val)
		curr = curr.Forward[0]
	}
	return slice
}

// NewSkipList 初始化跳表结构
func NewSkipList[T any](compare tile.Comparator[T]) *SkipList[T] {
	return &SkipList[T]{
		header: &skipListNode[T]{
			Forward: make([]*skipListNode[T], MaxLevel),
		},
		level:   1,
		compare: compare,
	}
}

// NewSkipListFromSlice 从切片中初始化跳表结构
func NewSkipListFromSlice[T any](slice []T, compare tile.Comparator[T]) *SkipList[T] {
	sl := NewSkipList[T](compare)
	for _, n := range slice {
		sl.Insert(n)
	}
	return sl
}

// randomLevel 表明 levels 的生成和跳表中的元素个数无关
func (sl *SkipList[T]) randomLevel() int {
	level := 1
	p := FactorP
	// 随机增加层数
	for (rand.Int31() & 0xFFFF) < int32(p*0xFFFF) {
		level++
	}
	if level < MaxLevel {
		return level
	}
	return MaxLevel
}

// Search 查找跳表元素
func (sl *SkipList[T]) Search(target T) bool {
	curr, _ := sl.traverse(target, sl.level)
	curr = curr.Forward[0] // 第一层包含所有元素
	return curr != nil && sl.compare(curr.Val, target) == 0
}

// traverse 遍历跳表元素
func (sl *SkipList[T]) traverse(Val T, level int) (*skipListNode[T], []*skipListNode[T]) {
	update := make([]*skipListNode[T], MaxLevel) // update[i] 包含位于level i 的插入/删除位置左侧的指针
	curr := sl.header
	for i := level - 1; i >= 0; i-- {
		for curr.Forward[i] != nil && sl.compare(curr.Forward[i].Val, Val) < 0 {
			curr = curr.Forward[i]
		}
		update[i] = curr
	}
	return curr, update
}

// Insert 插入元素
func (sl *SkipList[T]) Insert(Val T) {
	// 定位：找到每一层的前驱节点，存入 update
	_, update := sl.traverse(Val, sl.level)
	// 定高：随机生成新节点的层数
	level := sl.randomLevel()

	// 扩容：若新层数高于当前最高层，将高出的部分指向 Header
	if level > sl.level {
		for i := sl.level; i < level; i++ {
			update[i] = sl.header
		}
		sl.level = level // 更新跳表最大层数
	}

	// 链接：创建节点并逐层插入链表
	newNode := newSkipListNode[T](Val, level)
	for i := 0; i < level; i++ {
		newNode.Forward[i] = update[i].Forward[i] // 新节点接管后续
		update[i].Forward[i] = newNode            // 前驱指向新节点
	}

	// 计数：更新元素数量
	sl.size += 1
}

// Len 返回长度
func (sl *SkipList[T]) Len() int {
	return sl.size
}

// DeleteElement 删除元素
func (sl *SkipList[T]) DeleteElement(target T) bool {
	curr, update := sl.traverse(target, sl.level)
	node := curr.Forward[0]
	if node == nil || sl.compare(node.Val, target) != 0 {
		return true
	}
	// 删除 target 节点
	for i := 0; i < sl.level && update[i].Forward[i] == node; i++ {
		update[i].Forward[i] = node.Forward[i]
	}

	// 更新层级
	for sl.level > 1 && sl.header.Forward[sl.level-1] == nil {
		sl.level--
	}
	sl.size -= 1
	return true
}

// Peek 查看跳表中元素
func (sl *SkipList[T]) Peek() (T, error) {
	curr := sl.header
	curr = curr.Forward[0]
	var zero T
	if curr == nil {
		return zero, errors.New("跳表为空")
	}
	return curr.Val, nil
}

// Get 获取跳表中的值
func (sl *SkipList[T]) Get(index int) (T, error) {
	var zero T
	if index < 0 || index >= sl.size {
		return zero, errs.NewErrIndexOutOfRange(sl.size, index)
	}
	curr := sl.header
	for i := 0; i <= index; i++ {
		curr = curr.Forward[0]
	}
	return curr.Val, nil
}
