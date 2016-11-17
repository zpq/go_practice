package main

import (
	"fmt"
)

func bubble(nums []int) {
	ll := len(nums)
	for i := 0; i < ll; i++ {
		for j := i; j < ll; j++ {
			count++
			if nums[i] > nums[j] {
				nums[i] = nums[i] + nums[j]
				nums[j] = nums[i] - nums[j]
				nums[i] = nums[i] - nums[j]
			}
		}
	}
}

var count int

func main() {
	arr := []int{9, 8, 7, 6, 5, 4, 3, 2, 1}
	arr = []int{27, 28, 67, 16, 25, 14, 83, 103, 2, 10, 30}
	fmt.Println(arr)
	bubble(arr)
	fmt.Println(arr)
	fmt.Println(count)
}
