package syncpackage

import (
	"fmt"
	"sync"
	"time"
)

/*
	Cond :
		- Is the point where goroutines waiting for or announcing the occurence of an event

	Event :
		- any arbitary signal between two or more goroutines that carries no information
		   other than it's just occured

	Syntax:
		sync.NewCond(sync.Locker)

		sync.Locker is an interface so we can use any locker that implement Lock() and Unlock() function
*/

func ExampleCondUsage() {
	cond := sync.NewCond(&sync.Mutex{})

	queue := make([]interface{}, 0, 10)

	removeFromQueueFn := func(delay time.Duration) {
		time.Sleep(delay)
		cond.L.Lock()
		queue = queue[1:]
		fmt.Println("Removed from queue")
		defer cond.Signal()
		defer cond.L.Unlock()
	}

	for i := 0; i < 10; i++ {
		cond.L.Lock()

		if len(queue) == 2 {
			cond.Wait()
		}

		fmt.Println("Adding to queue")
		queue = append(queue, struct{}{})
		go removeFromQueueFn(time.Second)
		cond.L.Unlock()
	}
}

func ExampleUseBroadcastCond() {
	type Button struct {
		Clicked *sync.Cond
	}

	button := Button{
		Clicked: sync.NewCond(&sync.Mutex{}),
	}

	subscribe := func(cond *sync.Cond, fn func()) {
		var goroutineRunning sync.WaitGroup

		goroutineRunning.Add(1)
		go func() {
			goroutineRunning.Done()
			cond.L.Lock()
			defer cond.L.Unlock()
			cond.Wait()
			fn()
		}()
		goroutineRunning.Wait()
	}

	var clickRegistered sync.WaitGroup

	clickRegistered.Add(3)
	subscribe(button.Clicked, func() {
		fmt.Println("Maximizing window")
		defer clickRegistered.Done()
	})

	subscribe(button.Clicked, func() {
		fmt.Println("Displaying annoying dialog box!")
		defer clickRegistered.Done()
	})

	subscribe(button.Clicked, func() {
		fmt.Println("Mouse clicked")
		defer clickRegistered.Done()
	})
	button.Clicked.Broadcast()
	clickRegistered.Wait()
}
