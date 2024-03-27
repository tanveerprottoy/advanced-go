package main

import (
	"fmt"

	"github.com/tanveerprottoy/concurrency-go/internal/structs"
)

func main() {
	c := structs.NewCar(
		structs.WithClass("coupe"),
		structs.WithMod(true),
		structs.WithSuperCharger("turbo"),
	)
	fmt.Printf("Car: %+v\n", c)
	fmt.Printf("Class: %s\n", c.Class())
	fmt.Printf("Class: %t\n", c.Mod())
	fmt.Printf("Class: %s\n", c.SuperCharger())
}
