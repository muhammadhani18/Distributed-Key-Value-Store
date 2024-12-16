package store

import (
	"sync"
)

type KeyValueStore struct {
	data map[string]string
	mu   sync.RWMutex
}

// NewKeyValueStore creates a new KeyValueStore
func NewKeyValueStore() *KeyValueStore {
	return &KeyValueStore{
		data: make(map[string]string),
	}
}


// Put adds a key-value pair to the store
func (kvs *KeyValueStore) Put(key string, value string) {
	kvs.mu.Lock()
	defer kvs.mu.Unlock()
	kvs.data[key] = value
}

func (kvs *KeyValueStore) Get(key string) (string, bool) {
	kvs.mu.RLock()
	defer kvs.mu.RUnlock()
	value, exists := kvs.data[key]
	return value,  exists
}

func (kvs *KeyValueStore) Delete(key string) {
	kvs.mu.Lock()
	defer kvs.mu.Unlock()
	delete(kvs.data, key)
}

