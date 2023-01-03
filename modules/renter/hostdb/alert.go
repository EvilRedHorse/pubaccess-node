package hostdb

import "github.com/EvilRedHorse/pubaccess-node/modules"

// Alerts implements the modules.Alerter interface for the hostdb. It returns
// all alerts of the hostdb.
func (hdb *HostDB) Alerts() (crit, err, warn []modules.Alert) {
	return hdb.staticAlerter.Alerts()
}
