package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {

	in := 100

	c1 := incrementor(in)
	c2 := decrementor(in)

	for n := range merge(c1, c2) {
		fmt.Println(n)
	}

}

func incrementor(n int) <-chan string {
	out := make(chan string)
	go func() {
		for i := 0; i < n+1; i++ {
			time.Sleep(500 * time.Millisecond)
			fmt.Println("Foo: ", i)
		}
		close(out)
	}()
	return out
}

func decrementor(n int) <-chan string {
	out := make(chan string)
	go func() {
		for i := n; i > -1; i-- {
			time.Sleep(200 * time.Millisecond)
			out <- fmt.Sprint("Bar:", i)
		}
		close(out)
	}()
	return out
}

func merge(cs ...<-chan string) <-chan string {
	out := make(chan string)
	var wg sync.WaitGroup

	output := func(ch <-chan string) {
		for n := range ch {
			out <- n
		}
		wg.Done()
	}

	wg.Add(len(cs))
	for _, n := range cs {
		go output(n)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
