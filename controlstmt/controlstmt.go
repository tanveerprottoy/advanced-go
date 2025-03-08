package controlstmt

import "log"

func label() {
	x := 0

	// for loop work as a while loop
Lable1:
	for x < 8 {

		if x == 5 {
			// using goto statement
			x = x + 1

			goto Lable1
		}

		log.Printf("value: %d\n", x)

		x++
	}
}
