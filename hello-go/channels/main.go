package main

import (
	"fmt"
	"time"
)

func sayHello(source string) {
	fmt.Printf("Hello World from %s!\n", source)
}

// Note that the following method receives a channel through which it can
// communicate it's result once it becomes available.

func getValueAsync(result chan int) {
	// Simulate long-running operation (e.g. read something from disk)
	time.Sleep(10 * time.Millisecond)

	// Send result to calling goroutine through channel
	result <- 42
}

// Note that the following method does not return a value. The channel is
// just used to indicate the completion of the asynchronous work.

func doSomethingComplex(done chan bool) {
	// Simulate long-running operation
	time.Sleep(10 * time.Millisecond)

	done <- true
}

func main() {
	// Call sayHello directly and using the *go* keyword on a separate goroutine.
	sayHello("direct call")
	go sayHello("goroutine")

	// Call method on a different goroutine and give it a channel through which
	// it can send back the result.
	result := make(chan int)
	go getValueAsync(result)
	// Wait until result is available and print it.
	fmt.Println(<-result)

	// Do something asynchronously and wait for it to finish
	done := make(chan bool)
	go doSomethingComplex(done)
	// Wait until a message is available in the channel
	<-done
	fmt.Println("Complex operation is done")

	// Note the select statement here. You can use it to wait on
	// multiple channels. In this case, we use a channel from a
	// timer to implement a timeout functionality.
	go getValueAsync(result)
	select {
	case m := <-result:
		fmt.Println(m)
	case <-time.After(5 * time.Millisecond):
		fmt.Println("timed out")
	}

	// Let us print a status message for a certain amount of time.
	ticker := time.NewTicker(100 * time.Millisecond)
	// Note the anonymous function here.
	go func() {
		// Note that Go's range operator supports looping over
		// values received through a channel.
		for range ticker.C {
			fmt.Println("Tick")
		}
	}()
	// Wait for some time and then stop timer.
	<-time.After(500 * time.Millisecond)
	ticker.Stop()
}
