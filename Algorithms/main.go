package main

import (
	"fmt"
	"time"
	"math/rand"
)

func generateSlice(size int) []int {
    slice := make([]int, size, size)
    rand.Seed(time.Now().UnixNano())
    for i := 0; i < size; i++ {
        slice[i] = rand.Intn(9999999) - rand.Intn(9999999)
    }
    return slice
}

func main() {

	fmt.Println(generateSlice(500))

}
