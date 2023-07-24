package main

import "github.com/tanveerprottoy/concurrency-go/internal/concurrency"

func main() {
	concurrency.InitGoRoutineEx()
	concurrency.InitChannelEx()
	concurrency.InitChannelBufferEx()
	concurrency.InitChannelCloseEx()
	concurrency.InitChannelSelectEx()
	concurrency.InitDefaultSelectionEx()
	concurrency.InitMutexEx()
	concurrency.InitWaitGroupEx()
	go concurrency.ChannelBuffer()
}
