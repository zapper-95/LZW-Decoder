package mergesort

import (
	"sort"
	"testing"
)

func TestMergeSortNormal(t *testing.T) {
	nums := []int{9, 8, 7, 6, 5, 4, 3, 2, 1}
	sorted := MergeSort(nums)

	if !sort.IntsAreSorted(sorted) {
		t.Errorf("got %v", sorted)
	}

}

func TestMergeSortStrange(t *testing.T) {
	nums := []int{9, -8, 7, -6, 5, 4, 3, 2, 100000000000000000}
	sorted := MergeSort(nums)

	if !sort.IntsAreSorted(sorted) {
		t.Errorf("got %v", sorted)
	}

}
