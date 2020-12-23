package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	hash Hash
	//虚拟节点倍数
	replicas int
	keys     []int // Sorted
	hashMap  map[int]string
}

// New creates a Map instance
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (m *Map) Add(nodeKeys ...string) {
	for _, nodeKey := range nodeKeys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + nodeKey)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = nodeKey
		}

	}
	sort.Ints(m.keys)

}

func (m *Map) Get(nodeKey string) string {
	if len(nodeKey) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(nodeKey)))
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
