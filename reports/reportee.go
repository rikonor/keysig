package reports

// Reportee supports giving back data for reports
type Reportee interface {
	Data() [][]string
	WriteToCSV()
	RegisterWithReporter(r *Reporter)
}
