package cache

/// for saving the actual data in the cache


// for saving the data in the cache
// byte for allowing all kinds of data
type ByteView struct {
	b []byte
}

// get the length
func (v ByteView) Len() int {
	return len(v.b)
}

// get the byte slice as a copy
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

// get the byte slice as a string
func (v ByteView) String() string {
	return string(v.b)
}

// helper function for copying bytes, to avoid the modification of the original data
func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}