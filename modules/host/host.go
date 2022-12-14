// Package host is an implementation of the host module, and is responsible for
// participating in the storage ecosystem, turning available disk space an
// internet bandwidth into profit for the user.
package host

// TODO: what happens if the renter submits the revision early, before the
// final revision. Will the host mark the contract as complete?

// TODO: Host and renter are reporting errors where the renter is not adding
// enough fees to the file contract.

// TODO: Test the safety of the builder, it should be okay to have multiple
// builders open for up to 600 seconds, which means multiple blocks could be
// received in that time period. Should also check what happens if a parent
// gets confirmed on the blockchain before the builder is finished.

// TODO: Double check that any network connection has a finite deadline -
// handling action items properly requires that the locks held on the
// obligations eventually be released. There's also some more advanced
// implementation that needs to happen with the storage obligation locks to
// make sure that someone who wants a lock is able to get it eventually.

// TODO: Add contract compensation from form contract to the storage obligation
// financial metrics, and to the host's tracking.

// TODO: merge the network interfaces stuff, don't forget to include the
// 'announced' variable as one of the outputs.

// TODO: 'announced' doesn't tell you if the announcement made it to the
// blockchain.

// TODO: Need to make sure that the revision exchange for the renter and the
// host is being handled correctly. For the host, it's not so difficult. The
// host need only send the most recent revision every time. But, the host
// should not sign a revision unless the renter has explicitly signed such that
// the 'WholeTransaction' fields cover only the revision and that the
// signatures for the revision don't depend on anything else. The renter needs
// to verify the same when checking on a file contract revision from the host.
// If the host has submitted a file contract revision where the signatures have
// signed the whole file contract, there is an issue.

// TODO: there is a mistake in the file contract revision rpc, the host, if it
// does not have the right file contract id, should be returning an error there
// to the renter (and not just to it's calling function without informing the
// renter what's up).

// TODO: Need to make sure that the correct height is being used when adding
// sectors to the storage manager - in some places right now WindowStart is
// being used but really it's WindowEnd that should be in use.

// TODO: The host needs some way to blacklist file contracts that are being
// abusive by repeatedly getting free download batches.

// TODO: clean up all of the magic numbers in the host.

// TODO: revamp the finances for the storage obligations.

// TODO: host_test.go has commented out tests.

// TODO: network_test.go has commented out tests.

// TODO: persist_test.go has commented out tests.

// TODO: update_test.go has commented out tests.

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"gitlab.com/NebulousLabs/encoding"
	connmonitor "gitlab.com/NebulousLabs/monitor"
	"gitlab.com/NebulousLabs/siamux"
	"github.com/EvilRedHorse/pubaccess-node/build"
	"github.com/EvilRedHorse/pubaccess-node/crypto"
	"github.com/EvilRedHorse/pubaccess-node/modules"
	"github.com/EvilRedHorse/pubaccess-node/modules/host/api"
	"github.com/EvilRedHorse/pubaccess-node/modules/host/contractmanager"
	"github.com/EvilRedHorse/pubaccess-node/modules/host/mdm"
	"github.com/EvilRedHorse/pubaccess-node/modules/host/tokenstorage"
	"github.com/EvilRedHorse/pubaccess-node/persist"
	siasync "github.com/EvilRedHorse/pubaccess-node/sync"
	"github.com/EvilRedHorse/pubaccess-node/types"
)

const (
	// Directory name for token storage.
	tokenStorDir = "token_storage"
	// Names of the various persistent files in the host.
	dbFilename   = modules.HostDir + ".db"
	logFile      = modules.HostDir + ".log"
	settingsFile = modules.HostDir + ".json"
)

var (
	// dbMetadata is a header that gets put into the database to identify a
	// version and indicate that the database holds host information.
	dbMetadata = persist.Metadata{
		Header:  "Sia Host DB",
		Version: "0.5.2",
	}

	// Nil dependency errors.
	errNilCS      = errors.New("host cannot use a nil state")
	errNilTpool   = errors.New("host cannot use a nil transaction pool")
	errNilWallet  = errors.New("host cannot use a nil wallet")
	errNilGateway = errors.New("host cannot use nil gateway")

	// rpcPriceGuaranteePeriod defines the amount of time a host will guarantee
	// its prices to the renter.
	rpcPriceGuaranteePeriod = build.Select(build.Var{
		Standard: 10 * time.Minute,
		Dev:      5 * time.Minute,
		Testing:  1 * time.Minute,
	}).(time.Duration)

	// pruneExpiredRPCPriceTableFrequency is the frequency at which the host
	// checks if it can expire price tables that have an expiry in the past.
	pruneExpiredRPCPriceTableFrequency = build.Select(build.Var{
		Standard: 15 * time.Minute,
		Dev:      10 * time.Minute,
		Testing:  30 * time.Second,
	}).(time.Duration)

	storageObligationAuditInteval = types.BlocksPerDay // Audit contracts once a day
)

// A Host contains all the fields necessary for storing files for clients and
// performing the storage proofs on the received files.
type Host struct {
	// RPC Metrics - atomic variables need to be placed at the top to preserve
	// compatibility with 32bit systems. These values are not persistent.
	atomicDownloadCalls     uint64
	atomicErroredCalls      uint64
	atomicFormContractCalls uint64
	atomicRenewCalls        uint64
	atomicReviseCalls       uint64
	atomicSettingsCalls     uint64
	atomicUnrecognizedCalls uint64

	// Error management. There are a few different types of errors returned by
	// the host. These errors intentionally not persistent, so that the logging
	// limits of each error type will be reset each time the host is reset.
	// These values are not persistent.
	atomicCommunicationErrors uint64
	atomicConnectionErrors    uint64
	atomicConsensusErrors     uint64
	atomicInternalErrors      uint64
	atomicNormalErrors        uint64

	// Dependencies.
	cs            modules.ConsensusSet
	g             modules.Gateway
	tpool         modules.TransactionPool
	wallet        modules.Wallet
	staticAlerter *modules.GenericAlerter
	staticMux     *siamux.SiaMux
	dependencies  modules.Dependencies
	modules.StorageManager

	// Subsystems
	staticAccountManager *accountManager
	staticMDM            *mdm.MDM

	// Host ACID fields - these fields need to be updated in serial, ACID
	// transactions.
	announced    bool
	blockHeight  types.BlockHeight
	publicKey    types.SiaPublicKey
	secretKey    crypto.SecretKey
	recentChange modules.ConsensusChangeID
	unlockHash   types.UnlockHash // A wallet address that can receive coins.

	// Host transient fields - these fields are either determined at startup or
	// otherwise are not critical to always be correct.
	autoAddress          modules.NetAddress // Determined using automatic tooling in network.go
	financialMetrics     modules.HostFinancialMetrics
	settings             modules.HostInternalSettings
	revisionNumber       uint64
	workingStatus        modules.HostWorkingStatus
	connectabilityStatus modules.HostConnectabilityStatus
	// scheduledAuditBlockheight is the blockheight when the next AuditStorageObligations()
	// should be started
	scheduledAuditBlockheight uint64

	// A map of storage obligations that are currently being modified. Locks on
	// storage obligations can be long-running, and each storage obligation can
	// be locked separately.
	lockedStorageObligations map[types.FileContractID]*lockedObligation

	// Storage of tokens for prepaid downloads.
	tokenStor *tokenstorage.TokenStorage

	// A collection of rpc price tables, covered by its own RW mutex. It
	// contains the host's current price table and the set of price tables the
	// host has communicated to all renters, thus guaranteeing a set of prices
	// for a fixed period of time. The host's RPC prices are dynamic, and are
	// subject to various conditions specific to the RPC in question. Examples
	// of such conditions are congestion, load, liquidity, etc.
	staticPriceTables *hostPrices

	// Fields related to RHP3 bandwidhth.
	atomicStreamUpload   uint64
	atomicStreamDownload uint64

	// Misc state.
	db            *persist.BoltDatabase
	listener      net.Listener
	log           *persist.Logger
	mu            sync.RWMutex
	staticMonitor *connmonitor.Monitor
	persistDir    string
	port          string
	apiPort       string
	tg            siasync.ThreadGroup
}

// hostPrices is a helper type that wraps both the host's RPC price table and
// the set of price tables containing prices it has guaranteed to all renters,
// covered by a read write mutex to help lock contention. It contains a separate
// minheap that enables efficiently purging expired price tables from the
// 'guaranteed' map.
type hostPrices struct {
	current       modules.RPCPriceTable
	guaranteed    map[modules.UniqueID]*hostRPCPriceTable
	staticMinHeap priceTableHeap
	mu            sync.RWMutex
}

// managedCurrent returns the host's current price table
func (hp *hostPrices) managedCurrent() modules.RPCPriceTable {
	hp.mu.RLock()
	defer hp.mu.RUnlock()
	return hp.current
}

// managedGet returns the price table with given uid
func (hp *hostPrices) managedGet(uid modules.UniqueID) (pt *hostRPCPriceTable, found bool) {
	hp.mu.RLock()
	defer hp.mu.RUnlock()
	pt, found = hp.guaranteed[uid]
	return
}

// managedSetCurrent overwrites the current price table with the one that's
// given
func (hp *hostPrices) managedSetCurrent(pt modules.RPCPriceTable) {
	hp.mu.Lock()
	defer hp.mu.Unlock()
	hp.current = pt
}

// managedTrack adds the given price table to the 'guaranteed' map, that holds
// all of the price tables the host has recently guaranteed to renters. It will
// also add it to the heap which facilates efficient pruning of that map.
func (hp *hostPrices) managedTrack(pt *hostRPCPriceTable) {
	hp.mu.Lock()
	hp.guaranteed[pt.UID] = pt
	hp.mu.Unlock()
	hp.staticMinHeap.Push(pt)
}

// managedPruneExpired removes all of the price tables that have expired from
// the 'guaranteed' map.
func (hp *hostPrices) managedPruneExpired() {
	current := hp.managedCurrent()
	expired := hp.staticMinHeap.PopExpired()
	if len(expired) == 0 {
		return
	}
	hp.mu.Lock()
	for _, uid := range expired {
		// Sanity check to never prune the host's current price table. This can
		// never occur because the host's price table UID is not added to the
		// minheap.
		if uid == current.UID {
			build.Critical("The host's current price table should not be pruned")
			continue
		}
		delete(hp.guaranteed, uid)
	}
	hp.mu.Unlock()
}

// lockedObligation is a helper type that locks a TryMutex and a counter to
// indicate how many times the locked obligation has been fetched from the
// lockedStorageObligations map already.
type lockedObligation struct {
	mu siasync.TryMutex
	n  uint
}

// checkUnlockHash will check that the host has an unlock hash. If the host
// does not have an unlock hash, an attempt will be made to get an unlock hash
// from the wallet. That may fail due to the wallet being locked, in which case
// an error is returned.
func (h *Host) checkUnlockHash() error {
	addrs, err := h.wallet.AllAddresses()
	if err != nil {
		return err
	}
	hasAddr := false
	for _, addr := range addrs {
		if h.unlockHash == addr {
			hasAddr = true
			break
		}
	}
	if !hasAddr || h.unlockHash == (types.UnlockHash{}) {
		uc, err := h.wallet.NextAddress()
		if err != nil {
			return err
		}

		// Set the unlock hash and save the host. Saving is important, because
		// the host will be using this unlock hash to establish identity, and
		// losing it will mean silently losing part of the host identity.
		h.unlockHash = uc.UnlockHash()
		err = h.saveSync()
		if err != nil {
			return err
		}
	}
	return nil
}

// managedInternalSettings returns the settings of a host.
func (h *Host) managedInternalSettings() modules.HostInternalSettings {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.settings
}

// managedUpdatePriceTable will recalculate the RPC costs and update the host's
// price table accordingly.
func (h *Host) managedUpdatePriceTable() {
	// create a new RPC price table
	hes := h.ManagedExternalSettings()
	priceTable := modules.RPCPriceTable{
		// TODO: hardcoded cost should be updated to use a better value.
		AccountBalanceCost:   types.NewCurrency64(1),
		FundAccountCost:      types.NewCurrency64(1),
		UpdatePriceTableCost: types.NewCurrency64(1),

		// TODO: hardcoded MDM costs should be updated to use better values.
		HasSectorBaseCost:   types.NewCurrency64(1),
		MemoryTimeCost:      types.NewCurrency64(1),
		DropSectorsBaseCost: types.NewCurrency64(1),
		DropSectorsUnitCost: types.NewCurrency64(1),
		SwapSectorCost:      types.NewCurrency64(1),

		// Read related costs.
		ReadBaseCost:   hes.SectorAccessPrice,
		ReadLengthCost: types.NewCurrency64(1),

		// Write related costs.
		WriteBaseCost:   types.NewCurrency64(1),
		WriteLengthCost: types.NewCurrency64(1),
		WriteStoreCost:  hes.StoragePrice,

		// Init costs.
		InitBaseCost: hes.BaseRPCPrice,

		// LatestRevisionCost is set to a reasonable base + the estimated
		// bandwidth cost of downloading a filecontract. This isn't perfect but
		// at least scales a bit as the host updates their download bandwidth
		// prices.
		LatestRevisionCost: modules.DefaultBaseRPCPrice.Add(hes.DownloadBandwidthPrice.Mul64(modules.EstimatedFileContractTransactionSetSize)),

		// Bandwidth related fields.
		DownloadBandwidthCost: hes.DownloadBandwidthPrice,
		UploadBandwidthCost:   hes.UploadBandwidthPrice,
	}
	// update the pricetable
	h.staticPriceTables.managedSetCurrent(priceTable)
}

// threadedPruneExpiredPriceTables will expire price tables which have an expiry
// in the past.
//
// Note: threadgroup counter must be inside for loop. If not, calling 'Flush'
// on the threadgroup would deadlock.
func (h *Host) threadedPruneExpiredPriceTables() {
	for {
		func() {
			if err := h.tg.Add(); err != nil {
				return
			}
			defer h.tg.Done()
			h.staticPriceTables.managedPruneExpired()
		}()

		// Block until next cycle.
		select {
		case <-h.tg.StopChan():
			return
		case <-time.After(pruneExpiredRPCPriceTableFrequency):
			continue
		}
	}
}

// newHost returns an initialized Host, taking a set of dependencies as input.
// By making the dependencies an argument of the 'new' call, the host can be
// mocked such that the dependencies can return unexpected errors or unique
// behaviors during testing, enabling easier testing of the failure modes of
// the Host.
func newHost(dependencies modules.Dependencies, smDeps modules.Dependencies, cs modules.ConsensusSet, g modules.Gateway,
	tpool modules.TransactionPool, wallet modules.Wallet, mux *siamux.SiaMux, listenerAddress, persistDir string,
	hostAPIListener net.Listener, checkTokenExpirationFrequency time.Duration) (*Host, error) {
	// Check that all the dependencies were provided.
	if cs == nil {
		return nil, errNilCS
	}
	if g == nil {
		return nil, errNilGateway
	}
	if tpool == nil {
		return nil, errNilTpool
	}
	if wallet == nil {
		return nil, errNilWallet
	}

	var hostAPIPort string
	var err error
	if hostAPIListener == nil {
		// Get free port.
		hostAPIListener, err = net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			return nil, fmt.Errorf("make listener for host api2 port: %w", err)
		}
	}
	hostAPIPort = ":" + strconv.Itoa(hostAPIListener.Addr().(*net.TCPAddr).Port)

	// Create the host object.
	h := &Host{
		cs:                       cs,
		g:                        g,
		tpool:                    tpool,
		wallet:                   wallet,
		staticAlerter:            modules.NewAlerter("host"),
		staticMux:                mux,
		dependencies:             dependencies,
		lockedStorageObligations: make(map[types.FileContractID]*lockedObligation),
		staticPriceTables: &hostPrices{
			guaranteed: make(map[modules.UniqueID]*hostRPCPriceTable),
			staticMinHeap: priceTableHeap{
				heap: make([]*hostRPCPriceTable, 0),
			},
		},
		persistDir: persistDir,
		apiPort:    hostAPIPort,
	}

	// Create MDM.
	h.staticMDM = mdm.New(h)

	// Call stop in the event of a partial startup.
	defer func() {
		if err != nil {
			err = composeErrors(h.tg.Stop(), err)
		}
	}()

	// Create the perist directory if it does not yet exist.
	err = dependencies.MkdirAll(h.persistDir, 0700)
	if err != nil {
		return nil, fmt.Errorf("Could not create directory %v: %w", h.persistDir, err)
	}

	// Initialize the logger, and set up the stop call that will close the
	// logger.
	h.log, err = dependencies.NewLogger(filepath.Join(h.persistDir, logFile))
	if err != nil {
		return nil, fmt.Errorf("Error creating logger: %w", err)
	}

	h.tg.AfterStop(func() {
		err = h.log.Close()
		if err != nil {
			// State of the logger is uncertain, a Println will have to
			// suffice.
			fmt.Println("Error when closing the logger:", err)
		}
	})

	tokenStorageDir := filepath.Join(h.persistDir, tokenStorDir)
	if _, err = os.Stat(tokenStorageDir); os.IsNotExist(err) {
		// Create the token storage directory if it does not yet exist.
		if err = os.Mkdir(tokenStorageDir, 0755); err != nil {
			return nil, fmt.Errorf("error creating token storage directory: %w", err)
		}
	}
	// Add the storage manager to the host, and set up the stop call that will
	// close the storage manager.
	stManager, err := contractmanager.NewCustomContractManager(smDeps, filepath.Join(persistDir, "contractmanager"))
	if err != nil {
		h.log.Println("Could not open the storage manager:", err)
		return nil, fmt.Errorf("error creating contract manager: %w", err)
	}
	h.StorageManager = stManager

	h.tg.AfterStop(func() {
		err = h.StorageManager.Close()
		if err != nil {
			h.log.Println("Could not close storage manager:", err)
			err = fmt.Errorf("error closing the storage manager: %w", err)
		}
	})

	// Initialize token storage.
	// TODO: the current version of tokens storage does not support reverting blocks.
	// If contracts related to `TopUpToken` RPC are reverted, all the tokens resources remain in the storage.
	h.tokenStor, err = tokenstorage.NewTokenStorage(stManager, tokenStorageDir)
	if err != nil {
		return nil, fmt.Errorf("error initializing token storage: %w", err)
	}
	updateTokenSectorsChan := make(chan bool)
	h.tg.AfterStop(func() {
		updateTokenSectorsChan <- true
		err = h.tokenStor.Close(context.Background())
		if err != nil {
			h.log.Errorf("Error when closing token storage: %v", err)
		}
	})
	// Remove sectors from token when token storage resource ends.
	go h.tokenStor.CheckExpiration(checkTokenExpirationFrequency, updateTokenSectorsChan)

	// Load the prior persistence structures, and configure the host to save
	// before shutting down.
	err = h.load()
	if err != nil {
		return nil, fmt.Errorf("error loading persistence: %w", err)
	}
	h.tg.AfterStop(func() {
		err = h.saveSync()
		if err != nil {
			h.log.Errorf("Error saving host upon shutdown: %v", err)
		}
	})

	// Add the account manager subsystem
	h.staticAccountManager, err = h.newAccountManager()
	if err != nil {
		return nil, fmt.Errorf("error creating account manager: %w", err)
	}

	atomic.StoreUint64(&h.scheduledAuditBlockheight, uint64(h.blockHeight))

	// Subscribe to the consensus set.
	err = h.initConsensusSubscription()
	if err != nil {
		return nil, fmt.Errorf("error subscribing to consensus: %w", err)
	}

	// Create bandwidth monitor
	h.staticMonitor = connmonitor.NewMonitor()

	// Initialize the networking. We need to hold the lock while doing so since
	// the previous load subscribed the host to the consensus set.
	h.mu.Lock()
	err = h.initNetworking(listenerAddress)
	h.mu.Unlock()
	if err != nil {
		h.log.Errorf("Could not initialize host networking: %v", err)
		return nil, err
	}

	// Initialize the RPC price table
	h.managedUpdatePriceTable()

	// Ensure the expired RPC tables get pruned as to not leak memory
	go h.threadedPruneExpiredPriceTables()

	//	Initialize and run host API
	hostApi := api.NewAPI(h.tokenStor, h.secretKey, h)
	err = hostApi.Start(hostAPIListener)
	if err != nil {
		err = fmt.Errorf("error starting host api: %w", err)
		h.log.Errorln(err)
		return nil, err
	}
	h.tg.AfterStop(func() {
		err = hostApi.Close()
		if err != nil {
			h.log.Errorln("Could not close host API:", err)
			err = fmt.Errorf("error closing host API: %w", err)
		}
	})
	return h, nil
}

// New returns an initialized Host.
func New(cs modules.ConsensusSet, g modules.Gateway, tpool modules.TransactionPool, wallet modules.Wallet, mux *siamux.SiaMux, address, persistDir string, hostAPIListener net.Listener, checkTokenExpirationFrequency time.Duration) (*Host, error) {
	return newHost(modules.ProdDependencies, new(modules.ProductionDependencies), cs, g, tpool, wallet, mux, address, persistDir, hostAPIListener, checkTokenExpirationFrequency)
}

// NewCustomHost returns an initialized Host using the provided dependencies.
func NewCustomHost(deps modules.Dependencies, cs modules.ConsensusSet, g modules.Gateway, tpool modules.TransactionPool, wallet modules.Wallet, mux *siamux.SiaMux, address, persistDir string, hostAPIListener net.Listener, checkTokenExpirationFrequency time.Duration) (*Host, error) {
	return newHost(deps, new(modules.ProductionDependencies), cs, g, tpool, wallet, mux, address, persistDir, hostAPIListener, checkTokenExpirationFrequency)
}

// NewCustomTestHost allows passing in both host dependencies and storage
// manager dependencies. Used solely for testing purposes, to allow dependency
// injection into the host's submodules.
func NewCustomTestHost(deps modules.Dependencies, smDeps modules.Dependencies, cs modules.ConsensusSet, g modules.Gateway, tpool modules.TransactionPool, wallet modules.Wallet, mux *siamux.SiaMux, address, persistDir string, hostAPIListener net.Listener, checkTokenExpirationFrequency time.Duration) (*Host, error) {
	return newHost(deps, smDeps, cs, g, tpool, wallet, mux, address, persistDir, hostAPIListener, checkTokenExpirationFrequency)
}

// Close shuts down the host.
func (h *Host) Close() error {
	return h.tg.Stop()
}

// Announcement returns host announcement.
// For use as ghost.
func (h *Host) Announcement() []byte {
	ann := encoding.Marshal(modules.HostAnnouncement{
		Specifier:  modules.PrefixHostAnnouncement,
		NetAddress: modules.NetAddress(h.listener.Addr().String()),
		PublicKey:  h.PublicKey(),
	})
	sig := crypto.SignHash(crypto.HashBytes(ann), h.secretKey)
	return append(ann, sig[:]...)
}

// ExternalSettings returns the hosts external settings. These values cannot be
// set by the user (host is configured through InternalSettings), and are the
// values that get displayed to other hosts on the network.
func (h *Host) ExternalSettings() modules.HostExternalSettings {
	err := h.tg.Add()
	if err != nil {
		build.Critical("Call to ExternalSettings after close")
	}
	defer h.tg.Done()
	return h.ManagedExternalSettings()
}

// BandwidthCounters returns the Hosts's upload and download bandwidth
func (h *Host) BandwidthCounters() (uint64, uint64, time.Time, error) {
	if err := h.tg.Add(); err != nil {
		return 0, 0, time.Time{}, err
	}
	defer h.tg.Done()

	// Get the bandwidth usage for RHP1 & RHP2 connections.
	readBytes, writeBytes := h.staticMonitor.Counts()

	// Get the bandwidth usage for RHP3 connections. Unfortunately we can't just
	// wrap the siamux streams since that wouldn't give us the raw data sent over
	// the TCP connection. Since we want this to be as accurate as possible, we
	// use the `Limit` method on the streams before closing them to get the
	// accurate amount of data sent and received. This includes overhead such as
	// frame headers and encryption.
	readBytes += atomic.LoadUint64(&h.atomicStreamDownload)
	writeBytes += atomic.LoadUint64(&h.atomicStreamUpload)

	startTime := h.staticMonitor.StartTime()
	return writeBytes, readBytes, startTime, nil
}

// WorkingStatus returns the working state of the host, where working is
// defined as having received more than workingStatusThreshold settings calls
// over the period of workingStatusFrequency.
func (h *Host) WorkingStatus() modules.HostWorkingStatus {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.workingStatus
}

// ConnectabilityStatus returns the connectability state of the host, whether
// the host can connect to itself on its configured netaddress.
func (h *Host) ConnectabilityStatus() modules.HostConnectabilityStatus {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.connectabilityStatus
}

// FinancialMetrics returns information about the financial commitments,
// rewards, and activities of the host.
func (h *Host) FinancialMetrics() modules.HostFinancialMetrics {
	err := h.tg.Add()
	if err != nil {
		build.Critical("Call to FinancialMetrics after close")
	}
	defer h.tg.Done()
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.financialMetrics
}

// PublicKey returns the public key of the host that is used to facilitate
// relationships between the host and renter.
func (h *Host) PublicKey() types.SiaPublicKey {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.publicKey
}

// SetInternalSettings updates the host's internal HostInternalSettings object.
func (h *Host) SetInternalSettings(settings modules.HostInternalSettings) error {
	err := h.tg.Add()
	if err != nil {
		return err
	}
	defer h.tg.Done()

	h.mu.Lock()
	// By updating the internal settings the user might influence the host's
	// price table, we defer a call to update the price table to ensure it
	// reflects the updated settings.
	defer h.managedUpdatePriceTable()
	defer h.mu.Unlock()

	// The host should not be accepting file contracts if it does not have an
	// unlock hash.
	if settings.AcceptingContracts {
		err := h.checkUnlockHash()
		if err != nil {
			return errors.New("internal settings not updated, no unlock hash: " + err.Error())
		}
	}

	if settings.NetAddress != "" {
		err := settings.NetAddress.IsValid()
		if err != nil {
			return errors.New("internal settings not updated, invalid NetAddress: " + err.Error())
		}
	}

	// Check if the net address for the host has changed. If it has, and it's
	// not equal to the auto address, then the host is going to need to make
	// another blockchain announcement.
	if h.settings.NetAddress != settings.NetAddress && settings.NetAddress != h.autoAddress {
		h.announced = false
	}

	h.settings = settings
	h.revisionNumber++

	// The locked storage collateral was altered, we potentially want to
	// unregister the insufficient collateral budget alert
	h.tryUnregisterInsufficientCollateralBudgetAlert()

	err = h.saveSync()
	if err != nil {
		return errors.New("internal settings updated, but failed saving to disk: " + err.Error())
	}
	return nil
}

// InternalSettings returns the settings of a host.
func (h *Host) InternalSettings() modules.HostInternalSettings {
	err := h.tg.Add()
	if err != nil {
		return modules.HostInternalSettings{}
	}
	defer h.tg.Done()
	return h.managedInternalSettings()
}

// BlockHeight returns the host's current blockheight.
func (h *Host) BlockHeight() types.BlockHeight {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.blockHeight
}

// ManagedExternalSettings returns the host's external settings. These values
// cannot be set by the user (host is configured through InternalSettings), and
// are the values that get displayed to other hosts on the network.
func (h *Host) ManagedExternalSettings() modules.HostExternalSettings {
	_, maxFee := h.tpool.FeeEstimation()
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.externalSettings(maxFee)
}
