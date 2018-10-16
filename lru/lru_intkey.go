package lru

import (
	"sync"
	"time"
)

type IntLRUNode struct {
	data      interface{}
	dataDirty bool
	next      *IntLRUNode
	prev      *IntLRUNode
	time      int64
	key       int
}

type IntKeyLRU struct {
	head       *IntLRUNode
	tail       *IntLRUNode
	len        int
	strMapNode map[int]*IntLRUNode
	lock       sync.Mutex
}

func (lru *IntKeyLRU) getLRUNode(key int) *IntLRUNode {
	lru.lock.Lock()
	defer lru.lock.Unlock()
	n, ok := lru.strMapNode[key]
	if !ok {
		return nil
	}
	if n != nil {
		n.time = time.Now().Unix()
		if lru.head != lru.tail {
			if lru.head != n {
				n.prev.next = n.next
				if lru.tail == n {
					lru.tail = n.prev
					lru.tail.next = nil
				} else {
					n.next.prev = n.prev
				}
				n.next = lru.head
				lru.head.prev = n
				lru.head = n
				lru.head.prev = nil
			}
		}
	}
	return n
}

func (lru *IntKeyLRU) GetValue(key int) interface{} {
	n := lru.getLRUNode(key)
	if n != nil {
		return n.data
	}
	return nil
}

func (lru *IntKeyLRU) RmValue(key int) interface{} {
	lru.lock.Lock()
	defer lru.lock.Unlock()
	n, ok := lru.strMapNode[key]
	if !ok {
		return nil
	}

	if n != nil {
		if lru.head == lru.tail {
			lru.head = nil
			lru.tail = nil
		} else {
			if lru.head == n {
				lru.head = n.next
				lru.head.prev = nil
			} else {
				n.prev.next = n.next
				if lru.tail == n {
					lru.tail = n.prev
					lru.tail.next = nil
				} else {
					n.next.prev = n.prev
				}
			}
		}
		lru.len--
		delete(lru.strMapNode, key)
		return n.data
	}
	return nil
}

func (lru *IntKeyLRU) AddValue(str int, data interface{}) bool {
	n := lru.getLRUNode(str)
	lru.lock.Lock()
	defer lru.lock.Unlock()
	if n != nil {
		n.data = data
		return false
	} else {
		n = &IntLRUNode{
			data: data,
			key:  str,
			time: time.Now().Unix(),
		}
		if lru.head == nil {
			lru.head = n
			lru.tail = n
		} else {
			lru.head.prev = n
			n.next = lru.head
			lru.head = n
		}
		lru.len++
		lru.strMapNode[str] = n
		return true
	}
}

func (lru *IntKeyLRU) RangeReverse(iter *IntLRUNode) *IntLRUNode {
	if iter == nil {
		return lru.tail
	}
	return iter.prev
}

func (lru *IntKeyLRU) Range(iter *IntLRUNode) *IntLRUNode {
	if iter == nil {
		return lru.head
	}
	return iter.next
}

func NewIntLRU(initCapicity int) *IntKeyLRU {
	return &IntKeyLRU{
		strMapNode: make(map[int]*IntLRUNode, initCapicity),
		lock:       sync.Mutex{},
	}
}
