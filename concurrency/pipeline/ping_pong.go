package pipeline

import (
	"log"
	"time"
)

type Ball struct {
	hits int
}

func player(table chan *Ball, name string) {
	for {
		ball, ok := <-table

		if !ok {
			// channel closed
			log.Printf("channel was closed")
			return
		}

		ball.hits++

		log.Println(name, ball.hits)

		time.Sleep(100 * time.Millisecond)

		table <- ball
	}
}

func ExecuterPingPong() {
	table := make(chan *Ball)

	// close the table when this
	// function returns
	defer close(table)

	// run player in separate goroutines
	go player(table, "ping")

	go player(table, "pong")

	// start the game
	table <- new(Ball)

	// this will block the main goroutine
	// which will cause the player function
	// to keep running one after another
	// through the table channel
	time.Sleep(1 * time.Second)

	// log.Println("main goroutine resumed")

	// when sleep is done
	// stop the game
	// this will exit the main goroutine
	// but the other two will not exit
	// which is a leak
	<-table
}
