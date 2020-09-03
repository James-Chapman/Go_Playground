package search

// BinarySearch
func BinarySearch(item int, arr []int) int {
	arrLen := len(arr)
	n := 0
	if arrLen == 1 {
		n = arrLen - 1
	} else {
		n = (arrLen / 2)
	}

	if arr[n] == item {
		return n
	} else if arr[n] > item {
		n = BinarySearch(item, arr[:n])
	} else if arr[n] < item {
		n += BinarySearch(item, arr[n:])
	}

	return n
}
