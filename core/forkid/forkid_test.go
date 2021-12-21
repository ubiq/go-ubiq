// Copyright 2019 The go-ethereum Authors
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

package forkid

import (
	"bytes"
	"math"
	"testing"

	"github.com/ubiq/go-ubiq/v6/common"
	"github.com/ubiq/go-ubiq/v6/params"
	"github.com/ubiq/go-ubiq/v6/rlp"
)

// TestCreation tests that different genesis and fork rule combinations result in
// the correct fork ID.
func TestCreation(t *testing.T) {
	type testcase struct {
		head uint64
		want ID
	}
	tests := []struct {
		config  *params.ChainConfig
		genesis common.Hash
		cases   []testcase
	}{
		// Mainnet test cases
		{
			params.MainnetChainConfig,
			params.MainnetGenesisHash,
			[]testcase{
				{0, ID{Hash: checksumToBytes(0xf3073ee8), Next: 10}},            // Unsynced, last Frontier, Homestead,
				{10, ID{Hash: checksumToBytes(0x3f0fcc99), Next: 1075090}},      // First Spurious block
				{1075089, ID{Hash: checksumToBytes(0x3f0fcc99), Next: 1075090}}, // Last Spurious block
				{1075090, ID{Hash: checksumToBytes(0xa4ecb4b6), Next: 1500000}}, // First Byzantium, Constantinople, Petersbug, (andromeda)
				{1499999, ID{Hash: checksumToBytes(0xa4ecb4b6), Next: 1500000}}, // Last Byzantium, Constantinople, Petersbug, (andromeda)
				{1500000, ID{Hash: checksumToBytes(0x65ea97e0), Next: 1791793}}, // First Istanbul (taurus)
				{1791792, ID{Hash: checksumToBytes(0x65ea97e0), Next: 1791793}}, // Last Istanbul (taurus)
				{1791793, ID{Hash: checksumToBytes(0x9ec7f55b), Next: 0}},       // First Berlin, London (orion)
				// {5000000, ID{Hash: checksumToBytes(0x757a1c47), Next: 0}}, // Future Berlin block
			},
		},
		// Rinkeby test cases
		{
			params.RinkebyChainConfig,
			params.RinkebyGenesisHash,
			[]testcase{
				{0, ID{Hash: checksumToBytes(0x3b8e0691), Next: 1}},             // Unsynced, last Frontier block
				{1, ID{Hash: checksumToBytes(0x60949295), Next: 2}},             // First and last Homestead block
				{2, ID{Hash: checksumToBytes(0x8bde40dd), Next: 3}},             // First and last Tangerine block
				{3, ID{Hash: checksumToBytes(0xcb3a64bb), Next: 1035301}},       // First Spurious block
				{1035300, ID{Hash: checksumToBytes(0xcb3a64bb), Next: 1035301}}, // Last Spurious block
				{1035301, ID{Hash: checksumToBytes(0x8d748b57), Next: 3660663}}, // First Byzantium block
				{3660662, ID{Hash: checksumToBytes(0x8d748b57), Next: 3660663}}, // Last Byzantium block
				{3660663, ID{Hash: checksumToBytes(0xe49cab14), Next: 4321234}}, // First Constantinople block
				{4321233, ID{Hash: checksumToBytes(0xe49cab14), Next: 4321234}}, // Last Constantinople block
				{4321234, ID{Hash: checksumToBytes(0xafec6b27), Next: 5435345}}, // First Petersburg block
				{5435344, ID{Hash: checksumToBytes(0xafec6b27), Next: 5435345}}, // Last Petersburg block
				{5435345, ID{Hash: checksumToBytes(0xcbdb8838), Next: 8290928}}, // First Istanbul block
				{8290927, ID{Hash: checksumToBytes(0xcbdb8838), Next: 8290928}}, // Last Istanbul block
				{8290928, ID{Hash: checksumToBytes(0x6910c8bd), Next: 8897988}}, // First Berlin block
				{8897987, ID{Hash: checksumToBytes(0x6910c8bd), Next: 8897988}}, // Last Berlin block
				{8897988, ID{Hash: checksumToBytes(0x8E29F2F3), Next: 0}},       // First London block
				{10000000, ID{Hash: checksumToBytes(0x8E29F2F3), Next: 0}},      // Future London block
			},
		},
		// Goerli test cases
		{
			params.GoerliChainConfig,
			params.GoerliGenesisHash,
			[]testcase{
				{0, ID{Hash: checksumToBytes(0xa3f5ab08), Next: 1561651}},       // Unsynced, last Frontier, Homestead, Tangerine, Spurious, Byzantium, Constantinople and first Petersburg block
				{1561650, ID{Hash: checksumToBytes(0xa3f5ab08), Next: 1561651}}, // Last Petersburg block
				{1561651, ID{Hash: checksumToBytes(0xc25efa5c), Next: 4460644}}, // First Istanbul block
				{4460643, ID{Hash: checksumToBytes(0xc25efa5c), Next: 4460644}}, // Last Istanbul block
				{4460644, ID{Hash: checksumToBytes(0x757a1c47), Next: 5062605}}, // First Berlin block
				{5000000, ID{Hash: checksumToBytes(0x757a1c47), Next: 5062605}}, // Last Berlin block
				{5062605, ID{Hash: checksumToBytes(0xB8C6299D), Next: 0}},       // First London block
				{6000000, ID{Hash: checksumToBytes(0xB8C6299D), Next: 0}},       // Future London block
			},
		},
	}
	for i, tt := range tests {
		for j, ttt := range tt.cases {
			if have := NewID(tt.config, tt.genesis, ttt.head); have != ttt.want {
				t.Errorf("test %d, case %d: fork ID mismatch: have %x, want %x", i, j, have, ttt.want)
			}
		}
	}
}

// TestValidation tests that a local peer correctly validates and accepts a remote
// fork ID.
func TestValidation(t *testing.T) {
	tests := []struct {
		head uint64
		id   ID
		err  error
	}{
		// Local is mainnet Andromeda, remote announces the same. No future fork is announced.
		{1075099, ID{Hash: checksumToBytes(0xa4ecb4b6), Next: 0}, nil},

		// Local is mainnet Andromeda, remote announces the same. Remote also announces a next fork
		// at block 0xffffffff, but that is uncertain.
		{1075102, ID{Hash: checksumToBytes(0xa4ecb4b6), Next: math.MaxUint64}, nil},

		// Local is mainnet currently in Spurious only (so it's aware of Andromeda), remote announces
		// also Spurious, but it's not yet aware of Andromeda (e.g. non updated node before the fork).
		// In this case we don't know if Andromeda passed yet or not.
		{1075000, ID{Hash: checksumToBytes(0x3f0fcc99), Next: 0}, nil},

		// Local is mainnet currently in Spurious only (so it's aware of Andromeda), remote announces
		// also Spurious, and it's also aware of Andromeda (e.g. updated node before the fork). We
		// don't know if Andromeda passed yet (will pass) or not.
		{1075089, ID{Hash: checksumToBytes(0x3f0fcc99), Next: 1075090}, nil},

		// Local is mainnet currently in Spurious only (so it's aware of Andromeda), remote announces
		// also Spurious, and it's also aware of some random fork (e.g. misconfigured Andromeda). As
		// neither forks passed at neither nodes, they may mismatch, but we still connect for now.
		{1075089, ID{Hash: checksumToBytes(0x3f0fcc99), Next: math.MaxUint64}, nil},

		// Local is mainnet Andromeda, remote announces Spurious + knowledge about Andromeda. Remote
		// is simply out of sync, accept.
		{1075200, ID{Hash: checksumToBytes(0x3f0fcc99), Next: 1075090}, nil},

		// Local is mainnet Taurus, remote announces Spurious + knowledge about Andromeda. Remote
		// is definitely out of sync. It may or may not need the Taurus update, we don't know yet.
		{1500200, ID{Hash: checksumToBytes(0x3f0fcc99), Next: 1075090}, nil},

		// Local is mainnet Andromeda, remote announces Taurus. Local is out of sync, accept.
		{1400000, ID{Hash: checksumToBytes(0x65ea97e0), Next: 0}, nil},

		// Local is mainnet Spurious, remote announces Andromeda, but is not aware of Taurus. Local
		// out of sync. Local also knows about a future fork, but that is uncertain yet.
		{1075089, ID{Hash: checksumToBytes(0xa4ecb4b6), Next: 0}, nil},

		// Local is mainnet Taurus. remote announces Andromeda but is not aware of further forks.
		// Remote needs software update.
		{1500500, ID{Hash: checksumToBytes(0xa4ecb4b6), Next: 0}, ErrRemoteStale},

		// Local is mainnet Petersburg, and isn't aware of more forks. Remote announces Petersburg +
		// 0xffffffff. Local needs software update, reject.
		{7987396, ID{Hash: checksumToBytes(0x5cddc0e1), Next: 0}, ErrLocalIncompatibleOrStale},

		// Local is mainnet Byzantium, and is aware of Petersburg. Remote announces Petersburg +
		// 0xffffffff. Local needs software update, reject.
		{7279999, ID{Hash: checksumToBytes(0x5cddc0e1), Next: 0}, ErrLocalIncompatibleOrStale},

		// Local is mainnet Petersburg, remote is Rinkeby Petersburg.
		{7987396, ID{Hash: checksumToBytes(0xafec6b27), Next: 0}, ErrLocalIncompatibleOrStale},

		// Local is mainnet Arrow Glacier, far in the future. Remote announces Gopherium (non existing fork)
		// at some future block 88888888, for itself, but past block for local. Local is incompatible.
		//
		// This case detects non-upgraded nodes with majority hash power (typical Ropsten mess).
		{88888888, ID{Hash: checksumToBytes(0x20c327fc), Next: 88888888}, ErrLocalIncompatibleOrStale},

		// Local is mainnet Byzantium. Remote is also in Byzantium, but announces Gopherium (non existing
		// fork) at block 7279999, before Petersburg. Local is incompatible.
		{7279999, ID{Hash: checksumToBytes(0xa00bc324), Next: 7279999}, ErrLocalIncompatibleOrStale},
	}
	for i, tt := range tests {
		filter := newFilter(params.MainnetChainConfig, params.MainnetGenesisHash, func() uint64 { return tt.head })
		if err := filter(tt.id); err != tt.err {
			t.Errorf("test %d: validation error mismatch: have %v, want %v", i, err, tt.err)
		}
	}
}

// Tests that IDs are properly RLP encoded (specifically important because we
// use uint32 to store the hash, but we need to encode it as [4]byte).
func TestEncoding(t *testing.T) {
	tests := []struct {
		id   ID
		want []byte
	}{
		{ID{Hash: checksumToBytes(0), Next: 0}, common.Hex2Bytes("c6840000000080")},
		{ID{Hash: checksumToBytes(0xdeadbeef), Next: 0xBADDCAFE}, common.Hex2Bytes("ca84deadbeef84baddcafe,")},
		{ID{Hash: checksumToBytes(math.MaxUint32), Next: math.MaxUint64}, common.Hex2Bytes("ce84ffffffff88ffffffffffffffff")},
	}
	for i, tt := range tests {
		have, err := rlp.EncodeToBytes(tt.id)
		if err != nil {
			t.Errorf("test %d: failed to encode forkid: %v", i, err)
			continue
		}
		if !bytes.Equal(have, tt.want) {
			t.Errorf("test %d: RLP mismatch: have %x, want %x", i, have, tt.want)
		}
	}
}
