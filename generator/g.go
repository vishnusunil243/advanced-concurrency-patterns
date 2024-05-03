package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func repeatFunc[T any, K any](done <-chan K, fn func() T) <-chan T {
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
func take[T any, k any](done <-chan k, stream <-chan T, n int) <-chan T {
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
	isPrime := func(random int) bool {
		for i := random - 1; i > 1; i-- {
			if random%i == 0 {
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
			case random := <-randIntStream:
				if isPrime(random) {
					primes <- random
				}
			}
		}
	}()
	return primes
}
func main() {
	start := time.Now()
	done := make(chan int)
	defer close(done)
	rand := func() int { return rand.Intn(40000000) }
	randIntStream := repeatFunc(done, rand)
	// primeStream := primeFinder(done, randIntStream)
	//naive-slow-approach
	// for r := range take(done, primeStream, 10) {
	// 	fmt.Println(r)
	// }
	//fan-out
	CPUcount := runtime.NumCPU()
	primeFinderChannels := make([]<-chan int, CPUcount)
	for i := 0; i < CPUcount; i++ {
		primeFinderChannels[i] = primeFinder(done, randIntStream)
	}

	//fan-in
	fannedInStream := fanIn(done, primeFinderChannels...)
	for r := range take(done, fannedInStream, 10) {
		fmt.Println(r)
	}
	fmt.Println(time.Since(start))
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
