package main

import (
	"fmt"
)

func main() {
	var nums = []int{5, 3, 4, 2, 7, 6, 10, 3, 4, -1}

	for i := 0; i < len(nums); i++ {
		for j := 0; j < len(nums)-i-1; j++ {
			if nums[j] > nums[j+1] {
				nums[j], nums[j+1] = nums[j+1], nums[j]
			}
		}
	}

	fmt.Println(nums)

}
