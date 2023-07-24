package main

import "github.com/tanveerprottoy/concurrency-go/internal/app/concurrency/misc"

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
