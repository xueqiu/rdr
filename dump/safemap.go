// Copyright 2017 XUEQIU.COM
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This code is mostly copied from Beego BeeMap
package dump

import (
	"encoding/json"
	"sync"
)

// SafeMap is a map which concurrency safe
type SafeMap struct {
	//TODO mutex with key
	lock *sync.RWMutex
	sm   map[interface{}]interface{}
}

// NewSafeMap return new safemap
func NewSafeMap() *SafeMap {
	return &SafeMap{
		lock: new(sync.RWMutex),
		sm:   make(map[interface{}]interface{}),
	}
}

// Get from maps return the k's value
func (m *SafeMap) Get(k interface{}) interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if val, ok := m.sm[k]; ok {
		return val
	}
	return nil
}

// Set the given key and value. Returns false
// if the key is already in the map and changes nothing.
func (m *SafeMap) Set(k interface{}, v interface{}) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	if _, ok := m.sm[k]; !ok {
		m.sm[k] = v
	}
	m.sm[k] = v
	return true
}

// Check seturns true if k is exist in the map.
func (m *SafeMap) Check(k interface{}) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if _, ok := m.sm[k]; !ok {
		return false
	}
	return true
}

// Delete the given key and value.
func (m *SafeMap) Delete(k interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.sm, k)
}

// Items returns all items in safemap.
func (m *SafeMap) Items() map[interface{}]interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	r := make(map[interface{}]interface{})
	for k, v := range m.sm {
		r[k] = v
	}
	return r
}

// String transfer map into json
func (m *SafeMap) String() (s string) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	r := make(map[string]interface{})
	for k, v := range m.sm {
		r[k.(string)] = v
	}
	b, err := json.Marshal(r)
	if err != nil {
		s = ""
		return
	}
	s = string(b)
	return
}
