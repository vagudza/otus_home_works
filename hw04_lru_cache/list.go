package hw04lrucache

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
	List  // Remove me after realization.
	front *ListItem
	back  *ListItem
	len   int
}

func NewList() List {
	return new(list)
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.front
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	defer func() {
		l.len++
	}()

	newItem := &ListItem{
		Value: v,
	}

	// in case of empty list
	if l.front == nil {
		l.front = newItem
		l.back = newItem
		return l.front
	}

	newItem.Next = l.front
	l.front.Prev = newItem
	l.front = newItem

	return l.front
}

func (l *list) PushBack(v interface{}) *ListItem {
	defer func() {
		l.len++
	}()

	newItem := &ListItem{
		Value: v,
	}

	// in case of empty list
	if l.back == nil {
		l.front = newItem
		l.back = newItem
		return l.back
	}

	newItem.Prev = l.back
	l.back.Next = newItem
	l.back = newItem

	return l.back
}

func (l *list) Remove(i *ListItem) {
	defer func() {
		l.len--
	}()

	if i.Prev == nil {
		l.front = i.Next
	} else {
		i.Prev.Next = i.Next
	}

	if i.Next == nil {
		l.back = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}

	i.Prev = nil
	i.Next = nil
}
