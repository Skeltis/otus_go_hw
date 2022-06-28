package hw04lrucache

import (
	"errors"
)

type pushPosition int

const (
	toHead pushPosition = iota
	toTail
)

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	head   *ListItem
	tail   *ListItem
	length int
}

func NewList() List {
	return new(list)
}

func (input *list) Len() int {
	return input.length
}

func (input *list) Front() *ListItem {
	return input.head
}

func (input *list) Back() *ListItem {
	return input.tail
}

func (input *list) push(newItem *ListItem, pposition pushPosition) {
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

func (input *list) PushFront(v interface{}) *ListItem {
	newItem := ListItem{Value: v}
	input.push(&newItem, toHead)
	return &newItem
}

func (input *list) PushBack(v interface{}) *ListItem {
	newItem := ListItem{Value: v}
	input.push(&newItem, toTail)
	return &newItem
}

// We suppose that Remove and MoveToFront only called for items
// that exist in list, it's called on, but we can do at least those
// consistency checks, cause they don't add asymptotic difficulty.
func (input *list) shallowCheckIfItemNotInList(i *ListItem) {
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

func (input *list) removeInner(i *ListItem) {
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

func (input *list) Remove(i *ListItem) {
	input.shallowCheckIfItemNotInList(i)
	input.removeInner(i)
}

func (input *list) MoveToFront(i *ListItem) {
	if input.head == i {
		return
	}

	input.shallowCheckIfItemNotInList(i)
	input.removeInner(i)
	input.push(i, toHead)
}
