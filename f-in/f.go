package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func generator[T any](done <-chan int, fn func() T) <-chan T {
	randStream := make(chan T)
	go func() {
		defer close(randStream)
		for {
			select {
			case <-done:
				return
			case randStream <- fn():
			}
		}
	}()
	return randStream
}
func take[T any](done <-chan int, stream <-chan T, n int) <-chan T {
	s := make(chan T)
	go func() {
		defer close(s)
		for i := 0; i < n; i++ {
			select {
			case <-done:
				return
			case s <- <-stream:
			}
		}
	}()
	return s
}
func primeFinder(done chan int, stream <-chan int) <-chan int {
	primes := make(chan int)
	isPrime := func(rand int) bool {
		for i := rand - 1; i > 1; i-- {
			if rand%i == 0 {
				return false
			}
		}
		return true
	}
	go func() {
		defer close(primes)
		for {
			select {
			case <-done:
				return
			case rand := <-stream:
				if isPrime(rand) {
					primes <- rand
				}
			}
		}
	}()
	return primes
}
func fanIn(done <-chan int, channels ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	fannedInStream := make(chan int)
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
	done := make(chan int)
	randFunc := func() int { return rand.Intn(40000000) }
	randStream := generator(done, randFunc)
	// primes := primeFinder(done, randStream)
	// for r := range take(done, primes, 10) {
	// 	fmt.Println(r)
	// }
	cpu := runtime.NumCPU()
	primes := make([]<-chan int, cpu)
	for i := 0; i < cpu; i++ {
		primes[i] = primeFinder(done, randStream)
	}
	fannedIn := fanIn(done, primes...)
	for r := range take(done, fannedIn, 10) {
		fmt.Println(r)
	}
	fmt.Println(time.Since(start))
}
