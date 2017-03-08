package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type Cache struct {
	items        map[string]Item
	lock         sync.RWMutex
	gcInterval   time.Duration
	stopGcSignal chan bool
}

func (c *Cache) gcLoop() {
	ticker := time.NewTicker(c.gcInterval)
	for {
		select {
		case <-ticker.C:
			c.DeleteExpiredDatas()
		case <-c.stopGcSignal:
			ticker.Stop()
			return
		}
	}
}

func (c *Cache) delete(k string) {
	delete(c.items, k)
}

func (c *Cache) DeleteExpiredDatas() {
	c.lock.Lock()
	defer c.lock.Unlock()
	for k, v := range c.items {
		if v.Expire() {
			c.delete(k)
		}
	}
}

func (c *Cache) Get(k string) (interface{}, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	item, ok := c.items[k]
	if ok && !item.Expire() {
		return item.v, true
	}
	return nil, false
}

func (c *Cache) Set(k string, v interface{}, d time.Duration) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	item, ok := c.items[k]
	if ok && !item.Expire() {
		c.items[k] = Item{v, time.Now().Add(d).UnixNano()}
		return true
	}
	return false
}

func (c *Cache) Add(k string, v interface{}, d time.Duration) error {
	c.lock.Lock()
	_, ok := c.items[k]
	if ok {
		c.lock.Unlock()
		return fmt.Errorf("Item %s already exists!", k)
	}
	c.items[k] = Item{v, time.Now().Add(d).UnixNano()}
	c.lock.Unlock()
	return nil
}

func (c *Cache) Del(k string) {
	c.lock.Lock()
	c.delete(k)
	c.lock.Unlock()
}

func (c *Cache) save(w io.Writer) (err error) {
	gb := gob.NewEncoder(w)
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("Error registering item types with Gob library")
		}
	}()
	c.lock.Lock()
	defer c.lock.Unlock()
	for _, v := range c.items {
		gob.Register(v.v)
	}
	err = gb.Encode(&c.items)
	return
}

func (c *Cache) SaveToFile(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	if err = c.save(f); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

func (c *Cache) load(r io.Reader) (err error) {
	gb := gob.NewDecoder(r)
	items := map[string]Item{}
	err = gb.Decode(&items)
	if err == nil {
		c.lock.Lock()
		defer c.lock.Unlock()
		for k, v := range items {
			if !v.Expire() {
				c.items[k] = v
			}
		}
	}
	return err
}

func (c *Cache) LoadFromFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	if err = c.load(f); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

func (c *Cache) Count() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return len(c.items)
}

func (c *Cache) Flush() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.items = map[string]Item{}
}

func (c *Cache) StopGC() {
	c.stopGcSignal <- true
}

func NewCache(gcInterval time.Duration) *Cache {
	c := &Cache{
		gcInterval:   gcInterval,
		items:        map[string]Item{},
		stopGcSignal: make(chan bool),
	}
	go c.gcLoop()
	return c
}
