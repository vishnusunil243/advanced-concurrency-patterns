package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	generator := func(data string, stream chan interface{}) {
		for {
			select {
			case <-ctx.Done():
				return
			case stream <- data:
			}
		}
	}
	infiniteApple := make(chan interface{})
	go generator("apple", infiniteApple)
	infiniteOrange := make(chan interface{})
	go generator("orange", infiniteOrange)
	wg.Add(2)
	go func1(ctx, &wg, infiniteApple)
	go func2(&wg, infiniteOrange, ctx)
	wg.Wait()

}
func func1(ctx context.Context, Pwg *sync.WaitGroup, stream <-chan interface{}) {
	defer Pwg.Done()
	var wg sync.WaitGroup
	doWork := func(ctx context.Context) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-stream:
				if !ok {
					fmt.Println("channel closed")
					return
				}
				fmt.Println(v)
			}
		}
	}
	newCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go doWork(newCtx)
	}
	wg.Wait()
}
func func2(wg *sync.WaitGroup, stream <-chan interface{}, ctx context.Context) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case v, ok := <-stream:
			if !ok {
				fmt.Println("channel closed")
				return
			}
			fmt.Println(v)
		}
	}
}
