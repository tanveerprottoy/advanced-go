package std

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

func scan() {
	fmt.Print("Enter text: ")

	var input string

	fmt.Scanln(&input)

	fmt.Print(input)
}

func read() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text: ")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)

	fmt.Print("Enter text: ")
	var input string
	fmt.Scanln(&input)
	fmt.Print(input)
}

func flagEx() {
	fl := flag.Int("n", 1234, "help message for flag n")

	flBool := flag.Bool("b", false, "help message for flagname")

	flag.Parse()

	log.Println("ip has value ", *fl)
	log.Println("flagvar has value ", flBool)
}

func inputEx() {
	lnNumber := flag.Bool("n", false, "ln number")
	flag.Parse()

	input := os.Args[2]
	file, _ := os.Open("./" + input)
	scanner := bufio.NewScanner(file)
	ln := 1
	for scanner.Scan() {
		if *lnNumber {
			fmt.Print(ln, ": ")
		}

		fmt.Println(scanner.Text())
	}
}


