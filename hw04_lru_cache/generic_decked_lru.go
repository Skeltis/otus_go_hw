package hw04lrucache

import (
	"sync"
)

type GenericCache[T any] interface {
	Set(key Key, value T) bool
	Get(key Key) (T, bool)
	Clear()
}

type genericLruCache[T any] struct {
	capacity    int
	structBlock sync.Mutex
	queue       Deck[genericCacheItem[T]]
	items       map[Key]*DeckItem[genericCacheItem[T]]
}

type genericCacheItem[T any] struct {
	key          Key
	wrappedValue T
}

func NewGenericCache[T any](capacity int) (GenericCache[T], error) {
	if capacity < 1 {
		return nil, ErrCapacityMustBeGreaterThanZero
	}

	return &genericLruCache[T]{
		capacity: capacity,
		queue:    NewDeck[genericCacheItem[T]](),
		items:    make(map[Key]*DeckItem[genericCacheItem[T]], capacity),
	}, nil
}

func (input *genericLruCache[T]) Set(key Key, value T) bool {
	input.structBlock.Lock()
	defer input.structBlock.Unlock()
	element, exists := input.items[key]
	if exists {
		input.queue.MoveToFront(element)
		element.Value.wrappedValue = value
		return true
	}
	newCacheElement := genericCacheItem[T]{key: key, wrappedValue: value}
	newListElement := input.queue.PushFront(newCacheElement)
	if input.queue.Len() > input.capacity {
		removedCacheElement := input.queue.Back()
		input.queue.Remove(removedCacheElement)
		delete(input.items, removedCacheElement.Value.key)
	}
	input.items[key] = newListElement
	return false
}

func (input *genericLruCache[T]) Get(key Key) (T, bool) {
	var result T
	input.structBlock.Lock()
	defer input.structBlock.Unlock()
	listElement, exists := input.items[key]
	if exists {
		input.queue.MoveToFront(listElement)
		result = listElement.Value.wrappedValue
	}

	return result, exists
}

func (input *genericLruCache[T]) Clear() {
	input.structBlock.Lock()
	input.queue = NewDeck[genericCacheItem[T]]()
	input.items = make(map[Key]*DeckItem[genericCacheItem[T]], input.capacity)
	input.structBlock.Unlock()
}
