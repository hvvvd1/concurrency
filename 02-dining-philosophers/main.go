package main

import (
	"fmt"
	"sync"
	"time"
)

const hunger = 3

var eatTime = 3 * time.Second
var sleepTime = 1 * time.Second
var thinkTime = 1 * time.Second

var wg sync.WaitGroup

var philosophers = []string{
	"Аристович", "Платон", "Рене Декарт",
	"Алексей Шевцов", "Будда"}

var philosophersFinishEatingOrder []string
var orderMutex sync.Mutex

func diningProblem(philosopher string, leftFork, rightFork *sync.Mutex) {
	defer wg.Done()
	fmt.Println(philosopher, "is seated")
	time.Sleep(sleepTime)

	for i := hunger; i > 0; i-- {
		fmt.Println(philosopher, "is hungry")
		time.Sleep(sleepTime)

		leftFork.Lock()
		fmt.Printf("\t%s picked up the fork to his left\n", philosopher)
		rightFork.Lock()
		fmt.Printf("\t%s picked up the fork to his right\n", philosopher)

		fmt.Println(philosopher, "has both forks and is eating")
		time.Sleep(eatTime)

		fmt.Println(philosopher, "is thinking.")
		time.Sleep(thinkTime)

		fmt.Printf("\t%s put down the fork on his right.\n", philosopher)
		rightFork.Unlock()
		fmt.Printf("\t%s put down the fork on his left.\n", philosopher)
		leftFork.Unlock()
		time.Sleep(sleepTime)
	}

	fmt.Println(philosopher, "is satisfied")

	orderMutex.Lock()
	philosophersFinishEatingOrder = append(philosophersFinishEatingOrder, philosopher)
	orderMutex.Unlock()

	time.Sleep(sleepTime)

	fmt.Println(philosopher, "has left the table")
}

func main() {
	fmt.Println("The dining Philosophers problem ")

	wg.Add(len(philosophers))

	forkLeft := &sync.Mutex{}

	for i := 0; i < len(philosophers); i++ {
		forkRight := &sync.Mutex{}

		go diningProblem(philosophers[i], forkLeft, forkRight)

		forkLeft = forkRight
	}

	wg.Wait()

	fmt.Println("order in which philosophers left table:", philosophersFinishEatingOrder)
	fmt.Println("the table is empty")
}
