package nonce

import "time"

type node struct {
	key  string
	exp  time.Time
	prev *node
	next *node
}

type LRU struct {
	cap  int
	size int
	head *node
	tail *node
	idx  map[string]*node
}

func NewLRU(cap int) *LRU {
	if cap <= 0 {
		cap = 1024
	}
	h := &node{}
	t := &node{}
	h.next = t
	t.prev = h
	return &LRU{cap: cap, head: h, tail: t, idx: make(map[string]*node)}
}

func (l *LRU) Len() int { return l.size }

func (l *LRU) Get(k string) (time.Time, bool) {
	n, ok := l.idx[k]
	if !ok {
		return time.Time{}, false
	}
	l.moveFront(n)
	return n.exp, true
}

func (l *LRU) Put(k string, exp time.Time) {
	if n, ok := l.idx[k]; ok {
		n.exp = exp
		l.moveFront(n)
		return
	}
	n := &node{key: k, exp: exp}
	l.idx[k] = n
	l.insertFront(n)
	l.size++
	if l.size > l.cap {
		l.evict()
	}
}

func (l *LRU) Delete(k string) {
	n, ok := l.idx[k]
	if !ok {
		return
	}
	l.remove(n)
	delete(l.idx, k)
	l.size--
}

func (l *LRU) insertFront(n *node) {
	n.prev = l.head
	n.next = l.head.next
	l.head.next.prev = n
	l.head.next = n
}

func (l *LRU) remove(n *node) {
	n.prev.next = n.next
	n.next.prev = n.prev
	n.prev = nil
	n.next = nil
}

func (l *LRU) moveFront(n *node) {
	l.remove(n)
	l.insertFront(n)
}

func (l *LRU) evict() {
	n := l.tail.prev
	if n == l.head {
		return
	}
	l.remove(n)
	delete(l.idx, n.key)
	l.size--
}

func (l *LRU) Sweep(now time.Time) int {
	var dead []string
	for k, n := range l.idx {
		if !n.exp.IsZero() && !now.Before(n.exp) {
			dead = append(dead, k)
		}
	}
	for _, k := range dead {
		l.Delete(k)
	}
	return len(dead)
}
