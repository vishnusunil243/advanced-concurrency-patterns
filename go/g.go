package main

import (
	"fmt"
	"sync"
	"time"
)

func sum(start, end int, res chan int, arr []int, wg *sync.WaitGroup) {
	defer wg.Done()
	count := 0
	for i := start; i <= end; i++ {
		count += arr[i]
	}
	res <- count
}
func main() {
	arr := []int{1, 2, 3, 4, 5, 6, 7, 9}
	start := time.Now()
	numWorkers := 4
	res := make(chan int, numWorkers)
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		start := i * len(arr) / numWorkers
		end := start + len(arr)/numWorkers - 1
		if end >= len(arr) {
			end = len(arr) - 1
		}
		go sum(start, end, res, arr, &wg)
	}
	wg.Wait()
	close(res)
	c := 0
	for i := range res {
		c += i
	}
	fmt.Println(c)
	// fmt.Println(sum1(arr))
	fmt.Println("time taken is ", time.Since(start))

}
