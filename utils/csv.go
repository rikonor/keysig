package utils

import (
	"encoding/csv"
	"log"
	"os"
	"time"
)

// WriteToCSV accepts a file description (one word) and records to write to file
// The first line in our data should be row headers
func WriteToCSV(fDesc string, records [][]string) {
	// Append current time to file name
	cTime := time.Now().Format(time.RFC3339)

	fName := "./" + fDesc + "_" + cTime + ".csv"
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
