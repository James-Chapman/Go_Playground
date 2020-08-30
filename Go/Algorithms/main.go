package main

import (
	"fmt"
	"Algorithms/algorithms"
)

func main() {

	var s1 = []int{56, 84, 23, 76, 22, 867, 23, 99, 9, 21, 61, 38, 37, 30, 12, 54, 28}
	fmt.Println(s1)
	algorithms.InsertionSort(s1)
	fmt.Println(s1)

	var s2 = []int{56, 84, 23, 76, 22, 867, 23, 99, 9, 21, 61, 38, 37, 30, 12, 54, 28}
	fmt.Println(s2)
	algorithms.InsertionSort2(s2)
	fmt.Println(s2)

	var s3 = []int{56, 84, 23, 76, 22, 867, 23, 99, 9, 21, 61, 38, 37, 30, 12, 54, 28}
	fmt.Println(s3)
	s3 = algorithms.MergeSort(s3)
	fmt.Println(s3)

}
