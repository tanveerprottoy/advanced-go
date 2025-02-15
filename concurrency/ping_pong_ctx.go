package concurrency

import (
	"context"
	"log"
	"time"
)

func playerCtx(ctx context.Context, table chan *Ball, name string) {
	for {
		ball := <-table
		ball.hits++

		log.Println(name, ball.hits)

		time.Sleep(100 * time.Millisecond)

		table <- ball

		// check for context deadline or cancelled
		// before starting the loop
		if ctx.Err() == context.DeadlineExceeded || ctx.Err() == context.Canceled {
			// break the loop
			break
		}
	}
}

func ExecuterPingPongCtx() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	table := make(chan *Ball)

	// run player in separate goroutines
	go playerCtx(ctx, table, "ping")

	go playerCtx(ctx, table, "pong")

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
