package hashmap

import "testing"

func TestSet(t *testing.T) {
	s := Set[Int]{}
	for i := 0; i < 100; i++ {
		s.Add(Int(i))
	}
	if s.Size() != 100 {
		t.Errorf("expected size 100, got %d", s.Size())
	}

	for i := 0; i < 100; i++ {
		if !s.Contains(Int(i)) {
			t.Errorf("expected to find key %d", i)
		}
	}

	for i := 0; i < 100; i += 2 {
		s.Remove(Int(i))
	}
	if s.Size() != 50 {
		t.Errorf("expected size 50, got %d", s.Size())
	}
	if s.Contains(Int(0)) {
		t.Errorf("expected to not find key 0")
	}
	if !s.Contains(Int(1)) {
		t.Errorf("expected to find key 1")
	}
}

func intSet(is ...int) *Set[Int] {
	s := new(Set[Int])
	for _, i := range is {
		s.Add(Int(i))
	}
	return s
}

func TestSetEquality(t *testing.T) {
	s1 := intSet(1, 2, 3)
	s2 := intSet(1, 2, 3)
	if !s1.Equals(s2) {
		t.Errorf("expected %v to equal %v", s1, s2)
	}
	if !s2.Equals(s1) {
		t.Errorf("expected %v to equal %v", s2, s1)
	}

	s3 := intSet(3, 2, 1)
	if !s1.Equals(s3) {
		t.Errorf("expected %v to equal %v", s1, s3)
	}

	s4 := intSet(1, 2)
	s4.Add(Int(3))
	if !s1.Equals(s4) {
		t.Errorf("expected %v to equal %v", s1, s4)
	}

	s5 := intSet(1, 2, 3, 4)
	if s1.Equals(s5) {
		t.Errorf("expected %v to not equal %v", s1, s5)
	}
	s5.Remove(Int(4))
	if !s1.Equals(s5) {
		t.Errorf("expected %v to equal %v", s1, s5)
	}

	s6 := s1.Copy()
	if !s1.Equals(s6) {
		t.Errorf("expected %v to equal %v", s1, s6)
	}
	s6.Remove(Int(1))
	if s1.Equals(s6) {
		t.Errorf("expected %v to not equal %v", s1, s6)
	}

	s7 := new(Set[Int])
	s8 := new(Set[Int])
	s8.Add(Int(1))
	s8.Remove(Int(1))
	if !s7.Equals(s8) {
		t.Errorf("expected %v to equal %v", s7, s8)
	}
	if s7.Hash() != s8.Hash() {
		t.Errorf("expected %v to equal %v", s7.Hash(), s8.Hash())
	}
}

func TestSetIsSubset(t *testing.T) {
	tests := []struct {
		a, b *Set[Int]
		r    bool
	}{
		{intSet(1, 2, 3), intSet(1, 2, 3), true},
		{intSet(1, 2, 3), intSet(1, 2), false},
		{intSet(1, 2, 3), intSet(1, 2, 3, 4), true},
		{intSet(1, 2, 3), intSet(2, 3, 4, 5), false},
		{intSet(), intSet(1, 2, 3), true},
		{intSet(1, 2, 3), intSet(), false},
		{intSet(), intSet(), true},
	}
	for _, test := range tests {
		if test.a.IsSubset(test.b) != test.r {
			if test.r {
				t.Errorf("expected %v to be subset of %v", test.a, test.b)
			} else {
				t.Errorf("expected %v to not be subset of %v", test.a, test.b)
			}
		}
	}
}

func TestSetIsDisjoint(t *testing.T) {
	tests := []struct {
		a, b *Set[Int]
		r    bool
	}{
		{intSet(1, 2, 3), intSet(1, 2, 3), false},
		{intSet(1, 2, 3), intSet(1, 2), false},
		{intSet(1, 2, 3), intSet(1, 2, 3, 4), false},
		{intSet(1, 2, 3), intSet(2, 3, 4, 5), false},
		{intSet(), intSet(1, 2, 3), true},
		{intSet(1, 2, 3), intSet(), true},
		{intSet(), intSet(), true},
		{intSet(1, 2, 3), intSet(4, 5, 6), true},
	}
	for _, test := range tests {
		if test.a.IsDisjoint(test.b) != test.r {
			if test.r {
				t.Errorf("expected %v to be disjoint of %v", test.a, test.b)
			} else {
				t.Errorf("expected %v to not be disjoint of %v", test.a, test.b)
			}
		}
	}
}

func TestSetUnion(t *testing.T) {
	tests := []struct {
		a, b, r *Set[Int]
	}{
		{intSet(1, 2, 3), intSet(1, 2, 3), intSet(1, 2, 3)},
		{intSet(1, 2, 3), intSet(1, 2), intSet(1, 2, 3)},
		{intSet(1, 2, 3), intSet(1, 2, 3, 4), intSet(1, 2, 3, 4)},
		{intSet(1, 2, 3), intSet(2, 3, 4, 5), intSet(1, 2, 3, 4, 5)},
		{intSet(), intSet(1, 2, 3), intSet(1, 2, 3)},
		{intSet(1, 2, 3), intSet(), intSet(1, 2, 3)},
		{intSet(), intSet(), intSet()},
	}
	for _, test := range tests {
		if !test.a.Union(test.b).Equals(test.r) {
			t.Errorf("expected %v union %v to equal %v", test.a, test.b, test.r)
		}
	}
}

func TestSetIntersection(t *testing.T) {
	tests := []struct {
		a, b, r *Set[Int]
	}{
		{intSet(1, 2, 3), intSet(1, 2, 3), intSet(1, 2, 3)},
		{intSet(1, 2, 3), intSet(1, 2), intSet(1, 2)},
		{intSet(1, 2, 3), intSet(1, 2, 3, 4), intSet(1, 2, 3)},
		{intSet(1, 2, 3), intSet(2, 3, 4, 5), intSet(2, 3)},
		{intSet(), intSet(1, 2, 3), intSet()},
		{intSet(1, 2, 3), intSet(), intSet()},
		{intSet(), intSet(), intSet()},
	}
	for _, test := range tests {
		if !test.a.Intersection(test.b).Equals(test.r) {
			t.Errorf("expected %v intersection %v to equal %v", test.a, test.b, test.r)
		}
	}
}

func TestSetDifference(t *testing.T) {
	tests := []struct {
		a, b, r *Set[Int]
	}{
		{intSet(1, 2, 3), intSet(1, 2, 3), intSet()},
		{intSet(1, 2, 3), intSet(1, 2), intSet(3)},
		{intSet(1, 2, 3), intSet(1, 2, 3, 4), intSet()},
		{intSet(1, 2, 3), intSet(2, 3, 4, 5), intSet(1)},
		{intSet(), intSet(1, 2, 3), intSet()},
		{intSet(1, 2, 3), intSet(), intSet(1, 2, 3)},
		{intSet(), intSet(), intSet()},
	}
	for _, test := range tests {
		if !test.a.Difference(test.b).Equals(test.r) {
			t.Errorf("expected %v difference %v to equal %v", test.a, test.b, test.r)
		}
	}
}
