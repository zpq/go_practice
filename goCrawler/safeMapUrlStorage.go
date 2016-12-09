package main

import "sync"

type mapStorage struct {
	syn  sync.Mutex
	urls map[string][]string
}

func newMapStorage() *mapStorage {
	return &mapStorage{
		urls: make(map[string][]string),
	}
}

func (m *mapStorage) getOneURL(k string) string {
	m.syn.Lock()
	defer m.syn.Unlock()
	// fmt.Println(len(m.urls[k]))
	if len(m.urls[k]) > 0 {
		str := m.urls[k][0]
		m.urls[k] = m.urls[k][1:]
		return str
	}
	return ""
}

func (m *mapStorage) addOneURL(k, v string) {
	m.syn.Lock()
	defer m.syn.Unlock()
	m.urls[k] = append(m.urls[k], v)
}
