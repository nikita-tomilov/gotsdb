package utils

type HashCode func(interface{}) uint32
type Equals func(a, b interface{}) bool

type Map interface {
	Get(key interface{}) (value interface{})
	Put(key, value interface{})
	Delete(key interface{})
	Size() int
	Values() []interface{}
	Keys() []interface{}
}

type node struct {
	next *node
	ent  entry
}

type entry struct {
	key   interface{}
	value interface{}
}

func NewHashMap(hashcode HashCode, equals Equals) *HashMap {
	return &HashMap{
		m:        make(map[uint32]*node),
		hashcode: hashcode,
		equals:   equals,
	}
}

type HashMap struct {
	m        map[uint32]*node
	hashcode HashCode
	equals   Equals
}

func (hm *HashMap) Get(key interface{}) interface{} {
	h := hm.hashcode(key)
	root := hm.m[h]

	if root == nil {
		return nil
	}
	for n := root; n != nil; n = n.next {
		if hm.equals(n.ent.key, key) {
			return n.ent.value
		}
	}
	return nil
}

func (hm *HashMap) Put(key, value interface{}) {
	h := hm.hashcode(key)
	root := hm.m[h]
	if root == nil {
		hm.m[h] = &node{nil, entry{key, value}}
		return
	}
	this := root
	for {
		if hm.equals(this.ent.key, key) {
			this.ent.value = value
			return
		}
		if this.next == nil {
			this.next = &node{nil, entry{key, value}}
			return
		}
		this = this.next
	}
}

func (hm *HashMap) Delete(key interface{}) {
	h := hm.hashcode(key)
	root := hm.m[h]
	if root == nil {
		return
	}
	if hm.equals(root.ent.key, key) {
		if root.next == nil {
			delete(hm.m, h)
		} else {
			hm.m[h] = root.next
		}
	}
	for n := root; n.next != nil; n = n.next {
		m := n.next
		if hm.equals(m.ent.key, key) {
			n.next = m.next
			return
		}
	}
}

func (hm *HashMap) Values() []interface{} {
	values := make([]interface{}, hm.Size())
	i := 0
	for _, x := range hm.m {
		for n := x; n != nil; n = n.next {
			values[i] = n.ent.value
			i++
		}
	}
	return values
}

func (hm *HashMap) Keys() []interface{} {
	values := make([]interface{}, hm.Size())
	i := 0
	for _, x := range hm.m {
		for n := x; n != nil; n = n.next {
			values[i] = n.ent.key
			i++
		}
	}
	return values
}

func (hm *HashMap) Size() (s int) {
	for _, x := range hm.m {
		for n := x; n != nil; n = n.next {
			s++
		}
	}
	return
}
