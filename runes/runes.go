package runes

import (
	"fmt"
	"log"
)

func Ex() {
	var r rune = 'T'
	log.Println(r)

	// both ASCII and non-ASCII characters.
	str := "Hello, 世界"
	// slice of runes.
	runes := []rune(str)
	log.Println(runes)

	u := '😀'

	fmt.Printf("Data type of %c is %T and the rune value is %U \n", u, u, u)
}
