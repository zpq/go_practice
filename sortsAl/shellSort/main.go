package main

import "fmt"

//希尔排序就是分组插入排序

func shellSort(arr []int) {
	ll, tmp := len(arr), 0
	for gap := ll / 2; gap > 0; gap /= 2 { //把数据每个n位分成n组
		for i := 0; i < gap; i++ { //开始插入排序
			for j := i + gap; j < ll; j += gap {
				if arr[j] < arr[j-gap] { //两个数相隔gap位，如果前一个大于后一个
					tmp = arr[j] //保存前一个数
					k := j - gap
					for k >= 0 && arr[k] > tmp { //同插入排序
						arr[k+gap] = arr[k]
						k = k - gap
						count++
					}
					arr[k+gap] = tmp
					count++
				}
			}
		}
	}
}

var count int

func main() {
	arr := []int{11, 10, 9, 8, 7, 6, 5, 4, 3, 3, 2, 1}
	arr = []int{27, 28, 67, 16, 25, 14, 83, 103, 2, 10, 30}
	fmt.Println(arr)
	shellSort(arr)
	fmt.Println(arr)
	fmt.Println(count)
}
