package main

import (
	"fmt"

	"github.com/tanveerprottoy/advanced-go/structs"
)

func main() {
	c := structs.NewCar(
		structs.WithClass("coupe"),
		structs.WithMod(true),
		structs.WithSuperCharger("turbo"),
	)
	// prints:
	fmt.Printf("Car: %+v\n", c)
	fmt.Printf("Class: %s\n", c.Class())
	fmt.Printf("Class: %t\n", c.Mod())
	fmt.Printf("Class: %s\n", c.SuperCharger())
}
