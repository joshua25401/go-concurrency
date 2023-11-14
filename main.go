package main

import (
	"fmt"
	syncpackage "go-concurrency/sync-package"
	"time"
)

func main() {
	var startExec = time.Now()
	// Uncoment this line to use the function

	// Race Condition
	// Please use this command to run the main.go => go run -race main.go
	// concurrencyproblem.RaceConditionDemo()``
	// concurrencyproblem.RaceConditionDemo2()

	// Without Solution (Uncomment the line below)
	// for i := 0; i < 1000; i++ {
	// 	concurrencyproblem.RaceConditionDemo(i)
	// }

	// With delay solution by adding 1 second delay (Uncomment the line below)
	// for i := 0; i < 5; i++ {
	// 	concurrencyproblem.SolutionAddSleepToRaceCondition(i)
	// }

	// With memory sync solution
	// for i := 0; i < 100; i++ {
	// 	concurrencyproblem.SolutionWithMemorySynchronizationWithMutex(i)
	// }

	// With wait group solution
	// for i := 0; i < 100; i++ {
	// 	concurrencyproblem.SolutionAddingWaitGroupRaceCondition(i)
	// }

	// Deadlock
	// concurrencyproblem.DeadLock()

	// Livelock
	// concurrencyproblem.LiveLock()

	// Sync Package

	// WaitGroup
	// syncpackage.ExampleWaitGroup()
	// syncpackage.ExampleWaitGroup2()

	// Mutex
	// syncpackage.ExampleMutex()
	// syncpackage.BenchMarkMutexAndRWMutex()
	syncpackage.ExampleProblemRaceConditionSolution()

	// Cond
	// syncpackage.ExampleCondUsage()
	// syncpackage.ExampleUseBroadcastCond()
	var totalExecutionTime = time.Since(startExec)
	fmt.Printf("Total Execution Time = %v ms\n", totalExecutionTime.Milliseconds())
}
