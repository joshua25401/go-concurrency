package concurrencyproblem

import (
	"fmt"
	"sync"
	"time"
)

/*
	Deadlock is a state when all goroutine is waiting on one another and this problem can't be solved without outside intervention

	TODO: Summarize the Coffman Law Below!
*/

type value struct {
	mu    sync.Mutex
	value int
}

func DeadLock() {
	var wg sync.WaitGroup
	printSum := func(v1, v2 *value) {
		defer wg.Done()

		v1.mu.Lock()
		defer v1.mu.Unlock()

		time.Sleep(2 * time.Second)

		v2.mu.Lock()
		defer v2.mu.Unlock()

		fmt.Printf("sum = %v\n", v1.value+v2.value)
	}

	var a, b value

	wg.Add(2)
	go printSum(&a, &b)
	go printSum(&b, &a)
	wg.Wait()
}
