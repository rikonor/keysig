package reports

import "time"

const defaultCollectionPeriod = 30 * time.Second

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

// TriggerPeriodically will periodically trigger collection of reports
func (r *Reporter) TriggerPeriodically(d time.Duration) {
	cp := defaultCollectionPeriod
	if d != 0 {
		cp = d
	}

	go func() {
		for {
			time.Sleep(cp)
			r.CollectReports()
		}
	}()
}
