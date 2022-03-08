package main

import (
	"fmt"
	"strconv"
	"sync"
)

func producer(thread int, ch chan string, wg *sync.WaitGroup) {
	for i := 0; i < 3; i++ {
		ch <- strconv.Itoa(thread) + "----" + strconv.Itoa(i)
	}
	wg.Done()
}

func consumer(ch chan string, wg *sync.WaitGroup) {
	for i := range ch {
		fmt.Println(i)
	}
	wg.Done()
}
func main() {
	ch := make(chan string, 10)
	wgPr := new(sync.WaitGroup)
	wgCo := new(sync.WaitGroup)
	for i := 0; i < 3; i++ {
		wgPr.Add(1)
		go producer(i, ch, wgPr)
	}

	for i := 0; i < 2; i++ {
		wgCo.Add(1)
		go consumer(ch, wgCo)
	}

	wgPr.Wait()
	close(ch)
	wgCo.Wait()
}
