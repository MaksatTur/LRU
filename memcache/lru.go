package memcache

import (
	"container/list"
	"sync"
	"time"
)

type Item struct {
	Key string
	Value interface{}
	Created time.Time
}

type LRU struct {
	sync.Mutex
	Capacity int
	Items map[string]*list.Element
	Queue *list.List
}

func NewLRU(capacity int) *LRU {
	return &LRU{
		Capacity:capacity,
		Items:make(map[string]*list.Element),
		Queue:list.New(),
	}
}

func (lru *LRU) Set(key string, value interface{}){
	if element, exist := lru.Items[key]; exist{
		// lock
		lru.Lock()
		defer lru.Unlock()
		lru.Queue.MoveToFront(element)
		element.Value.(*Item).Value = value
		element.Value.(*Item).Created = time.Now()
		return
	}
	if lru.Capacity == lru.Queue.Len(){
		//fmt.Println("full")
		lru.purge()
		//fmt.Println("done")
	}
	newItem := &Item{
		Key: key,
		Value:value,
		Created:time.Now(),
	}
	lru.Lock()
	defer lru.Unlock()
	element := lru.Queue.PushFront(newItem)
	lru.Items[newItem.Key] = element
}

func (lru *LRU) purge(){
	if back := lru.Queue.Back(); back != nil{
		lru.Lock()
		defer lru.Unlock()
		element := lru.Queue.Remove(back).(*Item)
		delete(lru.Items, element.Key)
	}
}

func (lru *LRU) GetElementValue(key string) interface{}  {
	element, exist := lru.Items[key]
	if exist{
		lru.Lock()
		lru.Queue.MoveToFront(element)
		element.Value.(*Item).Created = time.Now()
		lru.Unlock()
		return element.Value.(*Item).Value
	}
	return nil
}
