package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})

	t.Run("checking prev\\next links with elements adding and removal", func(t *testing.T) {
		testList := NewList()

		testList.PushFront(10) // [10]
		require.Equal(t, 1, testList.Len())
		require.Equal(t, 10, testList.Front().Value)
		require.Equal(t, 10, testList.Back().Value)
		require.Nil(t, testList.Front().Prev)
		require.Nil(t, testList.Front().Next)
		require.Nil(t, testList.Back().Prev)
		require.Nil(t, testList.Back().Next)

		testList.PushBack(20) // [10, 20]

		require.Equal(t, 2, testList.Len())
		require.Equal(t, 10, testList.Front().Value)
		require.Equal(t, 20, testList.Front().Next.Value)
		require.Nil(t, testList.Front().Prev)
		require.Nil(t, testList.Front().Next.Next)
		require.Equal(t, 20, testList.Back().Value)
		require.Equal(t, 10, testList.Back().Prev.Value)
		require.Nil(t, testList.Back().Next)
		require.Nil(t, testList.Back().Prev.Prev)

		testList.PushFront(30) // [30, 10, 20]

		require.Equal(t, 3, testList.Len())

		require.Equal(t, 30, testList.Front().Value)
		require.Equal(t, 10, testList.Front().Next.Value)
		require.Equal(t, 20, testList.Front().Next.Next.Value)
		require.Nil(t, testList.Front().Prev)
		require.Nil(t, testList.Front().Next.Next.Next)

		require.Equal(t, 20, testList.Back().Value)
		require.Equal(t, 10, testList.Back().Prev.Value)
		require.Equal(t, 30, testList.Back().Prev.Prev.Value)
		require.Nil(t, testList.Back().Next)
		require.Nil(t, testList.Back().Prev.Prev.Prev)

		testList.Remove(testList.Front()) // [10, 20]
		require.Equal(t, 2, testList.Len())
		require.Equal(t, 10, testList.Front().Value)
		require.Equal(t, 20, testList.Back().Value)
		require.Nil(t, testList.Front().Prev)
		require.Equal(t, 20, testList.Front().Next.Value)
		require.Equal(t, 10, testList.Back().Prev.Value)
		require.Nil(t, testList.Back().Next)

		testList.Remove(testList.Back()) // [10]
		require.Equal(t, 1, testList.Len())
		require.Equal(t, 10, testList.Front().Value)
		require.Equal(t, 10, testList.Back().Value)
		require.Nil(t, testList.Front().Prev)
		require.Nil(t, testList.Front().Next)
		require.Nil(t, testList.Back().Prev)
		require.Nil(t, testList.Back().Next)

		testList.Remove(testList.Front()) // []
		require.Equal(t, 0, testList.Len())
	})
}
