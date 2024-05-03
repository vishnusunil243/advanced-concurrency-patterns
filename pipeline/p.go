package main

import "fmt"

func sliceToData(data []int) <-chan int {
	ch := make(chan int)
	go func() {
		for _, n := range data {
			ch <- n
		}
		close(ch)
	}()
	return ch
}
func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for i := range in {
			out <- i * i
		}
		close(out)
	}()
	return out
}
func main() {
	data := []int{1, 2, 3, 4, 5}
	//stage 1
	datachannel := sliceToData(data)
	//stage 2
	finalChannel := sq(datachannel)
	//
	for r := range finalChannel {
		fmt.Println(r)
	}

}
