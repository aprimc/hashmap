package hashmap

const fullBit = 1 << 63
const initialCapacity = 16

type mapEntry[K Comparable[K], V any] struct {
	hash1 uint64
	key   K
	value V
}

// Map is a hash map that uses open addressing with linear probing.
// It is not thread-safe.
// The zero value is an empty map ready to use.
type Map[K Comparable[K], V any] struct {
	entries []mapEntry[K, V]
	size    int
}

func (m *Map[K, V]) init() {
	m.entries = make([]mapEntry[K, V], initialCapacity)
}

// Size returns the number of elements in the map.
func (m *Map[K, V]) Size() int {
	return m.size
}

// Get returns the value associated with the given key.
// The second return value indicates if the key was found.
// The value is the zero value for the value type if the key was not found.
func (m *Map[K, V]) Get(key K) (V, bool) {
	var zero V
	if m.size == 0 {
		return zero, false
	}
	hash1 := key.Hash() | fullBit
	index := hash1 & uint64(len(m.entries)-1)
	entry := m.entries[index]
	for entry.hash1 != 0 {
		if entry.hash1 == hash1 && entry.key.Equals(key) {
			return entry.value, true
		}
		index = (index + 1) & uint64(len(m.entries)-1)
		entry = m.entries[index]
	}
	return zero, false
}

// Put adds the given key/value pair to the map.
// If the key already exists, the value is updated.
func (m *Map[K, V]) Put(key K, value V) {
	if m.entries == nil {
		m.init()
	}
	hash1 := key.Hash() | fullBit
	m.putHash1(hash1, key, value)
}

func (m *Map[K, V]) putHash1(hash1 uint64, key K, value V) {
	index := hash1 & uint64(len(m.entries)-1)
	entry := m.entries[index]
	for entry.hash1 != 0 {
		if entry.hash1 == hash1 && entry.key.Equals(key) {
			entry.value = value
			return
		}
		index = (index + 1) & uint64(len(m.entries)-1)
		entry = m.entries[index]
	}
	m.entries[index] = mapEntry[K, V]{hash1, key, value}
	m.size++
	if m.size > 3*len(m.entries)/4 {
		m.resize(len(m.entries) * 2)
	}
}

func (m *Map[K, V]) resize(cap int) {
	entries := m.entries
	m.size = 0
	m.entries = make([]mapEntry[K, V], cap)
	for _, entry := range entries {
		if entry.hash1 != 0 {
			m.putHash1(entry.hash1, entry.key, entry.value)
		}
	}
}

// Remove removes the given key from the map.
func (m *Map[K, V]) Remove(key K) {
	if m.size == 0 {
		return
	}
	hash1 := key.Hash() | fullBit
	index := hash1 & uint64(len(m.entries)-1)
	entry := m.entries[index]
	for entry.hash1 != 0 {
		if entry.hash1 == hash1 && entry.key.Equals(key) {
			m.entries[index] = mapEntry[K, V]{}
			m.size--
			if len(m.entries) > initialCapacity && m.size < len(m.entries)/4 {
				m.resize(len(m.entries) / 2)
				return
			}
			index = (index + 1) & uint64(len(m.entries)-1)
			for m.entries[index].hash1 != 0 {
				entry := m.entries[index]
				m.entries[index] = mapEntry[K, V]{}
				m.size--
				m.putHash1(entry.hash1, entry.key, entry.value)
				index = (index + 1) & uint64(len(m.entries)-1)
			}
			return
		}
		index = (index + 1) & uint64(len(m.entries)-1)
		entry = m.entries[index]
	}
}

// ForEach calls the given function for each key/value pair in the map.
func (m *Map[K, V]) ForEach(f func(K, V) error) error {
	for _, entry := range m.entries {
		if entry.hash1 != 0 {
			if err := f(entry.key, entry.value); err != nil {
				return err
			}
		}
	}
	return nil
}

// Copy returns a copy of the map.
func (m *Map[K, V]) Copy() *Map[K, V] {
	if m.size == 0 {
		return &Map[K, V]{}
	}
	c := &Map[K, V]{}
	c.entries = make([]mapEntry[K, V], len(m.entries))
	copy(c.entries, m.entries)
	c.size = m.size
	return c
}
