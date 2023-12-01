package services

import "sync"

type QNode struct {
	value []byte
	next  *QNode
}

type Queue struct {
	head   *QNode
	tail   *QNode
	length int
	sync   sync.Mutex
}

func NewQueue() *Queue {
	return &Queue{
		length: 0,
		head:   nil,
		tail:   nil,
		sync:   sync.Mutex{},
	}
}

func (q *Queue) Enqueue(value []byte) {
	q.sync.Lock()
	defer q.sync.Unlock()
	node := &QNode{value: value}

	if q.length == 0 {
		q.head = node
		q.tail = node
	} else {
		q.tail.next = node
		q.tail = node
	}
	q.length++
}

func (q *Queue) Dequeue() (bool, []byte) {
	q.sync.Lock()
	defer q.sync.Unlock()
	if q.length == 0 {
		return false, nil
	}

	value := q.head.value
	if q.length > 1 {
		q.head = q.head.next
	} else {
		q.head = nil
		q.tail = nil
	}
	q.length--
	return true, value
}

func (q *Queue) DequeueMany(count int) (bool, *[][]byte) {

	if q.length == 0 {
		return false, nil
	}
	var data [][]byte
	isFound := false
	for i := 0; i < count; i++ {
		ok, value := q.Dequeue()
		if !ok {
			break
		}
		data = append(data, value)
		isFound = true

	}
	return isFound, &data
}
