package types

// constants.go contains the ScPrime constants. Depending on which build tags are
// used, the constants will be initialized to different values.
//
// CONTRIBUTE: We don't have way to check that the non-test constants are all
// sane, plus we have no coverage for them.

import (
	"math"
	"math/big"
	"time"

	"gitlab.com/scpcorp/ScPrime/build"
)

var (
	// ASICHardforkHeight is the height at which the hardfork targeting
	// selected ASICs was activated.
	ASICHardforkHeight BlockHeight

	// ASICHardforkTotalTarget is the initial target after the ASIC hardfork.
	// The actual target at ASICHardforkHeight is replaced with this value in
	// order to prevent intolerably slow block times post-fork.
	ASICHardforkTotalTarget Target

	// ASICHardforkTotalTime is the initial total time after the ASIC
	// hardfork. The actual total time at ASICHardforkHeight is replaced with
	// this value in order to prevent intolerably slow block times post-fork.
	ASICHardforkTotalTime int64

	// ASICHardforkFactor is the factor by which the hashrate of targeted
	// ASICs will be reduced.
	ASICHardforkFactor = uint64(1)

	// ASICHardforkReplayProtectionPrefix is a byte that prefixes
	// SiacoinInputs and SiafundInputs when calculating SigHashes to protect
	// against replay attacks.
	ASICHardforkReplayProtectionPrefix = []byte(nil)

	// BlockFrequency is the desired number of seconds that
	// should elapse, on average, between successive Blocks.
	BlockFrequency BlockHeight
	// BlockSizeLimit is the maximum size of a binary-encoded Block
	// that is permitted by the consensus rules.
	BlockSizeLimit = uint64(2e6)
	// BlocksPerHour is the number of blocks expected to be mined per hour.
	BlocksPerHour = BlockHeight(6)
	// BlocksPerDay is the number of blocks expected to be mined per day.
	BlocksPerDay = 24 * BlocksPerHour
	// BlocksPerWeek is the number of blocks expected to be mined per week.
	BlocksPerWeek = 7 * BlocksPerDay
	// BlocksPerMonth is the number of blocks expected to be mined per month.
	BlocksPerMonth = 30 * BlocksPerDay
	// BlocksPerYear is the number of blocks expected to be mined per year.
	BlocksPerYear = 365 * BlocksPerDay
	// BurnAddressBlockHeight is the height at which the dev fund will be burnt
	// instead of being claimed by the dev fund. Setting this value to 0 will
	// prevent the dev fund from being burnt at any height.
	BurnAddressBlockHeight = BlockHeight(105000)
	// BurnAddressUnlockHash is the unlock hash for where to send coins to burn.
	BurnAddressUnlockHash = UnlockHashFromAddrStr("000000000000000000000000000000000000000000000000000000000000000089eb0d6a8a69")
	// DevFundEnabled is a boolean that when set to true will enable the ability to
	// configure a dev fund
	DevFundEnabled = true
	// DevFundInitialBlockHeight is the height at which the dev fund became mandatory
	DevFundInitialBlockHeight = BlockHeight(1)
	// DevFundDecayStartBlockHeight is the height at which the DevFundInitialPercentage
	// begins to linearly decay to the DevFundFinalPercentage
	DevFundDecayStartBlockHeight = BlockHeight(30000)
	// DevFundDecayEndBlockHeight is the height at which the DevFundInitialPercentage
	// has fully decayed to the DevFundFinalPercentage
	DevFundDecayEndBlockHeight = BlockHeight(105000)
	// DevFundInitialPercentage is the initial percentage of the block reward that is
	// sent to the DevFundUnlockHash before any dev fund percentage decay happens
	DevFundInitialPercentage = uint64(20)
	// DevFundFinalPercentage is the final percentage of the block reward that is sent
	//  to the DevFundUnlockHash after the dev fund percentage is fully decayed
	DevFundFinalPercentage = uint64(10)
	// DevFundUnlockHash is the unlock hash for the dev fund subsidy
	// Do not set this to the Zero address as doing so will cause the test that
	// verifies that a dev fee is set to fail
	DevFundUnlockHash = UnlockHashFromAddrStr("aefe0af2713c112ba4d10dee7753726e5c4de3f237ea455151342615c95d0e797d7a8cce7b05")

	// EndOfTime is value to be used when a date in the future is needed for
	// validation
	EndOfTime = time.Unix(0, math.MaxInt64)

	// ExtremeFutureThreshold is a temporal limit beyond which Blocks are
	// discarded by the consensus rules. When incoming Blocks are processed, their
	// Timestamp is allowed to exceed the processor's current time by a small amount.
	// But if the Timestamp is further into the future than ExtremeFutureThreshold,
	// the Block is immediately discarded.
	ExtremeFutureThreshold Timestamp
	// FutureThreshold is a temporal limit beyond which Blocks are
	// discarded by the consensus rules. When incoming Blocks are processed, their
	// Timestamp is allowed to exceed the processor's current time by no more than
	// FutureThreshold. If the excess duration is larger than FutureThreshold, but
	// smaller than ExtremeFutureThreshold, the Block may be held in memory until
	// the Block's Timestamp exceeds the current time by less than FutureThreshold.
	FutureThreshold Timestamp
	// GenesisBlock is the first block of the block chain
	GenesisBlock Block

	// GenesisID is used in many places. Calculating it once saves lots of
	// redundant computation.
	GenesisID BlockID
	// GenesisAirdropAllocation is the output creating the initial coins allocated
	// for the airdrop at network launch
	GenesisAirdropAllocation []SiacoinOutput

	// GenesisSiacoinAllocation is the set of SiacoinOutputs created in the Genesis
	// block
	GenesisSiacoinAllocation []SiacoinOutput
	// GenesisSiafundAllocation is the set of SiafundOutputs created in the Genesis
	// block.
	GenesisSiafundAllocation []SiafundOutput
	// ForkedGenesisSiafundAllocation is the set of SiafundOutputs created in the Genesis
	// block.
	ForkedGenesisSiafundAllocation []SiafundOutput
	// SiafundHardforkAllocation is allocation of new Siafunds at various hardforks.
	SiafundHardforkAllocation map[BlockHeight][]SiafundOutput
	// GenesisTimestamp is the timestamp when genesis block was mined
	GenesisTimestamp Timestamp
	// InitialCoinbase is the coinbase reward of the Genesis block.
	InitialCoinbase = uint64(300e3)
	// AirdropCommunityValue is the total amount of coins the community members will split
	// from the genesis block airdrop.
	AirdropCommunityValue = NewCurrency64(10000000000).Mul(SiacoinPrecision)
	// AirdropPoolValue is the total amount of coins the pools get
	// airdrop so that they can pay out miners in the first 144 blocks
	AirdropPoolValue = NewCurrency64(51840000).Mul(SiacoinPrecision)
	// AirdropNebulousLabsValue is a gift to the NebulousLabs Team to acknowledge all their
	// effort and hard work. THANK YOU!
	AirdropNebulousLabsValue = NewCurrency64(300000000).Mul(SiacoinPrecision)
	// AirdropSiaPrimeValue is the total amount of coins ScPrime gets to help bootstrap
	// expenses
	AirdropSiaPrimeValue = NewCurrency64(200000000).Mul(SiacoinPrecision)
	// MaturityDelay specifies the number of blocks that a maturity-required output
	// is required to be on hold before it can be spent on the blockchain.
	// Outputs are maturity-required if they are highly likely to be altered or
	// invalidated in the event of a small reorg. One example is the block reward,
	// as a small reorg may invalidate the block reward. Another example is a siafund
	// payout, as a tiny reorg may change the value of the payout, and thus invalidate
	// any transactions spending the payout. File contract payouts also are subject to
	// a maturity delay.
	MaturityDelay BlockHeight
	// MaxTargetAdjustmentDown restrict how much the block difficulty is allowed to
	// change in a single step, which is important to limit the effect of difficulty
	// raising and lowering attacks.
	MaxTargetAdjustmentDown *big.Rat
	// MaxTargetAdjustmentUp restrict how much the block difficulty is allowed to
	// change in a single step, which is important to limit the effect of difficulty
	// raising and lowering attacks.
	MaxTargetAdjustmentUp *big.Rat
	// MedianTimestampWindow tells us how many blocks to look back when calculating
	// the median timestamp over the previous n blocks. The timestamp of a block is
	// not allowed to be less than or equal to the median timestamp of the previous n
	// blocks, where for ScPrime this number is typically 11.
	MedianTimestampWindow = uint64(11)
	// MinimumCoinbase is the minimum coinbase reward for a block.
	// The coinbase decreases in each block after the Genesis block,
	// but it will not decrease past MinimumCoinbase.
	MinimumCoinbase uint64

	// Oak hardfork constants. Oak is the name of the difficulty algorithm for
	// ScPrime following a hardfork at block 135e3.

	// OakDecayDenom is the denominator for how much the total timestamp is decayed
	// each step.
	OakDecayDenom int64
	// OakDecayNum is the numerator for how much the total timestamp is decayed each
	// step.
	OakDecayNum int64
	// OakHardforkBlock is the height at which the hardfork to switch to the oak
	// difficulty adjustment algorithm is triggered.
	OakHardforkBlock BlockHeight
	// OakHardforkFixBlock is the height at which the hardfork to switch from the broken
	// oak difficulty adjustment algorithm to the fixed oak difficulty adjustment
	// algorithm is triggered.
	OakHardforkFixBlock BlockHeight
	// OakHardforkTxnSizeLimit is the maximum size allowed for a transaction, a change
	// which I believe was implemented simultaneously with the oak hardfork.
	OakHardforkTxnSizeLimit = uint64(64e3) // 64 KB
	// OakMaxBlockShift is the maximum number of seconds that the oak algorithm will shift
	// the difficulty.
	OakMaxBlockShift int64
	// OakMaxDrop is the drop is the maximum amount that the difficulty will drop each block.
	OakMaxDrop *big.Rat
	// OakMaxRise is the maximum amount that the difficulty will rise each block.
	OakMaxRise *big.Rat

	// RootDepth is the cumulative target of all blocks. The root depth is essentially
	// the maximum possible target, there have been no blocks yet, so there is no
	// cumulated difficulty yet.
	RootDepth = Target{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}
	// RootTarget is the target for the genesis block - basically how much work needs
	// to be done in order to mine the first block. The difficulty adjustment algorithm
	// takes over from there.
	RootTarget Target
	// SiacoinPrecision is the number of base units in a siacoin. This constant is used
	// for mining rewards calculation and supported for compatibility with
	// existing 3rd party intergations.
	// DEPRECATED: Since February 2020 one scprimecoin equals 10^27 Hastings
	// Use the types.ScPrimecoinPrecision constant.
	//
	// The base unit for Bitcoin is called a satoshi. We call 10^8 satoshis a bitcoin,
	// even though the code itself only ever works with satoshis.
	SiacoinPrecision = NewCurrency(new(big.Int).Exp(big.NewInt(10), big.NewInt(24), nil))
	// ScPrimecoinPrecision is the number of base units in a scprimecoin that is used
	// by clients (1 SCP = 10^27 H).
	ScPrimecoinPrecision = NewCurrency(new(big.Int).Exp(big.NewInt(10), big.NewInt(27), nil))
	// OldSiafundCount is the total number of Siafunds in existence before the SPF hardfork.
	OldSiafundCount = NewCurrency64(10000)
	// NewSiafundCount is the total number of Siafunds in existence after the SPF hardfork.
	NewSiafundCount = NewCurrency64(30000)
	// NewerSiafundCount is the total number of Siafunds in existence after the second SPF hardfork.
	NewerSiafundCount = NewCurrency64(200000000)
	// FirstSiafundMul is multiplier for percentage of siacoins that is taxed from FileContracts
	// before the first SPF hardfork.
	FirstSiafundMul = int64(39)
	// FirstSiafundDiv is divider for percentage of siacoins that is taxed from FileContracts
	// before the first SPF hardfork.
	FirstSiafundDiv = int64(1000)
	// FirstSiafundPortion is the percentage of siacoins that is taxed from FileContracts before the first SPF hardfork.
	FirstSiafundPortion = big.NewRat(FirstSiafundMul, FirstSiafundDiv)
	// SecondSiafundMul is multiplier for percentage of siacoins that is taxed from FileContracts
	// between the first and the second SPF hardforks.
	SecondSiafundMul = int64(150)
	// SecondSiafundDiv is divider for percentage of siacoins that is taxed from FileContracts
	// between the first and the second SPF hardforks.
	SecondSiafundDiv = int64(1000)
	// SecondSiafundPortion is the percentage of siacoins that is taxed from FileContracts between the first and the second SPF hardforks.
	SecondSiafundPortion = big.NewRat(SecondSiafundMul, SecondSiafundDiv)
	// ThirdSiafundMul is multiplier for percentage of siacoins that is taxed from FileContracts
	// after the second SPF hardfork.
	ThirdSiafundMul = int64(100)
	// ThirdSiafundDiv is divider for percentage of siacoins that is taxed from FileContracts
	// after the second SPF hardfork.
	ThirdSiafundDiv = int64(1000)
	// ThirdSiafundPortion is the percentage of siacoins that is taxed from FileContracts after the second SPF hardfork.
	ThirdSiafundPortion = big.NewRat(ThirdSiafundMul, ThirdSiafundDiv)
	// TargetWindow is the number of blocks to look backwards when determining how much
	// time has passed vs. how many blocks have been created. It's only used in the old,
	// broken difficulty adjustment algorithm.
	TargetWindow BlockHeight
)

var (
	// TaxHardforkHeight is the height at which the tax hardfork occurred.
	TaxHardforkHeight = build.Select(build.Var{
		Dev:      BlockHeight(10),
		Standard: BlockHeight(21e3),
		Testing:  BlockHeight(10),
	}).(BlockHeight)

	// SpfAirdropHeight is the height of SPF airdrop.
	SpfAirdropHeight = build.Select(build.Var{
		Dev:      BlockHeight(20),
		Standard: BlockHeight(7200),
		Testing:  BlockHeight(7200),
	}).(BlockHeight)

	// SpfHardforkHeight is the height of SPF hardfork.
	SpfHardforkHeight = build.Select(build.Var{
		Dev:      BlockHeight(100),
		Standard: BlockHeight(54550),
		Testing:  BlockHeight(10000),
	}).(BlockHeight)

	// SpfSecondHardforkHeight is the height of second SPF hardfork.
	SpfSecondHardforkHeight = build.Select(build.Var{
		Dev:      BlockHeight(200),
		Standard: BlockHeight(109000),
		Testing:  BlockHeight(20000),
	}).(BlockHeight)
)

// IsSpfHardfork returns true when one of Spf hardforks happens at given height.
func IsSpfHardfork(height BlockHeight) bool {
	if height == SpfHardforkHeight {
		return true
	}
	if height == SpfSecondHardforkHeight {
		return true
	}
	return false
}

// SiafundCount returns the total number of Siafunds by height.
func SiafundCount(height BlockHeight) Currency {
	if height > SpfSecondHardforkHeight {
		return NewerSiafundCount
	}
	if height > SpfHardforkHeight {
		return NewSiafundCount
	}
	return OldSiafundCount
}

// SiafundPortion returns SPF percentage by height.
func SiafundPortion(height BlockHeight) *big.Rat {
	if height > SpfSecondHardforkHeight {
		return ThirdSiafundPortion
	}
	if height > SpfHardforkHeight {
		return SecondSiafundPortion
	}
	return FirstSiafundPortion
}

// SiafundMul returns SPF percentage multiplier by height.
func SiafundMul(height BlockHeight) int64 {
	if height > SpfSecondHardforkHeight {
		return ThirdSiafundMul
	}
	if height > SpfHardforkHeight {
		return SecondSiafundMul
	}
	return FirstSiafundMul
}

// SiafundDiv returns SPF percentage divider by height.
func SiafundDiv(height BlockHeight) int64 {
	if height > SpfSecondHardforkHeight {
		return ThirdSiafundDiv
	}
	if height > SpfHardforkHeight {
		return SecondSiafundDiv
	}
	return FirstSiafundDiv
}

// scanAddress scans a types.UnlockHash from a string.
func scanAddress(addrStr string) (addr UnlockHash, err error) {
	err = addr.LoadString(addrStr)
	if err != nil {
		return UnlockHash{}, err
	}
	return addr, nil
}

// UnlockHashFromAddrStr convert string address to UnlockHash
func UnlockHashFromAddrStr(addrStr string) (addr UnlockHash) {
	dest, err := scanAddress(addrStr)
	if err != nil {
		return UnlockHash{}
	}
	return dest
}

// init checks which build constant is in place and initializes the variables
// accordingly.
func init() {
	if build.Release == "dev" {
		// 'dev' settings are for small developer testnets, usually on the same
		// computer. Settings are slow enough that a small team of developers
		// can coordinate their actions over a the developer testnets, but fast
		// enough that there isn't much time wasted on waiting for things to
		// happen.
		ASICHardforkHeight = math.MaxInt64
		ASICHardforkTotalTarget = Target{0, 0, 0, 8}
		ASICHardforkTotalTime = 800

		BlockFrequency = 12                      // 12 seconds: slow enough for developers to see ~each block, fast enough that blocks don't waste time.
		MaturityDelay = 10                       // 60 seconds before a delayed output matures.
		GenesisTimestamp = Timestamp(1528293910) // Change as necessary.
		RootTarget = Target{0, 0, 2}             // Standard developer CPUs will be able to mine blocks with the race library activated.

		TargetWindow = 20                              // Difficulty is adjusted based on prior 20 blocks.
		MaxTargetAdjustmentUp = big.NewRat(120, 100)   // Difficulty adjusts quickly.
		MaxTargetAdjustmentDown = big.NewRat(100, 120) // Difficulty adjusts quickly.
		FutureThreshold = 2 * 60                       // 2 minutes.
		ExtremeFutureThreshold = 4 * 60                // 4 minutes.

		MinimumCoinbase = 30e3

		OakHardforkBlock = 100
		OakHardforkFixBlock = 105
		OakDecayNum = 985
		OakDecayDenom = 1000
		OakMaxBlockShift = 3
		OakMaxRise = big.NewRat(102, 100)
		OakMaxDrop = big.NewRat(100, 102)

		GenesisAirdropAllocation = []SiacoinOutput{
			{
				Value:      AirdropCommunityValue,
				UnlockHash: UnlockHashFromAddrStr("436890aacc53f93f9cc4538d9b4abba27dd5be6ff8a064fae7b78a67809db5e210819ffc4a21"),
			},
			{
				Value:      AirdropPoolValue,
				UnlockHash: UnlockHashFromAddrStr("78054218b7d0bc04929e5a3e6a2ac5fed29b98898cba3d740dd31a1aae6e8c8b3ce7467d4e8f"),
			},
			{
				Value:      AirdropNebulousLabsValue,
				UnlockHash: UnlockHashFromAddrStr("7d0c44f7664e2d34e53efde0661a6f628ec9264785ae8e3cd7c973e8d190c3c97b5e3ecbc567"),
			},
			{
				Value:      AirdropSiaPrimeValue,
				UnlockHash: UnlockHashFromAddrStr("aefe0af2713c112ba4d10dee7753726e5c4de3f237ea455151342615c95d0e797d7a8cce7b05"),
			},
		}

		GenesisSiafundAllocation = []SiafundOutput{
			{
				Value:      NewCurrency64(2000),
				UnlockHash: UnlockHash{214, 166, 197, 164, 29, 201, 53, 236, 106, 239, 10, 158, 127, 131, 20, 138, 63, 221, 230, 16, 98, 247, 32, 77, 210, 68, 116, 12, 241, 89, 27, 223},
			},
			{
				Value:      NewCurrency64(7000),
				UnlockHash: UnlockHash{209, 246, 228, 60, 248, 78, 242, 110, 9, 8, 227, 248, 225, 216, 163, 52, 142, 93, 47, 176, 103, 41, 137, 80, 212, 8, 132, 58, 241, 189, 2, 17},
			},
			{
				Value:      NewCurrency64(1000),
				UnlockHash: UnlockConditions{}.UnlockHash(),
			},
		}

		ForkedGenesisSiafundAllocation = []SiafundOutput{
			{
				Value:      NewCurrency64(10000),
				UnlockHash: UnlockHashFromAddrStr("436890aacc53f93f9cc4538d9b4abba27dd5be6ff8a064fae7b78a67809db5e210819ffc4a21"),
			},
		}
		SiafundHardforkAllocation = map[BlockHeight][]SiafundOutput{
			SpfHardforkHeight: {
				{
					Value:      NewCurrency64(20000),
					UnlockHash: UnlockHashFromAddrStr("c25a37aa55baf3131e4d9335373338018d71f35bc22cd55d6d983394611d09a1fbb3f6edd5cc"),
				},
			},
			SpfSecondHardforkHeight: {
				{
					Value:      NewCurrency64(200000000 - 30000),
					UnlockHash: UnlockHashFromAddrStr("f3f512b45c24b531e571d59d99b804b19244cecf1da8117217101fa722e40a18917c91ebb080"),
				},
			},
		}
	} else if build.Release == "testing" {
		// 'testing' settings are for automatic testing, and create much faster
		// environments than a human can interact with.
		ASICHardforkHeight = math.MaxInt64
		ASICHardforkTotalTarget = Target{255, 255}
		ASICHardforkTotalTime = 10e3

		BlockFrequency = 1 // As fast as possible
		MaturityDelay = 3
		GenesisTimestamp = CurrentTimestamp() - 1e6
		RootTarget = Target{128} // Takes an expected 2 hashes; very fast for testing but still probes 'bad hash' code.

		// A restrictive difficulty clamp prevents the difficulty from climbing
		// during testing, as the resolution on the difficulty adjustment is
		// only 1 second and testing mining should be happening substantially
		// faster than that.
		TargetWindow = 200
		MaxTargetAdjustmentUp = big.NewRat(10001, 10000)
		MaxTargetAdjustmentDown = big.NewRat(9999, 10000)
		FutureThreshold = 3        // 3 seconds
		ExtremeFutureThreshold = 6 // 6 seconds

		MinimumCoinbase = 299990 // Minimum coinbase is hit after 10 blocks to make testing minimum-coinbase code easier.

		// Do not let the difficulty change rapidly - blocks will be getting
		// mined far faster than the difficulty can adjust to.
		OakHardforkBlock = 20
		OakHardforkFixBlock = 23
		OakDecayNum = 9999
		OakDecayDenom = 10e3
		OakMaxBlockShift = 3
		OakMaxRise = big.NewRat(10001, 10e3)
		OakMaxDrop = big.NewRat(10e3, 10001)

		GenesisAirdropAllocation = []SiacoinOutput{
			{
				Value:      AirdropCommunityValue,
				UnlockHash: UnlockHashFromAddrStr("436890aacc53f93f9cc4538d9b4abba27dd5be6ff8a064fae7b78a67809db5e210819ffc4a21"),
			},
			{
				Value:      AirdropPoolValue,
				UnlockHash: UnlockHashFromAddrStr("78054218b7d0bc04929e5a3e6a2ac5fed29b98898cba3d740dd31a1aae6e8c8b3ce7467d4e8f"),
			},
			{
				Value:      AirdropNebulousLabsValue,
				UnlockHash: UnlockHashFromAddrStr("7d0c44f7664e2d34e53efde0661a6f628ec9264785ae8e3cd7c973e8d190c3c97b5e3ecbc567"),
			},
			{
				Value:      AirdropSiaPrimeValue,
				UnlockHash: UnlockHashFromAddrStr("aefe0af2713c112ba4d10dee7753726e5c4de3f237ea455151342615c95d0e797d7a8cce7b05"),
			},
		}

		GenesisSiafundAllocation = []SiafundOutput{
			{
				Value:      NewCurrency64(2000),
				UnlockHash: UnlockHash{214, 166, 197, 164, 29, 201, 53, 236, 106, 239, 10, 158, 127, 131, 20, 138, 63, 221, 230, 16, 98, 247, 32, 77, 210, 68, 116, 12, 241, 89, 27, 223},
			},
			{
				Value:      NewCurrency64(7000),
				UnlockHash: UnlockHash{209, 246, 228, 60, 248, 78, 242, 110, 9, 8, 227, 248, 225, 216, 163, 52, 142, 93, 47, 176, 103, 41, 137, 80, 212, 8, 132, 58, 241, 189, 2, 17},
			},
			{
				Value:      NewCurrency64(1000),
				UnlockHash: UnlockConditions{}.UnlockHash(),
			},
		}

		ForkedGenesisSiafundAllocation = []SiafundOutput{
			{
				Value:      NewCurrency64(10000),
				UnlockHash: UnlockHashFromAddrStr("436890aacc53f93f9cc4538d9b4abba27dd5be6ff8a064fae7b78a67809db5e210819ffc4a21"),
			},
		}
		SiafundHardforkAllocation = map[BlockHeight][]SiafundOutput{
			SpfHardforkHeight: {
				{
					Value:      NewCurrency64(20000),
					UnlockHash: UnlockHashFromAddrStr("c25a37aa55baf3131e4d9335373338018d71f35bc22cd55d6d983394611d09a1fbb3f6edd5cc"),
				},
			},
			SpfSecondHardforkHeight: {
				{
					Value:      NewCurrency64(200000000 - 30000),
					UnlockHash: UnlockHashFromAddrStr("f3f512b45c24b531e571d59d99b804b19244cecf1da8117217101fa722e40a18917c91ebb080"),
				},
			},
		}
	} else if build.Release == "standard" {
		// 'standard' settings are for the full network. They are slow enough
		// that the network is secure in a real-world byzantine environment.

		// A hardfork height of max int64 was chosen to clarify that the we
		// expect the hardfork to never happen on the ScPrime blockchain.
		// A total time of 120,000 is chosen because that represents the total
		// time elapsed at a perfect equilibrium, indicating a visible average
		// block time that perfectly aligns with what is expected. A total
		// target of 67 leading zeroes is chosen because that aligns with the
		// amount of hashrate that we expect to be on the network after the
		// hardfork.
		ASICHardforkHeight = math.MaxInt64
		ASICHardforkTotalTarget = Target{0, 0, 0, 0, 0, 0, 0, 0, 32}
		ASICHardforkTotalTime = 120e3

		// A block time of 1 block per 10 minutes is chosen to follow Bitcoin's
		// example. The security lost by lowering the block time is not
		// insignificant, and the convenience gained by lowering the blocktime
		// even down to 90 seconds is not significant. I do feel that 10
		// minutes could even be too short, but it has worked well for Bitcoin.
		BlockFrequency = 600

		// Payouts take 1 day to mature. This is to prevent a class of double
		// spending attacks parties unintentionally spend coins that will stop
		// existing after a blockchain reorganization. There are multiple
		// classes of payouts in ScPrime that depend on a previous block - if that
		// block changes, then the output changes and the previously existing
		// output ceases to exist. This delay stops both unintentional double
		// spending and stops a small set of long-range mining attacks.
		MaturityDelay = 144

		// The genesis timestamp is set to June 6th, because that is when the
		// 100-block developer premine started. The trailing zeroes are a
		// bonus, and make the timestamp easier to memorize.
		GenesisTimestamp = Timestamp(1540955779)

		// The RootTarget was set such that the developers could reasonable
		// premine 100 blocks in a day. It was known to the developers at launch
		// this this was at least one and perhaps two orders of magnitude too
		// small.
		RootTarget = Target{0, 0, 0, 0, 0, 0, 2}

		// When the difficulty is adjusted, it is adjusted by looking at the
		// timestamp of the 1000th previous block. This minimizes the abilities
		// of miners to attack the network using rogue timestamps.
		TargetWindow = 1e3

		// The difficulty adjustment is clamped to 2.5x every 500 blocks. This
		// corresponds to 6.25x every 2 weeks, which can be compared to
		// Bitcoin's clamp of 4x every 2 weeks. The difficulty clamp is
		// primarily to stop difficulty raising attacks. ScPrime's safety margin is
		// similar to Bitcoin's despite the looser clamp because ScPrime's
		// difficulty is adjusted four times as often. This does result in
		// greater difficulty oscillation, a tradeoff that was chosen to be
		// acceptable due to ScPrime's more vulnerable position as an altcoin.
		MaxTargetAdjustmentUp = big.NewRat(25, 10)
		MaxTargetAdjustmentDown = big.NewRat(10, 25)

		// Blocks will not be accepted if their timestamp is more than 3 hours
		// into the future, but will be accepted as soon as they are no longer
		// 3 hours into the future. Blocks that are greater than 5 hours into
		// the future are rejected outright, as it is assumed that by the time
		// 2 hours have passed, those blocks will no longer be on the longest
		// chain. Blocks cannot be kept forever because this opens a DoS
		// vector.
		FutureThreshold = 3 * 60 * 60        // 3 hours.
		ExtremeFutureThreshold = 5 * 60 * 60 // 5 hours.

		// The minimum coinbase is set to 10,000. Because the coinbase
		// decreases by 1 every time, it means that ScPrime's coinbase will have an
		// increasingly potent dropoff for about 5 years, until inflation more
		// or less permanently settles around 2%.
		MinimumCoinbase = 10e3

		// The oak difficulty adjustment hardfork is set to trigger at block
		// 135,000, which is just under 6 months after the hardfork was first
		// released as beta software to the network. This hopefully gives
		// everyone plenty of time to upgrade and adopt the hardfork, while also
		// being earlier than the most optimistic shipping dates for the miners
		// that would otherwise be very disruptive to the network.
		//
		// There was a bug in the original Oak hardfork that had to be quickly
		// followed up with another fix. The height of that fix is the
		// OakHardforkFixBlock.
		OakHardforkBlock = 1500
		OakHardforkFixBlock = 1500

		// The decay is kept at 995/1000, or a decay of about 0.5% each block.
		// This puts the halflife of a block's relevance at about 1 day. This
		// allows the difficulty to adjust rapidly if the hashrate is adjusting
		// rapidly, while still keeping a relatively strong insulation against
		// random variance.
		OakDecayNum = 995
		OakDecayDenom = 1e3

		// The block shift determines the most that the difficulty adjustment
		// algorithm is allowed to shift the target block time. With a block
		// frequency of 600 seconds, the min target block time is 200 seconds,
		// and the max target block time is 1800 seconds.
		OakMaxBlockShift = 3

		// The max rise and max drop for the difficulty is kept at 0.4% per
		// block, which means that in 1008 blocks the difficulty can move a
		// maximum of about 55x. This is significant, and means that dramatic
		// hashrate changes can be responded to quickly, while still forcing an
		// attacker to do a significant amount of work in order to execute a
		// difficulty raising attack, and minimizing the chance that an attacker
		// can get lucky and fake a ton of work.
		OakMaxRise = big.NewRat(1004, 1e3)
		OakMaxDrop = big.NewRat(1e3, 1004)

		GenesisAirdropAllocation = []SiacoinOutput{
			{
				Value:      AirdropCommunityValue,
				UnlockHash: UnlockHashFromAddrStr("436890aacc53f93f9cc4538d9b4abba27dd5be6ff8a064fae7b78a67809db5e210819ffc4a21"),
			},
			{
				Value:      AirdropPoolValue,
				UnlockHash: UnlockHashFromAddrStr("78054218b7d0bc04929e5a3e6a2ac5fed29b98898cba3d740dd31a1aae6e8c8b3ce7467d4e8f"),
			},
			{
				Value:      AirdropNebulousLabsValue,
				UnlockHash: UnlockHashFromAddrStr("7d0c44f7664e2d34e53efde0661a6f628ec9264785ae8e3cd7c973e8d190c3c97b5e3ecbc567"),
			},
			{
				Value:      AirdropSiaPrimeValue,
				UnlockHash: UnlockHashFromAddrStr("aefe0af2713c112ba4d10dee7753726e5c4de3f237ea455151342615c95d0e797d7a8cce7b05"),
			},
		}

		ForkedGenesisSiafundAllocation = []SiafundOutput{
			{
				Value:      NewCurrency64(10000),
				UnlockHash: UnlockHashFromAddrStr("436890aacc53f93f9cc4538d9b4abba27dd5be6ff8a064fae7b78a67809db5e210819ffc4a21"),
			},
		}
		SiafundHardforkAllocation = map[BlockHeight][]SiafundOutput{
			SpfHardforkHeight: {
				{
					Value:      NewCurrency64(20000),
					UnlockHash: UnlockHashFromAddrStr("c25a37aa55baf3131e4d9335373338018d71f35bc22cd55d6d983394611d09a1fbb3f6edd5cc"),
				},
			},
			SpfSecondHardforkHeight: {
				{
					Value:      NewCurrency64(200000000 - 30000),
					UnlockHash: UnlockHashFromAddrStr("f3f512b45c24b531e571d59d99b804b19244cecf1da8117217101fa722e40a18917c91ebb080"),
				},
			},
		}

		GenesisSiafundAllocation = []SiafundOutput{
			{
				Value:      NewCurrency64(2),
				UnlockHash: UnlockHash{4, 57, 229, 188, 127, 20, 204, 245, 211, 167, 232, 130, 208, 64, 146, 62, 69, 98, 81, 102, 221, 7, 123, 100, 70, 107, 199, 113, 121, 26, 198, 252},
			},
			{
				Value:      NewCurrency64(6),
				UnlockHash: UnlockHash{4, 158, 29, 42, 105, 119, 43, 5, 138, 72, 190, 190, 101, 114, 79, 243, 189, 248, 208, 151, 30, 187, 233, 148, 225, 233, 28, 159, 19, 232, 75, 244},
			},
			{
				Value:      NewCurrency64(7),
				UnlockHash: UnlockHash{8, 7, 66, 250, 25, 74, 247, 108, 162, 79, 220, 151, 202, 228, 241, 11, 130, 138, 13, 248, 193, 167, 136, 197, 65, 63, 234, 174, 205, 216, 71, 230},
			},
			{
				Value:      NewCurrency64(8),
				UnlockHash: UnlockHash{44, 106, 239, 51, 138, 102, 242, 19, 204, 197, 248, 178, 219, 122, 152, 251, 19, 20, 52, 32, 175, 32, 4, 156, 73, 33, 163, 165, 222, 184, 217, 218},
			},
			{
				Value:      NewCurrency64(3),
				UnlockHash: UnlockHash{44, 163, 31, 233, 74, 103, 55, 132, 230, 159, 97, 78, 149, 147, 65, 110, 164, 211, 105, 173, 158, 29, 202, 43, 85, 217, 85, 75, 83, 37, 205, 223},
			},
			{
				Value:      NewCurrency64(1),
				UnlockHash: UnlockHash{51, 151, 146, 84, 199, 7, 59, 89, 111, 172, 227, 200, 62, 55, 165, 253, 238, 186, 28, 145, 47, 137, 200, 15, 70, 199, 187, 125, 243, 104, 179, 240},
			},
			{
				Value:      NewCurrency64(10),
				UnlockHash: UnlockHash{53, 118, 253, 229, 254, 229, 28, 131, 233, 156, 108, 58, 197, 152, 17, 160, 74, 252, 11, 49, 112, 240, 66, 119, 40, 98, 114, 251, 5, 86, 233, 117},
			},
			{
				Value:      NewCurrency64(50),
				UnlockHash: UnlockHash{56, 219, 3, 50, 28, 3, 166, 95, 141, 163, 202, 35, 60, 199, 219, 10, 151, 176, 228, 97, 176, 133, 189, 33, 211, 202, 83, 197, 31, 208, 254, 193},
			},
			{
				Value:      NewCurrency64(75),
				UnlockHash: UnlockHash{68, 190, 140, 87, 96, 232, 150, 32, 161, 177, 204, 65, 228, 223, 87, 217, 134, 90, 25, 56, 51, 45, 72, 107, 129, 12, 29, 202, 6, 7, 50, 13},
			},
			{
				Value:      NewCurrency64(10),
				UnlockHash: UnlockHash{69, 14, 201, 200, 90, 73, 245, 45, 154, 94, 161, 19, 199, 241, 203, 56, 13, 63, 5, 220, 121, 245, 247, 52, 194, 181, 252, 76, 130, 6, 114, 36},
			},
			{
				Value:      NewCurrency64(10),
				UnlockHash: UnlockHash{72, 128, 253, 207, 169, 48, 1, 26, 237, 205, 169, 102, 196, 224, 42, 186, 95, 151, 59, 226, 203, 136, 251, 223, 165, 38, 88, 110, 47, 213, 121, 224},
			},
			{
				Value:      NewCurrency64(50),
				UnlockHash: UnlockHash{72, 130, 164, 227, 218, 28, 60, 15, 56, 151, 212, 242, 77, 131, 232, 131, 42, 57, 132, 173, 113, 118, 66, 183, 38, 79, 96, 178, 105, 108, 26, 247},
			},
			{
				Value:      NewCurrency64(10),
				UnlockHash: UnlockHash{74, 210, 58, 228, 111, 69, 253, 120, 53, 195, 110, 26, 115, 76, 211, 202, 199, 159, 204, 14, 78, 92, 14, 131, 250, 22, 141, 236, 154, 44, 39, 135},
			},
			{
				Value:      NewCurrency64(15),
				UnlockHash: UnlockHash{85, 198, 154, 41, 196, 116, 226, 114, 202, 94, 214, 147, 87, 84, 247, 164, 195, 79, 58, 123, 26, 33, 68, 65, 116, 79, 181, 241, 241, 208, 215, 184},
			},
			{
				Value:      NewCurrency64(121),
				UnlockHash: UnlockHash{87, 239, 83, 125, 152, 14, 19, 22, 203, 136, 46, 192, 203, 87, 224, 190, 77, 236, 125, 18, 142, 223, 146, 70, 16, 23, 252, 19, 100, 69, 91, 111},
			},
			{
				Value:      NewCurrency64(222),
				UnlockHash: UnlockHash{91, 201, 101, 11, 188, 40, 35, 111, 236, 133, 31, 124, 97, 246, 140, 136, 143, 245, 152, 174, 111, 245, 188, 124, 21, 125, 187, 192, 203, 92, 253, 57},
			},
			{
				Value:      NewCurrency64(10),
				UnlockHash: UnlockHash{110, 240, 238, 173, 78, 138, 185, 138, 179, 227, 135, 153, 54, 132, 46, 62, 226, 206, 204, 35, 174, 107, 156, 15, 142, 2, 93, 132, 163, 60, 50, 89},
			},
			{
				Value:      NewCurrency64(3),
				UnlockHash: UnlockHash{114, 58, 147, 44, 64, 69, 72, 184, 65, 178, 213, 94, 157, 44, 88, 106, 92, 31, 145, 193, 215, 200, 215, 233, 99, 116, 36, 197, 160, 70, 79, 153},
			},
			{
				Value:      NewCurrency64(1),
				UnlockHash: UnlockHash{123, 106, 229, 101, 220, 252, 50, 203, 38, 183, 133, 152, 250, 167, 210, 155, 252, 102, 150, 29, 187, 3, 178, 53, 11, 145, 143, 33, 166, 115, 250, 40},
			},
			{
				Value:      NewCurrency64(5),
				UnlockHash: UnlockHash{124, 101, 207, 175, 50, 119, 207, 26, 62, 15, 247, 141, 150, 174, 73, 247, 238, 28, 77, 255, 222, 104, 166, 244, 112, 86, 227, 80, 215, 45, 69, 143},
			},
			{
				Value:      NewCurrency64(10),
				UnlockHash: UnlockHash{130, 184, 72, 15, 227, 79, 217, 205, 120, 254, 67, 69, 10, 49, 76, 194, 222, 30, 242, 62, 88, 179, 51, 117, 27, 166, 140, 6, 7, 22, 222, 185},
			},
			{
				Value:      NewCurrency64(25),
				UnlockHash: UnlockHash{134, 137, 198, 172, 96, 54, 45, 10, 100, 128, 91, 225, 226, 134, 143, 108, 31, 70, 187, 228, 54, 212, 70, 229, 149, 57, 64, 166, 153, 123, 238, 180},
			},
			{
				Value:      NewCurrency64(1),
				UnlockHash: UnlockHash{143, 253, 118, 229, 109, 181, 141, 224, 91, 144, 123, 160, 203, 221, 119, 104, 172, 13, 105, 77, 171, 185, 122, 54, 229, 168, 6, 130, 160, 130, 182, 151},
			},
			{
				Value:      NewCurrency64(8),
				UnlockHash: UnlockHash{147, 108, 249, 16, 36, 249, 108, 184, 196, 212, 241, 120, 219, 63, 45, 184, 86, 53, 96, 207, 130, 96, 210, 251, 136, 9, 193, 160, 131, 198, 221, 185},
			},
			{
				Value:      NewCurrency64(58),
				UnlockHash: UnlockHash{155, 79, 89, 28, 69, 71, 239, 198, 246, 2, 198, 254, 92, 59, 192, 205, 229, 152, 36, 186, 110, 122, 233, 221, 76, 143, 3, 238, 89, 231, 192, 23},
			},
			{
				Value:      NewCurrency64(2),
				UnlockHash: UnlockHash{156, 32, 76, 105, 213, 46, 66, 50, 27, 85, 56, 9, 106, 193, 80, 145, 19, 101, 84, 177, 145, 4, 125, 28, 79, 252, 43, 83, 118, 110, 206, 247},
			},
			{
				Value:      NewCurrency64(23),
				UnlockHash: UnlockHash{157, 169, 134, 24, 254, 22, 58, 188, 119, 87, 201, 238, 55, 168, 194, 131, 88, 18, 39, 168, 37, 2, 198, 194, 93, 202, 116, 146, 189, 17, 108, 44},
			},
			{
				Value:      NewCurrency64(10),
				UnlockHash: UnlockHash{158, 51, 104, 36, 242, 114, 67, 16, 168, 230, 4, 111, 241, 72, 5, 14, 182, 102, 169, 156, 144, 220, 103, 117, 223, 8, 58, 187, 124, 102, 80, 44},
			},
			{
				Value:      NewCurrency64(1),
				UnlockHash: UnlockHash{160, 175, 59, 33, 223, 30, 82, 60, 34, 110, 28, 203, 249, 93, 3, 16, 218, 12, 250, 206, 138, 231, 85, 67, 69, 191, 68, 198, 160, 87, 154, 68},
			},
			{
				Value:      NewCurrency64(75),
				UnlockHash: UnlockHash{163, 94, 51, 220, 14, 144, 83, 112, 62, 10, 0, 173, 161, 234, 211, 176, 186, 84, 9, 189, 250, 111, 33, 231, 114, 87, 100, 75, 72, 217, 11, 26},
			},
			{
				Value:      NewCurrency64(3),
				UnlockHash: UnlockHash{170, 7, 138, 116, 205, 20, 132, 197, 166, 251, 75, 93, 69, 6, 109, 244, 212, 119, 173, 114, 34, 18, 25, 21, 111, 203, 203, 253, 138, 104, 27, 36},
			},
			{
				Value:      NewCurrency64(90),
				UnlockHash: UnlockHash{173, 120, 128, 104, 186, 86, 151, 140, 191, 23, 231, 193, 77, 245, 243, 104, 196, 55, 155, 243, 111, 15, 84, 139, 148, 187, 173, 47, 104, 69, 141, 39},
			},
			{
				Value:      NewCurrency64(20),
				UnlockHash: UnlockHash{179, 185, 228, 166, 139, 94, 13, 193, 255, 227, 174, 99, 120, 105, 109, 221, 247, 4, 155, 243, 229, 37, 26, 98, 222, 12, 91, 80, 223, 33, 61, 56},
			},
			{
				Value:      NewCurrency64(5),
				UnlockHash: UnlockHash{193, 49, 103, 20, 170, 135, 182, 85, 149, 18, 159, 194, 152, 120, 162, 208, 49, 158, 220, 188, 114, 79, 1, 131, 62, 27, 86, 57, 244, 46, 64, 66},
			},
			{
				Value:      NewCurrency64(1),
				UnlockHash: UnlockHash{196, 71, 45, 222, 0, 21, 12, 121, 197, 224, 101, 65, 40, 57, 19, 119, 112, 205, 166, 23, 2, 91, 75, 231, 69, 143, 221, 68, 245, 75, 7, 52},
			},
			{
				Value:      NewCurrency64(44),
				UnlockHash: UnlockHash{196, 214, 236, 211, 227, 216, 152, 127, 164, 2, 235, 14, 235, 46, 142, 231, 83, 38, 7, 131, 208, 29, 179, 189, 62, 88, 129, 180, 119, 158, 214, 97},
			},
			{
				Value:      NewCurrency64(23),
				UnlockHash: UnlockHash{206, 58, 114, 148, 131, 49, 87, 197, 86, 18, 216, 26, 62, 79, 152, 175, 33, 4, 132, 160, 108, 231, 53, 200, 48, 76, 125, 94, 156, 85, 32, 130},
			},
			{
				Value:      NewCurrency64(80),
				UnlockHash: UnlockHash{200, 103, 135, 126, 197, 2, 203, 63, 241, 6, 245, 195, 220, 102, 27, 74, 232, 249, 201, 86, 207, 34, 51, 26, 180, 151, 136, 108, 112, 56, 132, 72},
			},
			{
				Value:      NewCurrency64(2),
				UnlockHash: UnlockHash{200, 249, 245, 218, 58, 253, 76, 250, 88, 114, 70, 239, 14, 2, 250, 123, 10, 192, 198, 61, 187, 155, 247, 152, 165, 174, 198, 24, 142, 39, 177, 119},
			},
			{
				Value:      NewCurrency64(1),
				UnlockHash: UnlockHash{209, 1, 199, 184, 186, 57, 21, 137, 33, 252, 219, 184, 130, 38, 32, 98, 63, 252, 250, 79, 70, 146, 169, 78, 180, 161, 29, 93, 38, 45, 175, 176},
			},
			{
				Value:      NewCurrency64(2),
				UnlockHash: UnlockHash{212, 107, 233, 43, 185, 138, 79, 253, 12, 237, 214, 17, 219, 198, 151, 92, 81, 129, 17, 120, 139, 58, 66, 119, 126, 220, 132, 136, 3, 108, 57, 58},
			},
			{
				Value:      NewCurrency64(3),
				UnlockHash: UnlockHash{214, 244, 146, 173, 173, 80, 33, 185, 29, 133, 77, 167, 185, 1, 38, 23, 111, 179, 104, 150, 105, 162, 120, 26, 245, 63, 114, 119, 52, 1, 44, 222},
			},
			{
				Value:      NewCurrency64(1),
				UnlockHash: UnlockHash{217, 218, 172, 16, 53, 134, 160, 226, 44, 138, 93, 53, 181, 62, 4, 209, 190, 27, 0, 93, 105, 17, 169, 61, 98, 145, 131, 112, 121, 55, 97, 184},
			},
			{
				Value:      NewCurrency64(1),
				UnlockHash: UnlockHash{223, 162, 172, 55, 54, 193, 37, 142, 200, 213, 230, 48, 186, 145, 184, 206, 15, 225, 167, 19, 37, 70, 38, 48, 135, 87, 205, 81, 187, 237, 181, 180},
			},
			{
				Value:      NewCurrency64(1),
				UnlockHash: UnlockHash{241, 46, 139, 41, 40, 63, 47, 169, 131, 173, 124, 246, 228, 213, 102, 44, 100, 217, 62, 237, 133, 154, 248, 69, 228, 2, 36, 206, 47, 250, 249, 170},
			},
			{
				Value:      NewCurrency64(50),
				UnlockHash: UnlockHash{241, 50, 229, 211, 66, 32, 115, 241, 117, 87, 180, 239, 76, 246, 14, 129, 105, 181, 153, 105, 105, 203, 229, 237, 23, 130, 193, 170, 100, 201, 38, 71},
			},
			{
				Value:      NewCurrency64(8841),
				UnlockHash: UnlockHash{125, 12, 68, 247, 102, 78, 45, 52, 229, 62, 253, 224, 102, 26, 111, 98, 142, 201, 38, 71, 133, 174, 142, 60, 215, 201, 115, 232, 209, 144, 195, 201},
			},
		}
	}

	// Create the genesis block.
	GenesisBlock = Block{
		Timestamp: GenesisTimestamp,
		Transactions: []Transaction{
			{SiacoinOutputs: GenesisAirdropAllocation},
			{SiafundOutputs: GenesisSiafundAllocation},
		},
	}
	// Calculate the genesis ID.
	GenesisID = GenesisBlock.ID()
}