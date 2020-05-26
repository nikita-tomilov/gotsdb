package kvs

import (
	"github.com/btcsuite/btcutil/base58"
	log "github.com/jeanphorn/log4go"
	"github.com/programmer74/gotsdb/utils"
	"io/ioutil"
	"os"
	"sync"
)

type FileKVS struct {
	Path string `summer.property:"kvs.fileKVSPath|/tmp/gotsdb/kvs"`
	lock sync.Mutex
}

func (f *FileKVS) createKey(s []byte) string {
	return base58.Encode(s)
}

func (f *FileKVS) toFilename(key []byte) string {
	return f.Path + "/" + f.createKey(key)
}

func (f *FileKVS) InitStorage() {
	os.MkdirAll(f.Path, os.ModePerm)
	log.Warn("FILE-BASED KVS storage initialized at %s", f.Path)
}

func (f *FileKVS) Save(key []byte, value []byte) {
	f.lock.Lock()
	log.Warn("Request on setting value on key %s", string(key))
	err := ioutil.WriteFile(f.toFilename(key), value, 0644)
	utils.Check(err)
	f.lock.Unlock()
}

func (f *FileKVS) KeyExists(key []byte) bool {
	f.lock.Lock()
	fname := f.toFilename(key)
	ok := utils.FileExists(fname)
	f.lock.Unlock()
	return ok
}

func (f *FileKVS) Retrieve(key []byte) []byte {
	isPresent := f.KeyExists(key)
	log.Warn("Request on getting value on key %s", string(key))
	if !isPresent {
		return nil
	}
	fname := f.toFilename(key)
	f.lock.Lock()
	val, err := ioutil.ReadFile(fname)
	utils.Check(err)
	f.lock.Unlock()
	return val
}

func (f *FileKVS) Delete(key []byte) {
	f.lock.Lock()
	fname := f.toFilename(key)
	utils.DeleteFile(fname)
	f.lock.Unlock()
}
