package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func generator[T any, K any](done <-chan K, fn func() T) <-chan T {
	stream := make(chan T)
	go func() {
		defer close(stream)
		for {
			select {
			case <-done:
				return
			case stream <- fn():
			}
		}
	}()
	return stream
}
func take[T any, K any](done <-chan K, stream <-chan T, n int) <-chan T {
	taken := make(chan T)
	go func() {
		defer close(taken)
		for i := 0; i < n; i++ {
			select {
			case <-done:
				return
			case taken <- <-stream:
			}
		}
	}()
	return taken
}
func primeFinder(done <-chan int, randIntStream <-chan int) <-chan int {
	isPrime := func(rand int) bool {
		for i := rand - 1; i > 1; i-- {
			if rand%i == 0 {
				return false
			}
		}
		return true
	}
	primes := make(chan int)
	go func() {
		defer close(primes)
		for {
			select {
			case <-done:
				return
			case rand := <-randIntStream:
				if isPrime(rand) {
					primes <- rand
				}
			}
		}
	}()
	return primes
}
func fanIn[T any](done <-chan int, channels ...<-chan T) <-chan T {
	var wg sync.WaitGroup
	fannedInStream := make(chan T)
	transfer := func(c <-chan T) {
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
	fmt.Println()
	done := make(chan int)
	randomStream := func() int { return rand.Intn(500000) }
	randIntStream := generator(done, randomStream)
	// primeStream := primeFinder(done, randIntStream)
	// for r := range take(done, primeStream, 10) {
	// 	fmt.Println(r)
	// }
	cpu := runtime.NumCPU()
	primes := make([]<-chan int, cpu)
	for i := 0; i < cpu; i++ {
		primes[i] = primeFinder(done, randIntStream)
	}
	fannedInStream := fanIn(done, primes...)
	for r := range take(done, fannedInStream, 10) {
		fmt.Println(r)
	}
	fmt.Println(time.Since(start))
}
