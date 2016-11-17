package main

import (
	"fmt"
)

func mergeSort(arr []int) {
	splitList(arr, 0, len(arr)-1, []int{})
}

func splitList(arr []int, l, r int, res []int) {
	if l < r {
		m := (l + r) / 2
		splitList(arr, l, m, res)
		splitList(arr, m+1, r, res)
		mergeTwoList(arr, res, l, m, r)
	}
}

//合并两个有序数组为一个有序数组
func mergeTwoList(arr []int, res []int, l, m, r int) {
	a := arr[l : m+1]
	b := arr[m+1 : r+1]
	i, j, m, n := 0, 0, len(a), len(b)
	for i < m || j < n {
		if i == m { // a没有值可以比较了,直接合并b数组剩下的值
			res = append(res, b[j:n]...)
			// j = n
			break
		}
		if j == n { // b没有值可以比较了,直接合并a数组剩下的值
			res = append(res, a[i:m]...)
			// i = m
			break
		}
		if a[i] > b[j] { //a >= b,把b中的值合并进结果数组，b索引加一
			res = append(res, b[j])
			j++
		} else { //反之,把a中的值合并进结果数组，a索引加一
			res = append(res, a[i])
			i++
		}
	}
	for i := 0; i < len(res); i++ { //修改原数组的值（待排序的数组，res是零时存放的）
		arr[l+i] = res[i]
	}
	return
}

var count int

func main() {
	// a, b := []int{1, 3, 5, 6, 7, 8}, []int{2, 4, 6}
	arr := []int{11, 10, 9, 8, 7, 6, 5, 4, 3, 3, 2, 1}
	arr = []int{27, 28, 67, 16, 25, 14, 83, 103, 2, 10, 30}
	fmt.Println(arr)
	mergeSort(arr)
	fmt.Println(arr)
}
