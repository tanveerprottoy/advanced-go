package main

import "txp/concurrency/app"

func main() {
	app.InitGoRoutineEx()
	app.InitChannelEx()
	app.InitChannelBufferEx()
	app.InitChannelCloseEx()
	app.InitChannelSelectEx()
	app.InitDefaultSelectionEx()
	app.InitMutexEx()
	app.InitWaitGroupEx()
	go app.ChannelBuffer()
}