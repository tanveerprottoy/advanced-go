package main

import concurrency "github.com/tanveerprottoy/concurrency-go/internal/concurrency/misc"

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
	concurrency.Timeout()
}
