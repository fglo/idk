package virtualmachine

type stack[T any] []T

func newStack[T any]() stack[T] {
	return make([]T, 0)
}

func (s *stack[T]) push(val T) {
	*s = append(*s, val)
}

func (s *stack[T]) pop() T {
	stack := *s
	val := stack[len(stack)-1]
	*s = stack[:len(stack)-1]
	return val
}
