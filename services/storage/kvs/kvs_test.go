package kvs

import (
	"fmt"
	"github.com/nikita-tomilov/gotsdb/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKVS_SaveRetrieveWork(t *testing.T) {
	//given
	storages := buildStorages()
	for _, s := range storages {
		func() {
			defer s.CloseStorage()
			//when
			s.Save(testKey1, testValue1)
			testRetrievedValue1 := s.Retrieve(testKey1)
			//then
			assert.Equal(t, testValue1, testRetrievedValue1, "Value should be retrieved by key at %s", s.String())
		}()
	}
}

func TestKVS_GetKeysAndDeleteWork(t *testing.T) {
	//given
	storages := buildStorages()
	for _, s := range storages {
		func() {
			defer s.CloseStorage()
			//when
			s.Save(testKey1, testValue1)
			s.Save(testKey2, testValue2)
			allKeys1 := s.GetAllKeys()
			//then
			assert.Equal(t, toBArray(testKey1, testKey2), allKeys1, "Should return all keys at %s", s.String())
			//when
			s.Delete(testKey1)
			allKeys2 := s.GetAllKeys()
			//then
			assert.Equal(t, toBArray(testKey2), allKeys2, "Should return second key at %s", s.String())
			//when
			testRetrievedValue1 := s.Retrieve(testKey1)
			//then
			assert.Nil(t, testRetrievedValue1, "Value for first key should be empty at %s", s.String())
		}()
	}
}


func buildStorages() []KeyValueStorage {
	inMem := buildInMemStorage()
	qL := buildFileStorage()
	return toArray(inMem, qL)
}

func toArray(items ...KeyValueStorage) []KeyValueStorage {
	return items
}

func buildInMemStorage() *InMemKVS {
	s := InMemKVS{}
	s.InitStorage()
	return &s
}

func buildFileStorage() *FileKVS {
	idx += 1
	s := FileKVS{Path: fmt.Sprintf("/tmp/gotsdb_test/test%d%d", utils.GetNowMillis(), idx)}
	s.InitStorage()
	return &s
}

var testKey1 = []byte("test-key-one")
var testValue1 = []byte("test-value-one")
var testKey2 = []byte("test-key-two")
var testValue2 = []byte("test-value-two")
var idx = 0

func toBArray(x ...[]byte) [][]byte {
	return x
}