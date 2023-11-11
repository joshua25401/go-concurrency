package syncpackage

import (
	"fmt"
	"sync"
	"time"
)

/*
	WaitGroup
		- Is a way to wait for a set of concurrent operations to complete
		when :
			- We don't care about the result of the concurrent operation
			- or we have other means of collecting the result
		if :
			- either of the above condition is not true
			- use : channel and select statement instead of WaitGroup
*/

func ExampleWaitGroup() {
	var wg sync.WaitGroup // define a empty waitgroup is valid in golang

	wg.Add(1) // 1. add how much concurrent process that we must wait
	go func() {
		defer wg.Done() // 3. call the Done() method to decrease and tell the waitgroup if any concurrent process have done their job
		fmt.Println("1st goroutine is working.... and now sleeping")
		time.Sleep(1 * time.Nanosecond)
	}() // 2. start the concurrent process

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("2nd goroutine is working.... and now sleeping")
		time.Sleep(2 * time.Nanosecond)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("3rd goroutine is working..... and now sleeping")
		time.Sleep(3 * time.Nanosecond)
	}()

	wg.Wait() // 4. Waiting for all goroutine to complete their job
	fmt.Println("All goroutine has complete")
}

func ExampleWaitGroup2() {
	hello := func(wg *sync.WaitGroup, id int) {
		defer wg.Done()
		fmt.Println("Got hello message from person number", id, "!")
	}

	const numberPerson = 5
	var wg sync.WaitGroup

	wg.Add(numberPerson)
	for i := 0; i < numberPerson; i++ {
		go hello(&wg, i)
	}
	wg.Wait()
}
