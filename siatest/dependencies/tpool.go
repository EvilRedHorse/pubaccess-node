package dependencies

import (
	"github.com/EvilRedHorse/pubaccess-node/modules"
)

// DependencyDoNotAcceptTxnSet will not accept a transaction set.
type DependencyDoNotAcceptTxnSet struct {
	modules.ProductionDependencies
}

// Disrupt returns true if the correct string is provided.
func (d *DependencyDoNotAcceptTxnSet) Disrupt(s string) bool {
	return s == "DoNotAcceptTxnSet"
}
