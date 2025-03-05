package csvreader

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Read values
func readCSV(file string) (<-chan []string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("opening file %w", err)
	}

	ch := make(chan []string)

	go func() {
		cr := csv.NewReader(f)
		cr.FieldsPerRecord = 3

		for {
			record, err := cr.Read()
			if errors.Is(err, io.EOF) {
				close(ch)

				return
			}

			ch <- record
		}
	}()

	return ch, nil
}

// Remove "invalid" records
func sanitize(strC <-chan []string) <-chan []string {
	ch := make(chan []string)

	go func() {
		for val := range strC {
			if len(val[0]) > 3 {
				fmt.Println("skipped ", val)
				continue
			}

			ch <- val
		}

		close(ch)
	}()

	return ch
}

// Modify received values
func titleize(strC <-chan []string) <-chan []string {
	ch := make(chan []string)

	go func() {
		for val := range strC {
			c := cases.Title(language.English)

			val[0] = c.String(val[0])
			val[1], val[2] = val[2], val[1]

			ch <- val
		}

		close(ch)
	}()

	return ch
}

func ExecuterCSVReader() {
	recordsC, err := readCSV("data.csv")
	if err != nil {
		log.Fatalf("Could not read csv %v", err)
	}

	for val := range sanitize(titleize(recordsC)) {
		log.Printf("%v\n", val)
	}
}
