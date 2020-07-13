package kvs

type KeyValueStorage interface {
	InitStorage()
	Save(key []byte, value []byte)
	KeyExists(key []byte) bool
	Retrieve(key []byte) []byte
	Delete(key []byte)
	GetAllKeys() [][]byte
}