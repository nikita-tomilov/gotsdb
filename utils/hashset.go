package utils

type emptyItem struct {}

type Set interface {
	Put(key interface{})
	Delete(key interface{})
	Size() int
	Values() []interface{}
}

func NewHashSet(hashcode HashCode, equals Equals) *HashSet {
	return &HashSet{
		NewHashMap(hashcode, equals),
	}
}

type HashSet struct {
	m *HashMap
}

func (hs *HashSet) Put(key interface{}) {
	hs.m.Put(key, emptyItem{})
}

func (hs *HashSet) Delete(key interface{}) {
	hs.m.Delete(key)
}

func (hs *HashSet) Size() int {
	return hs.m.Size()
}

func (hs *HashSet) Values() []interface{} {
	return hs.m.Keys()
}

