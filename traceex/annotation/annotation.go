package traceex

import (
	"context"
	"log"
	"runtime/trace"
)

// User annotation
// Package trace provides user annotation APIs that can be used to log interesting
// events during execution.
// There are three types of user annotations: log messages, regions, and tasks.
// Log emits a timestamped message to the execution trace along with additional
// information such as the category of the message and which goroutine called Log.
// The execution tracer provides UIs to filter and group goroutines using the log
// category and the message supplied in Log.
// A region is for logging a time interval during a goroutine's execution. By
// definition, a region starts and ends in the same goroutine. Regions can be nested
// to represent subintervals. For example, the following code records four regions
// in the execution trace to trace the durations of sequential steps in a cappuccino
// making operation.
func withRegionEx(ctx context.Context, orderID string) {
	trace.WithRegion(ctx, "makeCappuccino", func() {

		// orderID allows to identify a specific order
		// among many cappuccino order region records.
		trace.Log(ctx, "orderID", orderID)

		trace.WithRegion(ctx, "steamMilk", func() {
			trace.WithRegion(ctx, "boilWater", func() {
				log.Println("boiling water")
			})
		})

		trace.WithRegion(ctx, "extractCoffee", func() {
			log.Println("extracting coffee")
		})

		trace.WithRegion(ctx, "mixMilkCoffee", func() {
			log.Println("mixing milk and coffee")
		})
	})
}

// A task is a higher-level component that aids tracing of logical operations such as
// an RPC request, an HTTP request, or an interesting local operation which may require
// multiple goroutines working together. Since tasks can involve multiple goroutines,
// they are tracked via a context.Context object. NewTask creates a new task and embeds
// it in the returned context.Context object. Log messages and regions are attached to
// the task, if any, in the Context passed to Log and WithRegion.
// For example, assume that we decided to froth milk, extract coffee, and mix milk and
// coffee in separate goroutines. With a task, the trace tool can identify the
// goroutines involved in a specific cappuccino order.
func taskEx(ctx context.Context, orderID string) {
	ctx, task := trace.NewTask(ctx, "makeCappuccino")
	trace.Log(ctx, "orderID", orderID)

	milk := make(chan bool)
	espresso := make(chan bool)

	go func() {
		trace.WithRegion(ctx, "steamMilk", func() {
			log.Println("boiling water")
		})

		milk <- true
	}()

	go func() {
		trace.WithRegion(ctx, "extractCoffee", func() {
			log.Println("extracting coffee")
		})

		espresso <- true
	}()

	go func() {
		defer task.End() // When assemble is done, the order is complete.

		<-espresso
		<-milk

		trace.WithRegion(ctx, "mixMilkCoffee", func() {
			log.Println("mixing milk and coffee")
		})
	}()
}

func Executer() {
	log.Printf("trace.IsEnabled() = %t\n", trace.IsEnabled())

	withRegionEx(context.Background(), "123")

	taskEx(context.Background(), "456")
}
