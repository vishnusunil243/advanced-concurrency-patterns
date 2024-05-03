package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	workers := 5
	total := 100
	ch := make(chan int, workers)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		start := i*(total/workers) + 1
		end := start + (total / workers) - 1
		fmt.Println("start is ", start, "end is ", end)
		go add(start, end, ch, &wg)
	}
	wg.Wait()
	close(ch)
	s := 0
	for i := range ch {
		s += i
	}
	fmt.Println("sum is ", s)

}
func add(start, end int, ch chan int, wg *sync.WaitGroup) {
	sum := 0
	for i := start; i <= end; i++ {
		sum += i
	}
	ch <- sum
	wg.Done()
}
