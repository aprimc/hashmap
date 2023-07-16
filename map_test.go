package hashmap

import (
	"fmt"
	"hash/maphash"
	"testing"
)

func TestMap(t *testing.T) {
	m := Map[Int, Int]{}
	for i := 0; i < 100; i++ {
		m.Put(Int(i), Int(i*2))
	}
	if m.Size() != 100 {
		t.Errorf("expected size 100, got %d", m.Size())
	}

	for i := 0; i < 100; i++ {
		v, ok := m.Get(Int(i))
		if !ok {
			t.Errorf("expected to find key %d", i)
		}
		if v != Int(i*2) {
			t.Errorf("expected value %d, got %d", i*2, v)
		}
	}

	for i := 0; i < 100; i += 2 {
		m.Remove(Int(i))
	}
	if m.Size() != 50 {
		t.Errorf("expected size 50, got %d", m.Size())
	}
	if _, ok := m.Get(Int(0)); ok {
		t.Errorf("expected to not find key 0")
	}
	if _, ok := m.Get(Int(1)); !ok {
		t.Errorf("expected to not find key 1")
	}
}

// The benchmarks below are meant to check whether the overhead of the Map type is acceptable.
// If we're within an order of magnitude of the native map, we're good.
// We are not benchmarking deletes because the native map doesn't shrink when deleting elements and our map does.

func BenchmarkNativeMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := make(map[int]int)
		for j := 0; j < 1000; j++ {
			m[j] = j
		}
		for j := 0; j < 1000; j++ {
			_, ok := m[j]
			if !ok {
				b.Fatal("expected to find key")
			}
		}
	}
}

func BenchmarkMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		m := Map[Int, Int]{}
		for j := 0; j < 1000; j++ {
			m.Put(Int(j), Int(j))
		}
		for j := 0; j < 1000; j++ {
			_, ok := m.Get(Int(j))
			if !ok {
				b.Fatal("expected to find key")
			}
		}
	}
}

type bigKey struct {
	a String
	b String
	c String
	d String
	e String
}

func (bk bigKey) Hash() uint64 {
	h := maphash.Hash{}
	h.SetSeed(seed)
	h.WriteString(string(bk.a))
	h.WriteString(string(bk.b))
	h.WriteString(string(bk.c))
	h.WriteString(string(bk.d))
	h.WriteString(string(bk.e))
	return h.Sum64()
}

func (bk bigKey) Equals(other bigKey) bool {
	return bk.a == other.a && bk.b == other.b && bk.c == other.c && bk.d == other.d && bk.e == other.e
}

func BenchmarkBigKeyNativeMap(b *testing.B) {
	keys := make([]bigKey, 1000)
	for i := 0; i < 1000; i++ {
		keys[i] = bigKey{
			a: String(fmt.Sprint(i)),
			b: String(fmt.Sprint(i * 2)),
			c: String(fmt.Sprint(i * 3)),
			d: String(fmt.Sprint(i * 5)),
			e: String(fmt.Sprint(i * 7)),
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := make(map[bigKey]int)
		for _, k := range keys {
			m[k] = 0
		}
		for _, k := range keys {
			_, ok := m[k]
			if !ok {
				b.Fatal("expected to find key")
			}
		}
	}
}

func BenchmarkBigKeyMap(b *testing.B) {
	keys := make([]bigKey, 1000)
	for i := 0; i < 1000; i++ {
		keys[i] = bigKey{
			a: String(fmt.Sprint(i)),
			b: String(fmt.Sprint(i * 2)),
			c: String(fmt.Sprint(i * 3)),
			d: String(fmt.Sprint(i * 5)),
			e: String(fmt.Sprint(i * 7)),
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m := Map[bigKey, Int]{}
		for _, k := range keys {
			m.Put(k, 0)
		}
		for _, k := range keys {
			_, ok := m.Get(k)
			if !ok {
				b.Fatal("expected to find key")
			}
		}
	}
}
