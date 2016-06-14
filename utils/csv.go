package utils

import (
	"encoding/csv"
	"log"
	"os"
)

// WriteToCSV accepts a file description (one word) and records to write to file
// The first line in our data should be row headers
func WriteToCSV(fDesc string, records [][]string) {
	fName := "./" + fDesc + ".csv"
	f, err := os.Create(fName)
	if err != nil {
		log.Fatalln("error creating file:", fName)
	}

	w := csv.NewWriter(f)

	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}
