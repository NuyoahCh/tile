package queue

import (
	"errors"
	"github.com/NuyoahCh/tile"
	"github.com/NuyoahCh/tile/internal/slice"
)

var (
	ErrOutOfCapacity = errors.New("tile: 超出最大容量限制")
	ErrEmptyQueue    = errors.New("tile: 队列为空")
)

// PriorityQueue 基于小顶堆的优先队列
// 当capacity <= 0时，为无界队列，切片容量会动态扩缩容
// 当capacity > 0 时，为有界队列，初始化后就固定容量，不会扩缩容
type PriorityQueue[T any] struct {
	compare  tile.Comparator[T] // 用于比较前一个元素是否小于后一个元素
	capacity int                // 队列容量
	data     []T                // 队列中的元素，根节点从 0 开始
}

// Len 计算长度
func (p *PriorityQueue[T]) Len() int {
	return len(p.data)
}

// Cap 无界队列返回0，有界队列返回创建队列时设置的值
func (p *PriorityQueue[T]) Cap() int {
	return p.capacity
}

// IsBoundless 小于边界值
func (p *PriorityQueue[T]) IsBoundless() bool {
	return p.capacity <= 0
}

// isFull 容量已满
func (p *PriorityQueue[T]) isFull() bool {
	return p.capacity > 0 && len(p.data) == p.capacity
}

// isEmpty 长度为空
func (p *PriorityQueue[T]) isEmpty() bool {
	return len(p.data) < 1
}

// Peek 检查队列首个元素
func (p *PriorityQueue[T]) Peek() (T, error) {
	if p.isEmpty() {
		var t T
		return t, ErrEmptyQueue
	}
	return p.data[0], nil
}

func (p *PriorityQueue[T]) Enqueue(t T) error {
	if p.isFull() {
		return ErrOutOfCapacity
	}

	p.data = append(p.data, t)
	node := len(p.data) - 1
	parent := (node - 1) / 2
	for parent >= 0 && p.compare(p.data[node], p.data[parent]) < 0 {
		p.data[parent], p.data[node] = p.data[node], p.data[parent]
		node = parent
		parent = (parent - 1) >> 1
	}

	return nil
}

func (p *PriorityQueue[T]) Dequeue() (T, error) {
	if p.isEmpty() {
		var t T
		return t, ErrEmptyQueue
	}

	pop := p.data[0]
	p.data[0] = p.data[len(p.data)-1]
	p.data = p.data[:len(p.data)-1]
	p.shrinkIfNecessary()
	p.heapify(p.data, len(p.data), 0)
	return pop, nil
}

func (p *PriorityQueue[T]) shrinkIfNecessary() {
	if p.IsBoundless() {
		p.data = slice.Shrink[T](p.data)
	}
}

func (p *PriorityQueue[T]) heapify(data []T, n, i int) {
	minPos := i
	for {
		if left := i*2 + 1; left < n && p.compare(data[left], data[minPos]) < 0 {
			minPos = left
		}
		if right := i*2 + 2; right < n && p.compare(data[right], data[minPos]) < 0 {
			minPos = right
		}
		if minPos == i {
			break
		}
		data[i], data[minPos] = data[minPos], data[i]
		i = minPos
	}
}

// NewPriorityQueue 创建优先队列 capacity <= 0 时，为无界队列，否则有有界队列
func NewPriorityQueue[T any](capacity int, compare tile.Comparator[T]) *PriorityQueue[T] {
	sliceCap := capacity
	if capacity < 1 {
		capacity = 0
		sliceCap = 64
	}
	return &PriorityQueue[T]{
		capacity: capacity,
		data:     make([]T, 0, sliceCap),
		compare:  compare,
	}
}
