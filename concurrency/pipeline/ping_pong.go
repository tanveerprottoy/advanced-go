package pipeline

import (
	"time"
)

type Ball struct {
	hits int
}

func player(table chan *Ball, name string) {
	for {
		ball := <-table
		ball.hits++

		println(name, ball.hits)

		time.Sleep(100 * time.Millisecond)

		table <- ball
	}
}

func ExecuterPingPong() {
	table := make(chan *Ball)

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
