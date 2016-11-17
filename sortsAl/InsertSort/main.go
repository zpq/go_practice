/*
 int i, j;
    for (i = 1; i < n; i++)
        if (a[i] < a[i - 1])
        {
            int temp = a[i];
            for (j = i - 1; j >= 0 && a[j] > temp; j--)
                a[j + 1] = a[j];
            a[j + 1] = temp;
        }

1,2,3,4,5,2
1,2,3,4,2,5
1,2,3,2,4,5
1,2,2,3,4,5
*/
package main

import "fmt"

func insertSort(arr []int) {
	ll, tmp, j := len(arr), 0, 0
	for i := 1; i < ll; i++ { // 正常遍历
		if arr[i] < arr[i-1] { //后一个数小于前一个数时
			tmp = arr[i] //保存较小的数
			//开始从当前位置向前寻找，直到找到比tmp小或等于的数，把tmp插到后面，比他大的数都向后移动一个位置
			for j = i - 1; j >= 0 && arr[j] > tmp; j-- {
				arr[j+1] = arr[j]
				count++
			}
			arr[j+1] = tmp
		}
		count++
	}
}

var count int

func main() {
	arr := []int{9, 8, 7, 6, 5, 4, 3, 3, 2, 1}
	arr = []int{27, 28, 67, 16, 25, 14, 83, 103, 2, 10, 30}
	fmt.Println(arr)
	insertSort(arr)
	fmt.Println(arr)
	fmt.Println(count)
}
