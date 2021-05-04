package dict

import (
	"fmt"
	"testing"

	"github.com/CalderLund/DynamicGo/tuple"
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

func printChecker(str string, items []string) bool {
	output := false
	perm(items, func(a []string) {
		potentialMatch := fmt.Sprint(a)
		if potentialMatch[1:len(potentialMatch)-1] == str[1:len(str)-1] {
			output = true
		}
	})
	return output
}

func TestMakeEmptyDict(t *testing.T) {
	s, err := MakeDict()
	if err != nil {
		t.Error(err)
	}

	values := []string{}
	if ok := printChecker(s.String(), values); !ok {
		t.Fatalf("Got %s, which was unexpected", s.String())
	}
}

type dictTester struct {
	a []interface{}
}

func (test *dictTester) Iterate() <-chan interface{} {
	c := make(chan interface{})
	go func() {
		for _, v := range test.a {
			c <- v
		}
		close(c)
	}()
	return c
}

func TestMakePopulatedDict(t *testing.T) {
	type tester struct {
		a int
		b int
	}
	var test tester
	test.a = 1
	test.b = 2

	values := []string{"(1 2.2)", "(hello {1 2})", "({1 2} 42)"}

	d := dictTester{[]interface{}{1, 2.2, "hello", test, test, 42}}

	s, err := MakeDict(&d)
	if err != nil {
		t.Error(err)
	}
	if ok := printChecker(s.String(), values); !ok {
		t.Fatalf("Got %s, which was unexpected", s.String())
	}

	var t1, t2, t3 tuple.TupleInterface
	if t1, err = tuple.MakeTupleFromValues(1, 2.2); err != nil {
		t.Error(err)
	}
	if t2, err = tuple.MakeTupleFromValues("hello", test); err != nil {
		t.Error(err)
	}
	if t3, err = tuple.MakeTupleFromValues(test, 42); err != nil {
		t.Error(err)
	}

	s, err = MakeDictFromItems(t1, t2, t3)
	if err != nil {
		t.Error(err)
	}
	if ok := printChecker(s.String(), values); !ok {
		t.Fatalf("Got %s, which was unexpected", s.String())
	}

	s, err = MakeDictFromKeyValues([]interface{}{1, "hello", test}, []interface{}{2.2, test, 42})
	if err != nil {
		t.Error(err)
	}
	if ok := printChecker(s.String(), values); !ok {
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

	s, err := MakeDictFromKeyValues([]interface{}{1, "hello", test}, []interface{}{2.2, test, 42})
	if err != nil {
		t.Error(err)
	}

	if err := s.Remove("hello"); err != nil {
		t.Error(err)
	}

	var test2 tester
	test2.a = 1
	test2.b = 2
	if err := s.Remove(test2); err != nil {
		t.Error(err)
	}

	vals := []string{"(1 2.2)"}
	if ok := printChecker(s.String(), vals); !ok {
		t.Fatalf("Got %s, expected some sort of permutation of %v", s.String(), vals)
	}
}

func TestGet(t *testing.T) {
	s, err1 := MakeDictFromKeyValues([]interface{}{1, "hello"}, []interface{}{2.2, 42})
	if err1 != nil {
		t.Error(err1)
	}

	if value, err := s.Get(1); err != nil {
		t.Error(err)
	} else if value != 2.2 {
		t.Fatalf("Expected value to be 2.2, got %s", value)
	}

	if value, err := s.Get("hello"); err != nil {
		t.Error(err)
	} else if value != 42 {
		t.Fatalf("Expected value to be 42, got %s", value)
	}

	if value, err := s.Get(2.2); err != nil {
		t.Error(err)
	} else if value != nil {
		t.Fatalf("Expected nil value, got %s", value)
	}
}

func TestSet(t *testing.T) {
	s, err1 := MakeDictFromKeyValues([]interface{}{1, "hello"}, []interface{}{2.2, 42})
	if err1 != nil {
		t.Error(err1)
	}

	vals := []string{"(1 2.2)", "(hello 42)", "(100 -200.2)"}
	if err := s.Set(100, -200.2); err != nil {
		t.Error(err)
	} else if ok := printChecker(s.String(), vals); !ok {
		t.Fatalf("Got %s, expected some sort of permutation of %v", s.String(), vals)
	}

	vals = []string{"(1 -2.2)", "(hello 42)", "(100 -200.2)"}
	if err := s.Set(1, -2.2); err != nil {
		t.Error(err)
	} else if ok := printChecker(s.String(), vals); !ok {
		t.Fatalf("Got %s, expected some sort of permutation of %v", s.String(), vals)
	}

	if err := s.Set(nil, -2.2); err == nil {
		t.Fatalf("Expecting error (no nil key)")
	}
}

func TestCombine(t *testing.T) {
	s1, err1 := MakeDictFromKeyValues([]interface{}{1, "hello"}, []interface{}{2.2, 42})
	if err1 != nil {
		t.Error(err1)
	}

	s2, err2 := MakeDictFromKeyValues([]interface{}{3}, []interface{}{3})
	if err2 != nil {
		t.Error(err2)
	}

	s3, err3 := MakeDict()
	if err3 != nil {
		t.Error(err3)
	}

	if err := s1.Combine(s2); err != nil {
		t.Error(err)
	}

	vals := []string{"(1 2.2)", "(hello 42)", "(3 3)"}
	if ok := printChecker(s1.String(), vals); !ok {
		t.Fatalf("Got %s, expected some sort of permutation of %v", s1.String(), vals)
	}

	if err := s2.Combine(s3); err != nil {
		t.Error(err)
	}

	vals = []string{"(3 3)"}
	if ok := printChecker(s2.String(), vals); !ok {
		t.Fatalf("Got %s, expected some sort of permutation of %v", s1.String(), vals)
	}
}

func TestPopKey(t *testing.T) {
	s, err1 := MakeDictFromKeyValues([]interface{}{1, 2, 3}, []interface{}{1, 2, 3})
	if err1 != nil {
		t.Error(err1)
	}

	key, err := s.PopKey()
	if err != nil {
		t.Error(err)
	}
	var vals []string
	switch key {
	case 1:
		vals = append(vals, "(2 2)")
		vals = append(vals, "(3 3)")
	case 2:
		vals = append(vals, "(1 1)")
		vals = append(vals, "(3 3)")
	case 3:
		vals = append(vals, "(1 1)")
		vals = append(vals, "(2 2)")
	default:
		t.Errorf("Pop popped an unexpected value: %s", key)
	}

	if ok := printChecker(s.String(), vals); !ok {
		t.Fatalf("Got %s, expected some sort of permutation of %v", s.String(), vals)
	}
}

func TestPopValue(t *testing.T) {
	s, err1 := MakeDictFromKeyValues([]interface{}{1, 2, 3}, []interface{}{1, 2, 3})
	if err1 != nil {
		t.Error(err1)
	}

	value, err := s.PopValue()
	if err != nil {
		t.Error(err)
	}

	var vals []string
	switch value {
	case 1:
		vals = append(vals, "(2 2)")
		vals = append(vals, "(3 3)")
	case 2:
		vals = append(vals, "(1 1)")
		vals = append(vals, "(3 3)")
	case 3:
		vals = append(vals, "(1 1)")
		vals = append(vals, "(2 2)")
	default:
		t.Errorf("Pop popped an unexpected value: %s", value)
	}

	if ok := printChecker(s.String(), vals); !ok {
		t.Fatalf("Got %s, expected some sort of permutation of %v", s.String(), vals)
	}
}

func TestPop(t *testing.T) {
	s, err1 := MakeDictFromKeyValues([]interface{}{1, 2, 3}, []interface{}{1, 2, 3})
	if err1 != nil {
		t.Error(err1)
	}

	key, value, err := s.Pop()
	if err != nil {
		t.Error(err)
	} else if key != value {
		t.Fatalf("Pop did not pop correctly")
	}
	var vals []string
	switch value {
	case 1:
		vals = append(vals, "(2 2)")
		vals = append(vals, "(3 3)")
	case 2:
		vals = append(vals, "(1 1)")
		vals = append(vals, "(3 3)")
	case 3:
		vals = append(vals, "(1 1)")
		vals = append(vals, "(2 2)")
	default:
		t.Errorf("Pop popped an unexpected value: %s", value)
	}

	if ok := printChecker(s.String(), vals); !ok {
		t.Fatalf("Got %s, expected some sort of permutation of %v", s.String(), vals)
	}
}

func TestClear(t *testing.T) {
	s, err1 := MakeDictFromKeyValues([]interface{}{1, 2, 3}, []interface{}{1, 2, 3})
	if err1 != nil {
		t.Error(err1)
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
	if err := s.Set(1, 1); err != nil {
		t.Error(err)
	}
}

func TestContains(t *testing.T) {
	s, err1 := MakeDictFromKeyValues([]interface{}{1, 2, 3}, []interface{}{1.0, 2.0, 3.0})
	if err1 != nil {
		t.Error(err1)
	}

	if contains, err := s.Contains(2); err != nil {
		t.Error(err)
	} else if contains == false {
		t.Fatal("The set does not contain 2 when it should")
	}
	if contains, err := s.Contains(2.0); err != nil {
		t.Error(err)
	} else if contains == true {
		t.Fatal("The set contains 2.0 when it shouldn't")
	}
}

func TestEquals(t *testing.T) {
	s1, err1 := MakeDictFromKeyValues([]interface{}{1, 2, 3}, []interface{}{1.0, 2.0, 3.0})
	if err1 != nil {
		t.Error(err1)
	}

	s2, err2 := MakeDictFromKeyValues([]interface{}{1, 2, 3}, []interface{}{1, 2, 3})
	if err2 != nil {
		t.Error(err2)
	}

	s3, err3 := MakeDictFromKeyValues([]interface{}{1.0, 2.0, 3.0}, []interface{}{1, 2, 3})
	if err3 != nil {
		t.Error(err3)
	}

	s4, err1 := s1.Copy()
	if err1 != nil {
		t.Error(err1)
	}

	s5, err5 := MakeDictFromKeyValues([]interface{}{1, 2, 3}, []interface{}{1.0, 2.0, 3.0})
	if err5 != nil {
		t.Error(err5)
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
	} else if equals == true {
		t.Fatal("s3 and s1 should not be equal, but they are")
	}

	if equals, err := s3.Equals(s1); err != nil {
		t.Error(err)
	} else if equals == true {
		t.Fatal("s3 and s1 should not be equal, but they are")
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

	if equals, err := s1.Equals(s5); err != nil {
		t.Error(err)
	} else if equals == false {
		t.Fatal("s1 and s5 should be equal, but they are not")
	}

	if equals, err := s5.Equals(s1); err != nil {
		t.Error(err)
	} else if equals == false {
		t.Fatal("s5 and s1 should be equal, but they are not")
	}
}

func TestKeys(t *testing.T) {
	s, err1 := MakeDictFromKeyValues([]interface{}{1, 2, 3}, []interface{}{1.0, 2.0, 3.0})
	if err1 != nil {
		t.Error(err1)
	}
	if keys, err := s.Keys(); err != nil {
		t.Error(err)
	} else if ok := printChecker(keys.String(), []string{"1", "2", "3"}); !ok {
		t.Fatalf("Got %s, which was unexpected", keys.String())
	}
}

func TestValues(t *testing.T) {
	s, err1 := MakeDictFromKeyValues([]interface{}{1, 2, 3}, []interface{}{1.0, 2.0, 3.0})
	if err1 != nil {
		t.Error(err1)
	}
	// Float was passed, but it is stored as int
	// Must do external work on client side to address this
	if values, err := s.Values(); err != nil {
		t.Error(err)
	} else if ok := printChecker(values.String(), []string{"1", "2", "3"}); !ok {
		t.Fatalf("Got %s, which was unexpected", values.String())
	}
}

func TestItems(t *testing.T) {
	s, err1 := MakeDictFromKeyValues([]interface{}{1, 2, 3}, []interface{}{1.0, 2.0, 3.0})
	if err1 != nil {
		t.Error(err1)
	}
	// Float was passed, but it is stored as int
	// Must do external work on client side to address this
	if items, err := s.Items(); err != nil {
		t.Error(err)
	} else if ok := printChecker(items.String(), []string{"(1 1)", "(2 2)", "(3 3)"}); !ok {
		t.Fatalf("Got %s, which was unexpected", items.String())
	}
}
