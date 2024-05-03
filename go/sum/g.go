package main

import (
	"fmt"
	"sync"
	"time"
)

func sum1(start, end int, arr []int, done chan int) <-chan int {
	res := make(chan int)
	s := 0
	for i := start; i <= end; i++ {
		s += arr[i]
	}
	go func() {
		defer close(res)
		select {
		case <-done:
			return
		case res <- s:
		}
	}()
	return res
}
func fanIn(done chan int, channels ...<-chan int) <-chan int {
	fannedInStream := make(chan int)
	var wg sync.WaitGroup
	transfer := func(c <-chan int) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case fannedInStream <- i:
			}
		}
	}
	for _, c := range channels {
		wg.Add(1)
		go transfer(c)
	}
	go func() {
		wg.Wait()
		close(fannedInStream)
	}()
	return fannedInStream
}
func main() {
	start := time.Now()
	arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	numWorkers := 5
	channels := make([]<-chan int, numWorkers)
	done := make(chan int)
	defer close(done)
	s, end := 0, len(arr)/numWorkers
	if len(arr)%numWorkers > 0 {
		end = end + len(arr)%numWorkers
	}
	for i := 0; i < numWorkers; i++ {
		if end >= len(arr) {
			end = len(arr) - 1
		}
		fmt.Println("start is ", s, "end is", end)
		channels[i] = sum1(s, end, arr, done)
		s = end + 1
		end = s + len(arr)/numWorkers
	}
	fannedInStream := fanIn(done, channels...)
	res := 0
	for i := range fannedInStream {
		res += i
	}
	fmt.Println(res)
	fmt.Println("time taken is ", time.Since(start))
}
