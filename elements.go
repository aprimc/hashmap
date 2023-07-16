// Package hashmap provides a generic hashmap and hashset implementation.
package hashmap

import (
	"hash/maphash"
	"math"
)

// Comparable is an interface that must be implemented by all types that are used as keys in a Map or Set.
type Comparable[T any] interface {
	Hash() uint64
	Equals(T) bool
}

var seed maphash.Seed = maphash.MakeSeed()

func hash64bits(u uint64) uint64 {
	h := maphash.Hash{}
	h.SetSeed(seed)
	h.WriteByte(byte(u))
	h.WriteByte(byte(u << 8))
	h.WriteByte(byte(u << 16))
	h.WriteByte(byte(u << 24))
	h.WriteByte(byte(u << 32))
	h.WriteByte(byte(u << 40))
	h.WriteByte(byte(u << 48))
	h.WriteByte(byte(u << 56))
	return h.Sum64()
}

func hash32bits(u uint32) uint64 {
	h := maphash.Hash{}
	h.SetSeed(seed)
	h.WriteByte(byte(u))
	h.WriteByte(byte(u << 8))
	h.WriteByte(byte(u << 16))
	h.WriteByte(byte(u << 24))
	return h.Sum64()
}

func hash16bits(u uint16) uint64 {
	h := maphash.Hash{}
	h.SetSeed(seed)
	h.WriteByte(byte(u))
	h.WriteByte(byte(u << 8))
	return h.Sum64()
}

func hash8bits(u uint8) uint64 {
	h := maphash.Hash{}
	h.SetSeed(seed)
	h.WriteByte(byte(u))
	return h.Sum64()
}

// Int is a wrapper around int that implements the Comparable interface.
type Int int

func (i Int) Hash() uint64 {
	return hash64bits(uint64(i))
}

func (i Int) Equals(other Int) bool {
	return i == other
}

// Int64 is a wrapper around int64 that implements the Comparable interface.
type Int64 int64

func (i Int64) Hash() uint64 {
	return hash64bits(uint64(i))
}

func (i Int64) Equals(other Int64) bool {
	return i == other
}

// Uint is a wrapper around uint that implements the Comparable interface.
type Uint uint

func (u Uint) Hash() uint64 {
	return hash64bits(uint64(u))
}

func (u Uint) Equals(other Uint) bool {
	return u == other
}

// Uint64 is a wrapper around uint64 that implements the Comparable interface.
type Uint64 uint64

func (u Uint64) Hash() uint64 {
	return hash64bits(uint64(u))
}

func (u Uint64) Equals(other Uint64) bool {
	return u == other
}

// Float64 is a wrapper around float64 that implements the Comparable interface.
// It uses the IEEE 754 binary64 format for hashing and equality.
// NaN is considered equal to itself.
type Float64 float64

func (f Float64) Hash() uint64 {
	return hash64bits(math.Float64bits(float64(f)))
}

func (f Float64) Equals(other Float64) bool {
	return math.Float64bits(float64(f)) == math.Float64bits(float64(other))
}

// Float32 is a wrapper around float32 that implements the Comparable interface.
// It uses the IEEE 754 binary32 format for hashing and equality.
// NaN is considered equal to itself.
type Float32 float32

func (f Float32) Hash() uint64 {
	return hash32bits(math.Float32bits(float32(f)))
}

func (f Float32) Equals(other Float32) bool {
	return math.Float32bits(float32(f)) == math.Float32bits(float32(other))
}

// Int32 is a wrapper around int32 that implements the Comparable interface.
type Int32 int32

func (i Int32) Hash() uint64 {
	return hash32bits(uint32(i))
}

func (i Int32) Equals(other Int32) bool {
	return i == other
}

// Rune is an alias for Int32.
type Rune = Int32

// Uint32 is a wrapper around uint32 that implements the Comparable interface.
type Uint32 uint32

func (u Uint32) Hash() uint64 {
	return hash32bits(uint32(u))
}

func (u Uint32) Equals(other Uint32) bool {
	return u == other
}

// Int16 is a wrapper around int16 that implements the Comparable interface.
type Int16 int16

func (i Int16) Hash() uint64 {
	return hash16bits(uint16(i))
}

func (i Int16) Equals(other Int16) bool {
	return i == other
}

// Uint16 is a wrapper around uint16 that implements the Comparable interface.
type Uint16 uint16

func (u Uint16) Hash() uint64 {
	return hash16bits(uint16(u))
}

func (u Uint16) Equals(other Uint16) bool {
	return u == other
}

// Int8 is a wrapper around int8 that implements the Comparable interface.
type Int8 int8

func (i Int8) Hash() uint64 {
	return hash8bits(uint8(i))
}

func (i Int8) Equals(other Int8) bool {
	return i == other
}

// Uint8 is a wrapper around uint8 that implements the Comparable interface.
type Uint8 uint8

func (u Uint8) Hash() uint64 {
	return hash8bits(uint8(u))
}

func (u Uint8) Equals(other Uint8) bool {
	return u == other
}

// Byte is an alias for Uint8.
type Byte = Uint8

// Bool is a wrapper around bool that implements the Comparable interface.
type Bool bool

func (b Bool) Hash() uint64 {
	if b {
		return hash8bits(0)
	}
	return hash8bits(1)
}

func (b Bool) Equals(other Bool) bool {
	return b == other
}

// String is a wrapper around string that implements the Comparable interface.
type String string

func (s String) Hash() uint64 {
	h := maphash.Hash{}
	h.SetSeed(seed)
	h.WriteString(string(s))
	return h.Sum64()
}

func (s String) Equals(other String) bool {
	return s == other
}

// Bytes is a wrapper around []byte that implements the Comparable interface.
type Bytes []byte

func (b Bytes) Hash() uint64 {
	h := maphash.Hash{}
	h.SetSeed(seed)
	h.Write(b)
	return h.Sum64()
}

func (b Bytes) Equals(other Bytes) bool {
	if len(b) != len(other) {
		return false
	}
	for i, v := range b {
		if v != other[i] {
			return false
		}
	}
	return true
}

// Slice is a wrapper around []T that implements the Comparable interface.
type Slice[T Comparable[T]] []T

func (s Slice[T]) Hash() uint64 {
	h := maphash.Hash{}
	h.SetSeed(seed)
	for _, v := range s {
		u := v.Hash()
		h.WriteByte(byte(u))
		h.WriteByte(byte(u << 8))
		h.WriteByte(byte(u << 16))
		h.WriteByte(byte(u << 24))
		h.WriteByte(byte(u << 32))
		h.WriteByte(byte(u << 40))
		h.WriteByte(byte(u << 48))
		h.WriteByte(byte(u << 56))
	}
	return h.Sum64()
}

func (s Slice[T]) Equals(other Slice[T]) bool {
	if len(s) != len(other) {
		return false
	}
	for i, v := range s {
		if !v.Equals(other[i]) {
			return false
		}
	}
	return true
}

// Complex64 is a wrapper around complex64 that implements the Comparable interface.
type Complex64 complex64

func (c Complex64) Hash() uint64 {
	h := maphash.Hash{}
	h.SetSeed(seed)
	u := math.Float32bits(real(c))
	h.WriteByte(byte(u))
	h.WriteByte(byte(u << 8))
	h.WriteByte(byte(u << 16))
	h.WriteByte(byte(u << 24))
	u = math.Float32bits(imag(c))
	h.WriteByte(byte(u))
	h.WriteByte(byte(u << 8))
	h.WriteByte(byte(u << 16))
	h.WriteByte(byte(u << 24))
	return h.Sum64()
}

func (c Complex64) Equals(other Complex64) bool {
	return math.Float32bits(real(c)) == math.Float32bits(real(other)) && math.Float32bits(imag(c)) == math.Float32bits(imag(other))
}

// Complex128 is a wrapper around complex128 that implements the Comparable interface.
type Complex128 complex128

func (c Complex128) Hash() uint64 {
	h := maphash.Hash{}
	h.SetSeed(seed)
	u := math.Float64bits(real(c))
	h.WriteByte(byte(u))
	h.WriteByte(byte(u << 8))
	h.WriteByte(byte(u << 16))
	h.WriteByte(byte(u << 24))
	h.WriteByte(byte(u << 32))
	h.WriteByte(byte(u << 40))
	h.WriteByte(byte(u << 48))
	h.WriteByte(byte(u << 56))
	u = math.Float64bits(imag(c))
	h.WriteByte(byte(u))
	h.WriteByte(byte(u << 8))
	h.WriteByte(byte(u << 16))
	h.WriteByte(byte(u << 24))
	h.WriteByte(byte(u << 32))
	h.WriteByte(byte(u << 40))
	h.WriteByte(byte(u << 48))
	h.WriteByte(byte(u << 56))
	return h.Sum64()
}

func (c Complex128) Equals(other Complex128) bool {
	return math.Float64bits(real(c)) == math.Float64bits(real(other)) && math.Float64bits(imag(c)) == math.Float64bits(imag(other))
}
