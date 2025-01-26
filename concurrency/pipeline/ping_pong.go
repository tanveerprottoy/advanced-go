package pipeline

import "time"

type Ball struct {
	hits int
}

func player(name string, table chan *Ball) {
	for {
		ball := <-table
		ball.hits++

		println(name, ball.hits)
		
		table <- ball
	}
}

func ExecuterPingPong() {
	table := make(chan *Ball)

	go player("ping", table)
	go player("pong", table)

	// start the game
	// table <- new(Ball)

	time.Sleep(1 * time.Second)

	// stop the game
	<-table
}
