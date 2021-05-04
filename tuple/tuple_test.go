package tuple

import (
	"fmt"
	"testing"
)

func printChecker(str string, values []string) bool {
	output := false
	potentialMatch := fmt.Sprint(values)
	if potentialMatch[1:len(potentialMatch)-1] == str[1:len(str)-1] {
		output = true
	}
	return output
}

func TestMakeEmptyTuple(t *testing.T) {
	s, err := MakeTuple()
	if err != nil {
		t.Error(err)
	}

	values := []string{}
	if ok := printChecker(s.String(), values); !ok {
		t.Fatalf("Got %s, which was unexpected", s.String())
	}
}

type tupleTester struct {
	a []interface{}
}

func (test *tupleTester) Iterate() <-chan interface{} {
	c := make(chan interface{})
	go func() {
		for _, v := range test.a {
			c <- v
		}
		close(c)
	}()
	return c
}

func TestMakePopulatedTuple(t *testing.T) {
	type tester struct {
		a int
		b int
	}
	var test tester
	test.a = 1
	test.b = 2

	values := []string{"1", "2.2", "hello", "{1 2}"}

	input := make([]interface{}, 4)
	input[0] = 1
	input[1] = 2.2
	input[2] = "hello"
	input[3] = test
	tuple := tupleTester{input}

	s, err := MakeTuple(&tuple)
	if err != nil {
		t.Error(err)
	}
	if ok := printChecker(s.String(), values); !ok {
		t.Fatalf("Got %s, which was unexpected", s.String())
	}

	s, err = MakeTupleFromValues(1, 2.2, "hello", test)
	if err != nil {
		t.Error(err)
	}
	if ok := printChecker(s.String(), values); !ok {
		t.Fatalf("Got %s, which was unexpected", s.String())
	}
}

func TestContains(t *testing.T) {
	s, err := MakeTupleFromValues(1, 2.2, "hello")
	if err != nil {
		t.Error(err)
	}

	if contains, err := s.Contains(2.2); err != nil {
		t.Error(err)
	} else if contains == false {
		t.Fatal("The set does not contain 2.2 when it should")
	}

	if contains, err := s.Contains("hello"); err != nil {
		t.Error(err)
	} else if contains == false {
		t.Fatal("The set does not contain \"hello\" when it should")
	}

	if contains, err := s.Contains(17); err != nil {
		t.Error(err)
	} else if contains == true {
		t.Fatal("The set does contain 17 when it should not")
	}
}

func TestEquals(t *testing.T) {
	s1, err1 := MakeTupleFromValues(1, 2.2, "hello")
	if err1 != nil {
		t.Error(err1)
	}

	s2, err2 := s1.Copy()
	if err2 != nil {
		t.Error(err2)
	}

	s3, err3 := MakeTupleFromValues(1, "hello", 2.2)
	if err3 != nil {
		t.Error(err3)
	}

	s4, err4 := MakeTupleFromValues(123, 456)
	if err4 != nil {
		t.Error(err4)
	}

	if equals, err := s1.Equals(s1); err != nil {
		t.Error(err)
	} else if equals == false {
		t.Fatal("s1 should be equal to itself, but it is not")
	}

	if equals, err := s1.Equals(s2); err != nil {
		t.Error(err)
	} else if equals == false {
		t.Fatal("s1 should be equal to s2, but it is not")
	}

	if equals, err := s1.Equals(s3); err != nil {
		t.Error(err)
	} else if equals == true {
		t.Fatal("s1 should not be equal to s3, but it is")
	}

	if equals, err := s1.Equals(s4); err != nil {
		t.Error(err)
	} else if equals == true {
		t.Fatal("s1 should not be equal to s3, but it is")
	}
}

func TestConcatenate(t *testing.T) {
	s1, err1 := MakeTupleFromValues(1, 2.2, "hello")
	if err1 != nil {
		t.Error(err1)
	}

	s2, err2 := MakeTupleFromValues(1, 2.2, "hello")
	if err2 != nil {
		t.Error(err2)
	}

	s3, err3 := MakeTupleFromValues()
	if err3 != nil {
		t.Error(err3)
	}

	values := []string{"1", "2.2", "hello", "1", "2.2", "hello"}
	if c, err := s1.Concatenate(s1); err != nil {
		t.Error(err)
	} else if ok := printChecker(c.String(), values); !ok {
		t.Fatalf("Got %s, which was unexpected", c.String())
	}

	values = []string{"1", "2.2", "hello", "1", "2.2", "hello"}
	if c, err := s1.Concatenate(s2); err != nil {
		t.Error(err)
	} else if ok := printChecker(c.String(), values); !ok {
		t.Fatalf("Got %s, which was unexpected", c.String())
	}

	values = []string{"1", "2.2", "hello"}
	if c, err := s1.Concatenate(s3); err != nil {
		t.Error(err)
	} else if ok := printChecker(c.String(), values); !ok {
		t.Fatalf("Got %s, which was unexpected", c.String())
	}
}

func TestMultiply(t *testing.T) {
	s1, err1 := MakeTupleFromValues(1, 2.2, "hello")
	if err1 != nil {
		t.Error(err1)
	}

	s2, err2 := MakeTupleFromValues()
	if err2 != nil {
		t.Error(err2)
	}

	values := []string{}
	if c, err := s1.Multiply(-1); err != nil {
		t.Error(err)
	} else if ok := printChecker(c.String(), values); !ok {
		t.Fatalf("Got %s, which was unexpected", c.String())
	}

	values = []string{}
	if c, err := s1.Multiply(0); err != nil {
		t.Error(err)
	} else if ok := printChecker(c.String(), values); !ok {
		t.Fatalf("Got %s, which was unexpected", c.String())
	}

	values = []string{"1", "2.2", "hello"}
	if c, err := s1.Multiply(1); err != nil {
		t.Error(err)
	} else if ok := printChecker(c.String(), values); !ok {
		t.Fatalf("Got %s, which was unexpected", c.String())
	}

	values = []string{"1", "2.2", "hello", "1", "2.2", "hello"}
	if c, err := s1.Multiply(2); err != nil {
		t.Error(err)
	} else if ok := printChecker(c.String(), values); !ok {
		t.Fatalf("Got %s, which was unexpected", c.String())
	}

	values = []string{}
	if c, err := s2.Multiply(2); err != nil {
		t.Error(err)
	} else if ok := printChecker(c.String(), values); !ok {
		t.Fatalf("Got %s, which was unexpected", c.String())
	}
}

func TestGet(t *testing.T) {
	s1, err1 := MakeTupleFromValues(1, 2.2, "hello")
	if err1 != nil {
		t.Error(err1)
	}

	s2, err2 := MakeTupleFromValues()
	if err2 != nil {
		t.Error(err2)
	}

	if _, err := s1.Get(-1); err == nil {
		t.Fatalf("Was expecting index to be out of range")
	}

	if val, err := s1.Get(0); err != nil {
		t.Error(err)
	} else if val != 1 {
		t.Fatalf("Got %s, was expecting 1", val)
	}

	if val, err := s1.Get(1); err != nil {
		t.Error(err)
	} else if val != 2.2 {
		t.Fatalf("Got %s, was expecting 2.2", val)
	}

	if val, err := s1.Get(2); err != nil {
		t.Error(err)
	} else if val != "hello" {
		t.Fatalf("Got %s, was expecting hello", val)
	}

	if _, err := s1.Get(3); err == nil {
		t.Fatalf("Was expecting index to be out of range")
	}

	if _, err := s2.Get(0); err == nil {
		t.Fatalf("Was expecting index to be out of range")
	}
}

func TestRange(t *testing.T) {
	s1, err1 := MakeTupleFromValues(1, 2.2, "hello")
	if err1 != nil {
		t.Error(err1)
	}

	s2, err2 := MakeTupleFromValues()
	if err2 != nil {
		t.Error(err2)
	}

	if s, err := s1.Range(0, 3); err != nil {
		t.Error(err)
	} else if ok := printChecker(s.String(), []string{"1", "2.2", "hello"}); !ok {
		t.Fatalf("Got %s, which was unexpected", s.String())
	}

	if s, err := s1.Range(0, 1); err != nil {
		t.Error(err)
	} else if ok := printChecker(s.String(), []string{"1"}); !ok {
		t.Fatalf("Got %s, which was unexpected", s.String())
	}

	if s, err := s1.Range(1, 3); err != nil {
		t.Error(err)
	} else if ok := printChecker(s.String(), []string{"2.2", "hello"}); !ok {
		t.Fatalf("Got %s, which was unexpected", s.String())
	}

	if s, err := s1.Range(-1, 4); err != nil {
		t.Error(err)
	} else if ok := printChecker(s.String(), []string{"1", "2.2", "hello"}); !ok {
		t.Fatalf("Got %s, which was unexpected", s.String())
	}

	if s, err := s1.Range(2, 0); err != nil {
		t.Error(err)
	} else if ok := printChecker(s.String(), []string{}); !ok {
		t.Fatalf("Got %s, which was unexpected", s.String())
	}

	if s, err := s2.Range(0, 1); err != nil {
		t.Error(err)
	} else if ok := printChecker(s.String(), []string{}); !ok {
		t.Fatalf("Got %s, which was unexpected", s.String())
	}
}

func TestIndex(t *testing.T) {
	s1, err1 := MakeTupleFromValues(1, 2.2, "hello", 2.2)
	if err1 != nil {
		t.Error(err1)
	}

	s2, err2 := MakeTupleFromValues()
	if err2 != nil {
		t.Error(err2)
	}

	if i, err := s1.Index(1); err != nil {
		t.Error(err)
	} else if i != 0 {
		t.Fatalf("Got %d, was expecting 0", i)
	}

	if i, err := s1.Index(2.2); err != nil {
		t.Error(err)
	} else if i != 1 {
		t.Fatalf("Got %d, was expecting 1", i)
	}

	if i, err := s1.Index("hello"); err != nil {
		t.Error(err)
	} else if i != 2 {
		t.Fatalf("Got %d, was expecting 2", i)
	}

	if i, err := s1.Index(42); err != nil {
		t.Error(err)
	} else if i != -1 {
		t.Fatalf("Got %d, was expecting -1", i)
	}

	if i, err := s2.Index("hello"); err != nil {
		t.Error(err)
	} else if i != -1 {
		t.Fatalf("Got %d, was expecting -1", i)
	}
}

func TestCount(t *testing.T) {
	s1, err1 := MakeTupleFromValues(1, 2.2, "hello", 2.2)
	if err1 != nil {
		t.Error(err1)
	}

	s2, err2 := MakeTupleFromValues()
	if err2 != nil {
		t.Error(err2)
	}

	if i, err := s1.Count(1); err != nil {
		t.Error(err)
	} else if i != 1 {
		t.Fatalf("Got %d, was expecting 1", i)
	}

	if i, err := s1.Count(2.2); err != nil {
		t.Error(err)
	} else if i != 2 {
		t.Fatalf("Got %d, was expecting 2", i)
	}

	if i, err := s1.Count("hello"); err != nil {
		t.Error(err)
	} else if i != 1 {
		t.Fatalf("Got %d, was expecting 1", i)
	}

	if i, err := s1.Count(42); err != nil {
		t.Error(err)
	} else if i != 0 {
		t.Fatalf("Got %d, was expecting 0", i)
	}

	if i, err := s2.Count("hello"); err != nil {
		t.Error(err)
	} else if i != 0 {
		t.Fatalf("Got %d, was expecting 0", i)
	}
}
