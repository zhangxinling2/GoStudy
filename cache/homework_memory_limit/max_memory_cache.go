package cache

import (
	"context"
	"errors"
	"sync"
	"time"
)

type ListNode struct {
	Pre  *ListNode
	Next *ListNode
	key  string
}
type LRUCache struct {
	data map[string]*ListNode
	head *ListNode
	tail *ListNode
}

func NewList() *ListNode {
	head := &ListNode{}
	tail := &ListNode{}
	head.Next = tail
	tail.Pre = head
	return head
}

// addHead 将节点加到头部
func (l *ListNode) addHead(node *ListNode) {
	l.Next.Pre = node
	node.Next = l.Next
	l.Next = node
	node.Pre = l
}

// removeNode 删除链表中的节点
func removeNode(node *ListNode) {
	node.Next.Pre = node.Pre
	node.Pre.Next = node.Next
	node.Next = nil
	node.Pre = nil
}
func NewLRUCache(capacity int64) LRUCache {
	list := NewList()
	return LRUCache{
		data: make(map[string]*ListNode),
		head: list,
		tail: list.Next,
	}
}
func (l *LRUCache) Get(key string) error {
	res, ok := l.data[key]
	if !ok {
		return errors.New("key不存在")
	}
	removeNode(res)
	l.head.addHead(res)
	return nil
}
func (l *LRUCache) Del(key string) error {
	res, ok := l.data[key]
	if !ok {
		return errors.New("key不存在")
	}
	removeNode(res)
	delete(l.data, key)
	return nil
}

//Put 向缓存中添加数据,在使用put前判断有没有超出最大值
func (l *LRUCache) Put(key string) {
	//tp := reflect.TypeOf(val)
	//看现有数据有没有此key
	res, ok := l.data[key]
	//如果有就改变值，后将节点移到最前面
	if ok {
		removeNode(res)
		l.head.addHead(res)
		return
	}
	//如果没有，判断是否超出了单个最大值,如果是那么删除最后近最少使用的
	//for int64(tp.Size())+l.used > l.capacity {
	//	node := l.tail.Pre
	//	removeNode(node)
	//	delete(l.data, key)
	//	l.used = l.used - int64(tp.Size())
	//}
	newNode := &ListNode{
		key: key,
	}
	l.data[key] = newNode
	l.head.addHead(newNode)
	//l.used = l.used + int64(tp.Size())
}

type MaxMemoryCache struct {
	Cache
	cache LRUCache
	lock  *sync.RWMutex
	used  int64
	max   int64
}

var _ Cache = &MaxMemoryCache{}

func NewMaxMemoryCache(max int64, cache Cache) *MaxMemoryCache {
	res := &MaxMemoryCache{
		Cache: cache,
		cache: NewLRUCache(max),
		lock:  &sync.RWMutex{},
		max:   max,
		used:  0,
	}
	return res
}
func (m *MaxMemoryCache) Delete(ctx context.Context, key string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	err := m.Cache.Delete(ctx, key)
	if err != nil {
		return err
	}
	return nil
}
func (m *MaxMemoryCache) LoadAndDelete(ctx context.Context, key string) ([]byte, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	res, err := m.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	err = m.Delete(ctx, key)
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (m *MaxMemoryCache) Get(ctx context.Context, key string) ([]byte, error) {
	res, err := m.Cache.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	//Get成功后更新缓存
	m.cache.Get(key)
	return res, err
}

//Set 设置key,使用LRUCache中的PUT
func (m *MaxMemoryCache) Set(ctx context.Context, key string, val []byte,
	expiration time.Duration) error {
	//tp := reflect.TypeOf(val)
	m.lock.Lock()
	defer m.lock.Unlock()
	err := m.Cache.Set(ctx, key, val, expiration)
	if err != nil {
		return err
	}
	//循环删除直到能放下这个val
	//for int64(tp.Size())+m.used > m.max {
	//	node := m.cache.tail.Pre
	//	nodeKey := node.key
	//	err = m.Delete(ctx, nodeKey)
	//	if err != nil {
	//		return err
	//	}
	//}
	for int64(len(val))+m.used > m.max {
		node := m.cache.tail.Pre
		nodeKey := node.key
		err = m.Cache.Delete(ctx, nodeKey)
		if err != nil {
			return err
		}
	}
	//Set成功后更新缓存
	m.cache.Put(key)
	m.used = m.used + int64(len(val))
	//m.used = m.used + int64(tp.Size())
	return nil
}

func (m *MaxMemoryCache) OnEvicted(fn func(key string, val []byte)) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.Cache.OnEvicted(func(key string, val []byte) {
		//在这里减去used
		//tp := reflect.TypeOf(val)
		//m.used = m.used - int64(tp.Size())
		m.used = m.used - int64(len(val))
		//size减掉之后，把key在map中删掉
		err := m.cache.Del(key)
		if err != nil {
			return
		}
		fn(key, val)
	})
}
