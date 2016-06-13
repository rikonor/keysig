package reports

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

type Reporter struct {
	reportees []Reportee
}

func New() *Reporter {
	return &Reporter{}
}

func (r *Reporter) WriteToCSV(p Reportee) {
	fmt.Println("Writing metrics to CSV file")

	records := p.Data()

	f, err := os.Create("./output.csv")
	if err != nil {
		log.Fatalln("error creating file: ./output.csv")
	}

	w := csv.NewWriter(f)

	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	// Write any buffered data to the underlying writer (standard output).
	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}
