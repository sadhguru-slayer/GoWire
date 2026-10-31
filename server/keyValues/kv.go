package keyvalues

import "sync"

type KeyValues struct {
	mu     sync.RWMutex
	data   map[string]string
	hashes map[string]map[string]string
}

func CreateKeyValue() *KeyValues {
	return &KeyValues{
		data:   map[string]string{},
		hashes: map[string]map[string]string{},
	}
}

func (kv *KeyValues) Set(key string, value string) {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.data[key] = value
}

func (kv *KeyValues) Get(key string) (string, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	value, exists := kv.data[key]
	return value, exists
}

func (kv *KeyValues) Del(key string) bool {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	_, exists := kv.data[key]
	if exists {
		delete(kv.data, key)
	}
	return exists
}

func (kv *KeyValues) Exists(key string) bool {
	kv.mu.RLock()
	defer kv.mu.RUnlock()
	_, exists := kv.data[key]
	return exists
}

// Hashes
func (kv *KeyValues) HSet(key, field, value string) bool {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	if _, exists := kv.hashes[key]; !exists {
		kv.hashes[key] = map[string]string{}
	}
	_, fieldExists := kv.hashes[key][field]

	kv.hashes[key][field] = value

	return !fieldExists
}

func (kv *KeyValues) HGet(key, field string) (string, bool) {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	hash, exists := kv.hashes[key]
	if !exists {
		return "", false
	}

	value, exists := hash[field]
	return value, exists
}

func (kv *KeyValues) HDel(key, field string) bool {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	hash, exists := kv.hashes[key]
	if !exists {
		return false
	}

	_, exists = hash[field]
	if !exists {
		return false
	}

	delete(hash, field)
	if len(hash) == 0 {
		delete(kv.hashes, key)
	}
	return true
}

func (kv *KeyValues) HExists(key, field string) bool {
	kv.mu.RLock()
	defer kv.mu.RUnlock()

	hash, exists := kv.hashes[key]
	if !exists {
		return false
	}

	_, exists = hash[field]
	return exists
}
