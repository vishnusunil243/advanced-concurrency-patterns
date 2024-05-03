package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	done := make(chan interface{})
	defer close(done)
	cows := make(chan interface{}, 100)
	pigs := make(chan interface{}, 100)
	go func() {
		for {
			select {
			case <-done:
				return
			case cows <- "moo":
			}
		}
	}()
	go func() {
		for {
			select {
			case <-done:
				return
			case pigs <- "oink":
			}
		}
	}()
	wg.Add(2)
	go cowProcess(done, cows, &wg)
	go pigProcess(done, pigs, &wg)
	wg.Wait()
}
func cowProcess(done, cows <-chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	// for {
	// 	select {
	// 	case <-done:
	// 		return
	// 	case cow, ok := <-cows:
	// 		if !ok {
	// 			fmt.Println("cow channel closed")
	// 			return
	// 		}
	// 		fmt.Println(cow)

	// 	}
	// }
	for cow := range orDone(done, cows) {
		fmt.Println(cow)
	}
}
func pigProcess(done, pigs <-chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	// for {
	// 	select {
	// 	case <-done:
	// 		return
	// 	case pig, ok := <-pigs:
	// 		if !ok {
	// 			fmt.Println("pig channel closed")
	// 			return
	// 		}
	// 		fmt.Println(pig)

	// 	}
	// }
	for pig := range orDone(done, pigs) {
		fmt.Println(pig)
	}
}
func orDone(done, c <-chan interface{}) <-chan interface{} {
	relayStream := make(chan interface{})
	go func() {
		defer close(relayStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case relayStream <- v:
				}
			}
		}
	}()
	return relayStream
}
