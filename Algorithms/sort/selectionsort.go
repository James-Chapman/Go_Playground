package sort

// SelectionSort
func SelectionSort(arr []int) {
	var n = len(arr)
    for i, _ := range arr {
        m := i
        for j := i; j < n; j++ {
            if arr[j] < arr[m] {
                m = j
            }
        }
        arr[i], arr[m] = arr[m], arr[i]
    }
}