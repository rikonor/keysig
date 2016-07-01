package metrics

import (
	"github.com/rikonor/keysig/keylogger"
)

type Metric interface {
	RegisterWith(*keylogger.Manager) *Metric
}
