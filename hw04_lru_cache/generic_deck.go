package hw04lrucache

import "errors"

type Deck[T any] interface {
	Len() int
	Front() *DeckItem[T]
	Back() *DeckItem[T]
	PushFront(value T) *DeckItem[T]
	PushBack(value T) *DeckItem[T]
	Remove(i *DeckItem[T])
	MoveToFront(i *DeckItem[T])
}

type DeckItem[T any] struct {
	Value T
	Next  *DeckItem[T]
	Prev  *DeckItem[T]
}

type deck[T any] struct {
	head   *DeckItem[T]
	tail   *DeckItem[T]
	length int
}

func NewDeck[T any]() Deck[T] {
	return new(deck[T])
}

func (input *deck[T]) Len() int {
	return input.length
}

func (input *deck[T]) Front() *DeckItem[T] {
	return input.head
}

func (input *deck[T]) Back() *DeckItem[T] {
	return input.tail
}

func (input *deck[T]) push(newItem *DeckItem[T], pposition pushPosition) {
	if input.head == nil {
		input.head = newItem
		input.tail = newItem
		input.length++
		return
	}

	if pposition == toHead {
		input.head.Prev = newItem
		newItem.Next, input.head = input.head, newItem
	} else {
		input.tail.Next = newItem
		newItem.Prev, input.tail = input.tail, newItem
	}

	input.length++
}

func (input *deck[T]) PushFront(value T) *DeckItem[T] {
	newItem := DeckItem[T]{Value: value}
	input.push(&newItem, toHead)
	return &newItem
}

func (input *deck[T]) PushBack(value T) *DeckItem[T] {
	newItem := DeckItem[T]{Value: value}
	input.push(&newItem, toTail)
	return &newItem
}

// We suppose that Remove and MoveToFront only called for items
// that exist in list, it's called on, but we can do at least those
// consistency checks, cause they don't add asymptotic difficulty.
func (input *deck[T]) shallowCheckIfItemNotInList(i *DeckItem[T]) {
	if input.length == 0 {
		panic(errors.New("can't remove from empty collection"))
	}

	if i.Prev == nil && i != input.head {
		panic(errors.New("item supposed to be a head, but it's not"))
	}

	if i.Next == nil && i != input.tail {
		panic(errors.New("item supposed to be a tail, but it's not'"))
	}
}

func (input *deck[T]) removeInner(i *DeckItem[T]) {
	if input.length == 1 {
		input.head = nil
		input.tail = nil
		input.length--
		return
	}

	if input.head == i {
		input.head = input.head.Next
		input.head.Prev = nil
		i.Next = nil
		input.length--
		return
	}

	if input.tail == i {
		input.tail = input.tail.Prev
		input.tail.Next = nil
		i.Prev = nil
		input.length--
		return
	}

	i.Prev.Next = i.Next
	i.Next.Prev = i.Prev
	i.Prev = nil
	i.Next = nil
	input.length--
}

func (input *deck[T]) Remove(i *DeckItem[T]) {
	input.shallowCheckIfItemNotInList(i)
	input.removeInner(i)
}

func (input *deck[T]) MoveToFront(i *DeckItem[T]) {
	if input.head == i {
		return
	}

	input.shallowCheckIfItemNotInList(i)
	input.removeInner(i)
	input.push(i, toHead)
}
