package waitgroup

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

// doWork simulates a time consuming task
func doWork(ch chan<- int) {
	time.Sleep(1 * time.Second)

	ch <- rand.Intn(100)
}

func ExecuterDoWork() {
	// unbuffered channel
	ch := make(chan int)

	wg := sync.WaitGroup{}

	go func() {
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				doWork(ch)
			}()
		}

		// wait for all goroutines to finish
		wg.Wait()

		// close the channel
		close(ch)
	}()

	for n := range ch {
		log.Printf("Received: %d\n", n)
	}
}
