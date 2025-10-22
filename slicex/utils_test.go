package slicex

import (
	"fmt"
	"testing"
)

func TestDistinct(t *testing.T) {
	s1 := []string{"a", "b", "c"}
	s2 := []string{"a", "b", "d"}
	fmt.Println(Distinct(s1, s2))
}

func TestIntersect(t *testing.T) {
	s1 := []string{"a", "b", "c"}
	s2 := []string{"a", "b", "d"}
	s3 := []string{"a", "b", "e"}
	fmt.Println(Intersect(s1, s2, s3))
}

func TestExclude(t *testing.T) {
	s1 := []string{"a", "b", "c"}
	s2 := []string{"a", "b", "d"}
	fmt.Println(Exclude(s1, s2))
	fmt.Println(Exclude(s2, s1))
}
