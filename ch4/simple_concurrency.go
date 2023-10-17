package main

import (
	"fmt"
	"sync"
)

func main() {
	wgWorkers := sync.WaitGroup{}
	wgWorkers.Add(1)

	intCh := make(chan int)

	go receive(intCh, &wgWorkers)
	go send(intCh)

	wgWorkers.Wait()

	println("All done. Going to rest after a had day at work")
}

func send(ch chan int) {
	for i := 0; i < 64; i++ {
		ch <- i
	}
	close(ch)
}

func receive(ch chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for value := range ch {
		fmt.Printf("Read %v\n", value)
	}
}
