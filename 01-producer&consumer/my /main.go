package main

import (
	"fmt"
	"github.com/fatih/color"
	"math/rand"
	"time"
)

const NumberOfPizzas = 10

var pizzasMade, pizzasFailed, pizzasTotal int

type Producer struct {
	data chan PizzaOrder
	quit chan chan error
}

type PizzaOrder struct {
	pizzaNumber int
	message     string
	success     bool
}

func (p *Producer) Close() error {
	ch := make(chan error)
	p.quit <- ch
	return <-ch
}

func makePizza(pizzaNumber int) *PizzaOrder {
	pizzaNumber++
	if pizzaNumber <= NumberOfPizzas {
		delay := rand.Intn(5) + 1
		fmt.Printf("Received order #%d!\n", pizzaNumber)

		rnd := rand.Intn(12) + 1
		msg := ""
		success := false

		if rnd < 5 {
			pizzasFailed++
		} else {
			pizzasMade++
		}
		pizzasTotal++

		fmt.Printf("Making pizza number #%d, it will take %d seconds...\n", pizzaNumber, delay)
		//delay for a bit
		time.Sleep(time.Duration(delay) * time.Second)

		if rnd <= 2 {
			msg = fmt.Sprintf("*** We ran out of ingridients for pizza #%d", pizzaNumber)
		} else if rnd <= 4 {
			msg = fmt.Sprintf("*** The cook quit while making pizza #%d", pizzaNumber)
		} else {
			success = true
			msg = fmt.Sprintf("Pizza order #%d is ready!", pizzaNumber)
		}

		p := PizzaOrder{
			pizzaNumber: pizzaNumber,
			message:     msg,
			success:     success,
		}
		return &p
	}

	return &PizzaOrder{
		pizzaNumber: pizzaNumber,
	}
}

func pizzeria(pizzaMaker *Producer) {
	// keep track of which pizza we are making
	var i = 0
	// run forever or until we receive a quit notification
	// try to make pizzas

	for {
		currentPizza := makePizza(i)
		if currentPizza != nil {
			i = currentPizza.pizzaNumber
			select {
			// we try to make a pizza (we sent something to the data chanel)
			case pizzaMaker.data <- *currentPizza:
			case quitChan := <-pizzaMaker.quit:
				// close channel
				close(pizzaMaker.data)
				close(quitChan)
				return
			}
		}
	}

}

func main() {
	// seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// print out a message
	color.Cyan("The pizzeria is open for business!\n")

	// create a producer
	pizzaJob := Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}

	//run the producer on the background
	go pizzeria(&pizzaJob)

	for i := range pizzaJob.data {
		if i.pizzaNumber <= NumberOfPizzas {
			if i.success {
				color.Green(i.message)
				color.Green("Order #%d is out for delivery", i.pizzaNumber)
			} else {
				color.Red(i.message)
				color.Red("The customer is really mad")
			}
		} else {
			color.Cyan("Done making pizzas...")
			err := pizzaJob.Close()
			if err != nil {
				color.Red("*** Error clowing channel!", err)
			}
		}

		// print out the ending message
		color.Cyan("\nDone for the day.")

		color.Cyan("We made %d pizzas, but failed to make %d, with %d attempts", pizzasMade, pizzasFailed, pizzasTotal)
		switch {
		case pizzasFailed > 9:
			color.Red("Awful day")
		case pizzasFailed >= 6:
			color.Red("It was not a very good day...")
		case pizzasFailed >= 4:
			color.Yellow("It was an okay day...")
		case pizzasFailed >= 2:
			color.Yellow("It was a pretty good day")
		default:
			color.Green("It was a great day")
		}

	}

}
