package context

import "sync"

// Dictionary exposes a sync locked collection of keys and values.
type Dictionary struct {
	store map[interface{}]interface{}
	sync.RWMutex
}

// NewDictionary returns an initialized Dictionary store.
func NewDictionary() *Dictionary {
	return &Dictionary{
		store: make(map[interface{}]interface{}),
	}
}

// Get returns the value for the specified key or nil if the key does not exist.
func (d *Dictionary) Get(key interface{}) (result interface{}) {
	d.Lock()
	result = d.store[key]
	d.Unlock()
	return
}

// Set adds or updates the specified key/value pair.
func (d *Dictionary) Set(key interface{}, value interface{}) {
	d.Lock()
	d.store[key] = value
	d.Unlock()
}

// Remove removes the specifed key.
func (d *Dictionary) Remove(key interface{}) {
	d.Lock()
	delete(d.store, key)
	d.Unlock()
}
