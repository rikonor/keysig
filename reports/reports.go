package reports

// Reporter tracks reportees and instructs them to report
type Reporter struct {
	reportees map[string]Reportee
}

// New returns a new Reporter pointer
func New() *Reporter {
	return &Reporter{
		reportees: make(map[string]Reportee),
	}
}

// Register adds a Reportee to the reporter's set of reportees
func (r *Reporter) Register(pName string, p Reportee) {
	r.reportees[pName] = p
}

// CollectReports iterates over registered reportees and instructs them
// to write their reports to file
func (r *Reporter) CollectReports() {
	for _, p := range r.reportees {
		p.WriteToCSV()
	}
}
