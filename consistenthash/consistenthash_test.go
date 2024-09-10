package consistenthash

import (
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	hash := NewConsistentHashMap(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})

	// Given the above hash function, this will give replicas with "hashes":
	// when add 6 with three replicas, the virtual nodes will be 6, 16, 26
	// when add 4 with three replicas, the virtual nodes will be 4, 14, 24
	// when add 2 with three replicas, the virtual nodes will be 2, 12, 22
	hash.AddNode("6", "4", "2")

	// for _, key := range hash.keys {
	// 	t.Logf("%v", key)
	// }
	// print(hash.hashMap)

	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s. But get %s", k, v, hash.Get(k))
		}
	}

	// Adds 8, 18, 28
	hash.AddNode("8")

	// 27 should now map to 8.
	testCases["27"] = "8"

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s. But get %s", k, v, hash.Get(k))
		}
	}

}
