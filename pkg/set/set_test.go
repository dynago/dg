// Set implementation by Calder Lund

package set

import (
	"fmt"
	"testing"
)

/* permutation code from https://yourbasic.org/golang/generate-permutation-slice-string/ */

// Perm calls f with each permutation of a.
func perm(a []string, f func([]string)) {
	permHelper(a, f, 0)
}

// Permute the values at index i to len(a)-1.
func permHelper(a []string, f func([]string), i int) {
	if i > len(a) {
		f(a)
		return
	}
	permHelper(a, f, i+1)
	for j := i + 1; j < len(a); j++ {
		a[i], a[j] = a[j], a[i]
		permHelper(a, f, i+1)
		a[i], a[j] = a[j], a[i]
	}
}

func printChecker(str string, values []string) bool {
	output := false
	perm(values, func(a []string) {
		potentialMatch := fmt.Sprint(a)
		if potentialMatch[1:len(potentialMatch)-1] == str[1:len(str)-1] {
			output = true
		}
	})
	return output
}

func TestAdd(t *testing.T) {
	s, err := MakeSet()
	if err != nil {
		t.Error(err)
	}
	if err := s.Add(1); err != nil {
		t.Error(err)
	}
	if err := s.Add(2.2); err != nil {
		t.Error(err)
	}
	if err := s.Add("hello"); err != nil {
		t.Error(err)
	}

	type tester struct {
		a int
		b int
	}
	var test tester
	test.a = 1
	test.b = 2
	if err := s.Add(test); err != nil {
		t.Error(err)
	}

	if ok := printChecker(s.String(), []string{"1", "2.2", "hello", "{1 2}"}); !ok {
		t.Fatalf("Got %s, which was unexpected", s.String())
	}
}

func TestRemove(t *testing.T) {
	type tester struct {
		a int
		b int
	}
	var test tester
	test.a = 1
	test.b = 2

	s, err := MakeSetFromValues(1, 2.2, "hello", test)
	if err != nil {
		t.Error(err)
	}

	if err := s.Remove(2.2); err != nil {
		t.Error(err)
	}
	var test2 tester
	test2.a = 1
	test2.b = 2
	if err := s.Remove(test2); err != nil {
		t.Error(err)
	}

	vals := []string{"1", "hello"}
	if ok := printChecker(s.String(), vals); !ok {
		t.Fatalf("Got %s, expected some sort of permutation of %v", s.String(), vals)
	}
}

func TestCombine(t *testing.T) {
	s1, err := MakeSetFromValues(1, 2.2)
	if err != nil {
		t.Error(err)
	}

	s2, err := MakeSetFromValues("hello", 2.2)
	if err != nil {
		t.Error(err)
	}

	if err := s1.Combine(s2); err != nil {
		t.Error(err)
	}

	vals := []string{"1", "2.2", "hello"}
	if ok := printChecker(s1.String(), vals); !ok {
		t.Fatalf("Got %s, expected some sort of permutation of %v", s1.String(), vals)
	}
}

func TestPop(t *testing.T) {
	s, err := MakeSetFromValues(1, 2.2, "hello")
	if err != nil {
		t.Error(err)
	}

	value, err := s.Pop()
	if err != nil {
		t.Error(err)
	}

	var vals []string
	switch value {
	case 1:
		vals = append(vals, "2.2")
		vals = append(vals, "hello")
	case 2.2:
		vals = append(vals, "1")
		vals = append(vals, "hello")
	case "hello":
		vals = append(vals, "1")
		vals = append(vals, "2.2")
	default:
		t.Errorf("Pop popped an unexpected value: %s", value)
	}

	if ok := printChecker(s.String(), vals); !ok {
		t.Fatalf("Got %s, expected some sort of permutation of 2 of %v", s.String(), vals)
	}
}

func TestClear(t *testing.T) {
	s, err := MakeSetFromValues(1, 2.2, "hello")
	if err != nil {
		t.Error(err)
	}

	if err := s.Clear(); err != nil {
		t.Error(err)
	}
	if s.Length() != 0 {
		t.Fatalf("Expected the length of the set to be 0, got %d", s.Length())
	}

	vals := []string{}
	if ok := printChecker(s.String(), vals); !ok {
		t.Fatalf("Got %s, expected some sort of permutation of %v", s.String(), vals)
	}

	// Check that we can still add
	if err := s.Add("hello"); err != nil {
		t.Error(err)
	}
}

func TestContains(t *testing.T) {
	s, err := MakeSetFromValues(1, 2.2, "hello")
	if err != nil {
		t.Error(err)
	}

	if contains, err := s.Contains(2.2); err != nil {
		t.Error(err)
	} else if contains == false {
		t.Fatal("The set does not contain 2.2 when it should")
	}
	if contains, err := s.Contains(42); err != nil {
		t.Error(err)
	} else if contains == true {
		t.Fatal("The set contains 42 when it shouldn't")
	}
}

func TestDisjoint(t *testing.T) {
	s1, err := MakeSetFromValues(1, 2.2)
	if err != nil {
		t.Error(err)
	}

	s2, err := MakeSetFromValues(1, "hello")
	if err != nil {
		t.Error(err)
	}

	s3, err := MakeSetFromValues("hello")
	if err != nil {
		t.Error(err)
	}
	if err := s3.Add("hello"); err != nil {
		t.Error(err)
	}

	if disjoint, err := s1.Disjoint(s2); err != nil {
		t.Error(err)
	} else if disjoint == true {
		t.Fatal("s1 and s2 should not be disjoint, but they are")
	}

	if disjoint, err := s2.Disjoint(s1); err != nil {
		t.Error(err)
	} else if disjoint == true {
		t.Fatal("s2 and s1 should not be disjoint, but they are")
	}

	if disjoint, err := s1.Disjoint(s3); err != nil {
		t.Error(err)
	} else if disjoint == false {
		t.Fatal("s1 and s3 should be disjoint, but they are not")
	}

	if disjoint, err := s3.Disjoint(s1); err != nil {
		t.Error(err)
	} else if disjoint == false {
		t.Fatal("s3 and s1 should be disjoint, but they are not")
	}
}

func TestEquals(t *testing.T) {
	s1, err := MakeSetFromValues(1, 2.2)
	if err != nil {
		t.Error(err)
	}

	s2, err := MakeSetFromValues(1, "2.2")
	if err != nil {
		t.Error(err)
	}

	s3, err := MakeSetFromValues(1, 2.2)
	if err != nil {
		t.Error(err)
	}

	s4, err1 := s1.Copy()
	if err1 != nil {
		t.Error(err1)
	}

	if equals, err := s1.Equals(s2); err != nil {
		t.Error(err)
	} else if equals == true {
		t.Fatal("s1 and s2 should not be equal, but they are")
	}

	if equals, err := s2.Equals(s1); err != nil {
		t.Error(err)
	} else if equals == true {
		t.Fatal("s2 and s1 should not be equal, but they are")
	}

	if equals, err := s1.Equals(s3); err != nil {
		t.Error(err)
	} else if equals == false {
		t.Fatal("s1 and s3 should be equal, but they are not")
	}

	if equals, err := s3.Equals(s1); err != nil {
		t.Error(err)
	} else if equals == false {
		t.Fatal("s3 and s1 should be equal, but they are not")
	}

	if equals, err := s1.Equals(s4); err != nil {
		t.Error(err)
	} else if equals == false {
		t.Fatal("s1 and s4 should be equal, but they are not")
	}

	if equals, err := s4.Equals(s1); err != nil {
		t.Error(err)
	} else if equals == false {
		t.Fatal("s4 and s1 should be equal, but they are not")
	}
}

func TestSupersetAndSubset(t *testing.T) {
	s1, err := MakeSetFromValues(1, 2.2, "hello")
	if err != nil {
		t.Error(err)
	}

	s2, err := MakeSetFromValues(1, "2.2")
	if err != nil {
		t.Error(err)
	}

	s3, err := MakeSetFromValues(1, 2.2)
	if err != nil {
		t.Error(err)
	}

	if superset, err := s1.SupersetOf(s2); err != nil {
		t.Error(err)
	} else if superset == true {
		t.Fatal("s1 should not be a superset of s2, but is")
	}

	if superset, err := s1.SupersetOf(s3); err != nil {
		t.Error(err)
	} else if superset == false {
		t.Fatal("s1 should be a superset of s3, but is not")
	}

	if subset, err := s3.SubsetOf(s1); err != nil {
		t.Error(err)
	} else if subset == false {
		t.Fatal("s3 should be a subset of s1, but is not")
	}

	if subset, err := s3.SubsetOf(s2); err != nil {
		t.Error(err)
	} else if subset == true {
		t.Fatal("s3 should not be a subset of s2, but is")
	}
}

func TestIntersection(t *testing.T) {
	s1, err := MakeSetFromValues(1, 2.2)
	if err != nil {
		t.Error(err)
	}

	s2, err := MakeSetFromValues(1, "hello")
	if err != nil {
		t.Error(err)
	}

	vals := []string{"1"}

	s, err := s1.Intersection(s2)
	if err != nil {
		t.Error(err)
	}
	if ok := printChecker(s.String(), vals); !ok {
		t.Fatalf("Got %s, expected some sort of permutation of %v", s.String(), vals)
	}

	s, err = s2.Intersection(s1)
	if err != nil {
		t.Error(err)
	}
	if ok := printChecker(s.String(), vals); !ok {
		t.Fatalf("Got %s, expected some sort of permutation of %v", s.String(), vals)
	}
}

func TestDifference(t *testing.T) {
	s1, err := MakeSetFromValues(1, 2.2)
	if err != nil {
		t.Error(err)
	}

	s2, err := MakeSetFromValues(1, "hello")
	if err != nil {
		t.Error(err)
	}

	diff12 := []string{"2.2"}
	diff21 := []string{"hello"}
	symm := []string{"2.2", "hello"}

	s, err := s1.Difference(s2)
	if err != nil {
		t.Error(err)
	}
	if ok := printChecker(s.String(), diff12); !ok {
		t.Fatalf("Got %s, expected some sort of permutation of %v", s.String(), diff12)
	}

	s, err = s2.Difference(s1)
	if err != nil {
		t.Error(err)
	}
	if ok := printChecker(s.String(), diff21); !ok {
		t.Fatalf("Got %s, expected some sort of permutation of %v", s.String(), diff21)
	}

	s, err = s1.SymmetricDifference(s2)
	if err != nil {
		t.Error(err)
	}
	if ok := printChecker(s.String(), symm); !ok {
		t.Fatalf("Got %s, expected some sort of permutation of %v", s.String(), symm)
	}

	s, err = s2.SymmetricDifference(s1)
	if err != nil {
		t.Error(err)
	}
	if ok := printChecker(s.String(), symm); !ok {
		t.Fatalf("Got %s, expected some sort of permutation of %v", s.String(), symm)
	}
}

func TestUnion(t *testing.T) {
	s1, err := MakeSetFromValues(1, 2.2)
	if err != nil {
		t.Error(err)
	}

	s2, err := MakeSetFromValues(1, "hello")
	if err != nil {
		t.Error(err)
	}

	s3, err := MakeSet()
	if err != nil {
		t.Error(err)
	}

	vals := []string{"1", "2.2", "hello"}

	s, err := s1.Union(s2)
	if err != nil {
		t.Error(err)
	}
	if ok := printChecker(s.String(), vals); !ok {
		t.Fatalf("Got %s, expected some sort of permutation of %v", s.String(), vals)
	}

	s, err = s.Union(s3)
	if err != nil {
		t.Error(err)
	}
	if ok := printChecker(s.String(), vals); !ok {
		t.Fatalf("Got %s, expected some sort of permutation of %v", s.String(), vals)
	}
}

func TestIterable(t *testing.T) {
	s1, err := MakeSetFromValues(1, 2.2)
	if err != nil {
		t.Error(err)
	}

	s, err := MakeSet(s1)
	if err != nil {
		t.Error(err)
	}

	vals := []string{"1", "2.2"}
	if ok := printChecker(s.String(), vals); !ok {
		t.Fatalf("Got %s, expected some sort of permutation of %v", s.String(), vals)
	}
}
