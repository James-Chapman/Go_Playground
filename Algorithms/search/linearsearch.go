package search

// LinearSearch
func LinearSearch(item int, arr []int) int {
	for i, v := range arr {
		if v == item {
			return i
		}
	}
	return -1
}