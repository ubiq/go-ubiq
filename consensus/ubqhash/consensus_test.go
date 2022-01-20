// Copyright 2017 The go-ethereum Authors
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

package ubqhash

import (
	"encoding/binary"
	"encoding/json"
	"math/big"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/ubiq/go-ubiq/v7/common"
	"github.com/ubiq/go-ubiq/v7/common/math"
	"github.com/ubiq/go-ubiq/v7/core/types"
	"github.com/ubiq/go-ubiq/v7/params"
)

type diffTest struct {
	ParentTimestamp    uint64
	ParentDifficulty   *big.Int
	CurrentTimestamp   uint64
	CurrentBlocknumber *big.Int
	CurrentDifficulty  *big.Int
}

func (d *diffTest) UnmarshalJSON(b []byte) (err error) {
	var ext struct {
		ParentTimestamp    string
		ParentDifficulty   string
		CurrentTimestamp   string
		CurrentBlocknumber string
		CurrentDifficulty  string
	}
	if err := json.Unmarshal(b, &ext); err != nil {
		return err
	}

	d.ParentTimestamp = math.MustParseUint64(ext.ParentTimestamp)
	d.ParentDifficulty = math.MustParseBig256(ext.ParentDifficulty)
	d.CurrentTimestamp = math.MustParseUint64(ext.CurrentTimestamp)
	d.CurrentBlocknumber = math.MustParseBig256(ext.CurrentBlocknumber)
	d.CurrentDifficulty = math.MustParseBig256(ext.CurrentDifficulty)

	return nil
}

func TestCalcDifficulty(t *testing.T) {
	file, err := os.Open(filepath.Join("..", "..", "tests", "testdata", "BasicTests", "difficulty.json"))
	if err != nil {
		t.Skip(err)
	}
	defer file.Close()

	tests := make(map[string]diffTest)
	err = json.NewDecoder(file).Decode(&tests)
	if err != nil {
		t.Fatal(err)
	}

	config := &params.ChainConfig{HomesteadBlock: big.NewInt(1150000)}

	for name, test := range tests {
		number := new(big.Int).Sub(test.CurrentBlocknumber, big.NewInt(1))
		diff := CalcDifficulty(nil, config, test.CurrentTimestamp, &types.Header{
			Number:     number,
			Time:       test.ParentTimestamp,
			Difficulty: test.ParentDifficulty,
		})
		if diff.Cmp(test.CurrentDifficulty) != 0 {
			t.Error(name, "failed. Expected", test.CurrentDifficulty, "and calculated", diff)
		}
	}
}

func randSlice(min, max uint32) []byte {
	var b = make([]byte, 4)
	rand.Read(b)
	a := binary.LittleEndian.Uint32(b)
	size := min + a%(max-min)
	out := make([]byte, size)
	rand.Read(out)
	return out
}

func TestDifficultyCalculators(t *testing.T) {
	rand.Seed(2)
	for i := 0; i < 5000; i++ {
		// 1 to 300 seconds diff
		var timeDelta = uint64(1 + rand.Uint32()%3000)
		diffBig := big.NewInt(0).SetBytes(randSlice(2, 10))
		if diffBig.Cmp(params.MinimumDifficulty) < 0 {
			diffBig.Set(params.MinimumDifficulty)
		}
		//rand.Read(difficulty)
		header := &types.Header{
			Difficulty: diffBig,
			Number:     new(big.Int).SetUint64(rand.Uint64() % 50_000_000),
			Time:       rand.Uint64() - timeDelta,
		}
		if rand.Uint32()&1 == 0 {
			header.UncleHash = types.EmptyUncleHash
		}
		bombDelay := new(big.Int).SetUint64(rand.Uint64() % 50_000_000)
		for i, pair := range []struct {
			bigFn  func(time uint64, parent *types.Header) *big.Int
			u256Fn func(time uint64, parent *types.Header) *big.Int
		}{
			{FrontierDifficultyCalulator, CalcDifficultyFrontierU256},
			{HomesteadDifficultyCalulator, CalcDifficultyHomesteadU256},
			{DynamicDifficultyCalculator(bombDelay), MakeDifficultyCalculatorU256(bombDelay)},
			{DynamicDifficultyCalculator(nil), MakeDifficultyCalculatorU256(nil)},
		} {
			time := header.Time + timeDelta
			want := pair.bigFn(time, header)
			have := pair.u256Fn(time, header)
			if want.BitLen() > 256 {
				continue
			}
			if want.Cmp(have) != 0 {
				t.Fatalf("pair %d: want %x have %x\nparent.Number: %x\np.Time: %x\nc.Time: %x\nBombdelay: %v\n", i, want, have,
					header.Number, header.Time, time, bombDelay)
			}
		}
	}
}

func BenchmarkDifficultyCalculator(b *testing.B) {
	x1 := makeDifficultyCalculator(big.NewInt(1000000))
	x2 := MakeDifficultyCalculatorU256(big.NewInt(1000000))
	h := &types.Header{
		ParentHash: common.Hash{},
		UncleHash:  types.EmptyUncleHash,
		Difficulty: big.NewInt(0xffffff),
		Number:     big.NewInt(500000),
		Time:       1000000,
	}
	b.Run("big-frontier", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			calcDifficultyFrontier(1000014, h)
		}
	})
	b.Run("u256-frontier", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			CalcDifficultyFrontierU256(1000014, h)
		}
	})
	b.Run("big-homestead", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			calcDifficultyHomestead(1000014, h)
		}
	})
	b.Run("u256-homestead", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			CalcDifficultyHomesteadU256(1000014, h)
		}
	})
	b.Run("big-generic", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			x1(1000014, h)
		}
	})
	b.Run("u256-generic", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			x2(1000014, h)
		}
	})
}

func TestCalcBaseBlockReward(t *testing.T) {
	config := *params.MainnetChainConfig
	_, reward := CalcBaseBlockReward(config.Ubqhash, big.NewInt(1), false)
	if reward.Cmp(big.NewInt(8e+18)) != 0 {
		t.Error("TestCalcBaseBlockReward 8 (start)", "failed. Expected", big.NewInt(8e+18), "and calculated", reward)
	}
	_, reward = CalcBaseBlockReward(config.Ubqhash, big.NewInt(358363), false)
	if reward.Cmp(big.NewInt(8e+18)) != 0 {
		t.Error("TestCalcBaseBlockReward 8 (end)", "failed. Expected", big.NewInt(8e+18), "and calculated", reward)
	}

	_, reward = CalcBaseBlockReward(config.Ubqhash, big.NewInt(358364), false)
	if reward.Cmp(big.NewInt(7e+18)) != 0 {
		t.Error("TestCalcBaseBlockReward 7 (start)", "failed. Expected", big.NewInt(7e+18), "and calculated", reward)
	}
	_, reward = CalcBaseBlockReward(config.Ubqhash, big.NewInt(716727), false)
	if reward.Cmp(big.NewInt(7e+18)) != 0 {
		t.Error("TestCalcBaseBlockReward 7 (end)", "failed. Expected", big.NewInt(7e+18), "and calculated", reward)
	}

	_, reward = CalcBaseBlockReward(config.Ubqhash, big.NewInt(716728), false)
	if reward.Cmp(big.NewInt(6e+18)) != 0 {
		t.Error("TestCalcBaseBlockReward 6 (start)", "failed. Expected", big.NewInt(6e+18), "and calculated", reward)
	}
	_, reward = CalcBaseBlockReward(config.Ubqhash, big.NewInt(1075090), false)
	if reward.Cmp(big.NewInt(6e+18)) != 0 {
		t.Error("TestCalcBaseBlockReward 6 (end)", "failed. Expected", big.NewInt(6e+18), "and calculated", reward)
	}

	_, reward = CalcBaseBlockReward(config.Ubqhash, big.NewInt(1075091), false)
	if reward.Cmp(big.NewInt(5e+18)) != 0 {
		t.Error("TestCalcBaseBlockReward 5 (start)", "failed. Expected", big.NewInt(5e+18), "and calculated", reward)
	}
	_, reward = CalcBaseBlockReward(config.Ubqhash, big.NewInt(1433454), false)
	if reward.Cmp(big.NewInt(5e+18)) != 0 {
		t.Error("TestCalcBaseBlockReward 5 (end)", "failed. Expected", big.NewInt(5e+18), "and calculated", reward)
	}

	_, reward = CalcBaseBlockReward(config.Ubqhash, big.NewInt(1433455), false)
	if reward.Cmp(big.NewInt(4e+18)) != 0 {
		t.Error("TestCalcBaseBlockReward 4 (start)", "failed. Expected", big.NewInt(4e+18), "and calculated", reward)
	}
	_, reward = CalcBaseBlockReward(config.Ubqhash, big.NewInt(1791818), false)
	if reward.Cmp(big.NewInt(4e+18)) != 0 {
		t.Error("TestCalcBaseBlockReward 4 (end)", "failed. Expected", big.NewInt(4e+18), "and calculated", reward)
	}

	_, reward = CalcBaseBlockReward(config.Ubqhash, big.NewInt(1791819), false)
	if reward.Cmp(big.NewInt(3e+18)) != 0 {
		t.Error("TestCalcBaseBlockReward 3 (start)", "failed. Expected", big.NewInt(3e+18), "and calculated", reward)
	}
	_, reward = CalcBaseBlockReward(config.Ubqhash, big.NewInt(2150181), false)
	if reward.Cmp(big.NewInt(3e+18)) != 0 {
		t.Error("TestCalcBaseBlockReward 3 (end)", "failed. Expected", big.NewInt(3e+18), "and calculated", reward)
	}

	_, reward = CalcBaseBlockReward(config.Ubqhash, big.NewInt(2150182), false)
	if reward.Cmp(big.NewInt(2e+18)) != 0 {
		t.Error("TestCalcBaseBlockReward 2 (start)", "failed. Expected", big.NewInt(2e+18), "and calculated", reward)
	}
	_, reward = CalcBaseBlockReward(config.Ubqhash, big.NewInt(2508545), false)
	if reward.Cmp(big.NewInt(2e+18)) != 0 {
		t.Error("TestCalcBaseBlockReward 2 (end)", "failed. Expected", big.NewInt(2e+18), "and calculated", reward)
	}

	_, reward = CalcBaseBlockReward(config.Ubqhash, big.NewInt(2508546), false)
	if reward.Cmp(big.NewInt(1e+18)) != 0 {
		t.Error("TestCalcBaseBlockReward 1 (start)", "failed. Expected", big.NewInt(1e+18), "and calculated", reward)
	}

	_, reward = CalcBaseBlockReward(config.Ubqhash, big.NewInt(2000000), true)
	if reward.Cmp(big.NewInt(15e+17)) != 0 {
		t.Error("TestCalcBaseBlockReward (orion)", "failed. Expected", big.NewInt(15e+17), "and calculated", reward)
	}
}

func TestCalcUncleBlockReward(t *testing.T) {
	config := params.MainnetChainConfig
	reward := big.NewInt(8e+18)
	// depth 1
	u := CalcUncleBlockReward(config, big.NewInt(5), big.NewInt(4), reward)
	if u.Cmp(big.NewInt(4e+18)) != 0 {
		t.Error("TestCalcUncleBlockReward 8", "failed. Expected", big.NewInt(4e+18), "and calculated", u)
	}

	// depth 2
	u = CalcUncleBlockReward(config, big.NewInt(8), big.NewInt(6), reward)
	if u.Cmp(big.NewInt(0)) != 0 {
		t.Error("TestCalcUncleBlockReward 8", "failed. Expected", big.NewInt(0), "and calculated", u)
	}

	// depth 3 (before negative fix)
	u = CalcUncleBlockReward(config, big.NewInt(8), big.NewInt(5), reward)
	if u.Cmp(big.NewInt(-4e+18)) != 0 {
		t.Error("TestCalcUncleBlockReward 8", "failed. Expected", big.NewInt(-4e+18), "and calculated", u)
	}

	// depth 3 (after negative fix)
	u = CalcUncleBlockReward(config, big.NewInt(10), big.NewInt(7), reward)
	if u.Cmp(big.NewInt(0)) != 0 {
		t.Error("TestCalcUncleBlockReward 8", "failed. Expected", big.NewInt(0), "and calculated", u)
	}

	_, reward = CalcBaseBlockReward(config.Ubqhash, big.NewInt(358364), false)
	expected := big.NewInt(35e+17)
	// depth 1 (after stepdown)
	u = CalcUncleBlockReward(config, big.NewInt(8), big.NewInt(7), reward)
	if u.Cmp(expected) != 0 {
		t.Error("TestCalcUncleBlockReward 7", "failed. Expected", expected, "and calculated", u)
	}

	_, reward = CalcBaseBlockReward(config.Ubqhash, big.NewInt(1075091), false)
	expected = big.NewInt(25e+17)
	// depth 1 (after stepdown)
	u = CalcUncleBlockReward(config, big.NewInt(8), big.NewInt(7), reward)
	if u.Cmp(expected) != 0 {
		t.Error("TestCalcUncleBlockReward 5", "failed. Expected", expected, "and calculated", u)
	}
}
