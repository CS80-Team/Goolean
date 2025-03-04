package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStack(t *testing.T) {
	t.Run("Test NewStack with default size", func(t *testing.T) {
		stack := NewStack[int]()
		assert.Equal(t, 10, stack.GetCapacity(), "Expected capacity 10")
		assert.True(t, stack.IsEmpty(), "Expected stack to be empty")
	})

	t.Run("Test NewStack with custom size", func(t *testing.T) {
		stack := NewStack[int](5)
		assert.Equal(t, 5, stack.GetCapacity(), "Expected capacity 5")
	})

	t.Run("Test NewStack with predefined slice", func(t *testing.T) {
		stack := NewStack[int]([]int{1, 2, 3})
		assert.Equal(t, 3, stack.GetSize(), "Expected size 3")
		assert.Equal(t, 3, stack.Peek(), "Expected top element 3")
	})

	t.Run("Test NewStack with Predefined slice takes a copy from the slice", func(t *testing.T) {
		data := []int{1, 2, 3}
		stack := NewStack[int](data)
		data[2] = 100
		assert.Equal(t, 3, stack.Peek(), "Expected top element 3")
	})

	t.Run("Test NewStack with invalid parameter", func(t *testing.T) {
		require.Panics(t, func() { NewStack[int]("invalid") }, "Expected panic with invalid parameter")
	})

	t.Run("Test Push and Pop", func(t *testing.T) {
		stack := NewStack[int]()
		stack.Push(10)
		stack.Push(20)
		assert.Equal(t, 20, stack.Pop(), "Expected popped element 20")
		assert.Equal(t, 10, stack.Pop(), "Expected popped element 10")
	})

	t.Run("Test Peek", func(t *testing.T) {
		stack := NewStack[int]()
		stack.Push(100)
		assert.Equal(t, 100, stack.Peek(), "Expected top element 100")
	})

	t.Run("Test IsEmpty", func(t *testing.T) {
		stack := NewStack[int]()
		assert.True(t, stack.IsEmpty(), "Expected stack to be empty")
		stack.Push(5)
		assert.False(t, stack.IsEmpty(), "Expected stack to be non-empty")
	})

	t.Run("Test GetSize", func(t *testing.T) {
		stack := NewStack[int]()
		stack.Push(1)
		stack.Push(2)
		assert.Equal(t, 2, stack.GetSize(), "Expected size 2")
	})

	t.Run("Test Clear", func(t *testing.T) {
		stack := NewStack[int]()
		stack.Push(1)
		stack.Push(2)
		stack.Clear()
		assert.True(t, stack.IsEmpty(), "Expected stack to be empty after clear")
	})

	t.Run("Test Pop on empty stack", func(t *testing.T) {
		stack := NewStack[int]()
		require.Panics(t, func() { stack.Pop() }, "Expected panic when popping from empty stack")
	})

	t.Run("Test Peek on empty stack", func(t *testing.T) {
		stack := NewStack[int]()
		require.Panics(t, func() { stack.Peek() }, "Expected panic when peeking empty stack")
	})
}
