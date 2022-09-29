package ds

import (
	"encoding/json"
	"fmt"
	"sync"
)

type Map[K comparable, V any] struct {
	data  map[K]V
	mutex sync.Mutex
}

// Exists checks if the value exists in map
func (set *Map[K, V]) Exists(key K) bool {
	set.mutex.Lock()
	defer set.mutex.Unlock()

	if set.data == nil {
		set.data = make(map[K]V)
	}

	_, ok := set.data[key]
	return ok
}

// Insert adds value to map, replacing any existing value
func (set *Map[K, V]) Insert(key K, value V) {
	set.mutex.Lock()
	defer set.mutex.Unlock()

	if set.data == nil {
		set.data = make(map[K]V)
	}

	set.data[key] = value
}

// Delete deletes value from map if exists, does nothing if it does not exist
func (set *Map[K, V]) Delete(key K) {
	set.mutex.Lock()
	defer set.mutex.Unlock()

	if set.data == nil {
		set.data = make(map[K]V)
	}

	delete(set.data, key)
}

// Get returns value from map or panic if key does not exist.
// If panic is not desired, should always use Exists to verify key existence.
func (set *Map[K, V]) Get(key K) V {
	set.mutex.Lock()
	defer set.mutex.Unlock()

	if set.data == nil {
		set.data = make(map[K]V)
	}

	result, ok := set.data[key]
	if !ok {
		panic(fmt.Sprintf("value for key '%v' does not exist", key))
	}
	return result
}

// Range will iterate through all key/values in map, returning false will stop iterations
// other behavior is same as sync.Map
func (set *Map[K, V]) Range(f func(key K, value V) bool) {
	set.mutex.Lock()
	defer set.mutex.Unlock()

	if set.data == nil {
		set.data = make(map[K]V)
	}

	for k, v := range set.data {
		if !f(k, v) {
			break
		}
	}
}

// Keys will return a list of all values in map
func (set *Map[K, V]) Keys() []K {
	set.mutex.Lock()
	defer set.mutex.Unlock()

	if set.data == nil {
		set.data = make(map[K]V)
	}

	result := make([]K, 0)
	for k := range set.data {
		result = append(result, k)
	}
	return result
}

// Values will return a list of all values in map
func (set *Map[K, V]) Values() []V {
	set.mutex.Lock()
	defer set.mutex.Unlock()

	if set.data == nil {
		set.data = make(map[K]V)
	}

	result := make([]V, 0)
	for _, v := range set.data {
		result = append(result, v)
	}
	return result
}

// MarshalJSON is used for easy JSON marshalling of this map
func (set *Map[K, V]) MarshalJSON() ([]byte, error) {
	set.mutex.Lock()
	defer set.mutex.Unlock()

	return json.Marshal(set.data)
}
