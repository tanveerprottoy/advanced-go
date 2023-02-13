package async

import (
	"fmt"
	"time"
)

func SleeperFunc(ch chan int, d time.Duration) {
	fmt.Println("sleeping for: ", d)
	time.Sleep(d)
	ch <- 0
}

func CallerFunc() {
	ch := make(chan int)
	go SleeperFunc(ch, time.Second*5)
	res := <-ch
	fmt.Println(res)
}
