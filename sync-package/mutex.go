package syncpackage

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"sync"
	"text/tabwriter"
	"time"
)

/*
	Mutex
		- Stands for "Mutual Exclusion"
		- A way to guard critical section in program

	Critical Section
		- is an area of your program requires exclusive access to a shared resources

	Mutex vs Channel :
		- Channel is share memory by communicating
		- Mutex is shares memory by creating a convention developers must follow to synchronize access to memory
*/

func ExampleMutex() {
	var count int       // Here is our critical section. Because, we share the count memory to increment() and decrement() method to modify
	var lock sync.Mutex // Is a valid way to instantiate mutex

	increment := func() {
		lock.Lock() // So, we lock here
		defer lock.Unlock()
		count++
		fmt.Printf("Incrementing count : %d\n", count)
	}

	decrement := func() {
		lock.Lock() // And here too to make sure the memory sync
		defer lock.Unlock()
		count--
		fmt.Printf("Decrementing count: %d\n", count)
	}

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			increment()
		}()
	}

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			decrement()
		}()
	}

	wg.Wait()

	fmt.Printf("Final count : %d\n", count)
}

// Here we do a little benchmarking the Mutex & RWMutex to know what exactly the cost of using them.
func BenchMarkMutexAndRWMutex() {
	producer := func(wg *sync.WaitGroup, l sync.Locker) {
		defer wg.Done()
		for i := 5; i > 0; i-- {
			l.Lock()
			l.Unlock()
			time.Sleep(1 * time.Nanosecond)
		}
	}

	observer := func(wg *sync.WaitGroup, l sync.Locker) {
		defer wg.Done()
		l.Lock()
		defer l.Unlock()
	}

	test := func(count int, mutex, rwMutex sync.Locker) time.Duration {
		var wg sync.WaitGroup
		wg.Add(count + 1)
		beginTestTime := time.Now()
		go producer(&wg, mutex)
		for i := count; i > 0; i-- {
			go observer(&wg, rwMutex)
		}

		wg.Wait()
		return time.Since(beginTestTime)
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 1, 2, ' ', 0)
	defer tw.Flush()

	var m sync.RWMutex
	fmt.Fprintf(tw, "Readers\tRWMutex\tMutex\n")
	for i := 0; i < 20; i++ {
		count := int(math.Pow(2, float64(i)))
		fmt.Fprintf(
			tw,
			"%d\t%v\t%v\n",
			count,
			test(count, &m, m.RLocker()),
			test(count, &m, &m),
		)
	}
}

func ExampleProblemRaceCondition() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	doneChannel := make(chan bool)

	// Shared Data
	var sharedData int

	incrementor := func() {
		sharedData++
	}

	decrementor := func() {
		sharedData--
	}

	printer := func() {
		fmt.Println("Real Data:", sharedData)
		fmt.Println()
	}

	go func() {
		for {
			incrementor()
			fmt.Println("Current value (incr):", sharedData)
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			decrementor()
			fmt.Println("Current value (decr):", sharedData)
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		time.Sleep(20 * time.Second)
		doneChannel <- true
	}()

	go func() {
		for {
			select {
			case <-doneChannel:
				fmt.Println("Program Exited")
				return

			case <-ticker.C:
				printer()
			}
		}
	}()

	<-doneChannel
}

func ExampleProblemRaceConditionSolution() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	doneChannel := make(chan bool)

	var mu sync.Mutex
	// Shared Data
	var sharedData int

	incrementor := func(lock *sync.Mutex) {
		lock.Lock()
		defer lock.Unlock()
		sharedData++
		fmt.Println("Current value (incr):", sharedData)
	}

	decrementor := func(lock *sync.Mutex) {
		lock.Lock()
		defer lock.Unlock()
		sharedData--
		fmt.Println("Current value (decr):", sharedData)
	}

	printer := func(lock *sync.Mutex) {
		lock.Lock()
		defer lock.Unlock()
		fmt.Println("Real Data:", sharedData)
		fmt.Println()
	}

	go func() {
		for {
			incrementor(&mu)
			rand := rand.Int() % (100 * 10)
			time.Sleep(time.Duration(rand * int(time.Millisecond)))
		}
	}()

	go func() {
		for {
			decrementor(&mu)
			rand := rand.Int() % (100 * 100)
			time.Sleep(time.Duration(rand * int(time.Millisecond)))
		}
	}()

	go func() {
		time.Sleep(20 * time.Second)
		doneChannel <- true
	}()

	for {
		select {
		case <-doneChannel:
			fmt.Println("Program Exited")
			return

		case <-ticker.C:
			printer(&mu)
		}
	}
}
