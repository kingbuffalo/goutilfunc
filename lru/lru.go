package lru

import (
	"sync"
	"time"
)

type StrLRUNode struct {
	data      interface{}
	dataDirty bool
	next      *StrLRUNode
	prev      *StrLRUNode
	time      int64
	key       string
}

type StrKeyLRU struct {
	head       *StrLRUNode
	tail       *StrLRUNode
	len        int
	strMapNode map[string]*StrLRUNode
	lock       sync.Mutex
}

func (lru *StrKeyLRU) getLRUNode(key string) *StrLRUNode {
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

func (lru *StrKeyLRU) GetValue(key string) interface{} {
	n := lru.getLRUNode(key)
	if n != nil {
		return n.data
	}
	return nil
}

func (lru *StrKeyLRU) RmValue(key string) interface{} {
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

func (lru *StrKeyLRU) AddValue(str string, data interface{}) bool {
	n := lru.getLRUNode(str)
	lru.lock.Lock()
	defer lru.lock.Unlock()
	if n != nil {
		n.data = data
		return false
	} else {
		n = &StrLRUNode{
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

func (lru *StrKeyLRU) RangeReverse(iter *StrLRUNode) *StrLRUNode {
	if iter == nil {
		return lru.tail
	}
	return iter.prev
}

func (lru *StrKeyLRU) Range(iter *StrLRUNode) *StrLRUNode {
	if iter == nil {
		return lru.head
	}
	return iter.next
}

func NewStrLRU(initCapicity int) *StrKeyLRU {
	return &StrKeyLRU{
		strMapNode: make(map[string]*StrLRUNode, initCapicity),
		lock:       sync.Mutex{},
	}
}
