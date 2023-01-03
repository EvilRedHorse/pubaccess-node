package wallet

import (
	"os"

	"github.com/EvilRedHorse/pubaccess-node/persist"
	"github.com/EvilRedHorse/pubaccess-node/siatest"
)

// walletTestDir creates a temporary testing directory for a wallet test. This
// should only every be called once per test. Otherwise it will delete the
// directory again.
func walletTestDir(testName string) string {
	path := siatest.TestDir("wallet", testName)
	if err := os.MkdirAll(path, persist.DefaultDiskPermissionsTest); err != nil {
		panic(err)
	}
	return path
}
