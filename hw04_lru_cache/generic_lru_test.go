package hw04lrucache

import (
	"github.com/stretchr/testify/require"
	"math/rand"
	"strconv"
	"sync"
	"testing"
)

func TestCacheGeneric(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c, _ := NewGenericCache[int](10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c, _ := NewGenericCache[int](5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Zero(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c, _ := NewGenericCache[int](3)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		wasInCache = c.Set("ccc", 300)
		require.False(t, wasInCache)

		wasInCache = c.Set("ddd", 400)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.False(t, ok)
		require.Zero(t, val)

		wasInCache = c.Set("bbb", 500)
		require.True(t, wasInCache)

		wasInCache = c.Set("eee", 600)
		require.False(t, wasInCache)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Zero(t, val)
	})
}

func TestCacheMultithreadingGeneric(t *testing.T) {
	c, _ := NewGenericCache[int](10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}

func TestZeroCapacityGenericCacheReturnsError(t *testing.T) {
	c, err := NewGenericCache[int](0)
	require.Nil(t, c)
	require.ErrorIs(t, err, ErrCapacityMustBeGreaterThanZero)
}
