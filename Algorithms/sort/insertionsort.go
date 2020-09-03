package sort

// InsertionSort is InsertionSort without the need for a temporary variable
func InsertionSort(arr []int) {
	for i, _ := range arr {
		j := i
		for j > 0 {
			// If previous element is bigger than current
			if arr[j-1] > arr[j] {
				// Swap previous and current
				arr[j-1], arr[j] = arr[j], arr[j-1]
			}
			j--
		}
	}
}

// InsertionSort2 is InsertionSort with a temporary variable
func InsertionSort2(arr []int) {
	for i, _ := range arr {
		j := i
		for j > 0 {
			// If previous element is bigger than current
			if arr[j-1] > arr[j] {
				x := arr[j]
				// Swap previous and current
				arr[j] = arr[j-1]
				arr[j-1] = x
			}
			j--
		}
	}
}