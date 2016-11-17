/*
3,16,2,8,1  //min = 1 ; 1和3交换
1,16,2,8,3  //min = 2; 2和16交换
1,2,16,8,3  //min = 3; 3和16交换
1,2,3,8,16  //排序结束
*/

package main

import "fmt"

func selectSort(arr []int) {
	ll, minIndex := len(arr), 0
	for i := 0; i < ll; i++ {
		minIndex = i                  //把第i位数字当作最小数
		for j := i + 1; j < ll; j++ { //往后比较
			if arr[j] < arr[minIndex] { //找出第i位后最小的数,最小数与第i位可能需要交换
				minIndex = j
			}
			count++
		}
		if arr[i] != arr[minIndex] { //当需交换的两个数不相等时才交换
			arr[i] = arr[i] + arr[minIndex]
			arr[minIndex] = arr[i] - arr[minIndex]
			arr[i] = arr[i] - arr[minIndex]
		}
	}
}

var count int

func main() {
	arr := []int{11, 10, 9, 8, 7, 6, 5, 4, 3, 3, 2, 1}
	arr = []int{27, 28, 67, 16, 25, 14, 83, 103, 2, 10, 30}
	// arr = []int{3, 16, 2, 8, 1}
	fmt.Println(arr)
	selectSort(arr)
	fmt.Println(arr)
	fmt.Println(count)
}
