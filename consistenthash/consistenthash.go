package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash maps bytes to uint32
type Hash func(data []byte) uint32

// Map contains all hashed keys
type ConsistentHashMap struct {
	hash     Hash
	replicas int
	keys     []int
	hashMap  map[int]string
}

// New creates a Map instance
func NewConsistentHashMap(replicas int, fn Hash) *ConsistentHashMap {
	m := &ConsistentHashMap{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add adds some keys to the hash, add nodes to the ring
func (m *ConsistentHashMap) AddNode(peers ...string) {
	for _, peer := range peers {
		for i := 0; i < m.replicas; i++ {
			// name virtual node as index + key, and do the hash
			hash := int(m.hash([]byte(strconv.Itoa(i) + peer)))

			// add virtual node hash to the ring
			m.keys = append(m.keys, hash)

			// add virtual node read node mapping to the map
			m.hashMap[hash] = peer
		}
	}
	sort.Ints(m.keys)
}

// Get the node by a key
func (m *ConsistentHashMap) Get(key string) string {
	// no nodes
	if len(m.keys) == 0 {
		return ""
	}

	// hash the key
	hash := int(m.hash([]byte(key)))

	// find the closest virtual node in the ring
	// binary search to get the closest virtual node
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	// get the real node name
	// for handling edge case that the idx == len(m.keys), use mod to get the first virtual node
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
