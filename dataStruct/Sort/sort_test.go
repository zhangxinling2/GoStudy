package sort

import (
	"fmt"
	"testing"
)

//冒泡排序
func BubbleSort(arr *[]int) {
	for i := 0; i < len(*arr); i++ {
		flag := false
		for j := len(*arr) - 1; j > i; j-- {
			if (*arr)[j] < (*arr)[j-1] {
				flag = true
				// tmp := (*arr)[j]
				// (*arr)[j] = (*arr)[j-1]
				// (*arr)[j-1] = tmp
				(*arr)[j], (*arr)[j-1] = (*arr)[j-1], (*arr)[j]
			}
		}
		if flag == false {
			break
		}
	}
}

//选择排序
func SelectSort(arr []int) {
	for i := 0; i < len(arr); i++ {
		min := arr[i]
		minIndex := i
		for j := i; j < len(arr); j++ {
			if arr[j] < min {
				minIndex = j
				min = arr[j]
			}
		}
		arr[i], arr[minIndex] = arr[minIndex], arr[i]
	}
}

//插入排序，将数据分为两部分，一部分排好序，之后从乱序中逐个向前比较进入排好序的部分
func InsertSort(arr *[]int) {
	for i := 1; i < len(*arr); i++ {
		for j := i; j > 0; j-- {
			if (*arr)[j] < (*arr)[j-1] {
				(*arr)[j], (*arr)[j-1] = (*arr)[j-1], (*arr)[j]
			} else {
				break
			}
		}
	}
}

//折半插入排序
func BinaryInsertSort(arr []int) {
	for i := 1; i < len(arr); i++ {
		left, right := 0, i-1
		for left <= right {
			half := (left + right) / 2
			if arr[half] > arr[i] {
				right = half - 1
			} else {
				left = half + 1
			}
		}
		for i > left {
			arr[i], arr[i-1] = arr[i-1], arr[i]
			i--
		}
	}
}

//希尔排序
func ShellSort(arr []int) {
	l := len(arr)
	for k := l / 2; k > 0; k = k / 2 {
		for i := k; i < l; i++ {
			j := i
			tmp := arr[j]
			for j >= k && arr[j-k] > tmp {
				arr[j] = arr[j-k]
				j = j - k
			}
			arr[j] = tmp
		}
	}
}

//归并排序
func MergeSort(arr []int) []int {
	if len(arr) < 2 {
		return arr
	}
	mid := len(arr) / 2
	return Merge(MergeSort(arr[:mid]), MergeSort(arr[mid:]))
}
func Merge(left []int, right []int) []int {
	i, j := 0, 0
	res := []int{}
	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			res = append(res, left[i])
			i++
		} else {
			res = append(res, right[j])
			j++
		}
	}
	if i < len(left) {
		res = append(res, left[i:]...)
	}
	if j < len(right) {
		res = append(res, right[j:]...)
	}
	return res
}

//快速排序
func QuickSort(arr []int) {
	if len(arr) < 2 {
		return
	}
	pivot := arr[0]
	i, j := 0, len(arr)-1
	for i < j {
		for i < j && arr[j] >= pivot {
			j--
		}
		arr[i] = arr[j]
		for i < j && arr[i] <= pivot {
			i++
		}
		arr[j] = arr[i]
	}
	arr[i] = pivot
	QuickSort(arr[:i])
	QuickSort(arr[i+1:])
}

//堆排序
func HeapSort(arr []int) {
	buildHeap(arr)
	for i := len(arr) - 1; i > 0; i-- {
		arr[0], arr[i] = arr[i], arr[0]
		heapAdjust(arr, 0, i)
	}
}

func buildHeap(arr []int) {
	for i := len(arr)/2 - 1; i >= 0; i-- {
		heapAdjust(arr, i, len(arr))
	}
}

func heapAdjust(arr []int, i, len int) {
	tmp := arr[i]
	for k := 2*i + 1; k < len; k = 2*k + 1 {
		if k+1 < len && arr[k] < arr[k+1] {
			k++
		}
		if arr[k] > tmp {
			arr[i] = arr[k]
			i = k
		} else {
			break
		}
	}
	arr[i] = tmp
}
func TestSort(t *testing.T) {
	arr := []int{1, 6, 8, 3, 11, 56, 33, 23, 7, 9, 10, 3}
	//BubbleSort(&arr)
	//SelectSort(arr)
	//InsertSort(&arr)
	//BinaryInsertSort(arr)
	//ShellSort(arr)
	//res := MergeSort(arr)
	//fmt.Println(res)
	//QuickSort(arr)
	HeapSort(arr)
	fmt.Println(arr)
}
