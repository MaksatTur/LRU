package main

import (
	"container/list"
	"fmt"
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
	capacity int
	items map[string]*list.Element
	queue *list.List
}

func NewLRU(capacity int) *LRU {
	return &LRU{
		capacity:capacity,
		items:make(map[string]*list.Element),
		queue:list.New(),
	}
}

func (lru *LRU) Set(key string, value interface{}){
	if element, exist := lru.items[key]; exist{
		// lock
		lru.Lock()
		lru.queue.MoveToFront(element)
		element.Value.(*Item).Value = value
		element.Value.(*Item).Created = time.Now()
		// unlock
		lru.Unlock()
	}
	if lru.capacity == lru.queue.Len(){
		lru.purge()
	}
	newItem := &Item{
		Key: key,
		Value:value,
		Created:time.Now(),
	}
	lru.Lock()
	element := lru.queue.PushFront(newItem)
	lru.items[newItem.Key] = element
	lru.Unlock()
}

func (lru *LRU) purge(){
	if back := lru.queue.Back(); back != nil{
		lru.Lock()
		element := lru.queue.Remove(back)
		delete(lru.items, element.(Item).Key)
		lru.Unlock()
	}
}

func (lru *LRU) GetElementValue(key string) interface{}  {
	element, exist := lru.items[key]
	if exist{
		lru.Lock()
		lru.queue.MoveToFront(element)
		element.Value.(*Item).Created = time.Now()
		lru.Unlock()
		return element.Value.(*Item).Value
	}
	return nil
}

func main() {
	lru := NewLRU(5)
	lru.Set("1", "1")
	lru.Set("2", "2")
	lru.Set("3", 3)
	lru.Set("4", 4)
	lru.Set("5", nil)

	for key, _ := range lru.items{
		if element := lru.GetElementValue(key); element != nil{
			fmt.Printf("Key = %s Value = %v\n", key, element)
		} else {
			fmt.Printf("Key = %s Value = %s\n", key, "nil")
		}
		//fmt.Printf("Key = %s, value = %v\n", key, val.Value)
	}


}
