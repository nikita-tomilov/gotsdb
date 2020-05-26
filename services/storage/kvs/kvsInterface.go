package kvs

type Storage interface {
	InitStorage()
	Save(key []byte, value []byte)
	KeyExists(key []byte) bool
	Retrieve(key []byte) []byte
	Delete(key []byte)
}