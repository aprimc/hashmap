package hashmap

type setEntry[K Comparable[K]] struct {
	hash1 uint64
	key   K
}

// Set is a hash set that uses open addressing with linear probing.
// It is not thread-safe.
// The zero value is an empty set ready to use.
// Set implements Comparable, so it can be used as a key in a Map or an element in a Set.
type Set[K Comparable[K]] struct {
	entries []setEntry[K]
	size    int
	hash    uint64
}

var emptySetHash uint64 = Int(0).Hash()

func (s *Set[K]) init() {
	s.entries = make([]setEntry[K], initialCapacity)
	s.hash = emptySetHash
}

// Size returns the number of elements in the set.
func (s *Set[K]) Size() int {
	return s.size
}

// Contains returns true if the set contains the given key.
func (s *Set[K]) Contains(key K) bool {
	if s.size == 0 {
		return false
	}
	hash1 := key.Hash() | fullBit
	return s.containsHash1Key(hash1, key)
}

func (s Set[K]) containsHash1Key(hash1 uint64, key K) bool {
	index := hash1 & uint64(len(s.entries)-1)
	entry := s.entries[index]
	for entry.hash1 != 0 {
		if entry.hash1 == hash1 && entry.key.Equals(key) {
			return true
		}
		index = (index + 1) & uint64(len(s.entries)-1)
		entry = s.entries[index]
	}
	return false
}

// Add adds the given key to the set.
func (s *Set[K]) Add(key K) {
	hash1 := key.Hash() | fullBit
	s.addHash1Key(hash1, key)
}

func (s *Set[K]) addHash1Key(hash1 uint64, key K) {
	if s.entries == nil {
		s.init()
	}
	index := hash1 & uint64(len(s.entries)-1)
	entry := s.entries[index]
	for entry.hash1 != 0 {
		if entry.hash1 == hash1 && entry.key.Equals(key) {
			return
		}
		index = (index + 1) & uint64(len(s.entries)-1)
		entry = s.entries[index]
	}
	s.entries[index] = setEntry[K]{hash1, key}
	s.size++
	s.hash ^= hash1
	if s.size > 3*len(s.entries)/4 {
		s.resize(len(s.entries) * 2)
	}
}

func (s *Set[K]) resize(cap int) {
	entries := s.entries
	s.size = 0
	s.hash = emptySetHash
	s.entries = make([]setEntry[K], cap)
	for _, entry := range entries {
		if entry.hash1 != 0 {
			s.addHash1Key(entry.hash1, entry.key)
		}
	}
}

// Remove removes the given key from the set.
func (s *Set[K]) Remove(key K) {
	if s.size == 0 {
		return
	}
	hash1 := key.Hash() | fullBit
	s.removeHash1Key(hash1, key)
}

func (s *Set[K]) removeHash1Key(hash1 uint64, key K) {
	index := hash1 & uint64(len(s.entries)-1)
	entry := s.entries[index]
	for entry.hash1 != 0 {
		if entry.hash1 == hash1 && entry.key.Equals(key) {
			s.entries[index] = setEntry[K]{}
			s.size--
			s.hash ^= hash1
			if len(s.entries) > initialCapacity && s.size < len(s.entries)/4 {
				s.resize(len(s.entries) / 2)
			}
			index = (index + 1) & uint64(len(s.entries)-1)
			for s.entries[index].hash1 != 0 {
				entry := s.entries[index]
				s.entries[index] = setEntry[K]{}
				s.size--
				s.hash ^= entry.hash1
				s.addHash1Key(entry.hash1, entry.key)
				index = (index + 1) & uint64(len(s.entries)-1)
			}
			return
		}
		index = (index + 1) & uint64(len(s.entries)-1)
		entry = s.entries[index]
	}
}

// ForEach calls the given function for each key in the set.
func (s *Set[K]) ForEach(f func(K) error) error {
	for _, entry := range s.entries {
		if entry.hash1 != 0 {
			if err := f(entry.key); err != nil {
				return err
			}
		}
	}
	return nil
}

// Copy returns a copy of the set.
func (s *Set[K]) Copy() *Set[K] {
	if s.size == 0 {
		return &Set[K]{}
	}
	c := &Set[K]{}
	c.entries = make([]setEntry[K], len(s.entries))
	copy(c.entries, s.entries)
	c.size = s.size
	c.hash = s.hash
	return c
}

// Hash returns the hash code for the set.
func (s *Set[K]) Hash() uint64 {
	if s.size == 0 {
		return emptySetHash
	}
	return s.hash
}

// Equals returns true if the set is equal to the given set.
func (s *Set[K]) Equals(t *Set[K]) bool {
	if s.size != t.size {
		return false
	}
	if s.size == 0 {
		return true
	}
	if s.Hash() != t.Hash() {
		return false
	}
	if len(s.entries) > len(t.entries) {
		s, t = t, s
	}
	for _, entry := range s.entries {
		if entry.hash1 != 0 {
			if !t.containsHash1Key(entry.hash1, entry.key) {
				return false
			}
		}
	}
	return true
}

// IsSubset returns true if the set is a subset of the given set.
func (s *Set[K]) IsSubset(t *Set[K]) bool {
	if s.size > t.size {
		return false
	}
	if s.size == 0 {
		return true
	}
	for _, entry := range s.entries {
		if entry.hash1 != 0 {
			if !t.containsHash1Key(entry.hash1, entry.key) {
				return false
			}
		}
	}
	return true
}

// IsDisjoint returns true if the intersection of the set and the given set is empty.
func (s *Set[K]) IsDisjoint(t *Set[K]) bool {
	if s.size == 0 || t.size == 0 {
		return true
	}
	if len(s.entries) > len(t.entries) {
		s, t = t, s
	}
	for _, entry := range s.entries {
		if entry.hash1 != 0 {
			if t.containsHash1Key(entry.hash1, entry.key) {
				return false
			}
		}
	}
	return true
}

// Union returns a new set with all the elements in both sets.
func (s *Set[K]) Union(t *Set[K]) *Set[K] {
	if s.size == 0 {
		return t.Copy()
	}
	if t.size == 0 {
		return s.Copy()
	}
	if len(s.entries) > len(t.entries) {
		s, t = t, s
	}
	r := t.Copy()
	for _, entry := range s.entries {
		if entry.hash1 != 0 {
			r.addHash1Key(entry.hash1, entry.key)
		}
	}
	return r
}

// Intersection returns a new set with the elements that are in both sets.
func (s *Set[K]) Intersection(t *Set[K]) *Set[K] {
	r := new(Set[K])
	if s.size == 0 || t.size == 0 {
		return r
	}
	if len(s.entries) > len(t.entries) {
		s, t = t, s
	}
	for _, entry := range s.entries {
		if entry.hash1 != 0 {
			if t.containsHash1Key(entry.hash1, entry.key) {
				r.addHash1Key(entry.hash1, entry.key)
			}
		}
	}
	return r
}

// Difference returns a new set with the elements that are in the set but not in the given set.
func (s *Set[K]) Difference(t *Set[K]) *Set[K] {
	r := s.Copy()
	if s.size == 0 || t.size == 0 {
		return r
	}
	for _, entry := range t.entries {
		if entry.hash1 != 0 {
			r.removeHash1Key(entry.hash1, entry.key)
		}
	}
	return r
}
