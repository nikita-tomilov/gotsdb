package kvs

import (
	"github.com/btcsuite/btcutil/base58"
	"sync"
)

type InMemKVS struct {
	data map[string][]byte
	lock sync.Mutex
}

func (f *InMemKVS) createKey(s []byte) string {
	return base58.Encode(s)
}

func (f *InMemKVS) getKey(s string) []byte {
	return base58.Decode(s)
}

func (f *InMemKVS) InitStorage() {
	f.data = make(map[string][]byte)
}

func (f *InMemKVS) CloseStorage() {
	//nothing here
}

func (f *InMemKVS) Save(key []byte, value []byte) {
	f.lock.Lock()
	f.data[f.createKey(key)] = value
	f.lock.Unlock()
}

func (f *InMemKVS) KeyExists(key []byte) bool {
	f.lock.Lock()
	_, found := f.data[f.createKey(key)]
	f.lock.Unlock()
	return found
}

func (f *InMemKVS) Retrieve(key []byte) []byte {
	f.lock.Lock()
	value, _ := f.data[f.createKey(key)]
	f.lock.Unlock()
	return value
}

func (f *InMemKVS) Delete(key []byte) {
	f.lock.Lock()
	delete(f.data, f.createKey(key))
	f.lock.Unlock()
}

func (f *InMemKVS) GetAllKeys() [][]byte {
	f.lock.Lock()
	keys := make([][]byte, len(f.data))  //wow go has no mechanism of retrieving map's keys
	i := 0
	for k := range f.data {
		keys[i] = f.getKey(k)
		i++
	}
	f.lock.Unlock()
	return keys
}

func (f *InMemKVS) String() string {
	return "Simple in-memory map-based KVS"
}