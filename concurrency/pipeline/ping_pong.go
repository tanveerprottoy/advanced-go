package pipeline

import (
	"log"
	"time"
)

type Ball struct {
	hits int
}

func player(name string, table chan *Ball) {
	for {
		ball := <-table
		ball.hits++

		println(name, ball.hits)

		table <- ball

		log.Printf("running player in loop, table: %v", table)
	}
}

func ExecuterPingPong() {
	table := make(chan *Ball)

	// run player in separate goroutines
	go player("ping", table)

	go player("pong", table)

	// start the game
	table <- new(Ball)

	// this will block the main goroutine
	// which will cause the player function
	// to keep running one after another
	// through the table channel
	time.Sleep(1 * time.Microsecond)

	log.Println("main goroutine resumed")

	// when sleep is done
	// stop the game
	<-table
}
