// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package params

import (
	"encoding/binary"
	"fmt"
	"math"
	"math/big"

	"github.com/ubiq/go-ubiq/common"
	"github.com/ubiq/go-ubiq/crypto"
)

// Genesis hashes to enforce below configs on.
var (
	MainnetGenesisHash = common.HexToHash("0x406f1b7dd39fca54d8c702141851ed8b755463ab5b560e6f19b963b4047418af")
	TestnetGenesisHash = common.HexToHash("0x41941023680923e0fe4d74a34bdac8141f2540e3ae90623718e47d66d1ca4a2d")
	// RinkebyGenesisHash = common.HexToHash("0x6341fd3daf94b748c72ced5a5b26028f2474f5f00d824504e4fa37a75767e177")
)

// TrustedCheckpoints associates each known checkpoint with the genesis hash of
// the chain it belongs to.
var TrustedCheckpoints = map[common.Hash]*TrustedCheckpoint{
	MainnetGenesisHash: MainnetTrustedCheckpoint,
	TestnetGenesisHash: TestnetTrustedCheckpoint,
	//RinkebyGenesisHash: RinkebyTrustedCheckpoint,
}

// CheckpointOracles associates each known checkpoint oracles with the genesis hash of
// the chain it belongs to.
var CheckpointOracles = map[common.Hash]*CheckpointOracleConfig{
	MainnetGenesisHash: MainnetCheckpointOracle,
	TestnetGenesisHash: TestnetCheckpointOracle,
	//RinkebyGenesisHash: RinkebyCheckpointOracle,
}

var (
	// MainnetChainConfig is the chain parameters to run a node on the main network.
	MainnetChainConfig = &ChainConfig{
		ChainID:             big.NewInt(8),
		HomesteadBlock:      big.NewInt(0),
		EIP150Block:         big.NewInt(0),
		EIP150Hash:          common.HexToHash("0x2086799aeebeae135c246c65021c82b4e15a2c451340993aacfd2751886514f0"),
		EIP155Block:         big.NewInt(10),
		EIP158Block:         big.NewInt(10),
		ByzantiumBlock:      big.NewInt(1075090), // Andromeda
		ConstantinopleBlock: big.NewInt(1075090), // Andromeda
		PetersburgBlock:     big.NewInt(1075090), // Andromeda
		IstanbulBlock:       big.NewInt(9069000),
		Ubqhash: &UbqhashConfig{
			DigishieldModBlock: big.NewInt(4088),
			FluxBlock:          big.NewInt(8000),
			MonetaryPolicy: []UbqhashMPStep{
				UbqhashMPStep{
					Block:  big.NewInt(0),
					Reward: big.NewInt(8e+18),
				},
				UbqhashMPStep{
					Block:  big.NewInt(358363),
					Reward: big.NewInt(7e+18),
				},
				UbqhashMPStep{
					Block:  big.NewInt(716727),
					Reward: big.NewInt(6e+18),
				},
				UbqhashMPStep{
					Block:  big.NewInt(1075090),
					Reward: big.NewInt(5e+18),
				},
				UbqhashMPStep{
					Block:  big.NewInt(1433454),
					Reward: big.NewInt(4e+18),
				},
				UbqhashMPStep{
					Block:  big.NewInt(1791818),
					Reward: big.NewInt(3e+18),
				},
				UbqhashMPStep{
					Block:  big.NewInt(2150181),
					Reward: big.NewInt(2e+18),
				},
				UbqhashMPStep{
					Block:  big.NewInt(2508545),
					Reward: big.NewInt(1e+18),
				},
			},
		},
	}

	// MainnetTrustedCheckpoint contains the light client trusted checkpoint for the main network.
	MainnetTrustedCheckpoint = &TrustedCheckpoint{
		SectionIndex: 227,
		SectionHead:  common.HexToHash("0xa2e0b25d72c2fc6e35a7f853cdacb193b4b4f95c606accf7f8fa8415283582c7"),
		CHTRoot:      common.HexToHash("0xf69bdd4053b95b61a27b106a0e86103d791edd8574950dc96aa351ab9b9f1aa0"),
		BloomRoot:    common.HexToHash("0xec1b454d4c6322c78ccedf76ac922a8698c3cac4d98748a84af4995b7bd3d744"),
	}

	// MainnetCheckpointOracle contains a set of configs for the main network oracle.
	MainnetCheckpointOracle = &CheckpointOracleConfig{
		Address: common.HexToAddress("0x9a9070028361F7AAbeB3f2F2Dc07F82C4a98A02a"),
		Signers: []common.Address{
			common.HexToAddress("0x1b2C260efc720BE89101890E4Db589b44E950527"), // Peter
			common.HexToAddress("0x78d1aD571A1A09D60D9BBf25894b44e4C8859595"), // Martin
			common.HexToAddress("0x286834935f4A8Cfb4FF4C77D5770C2775aE2b0E7"), // Zsolt
			common.HexToAddress("0xb86e2B0Ab5A4B1373e40c51A7C712c70Ba2f9f8E"), // Gary
			common.HexToAddress("0x0DF8fa387C602AE62559cC4aFa4972A7045d6707"), // Guillaume
		},
		Threshold: 2,
	}

	// TestnetChainConfig contains the chain parameters to run a node on the Ropsten test network.
	TestnetChainConfig = &ChainConfig{
		ChainID:             big.NewInt(9),
		HomesteadBlock:      big.NewInt(0),
		EIP150Block:         big.NewInt(0),
		EIP150Hash:          common.HexToHash("0x41941023680923e0fe4d74a34bdac8141f2540e3ae90623718e47d66d1ca4a2d"),
		EIP155Block:         big.NewInt(10),
		EIP158Block:         big.NewInt(10),
		ByzantiumBlock:      big.NewInt(math.MaxInt64),
		ConstantinopleBlock: big.NewInt(math.MaxInt64),
		PetersburgBlock:     big.NewInt(math.MaxInt64),
		IstanbulBlock:       big.NewInt(6485846),
		Ubqhash: &UbqhashConfig{
			DigishieldModBlock: big.NewInt(4088),
			FluxBlock:          big.NewInt(8000),
			MonetaryPolicy: []UbqhashMPStep{
				UbqhashMPStep{
					Block:  big.NewInt(0),
					Reward: big.NewInt(8e+18),
				},
				UbqhashMPStep{
					Block:  big.NewInt(358363),
					Reward: big.NewInt(7e+18),
				},
				UbqhashMPStep{
					Block:  big.NewInt(716727),
					Reward: big.NewInt(6e+18),
				},
				UbqhashMPStep{
					Block:  big.NewInt(1075090),
					Reward: big.NewInt(5e+18),
				},
				UbqhashMPStep{
					Block:  big.NewInt(1433454),
					Reward: big.NewInt(4e+18),
				},
				UbqhashMPStep{
					Block:  big.NewInt(1791818),
					Reward: big.NewInt(3e+18),
				},
				UbqhashMPStep{
					Block:  big.NewInt(2150181),
					Reward: big.NewInt(2e+18),
				},
				UbqhashMPStep{
					Block:  big.NewInt(2508545),
					Reward: big.NewInt(1e+18),
				},
			},
		},
	}

	// TestnetTrustedCheckpoint contains the light client trusted checkpoint for the Ropsten test network.
	TestnetTrustedCheckpoint = &TrustedCheckpoint{
		SectionIndex: 161,
		SectionHead:  common.HexToHash("0x5378afa734e1feafb34bcca1534c4d96952b754579b96a4afb23d5301ecececc"),
		CHTRoot:      common.HexToHash("0x1cf2b071e7443a62914362486b613ff30f60cea0d9c268ed8c545f876a3ee60c"),
		BloomRoot:    common.HexToHash("0x5ac25c84bd18a9cbe878d4609a80220f57f85037a112644532412ba0d498a31b"),
	}

	// RopstenCheckpointOracle contains a set of configs for the Ropsten test network oracle.
	TestnetCheckpointOracle = &CheckpointOracleConfig{
		Address: common.HexToAddress("0xEF79475013f154E6A65b54cB2742867791bf0B84"),
		Signers: []common.Address{
			common.HexToAddress("0x32162F3581E88a5f62e8A61892B42C46E2c18f7b"), // Peter
			common.HexToAddress("0x78d1aD571A1A09D60D9BBf25894b44e4C8859595"), // Martin
			common.HexToAddress("0x286834935f4A8Cfb4FF4C77D5770C2775aE2b0E7"), // Zsolt
			common.HexToAddress("0xb86e2B0Ab5A4B1373e40c51A7C712c70Ba2f9f8E"), // Gary
			common.HexToAddress("0x0DF8fa387C602AE62559cC4aFa4972A7045d6707"), // Guillaume
		},
		Threshold: 2,
	}

	/*
		// RinkebyChainConfig contains the chain parameters to run a node on the Rinkeby test network.
		RinkebyChainConfig = &ChainConfig{
			ChainID:             big.NewInt(4),
			HomesteadBlock:      big.NewInt(1),
			EIP150Block:         big.NewInt(2),
			EIP150Hash:          common.HexToHash("0x9b095b36c15eaf13044373aef8ee0bd3a382a5abb92e402afa44b8249c3a90e9"),
			EIP155Block:         big.NewInt(3),
			EIP158Block:         big.NewInt(3),
			ByzantiumBlock:      big.NewInt(1035301),
			ConstantinopleBlock: big.NewInt(3660663),
			PetersburgBlock:     big.NewInt(4321234),
			IstanbulBlock:       big.NewInt(5435345),
			Clique: &CliqueConfig{
				Period: 15,
				Epoch:  30000,
			},
		}

		// RinkebyTrustedCheckpoint contains the light client trusted checkpoint for the Rinkeby test network.
		RinkebyTrustedCheckpoint = &TrustedCheckpoint{
			SectionIndex: 196,
			SectionHead:  common.HexToHash("0x66faad1de5cd0c9da5c4c0b0d4e2e86c2ed6a9cde7441a9211deb3b6d049a01e"),
			CHTRoot:      common.HexToHash("0x5752c6633b5d052298316a4d7dd9d2e931b83e3387584f82998a1f6f05b5e4c1"),
			BloomRoot:    common.HexToHash("0x6a2e14dc35d2b6e0361af41a0e28143b59a578a4458e58ca2fb2172b6688b963"),
		}

		// RinkebyCheckpointOracle contains a set of configs for the Rinkeby test network oracle.
		RinkebyCheckpointOracle = &CheckpointOracleConfig{
			Address: common.HexToAddress("0xebe8eFA441B9302A0d7eaECc277c09d20D684540"),
			Signers: []common.Address{
				common.HexToAddress("0xd9c9cd5f6779558b6e0ed4e6acf6b1947e7fa1f3"), // Peter
				common.HexToAddress("0x78d1aD571A1A09D60D9BBf25894b44e4C8859595"), // Martin
				common.HexToAddress("0x286834935f4A8Cfb4FF4C77D5770C2775aE2b0E7"), // Zsolt
				common.HexToAddress("0xb86e2B0Ab5A4B1373e40c51A7C712c70Ba2f9f8E"), // Gary
			},
			Threshold: 2,
		}*/

	// AllUbqhashProtocolChanges contains every protocol change (EIPs) introduced
	// and accepted by the Ethereum core developers into the Ubqhash consensus.
	//
	// This configuration is intentionally not using keyed fields to force anyone
	// adding flags to the config to also have to set these fields.
	AllUbqhashProtocolChanges = &ChainConfig{big.NewInt(1337), big.NewInt(0), big.NewInt(0), common.Hash{}, big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), nil, &UbqhashConfig{big.NewInt(0), big.NewInt(0), []UbqhashMPStep{}}, nil}

	// AllCliqueProtocolChanges contains every protocol change (EIPs) introduced
	// and accepted by the Ubiq core developers into the Clique consensus.
	//
	// This configuration is intentionally not using keyed fields to force anyone
	// adding flags to the config to also have to set these fields.
	AllCliqueProtocolChanges = &ChainConfig{big.NewInt(1337), big.NewInt(0), big.NewInt(0), common.Hash{}, big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), nil, nil, &CliqueConfig{Period: 0, Epoch: 30000}}

	TestChainConfig = &ChainConfig{big.NewInt(1), big.NewInt(0), big.NewInt(0), common.Hash{}, big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), nil, &UbqhashConfig{big.NewInt(0), big.NewInt(0), []UbqhashMPStep{}}, nil}
	TestRules       = TestChainConfig.Rules(new(big.Int))
)

// TrustedCheckpoint represents a set of post-processed trie roots (CHT and
// BloomTrie) associated with the appropriate section index and head hash. It is
// used to start light syncing from this checkpoint and avoid downloading the
// entire header chain while still being able to securely access old headers/logs.
type TrustedCheckpoint struct {
	SectionIndex uint64      `json:"sectionIndex"`
	SectionHead  common.Hash `json:"sectionHead"`
	CHTRoot      common.Hash `json:"chtRoot"`
	BloomRoot    common.Hash `json:"bloomRoot"`
}

// HashEqual returns an indicator comparing the itself hash with given one.
func (c *TrustedCheckpoint) HashEqual(hash common.Hash) bool {
	if c.Empty() {
		return hash == common.Hash{}
	}
	return c.Hash() == hash
}

// Hash returns the hash of checkpoint's four key fields(index, sectionHead, chtRoot and bloomTrieRoot).
func (c *TrustedCheckpoint) Hash() common.Hash {
	buf := make([]byte, 8+3*common.HashLength)
	binary.BigEndian.PutUint64(buf, c.SectionIndex)
	copy(buf[8:], c.SectionHead.Bytes())
	copy(buf[8+common.HashLength:], c.CHTRoot.Bytes())
	copy(buf[8+2*common.HashLength:], c.BloomRoot.Bytes())
	return crypto.Keccak256Hash(buf)
}

// Empty returns an indicator whether the checkpoint is regarded as empty.
func (c *TrustedCheckpoint) Empty() bool {
	return c.SectionHead == (common.Hash{}) || c.CHTRoot == (common.Hash{}) || c.BloomRoot == (common.Hash{})
}

// CheckpointOracleConfig represents a set of checkpoint contract(which acts as an oracle)
// config which used for light client checkpoint syncing.
type CheckpointOracleConfig struct {
	Address   common.Address   `json:"address"`
	Signers   []common.Address `json:"signers"`
	Threshold uint64           `json:"threshold"`
}

// ChainConfig is the core config which determines the blockchain settings.
//
// ChainConfig is stored in the database on a per block basis. This means
// that any network, identified by its genesis block, can have its own
// set of configuration options.
type ChainConfig struct {
	ChainID *big.Int `json:"chainId"` // chainId identifies the current chain and is used for replay protection

	HomesteadBlock *big.Int `json:"homesteadBlock,omitempty"` // Homestead switch block (nil = no fork, 0 = already homestead)

	// EIP150 implements the Gas price changes (https://github.com/ethereum/EIPs/issues/150)
	EIP150Block *big.Int    `json:"eip150Block,omitempty"` // EIP150 HF block (nil = no fork)
	EIP150Hash  common.Hash `json:"eip150Hash,omitempty"`  // EIP150 HF hash (needed for header only clients as only gas pricing changed)

	EIP155Block *big.Int `json:"eip155Block,omitempty"` // EIP155 HF block
	EIP158Block *big.Int `json:"eip158Block,omitempty"` // EIP158 HF block

	ByzantiumBlock      *big.Int `json:"byzantiumBlock,omitempty"`      // Byzantium switch block (nil = no fork, 0 = already on byzantium)
	ConstantinopleBlock *big.Int `json:"constantinopleBlock,omitempty"` // Constantinople switch block (nil = no fork, 0 = already activated)
	PetersburgBlock     *big.Int `json:"petersburgBlock,omitempty"`     // Petersburg switch block (nil = same as Constantinople)
	IstanbulBlock       *big.Int `json:"istanbulBlock,omitempty"`       // Istanbul switch block (nil = no fork, 0 = already on istanbul)
	EWASMBlock          *big.Int `json:"ewasmBlock,omitempty"`          // EWASM switch block (nil = no fork, 0 = already activated)

	// Various consensus engines
	Ubqhash *UbqhashConfig `json:"ubqhash,omitempty"`
	Clique  *CliqueConfig  `json:"clique,omitempty"`
}

// Ubqhash monetary policy reward step
type UbqhashMPStep struct {
	Block  *big.Int `json:"block"`
	Reward *big.Int `json:"reward"`
}

// UbqhashConfig is the consensus engine configs for proof-of-work based sealing.
type UbqhashConfig struct {
	DigishieldModBlock *big.Int        `json:"digishieldModBlock,omitempty"` // Block to activate the DigiShield V3 mod
	FluxBlock          *big.Int        `json:"fluxBlock"`                    // Block to activate the Flux difficulty algorithm
	MonetaryPolicy     []UbqhashMPStep `json:"monetaryPolicy"`               // Blocks to step the block reward down
}

// String implements the stringer interface, returning the consensus engine details.
func (c *UbqhashConfig) String() string {
	return "ubqhash"
}

// CliqueConfig is the consensus engine configs for proof-of-authority based sealing.
type CliqueConfig struct {
	Period uint64 `json:"period"` // Number of seconds between blocks to enforce
	Epoch  uint64 `json:"epoch"`  // Epoch length to reset votes and checkpoint
}

// String implements the stringer interface, returning the consensus engine details.
func (c *CliqueConfig) String() string {
	return "clique"
}

// String implements the fmt.Stringer interface.
func (c *ChainConfig) String() string {
	var engine interface{}
	switch {
	case c.Ubqhash != nil:
		engine = c.Ubqhash
	case c.Clique != nil:
		engine = c.Clique
	default:
		engine = "unknown"
	}
	return fmt.Sprintf("{ChainID: %v Homestead: %v EIP150: %v EIP155: %v EIP158: %v Byzantium: %v Constantinople: %v  Petersburg: %v Istanbul: %v Engine: %v}",
		c.ChainID,
		c.HomesteadBlock,
		c.EIP150Block,
		c.EIP155Block,
		c.EIP158Block,
		c.ByzantiumBlock,
		c.ConstantinopleBlock,
		c.PetersburgBlock,
		c.IstanbulBlock,
		engine,
	)
}

// IsHomestead returns whether num is either equal to the homestead block or greater.
func (c *ChainConfig) IsHomestead(num *big.Int) bool {
	return isForked(c.HomesteadBlock, num)
}

// IsEIP150 returns whether num is either equal to the EIP150 fork block or greater.
func (c *ChainConfig) IsEIP150(num *big.Int) bool {
	return isForked(c.EIP150Block, num)
}

// IsEIP155 returns whether num is either equal to the EIP155 fork block or greater.
func (c *ChainConfig) IsEIP155(num *big.Int) bool {
	return isForked(c.EIP155Block, num)
}

// IsEIP158 returns whether num is either equal to the EIP158 fork block or greater.
func (c *ChainConfig) IsEIP158(num *big.Int) bool {
	return isForked(c.EIP158Block, num)
}

// IsByzantium returns whether num is either equal to the Byzantium fork block or greater.
func (c *ChainConfig) IsByzantium(num *big.Int) bool {
	return isForked(c.ByzantiumBlock, num)
}

// IsConstantinople returns whether num is either equal to the Constantinople fork block or greater.
func (c *ChainConfig) IsConstantinople(num *big.Int) bool {
	return isForked(c.ConstantinopleBlock, num)
}

// IsPetersburg returns whether num is either
// - equal to or greater than the PetersburgBlock fork block,
// - OR is nil, and Constantinople is active
func (c *ChainConfig) IsPetersburg(num *big.Int) bool {
	return isForked(c.PetersburgBlock, num) || c.PetersburgBlock == nil && isForked(c.ConstantinopleBlock, num)
}

// IsIstanbul returns whether num is either equal to the Istanbul fork block or greater.
func (c *ChainConfig) IsIstanbul(num *big.Int) bool {
	return isForked(c.IstanbulBlock, num)
}

// IsEWASM returns whether num represents a block number after the EWASM fork
func (c *ChainConfig) IsEWASM(num *big.Int) bool {
	return isForked(c.EWASMBlock, num)
}

// CheckCompatible checks whether scheduled fork transitions have been imported
// with a mismatching chain configuration.
func (c *ChainConfig) CheckCompatible(newcfg *ChainConfig, height uint64) *ConfigCompatError {
	bhead := new(big.Int).SetUint64(height)

	// Iterate checkCompatible to find the lowest conflict.
	var lasterr *ConfigCompatError
	for {
		err := c.checkCompatible(newcfg, bhead)
		if err == nil || (lasterr != nil && err.RewindTo == lasterr.RewindTo) {
			break
		}
		lasterr = err
		bhead.SetUint64(err.RewindTo)
	}
	return lasterr
}

// CheckConfigForkOrder checks that we don't "skip" any forks, geth isn't pluggable enough
// to guarantee that forks can be implemented in a different order than on official networks
func (c *ChainConfig) CheckConfigForkOrder() error {
	type fork struct {
		name  string
		block *big.Int
	}
	var lastFork fork
	for _, cur := range []fork{
		{"homesteadBlock", c.HomesteadBlock},
		{"eip150Block", c.EIP150Block},
		{"eip155Block", c.EIP155Block},
		{"eip158Block", c.EIP158Block},
		{"byzantiumBlock", c.ByzantiumBlock},
		{"constantinopleBlock", c.ConstantinopleBlock},
		{"petersburgBlock", c.PetersburgBlock},
		{"istanbulBlock", c.IstanbulBlock},
	} {
		if lastFork.name != "" {
			// Next one must be higher number
			if lastFork.block == nil && cur.block != nil {
				return fmt.Errorf("unsupported fork ordering: %v not enabled, but %v enabled at %v",
					lastFork.name, cur.name, cur.block)
			}
			if lastFork.block != nil && cur.block != nil {
				if lastFork.block.Cmp(cur.block) > 0 {
					return fmt.Errorf("unsupported fork ordering: %v enabled at %v, but %v enabled at %v",
						lastFork.name, lastFork.block, cur.name, cur.block)
				}
			}
		}
		lastFork = cur
	}
	return nil
}

func (c *ChainConfig) checkCompatible(newcfg *ChainConfig, head *big.Int) *ConfigCompatError {
	if isForkIncompatible(c.HomesteadBlock, newcfg.HomesteadBlock, head) {
		return newCompatError("Homestead fork block", c.HomesteadBlock, newcfg.HomesteadBlock)
	}
	if isForkIncompatible(c.EIP150Block, newcfg.EIP150Block, head) {
		return newCompatError("EIP150 fork block", c.EIP150Block, newcfg.EIP150Block)
	}
	if isForkIncompatible(c.EIP155Block, newcfg.EIP155Block, head) {
		return newCompatError("EIP155 fork block", c.EIP155Block, newcfg.EIP155Block)
	}
	if isForkIncompatible(c.EIP158Block, newcfg.EIP158Block, head) {
		return newCompatError("EIP158 fork block", c.EIP158Block, newcfg.EIP158Block)
	}
	if c.IsEIP158(head) && !configNumEqual(c.ChainID, newcfg.ChainID) {
		return newCompatError("EIP158 chain ID", c.EIP158Block, newcfg.EIP158Block)
	}
	if isForkIncompatible(c.ByzantiumBlock, newcfg.ByzantiumBlock, head) {
		return newCompatError("Byzantium fork block", c.ByzantiumBlock, newcfg.ByzantiumBlock)
	}
	if isForkIncompatible(c.ConstantinopleBlock, newcfg.ConstantinopleBlock, head) {
		return newCompatError("Constantinople fork block", c.ConstantinopleBlock, newcfg.ConstantinopleBlock)
	}
	if isForkIncompatible(c.PetersburgBlock, newcfg.PetersburgBlock, head) {
		return newCompatError("Petersburg fork block", c.PetersburgBlock, newcfg.PetersburgBlock)
	}
	if isForkIncompatible(c.IstanbulBlock, newcfg.IstanbulBlock, head) {
		return newCompatError("Istanbul fork block", c.IstanbulBlock, newcfg.IstanbulBlock)
	}
	if isForkIncompatible(c.EWASMBlock, newcfg.EWASMBlock, head) {
		return newCompatError("ewasm fork block", c.EWASMBlock, newcfg.EWASMBlock)
	}
	return nil
}

// isForkIncompatible returns true if a fork scheduled at s1 cannot be rescheduled to
// block s2 because head is already past the fork.
func isForkIncompatible(s1, s2, head *big.Int) bool {
	return (isForked(s1, head) || isForked(s2, head)) && !configNumEqual(s1, s2)
}

// isForked returns whether a fork scheduled at block s is active at the given head block.
func isForked(s, head *big.Int) bool {
	if s == nil || head == nil {
		return false
	}
	return s.Cmp(head) <= 0
}

func configNumEqual(x, y *big.Int) bool {
	if x == nil {
		return y == nil
	}
	if y == nil {
		return x == nil
	}
	return x.Cmp(y) == 0
}

// ConfigCompatError is raised if the locally-stored blockchain is initialised with a
// ChainConfig that would alter the past.
type ConfigCompatError struct {
	What string
	// block numbers of the stored and new configurations
	StoredConfig, NewConfig *big.Int
	// the block number to which the local chain must be rewound to correct the error
	RewindTo uint64
}

func newCompatError(what string, storedblock, newblock *big.Int) *ConfigCompatError {
	var rew *big.Int
	switch {
	case storedblock == nil:
		rew = newblock
	case newblock == nil || storedblock.Cmp(newblock) < 0:
		rew = storedblock
	default:
		rew = newblock
	}
	err := &ConfigCompatError{what, storedblock, newblock, 0}
	if rew != nil && rew.Sign() > 0 {
		err.RewindTo = rew.Uint64() - 1
	}
	return err
}

func (err *ConfigCompatError) Error() string {
	return fmt.Sprintf("mismatching %s in database (have %d, want %d, rewindto %d)", err.What, err.StoredConfig, err.NewConfig, err.RewindTo)
}

// Rules wraps ChainConfig and is merely syntactic sugar or can be used for functions
// that do not have or require information about the block.
//
// Rules is a one time interface meaning that it shouldn't be used in between transition
// phases.
type Rules struct {
	ChainID                                                 *big.Int
	IsHomestead, IsEIP150, IsEIP155, IsEIP158               bool
	IsByzantium, IsConstantinople, IsPetersburg, IsIstanbul bool
}

// Rules ensures c's ChainID is not nil.
func (c *ChainConfig) Rules(num *big.Int) Rules {
	chainID := c.ChainID
	if chainID == nil {
		chainID = new(big.Int)
	}
	return Rules{
		ChainID:          new(big.Int).Set(chainID),
		IsHomestead:      c.IsHomestead(num),
		IsEIP150:         c.IsEIP150(num),
		IsEIP155:         c.IsEIP155(num),
		IsEIP158:         c.IsEIP158(num),
		IsByzantium:      c.IsByzantium(num),
		IsConstantinople: c.IsConstantinople(num),
		IsPetersburg:     c.IsPetersburg(num),
		IsIstanbul:       c.IsIstanbul(num),
	}
}
