package main

import (
	"fmt"
)

func findMid(l, r int, nums []int) int {
	i, j := l, r
	x := nums[i]
	for i < j {
		//先从右边开始找， 因为基准值选取的是最左边的值
		for nums[j] >= x && i < j { //当右边的值大于等于基准值时，右索引往左走
			j--
			count++
		}
		//一旦跳出循环，说明找到一个数比基准值小，应该移动该数字的位置， 并且左索引往右移动一个位置
		if i < j {
			nums[i] = nums[j]
			i++
			count++
		}

		//然后从左边开始找
		for nums[i] < x && i < j { //当左边的值小于基准值时， 左索引往右走
			i++
			count++
		}
		//一旦跳出循环，说明找到一个数大于等于基准值，应该移动该数字的位置，并且右索引往左移动一个位置
		if i < j {
			nums[j] = nums[i]
			j--
			count++
		}
	}
	nums[i] = x //把基准值放到中间位置
	return i    //返回中间位置的索引值
}

func quickSort(l, r int, arr []int) {
	if l < r {
		m := findMid(l, r, arr)
		quickSort(l, m-1, arr)
		quickSort(m+1, r, arr)
	}
}

var count int

func main() {
	arr := []int{9, 8, 7, 6, 5, 4, 3, 3, 2, 1}
	arr = []int{27, 28, 67, 16, 25, 14, 83, 103, 2, 10, 30}
	fmt.Println(arr)
	quickSort(0, len(arr)-1, arr)
	fmt.Println(arr)
	fmt.Println(count)
}
