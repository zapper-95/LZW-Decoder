package mergesort

func MergeSort(nums []int) []int {
	// base case
	if len(nums) <= 1 {
		return nums
	}

	nums1 := MergeSort(nums[:len(nums)/2])
	nums2 := MergeSort(nums[len(nums)/2:])

	return Merge(nums1, nums2)
}

func Merge(nums1 []int, nums2 []int) []int {
	nums3 := make([]int, 0)
	i, j := 0, 0
	for {
		if i == len(nums1) || j == len(nums2) {
			break
		}

		if nums1[i] < nums2[j] {
			nums3 = append(nums3, nums1[i])
			i++
		} else {
			nums3 = append(nums3, nums2[j])
			j++
		}

	}

	if i == len(nums1) {
		nums3 = append(nums3, nums2[j:]...)
	} else {
		nums3 = append(nums3, nums1[i:]...)
	}

	return nums3

}
