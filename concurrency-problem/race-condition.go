package concurrencyproblem

import (
	"fmt"
	"sync"
	"time"
)

func printData(iteration int, data int) {
	if data == 0 {
		fmt.Printf("Data %d = %v\n", iteration, data)
	} else {
		fmt.Printf("Changed Data %d = %d\n", iteration, data)
	}
}

func RaceConditionDemo(iteration int) {
	var data int

	go func() {
		data = iteration
	}()

	printData(iteration, data)
}

func RaceConditionDemo2() {
	var mu sync.Mutex
	var counter int

	incrementFn := func() {
		mu.Lock()
		defer mu.Unlock()
		time.Sleep(3 * time.Second)
		counter++
	}

	decrementFn := func() {
		mu.Lock()
		defer mu.Unlock()
		time.Sleep(2 * time.Second)
		counter--
	}

	var stop = false

	go func() {
		for {
			go incrementFn()
		}
	}()

	go func() {
		for {
			go decrementFn()
		}
	}()

	for !stop {
		if counter%2 == 0 && counter != 0 && counter > 10 {
			fmt.Println(counter, "= GENAP")
			stop = true
		}
	}
}

func SolutionAddSleepToRaceCondition(iteration int) {
	var data int

	go func() {
		data = iteration
	}()

	// Here we add sleep to wait for goroutine come back to this function
	time.Sleep(1 * time.Second)

	// This solution is not effective because it's add more drawback to performance app instead of efficiency
	// Adding more time = Adding more response time and we can't measure exact time the goroutine can done the job
	// This solution is not 100% solve the race condition for this sample
	// Summary : Don't use more delay to solve your race condition problem

	printData(iteration, data)
}

func SolutionAddingWaitGroupRaceCondition(iteration int) {
	var data int
	var wg sync.WaitGroup

	// This is a critical section
	/**
	Because :
		1. this line is try to update the data value with the data that this function receive from the iteration parameter
		2. the data value is used by printData() function.
		     So, if we execute the update value in other routine that the printData() function use it will become a race condition because we never know when exact time the data has been updated by other goroutine
	*/
	wg.Add(1)
	go func() {
		defer wg.Done()
		data = iteration
	}()
	wg.Wait()

	// This solution is quite effective to solve the race condition problem
	// Instead of add exact delay for waiting the goroutine do their task.
	// Here we just declare a watcher (waitgroup) that can hold the execution of the program.
	// Whenever the goroutine has done their job, the watcher open the gate to another line execution of the program
	printData(iteration, data)
}

func SolutionWithMemorySynchronizationWithMutex(iteration int) {
	var data int
	var mu sync.Mutex

	// Lock the part when we gain access to shared memory data
	// Here we lock the critical section when we trying to modify variable data that in this case is shared to other routine by goroutine
	go func() {
		mu.Lock()
		defer mu.Unlock()
		data++
	}()

	// After that, we lock other critical section
	// In this stage, we lock the operation to printData because they need to access the data variable which in this case the memory is shared

	// This kind of solution maybe can solve our data race issue (you can see by run the program with go run -race main.go)
	// But still it cannot define whether part of the program that has been executed first. You'll notice the data value is always 0 or sometimes it can be same like the iteration variable
	mu.Lock()
	defer mu.Unlock()
	printData(iteration, data)
}
