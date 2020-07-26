package kvs

type KeyValueStorage interface {
	InitStorage()
	CloseStorage()
	Save(key []byte, value []byte)
	KeyExists(key []byte) bool
	Retrieve(key []byte) []byte
	Delete(key []byte)
	GetAllKeys() [][]byte
	String() string
}