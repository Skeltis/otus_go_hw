package hw04lrucache

import (
	"errors"
	"sync"
)

var ErrCapacityMustBeGreaterThanZero = errors.New("cache capacity can't be less than 1")

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity    int
	structBlock sync.Mutex
	queue       List
	items       map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) (Cache, error) {
	if capacity < 1 {
		return nil, ErrCapacityMustBeGreaterThanZero
	}
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}, nil
}

func (input *lruCache) Set(key Key, value interface{}) bool {
	input.structBlock.Lock()
	defer input.structBlock.Unlock()
	element, exists := input.items[key]
	if exists {
		// type assertion won't return a pointer to the same internal struct
		input.queue.Remove(element)
	}
	newCacheElement := cacheItem{key: key, value: value}
	newListElement := input.queue.PushFront(newCacheElement)
	if !exists && input.queue.Len() > input.capacity {
		removedCacheElement := input.queue.Back()
		removedElement, ok := removedCacheElement.Value.(cacheItem)
		if !ok {
			panic(errors.New("invalid internal queue type"))
		}
		input.queue.Remove(removedCacheElement)
		delete(input.items, removedElement.key)
	}
	input.items[key] = newListElement
	return exists
}

func (input *lruCache) Get(key Key) (interface{}, bool) {
	input.structBlock.Lock()
	defer input.structBlock.Unlock()
	listElement, exists := input.items[key]
	if exists {
		input.queue.MoveToFront(listElement)
		foundElement, ok := listElement.Value.(cacheItem)
		if !ok {
			panic(errors.New("invalid internal queue type"))
		}
		return foundElement.value, exists
	}

	return nil, false
}

func (input *lruCache) Clear() {
	input.structBlock.Lock()
	input.queue = NewList()
	input.items = make(map[Key]*ListItem, input.capacity)
	input.structBlock.Unlock()
}
