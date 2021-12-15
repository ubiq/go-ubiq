// Copyright 2020 The go-ethereum Authors
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
	"fmt"
	"math/big"

	"github.com/holiman/uint256"
	"github.com/ubiq/go-ubiq/v5/consensus"
	"github.com/ubiq/go-ubiq/v5/core/types"
	"github.com/ubiq/go-ubiq/v5/log"
	"github.com/ubiq/go-ubiq/v5/params"
)

const (
	// frontierDurationLimit is for Frontier:
	// The decision boundary on the blocktime duration used to determine
	// whether difficulty should go up or down.
	frontierDurationLimit = 13
	// minimumDifficulty The minimum that the difficulty may ever be.
	minimumDifficulty = 131072
	// expDiffPeriod is the exponential difficulty period
	expDiffPeriodUint = 100000
	// difficultyBoundDivisorBitShift is the bound divisor of the difficulty (2048),
	// This constant is the right-shifts to use for the division.
	difficultyBoundDivisor = 11
)

// Diff algo constants.
var (
	big88 = big.NewInt(88)

	digishieldV3Config = &diffConfig{
		AveragingWindow: big.NewInt(21),
		MaxAdjustDown:   big.NewInt(16), // 16%
		MaxAdjustUp:     big.NewInt(8),  // 8%
		Factor:          big.NewInt(100),
	}

	digishieldV3ModConfig = &diffConfig{
		AveragingWindow: big.NewInt(88),
		MaxAdjustDown:   big.NewInt(3), // 3%
		MaxAdjustUp:     big.NewInt(2), // 2%
		Factor:          big.NewInt(100),
	}

	fluxConfig = &diffConfig{
		AveragingWindow: big.NewInt(88),
		MaxAdjustDown:   big.NewInt(5), // 0.5%
		MaxAdjustUp:     big.NewInt(3), // 0.3%
		Dampen:          big.NewInt(1), // 0.1%
		Factor:          big.NewInt(1000),
	}
)

type diffConfig struct {
	AveragingWindow *big.Int `json:"averagingWindow"`
	MaxAdjustDown   *big.Int `json:"maxAdjustDown"`
	MaxAdjustUp     *big.Int `json:"maxAdjustUp"`
	Dampen          *big.Int `json:"dampen,omitempty"`
	Factor          *big.Int `json:"factor"`
}

// Difficulty timespans
func averagingWindowTimespan(config *diffConfig) *big.Int {
	x := new(big.Int)
	return x.Mul(config.AveragingWindow, big88)
}

func minActualTimespan(config *diffConfig, dampen bool) *big.Int {
	x := new(big.Int)
	y := new(big.Int)
	z := new(big.Int)
	if dampen {
		x.Sub(config.Factor, config.Dampen)
		y.Mul(averagingWindowTimespan(config), x)
		z.Div(y, config.Factor)
	} else {
		x.Sub(config.Factor, config.MaxAdjustUp)
		y.Mul(averagingWindowTimespan(config), x)
		z.Div(y, config.Factor)
	}
	return z
}

func maxActualTimespan(config *diffConfig, dampen bool) *big.Int {
	x := new(big.Int)
	y := new(big.Int)
	z := new(big.Int)
	if dampen {
		x.Add(config.Factor, config.Dampen)
		y.Mul(averagingWindowTimespan(config), x)
		z.Div(y, config.Factor)
	} else {
		x.Add(config.Factor, config.MaxAdjustDown)
		y.Mul(averagingWindowTimespan(config), x)
		z.Div(y, config.Factor)
	}
	return z
}

// CalcDifficultyFrontierU256 is the difficulty adjustment algorithm. It returns the
// difficulty that a new block should have when created at time given the parent
// block's time and difficulty. The calculation uses the Frontier rules.
func CalcDifficultyFrontierU256(time uint64, parent *types.Header) *big.Int {
	/*
		Algorithm
		block_diff = pdiff + pdiff / 2048 * (1 if time - ptime < 13 else -1) + int(2^((num // 100000) - 2))

		Where:
		- pdiff  = parent.difficulty
		- ptime = parent.time
		- time = block.timestamp
		- num = block.number
	*/

	pDiff, _ := uint256.FromBig(parent.Difficulty) // pDiff: pdiff
	adjust := pDiff.Clone()
	adjust.Rsh(adjust, difficultyBoundDivisor) // adjust: pDiff / 2048

	if time-parent.Time < frontierDurationLimit {
		pDiff.Add(pDiff, adjust)
	} else {
		pDiff.Sub(pDiff, adjust)
	}
	if pDiff.LtUint64(minimumDifficulty) {
		pDiff.SetUint64(minimumDifficulty)
	}
	// 'pdiff' now contains:
	// pdiff + pdiff / 2048 * (1 if time - ptime < 13 else -1)

	if periodCount := (parent.Number.Uint64() + 1) / expDiffPeriodUint; periodCount > 1 {
		// diff = diff + 2^(periodCount - 2)
		expDiff := adjust.SetOne()
		expDiff.Lsh(expDiff, uint(periodCount-2)) // expdiff: 2 ^ (periodCount -2)
		pDiff.Add(pDiff, expDiff)
	}
	return pDiff.ToBig()
}

// CalcDifficultyHomesteadU256 is the difficulty adjustment algorithm. It returns
// the difficulty that a new block should have when created at time given the
// parent block's time and difficulty. The calculation uses the Homestead rules.
func CalcDifficultyHomesteadU256(time uint64, parent *types.Header) *big.Int {
	/*
		https://github.com/ethereum/EIPs/blob/master/EIPS/eip-2.md
		Algorithm:
		block_diff = pdiff + pdiff / 2048 * max(1 - (time - ptime) / 10, -99) + 2 ^ int((num / 100000) - 2))

		Our modification, to use unsigned ints:
		block_diff = pdiff - pdiff / 2048 * max((time - ptime) / 10 - 1, 99) + 2 ^ int((num / 100000) - 2))

		Where:
		- pdiff  = parent.difficulty
		- ptime = parent.time
		- time = block.timestamp
		- num = block.number
	*/

	pDiff, _ := uint256.FromBig(parent.Difficulty) // pDiff: pdiff
	adjust := pDiff.Clone()
	adjust.Rsh(adjust, difficultyBoundDivisor) // adjust: pDiff / 2048

	x := (time - parent.Time) / 10 // (time - ptime) / 10)
	var neg = true
	if x == 0 {
		x = 1
		neg = false
	} else if x >= 100 {
		x = 99
	} else {
		x = x - 1
	}
	z := new(uint256.Int).SetUint64(x)
	adjust.Mul(adjust, z) // adjust: (pdiff / 2048) * max((time - ptime) / 10 - 1, 99)
	if neg {
		pDiff.Sub(pDiff, adjust) // pdiff - pdiff / 2048 * max((time - ptime) / 10 - 1, 99)
	} else {
		pDiff.Add(pDiff, adjust) // pdiff + pdiff / 2048 * max((time - ptime) / 10 - 1, 99)
	}
	if pDiff.LtUint64(minimumDifficulty) {
		pDiff.SetUint64(minimumDifficulty)
	}
	// for the exponential factor, a.k.a "the bomb"
	// diff = diff + 2^(periodCount - 2)
	if periodCount := (1 + parent.Number.Uint64()) / expDiffPeriodUint; periodCount > 1 {
		expFactor := adjust.Lsh(adjust.SetOne(), uint(periodCount-2))
		pDiff.Add(pDiff, expFactor)
	}
	return pDiff.ToBig()
}

// MakeDifficultyCalculatorU256 creates a difficultyCalculator with the given bomb-delay.
// the difficulty is calculated with Byzantium rules, which differs from Homestead in
// how uncles affect the calculation
func MakeDifficultyCalculatorU256(bombDelay *big.Int) func(time uint64, parent *types.Header) *big.Int {
	// Note, the calculations below looks at the parent number, which is 1 below
	// the block number. Thus we remove one from the delay given
	medianBlockTime := 15 // Ubiq - Orion
	bombDelayFromParent := new(big.Int).Uint64()
	if bombDelay != nil {
		bombDelayFromParent = bombDelay.Uint64() - 1
		medianBlockTime = 9 // Ethereum
	}
	return func(time uint64, parent *types.Header) *big.Int {
		/*
			https://github.com/ethereum/EIPs/issues/100
			pDiff = parent.difficulty
			BLOCK_DIFF_FACTOR = 9
			a = pDiff + (pDiff // BLOCK_DIFF_FACTOR) * adj_factor
			b = min(parent.difficulty, MIN_DIFF)
			child_diff = max(a,b )
		*/
		x := (time - parent.Time) / uint64(medianBlockTime) // (block_timestamp - parent_timestamp) // 15
		c := uint64(1)                                      // if parent.unclehash == emptyUncleHashHash
		if parent.UncleHash != types.EmptyUncleHash {
			c = 2
		}
		xNeg := x >= c
		if xNeg {
			// x is now _negative_ adjustment factor
			x = x - c // - ( (t-p)/p -( 2 or 1) )
		} else {
			x = c - x // (2 or 1) - (t-p)/15
		}
		if x > 99 {
			x = 99 // max(x, 99)
		}
		// parent_diff + (parent_diff / 2048 * max((2 if len(parent.uncles) else 1) - ((timestamp - parent.timestamp) // 15), -99))
		y := new(uint256.Int)
		y.SetFromBig(parent.Difficulty)    // y: p_diff
		pDiff := y.Clone()                 // pdiff: p_diff
		z := new(uint256.Int).SetUint64(x) //z : +-adj_factor (either pos or negative)
		y.Rsh(y, difficultyBoundDivisor)   // y: p__diff / 2048
		z.Mul(y, z)                        // z: (p_diff / 2048 ) * (+- adj_factor)

		if xNeg {
			y.Sub(pDiff, z) // y: parent_diff + parent_diff/2048 * adjustment_factor
		} else {
			y.Add(pDiff, z) // y: parent_diff + parent_diff/2048 * adjustment_factor
		}
		// minimum difficulty can ever be (before exponential factor)
		if y.LtUint64(minimumDifficulty) {
			y.SetUint64(minimumDifficulty)
		}
		if bombDelay != nil {
			// calculate a fake block number for the ice-age delay
			// Specification: https://eips.ethereum.org/EIPS/eip-1234
			var pNum = parent.Number.Uint64()
			if pNum >= bombDelayFromParent {
				if fakeBlockNumber := pNum - bombDelayFromParent; fakeBlockNumber >= 2*expDiffPeriodUint {
					z.SetOne()
					z.Lsh(z, uint(fakeBlockNumber/expDiffPeriodUint-2))
					y.Add(z, y)
				}
			}
		}
		return y.ToBig()
	}
}

// calcDifficultyDigishieldV3 is the original ubiq difficulty adjustment algorithm.
// It returns the difficulty that a new block should have when created at time
// given the parent block's time and difficulty.
// Based on Digibyte's Digishield v3 retargeting
func CalcDifficultyDigishieldV3(chain consensus.ChainHeaderReader, parentNumber, parentDiff *big.Int, parent *types.Header, digishield *diffConfig) *big.Int {
	// holds intermediate values to make the algo easier to read & audit
	x := new(big.Int)
	nFirstBlock := new(big.Int)
	nFirstBlock.Sub(parentNumber, digishield.AveragingWindow)

	log.Debug(fmt.Sprintf("CalcDifficulty parentNumber: %v parentDiff: %v", parentNumber, parentDiff))

	// Check we have enough blocks
	if parentNumber.Cmp(digishield.AveragingWindow) < 1 {
		log.Debug(fmt.Sprintf("CalcDifficulty: parentNumber(%+x) < digishieldV3Config.AveragingWindow(%+x)", parentNumber, digishield.AveragingWindow))
		x.Set(parentDiff)
		return x
	}

	// Limit adjustment step
	// Use medians to prevent time-warp attacks
	nLastBlockTime := chain.CalcPastMedianTime(parentNumber.Uint64(), parent)
	nFirstBlockTime := chain.CalcPastMedianTime(nFirstBlock.Uint64(), parent)
	nActualTimespan := new(big.Int)
	nActualTimespan.Sub(nLastBlockTime, nFirstBlockTime)
	log.Debug(fmt.Sprintf("CalcDifficulty nActualTimespan = %v before dampening", nActualTimespan))

	y := new(big.Int)
	y.Sub(nActualTimespan, averagingWindowTimespan(digishield))
	y.Div(y, big.NewInt(4))
	nActualTimespan.Add(y, averagingWindowTimespan(digishield))
	log.Debug(fmt.Sprintf("CalcDifficulty nActualTimespan = %v before bounds", nActualTimespan))

	if nActualTimespan.Cmp(minActualTimespan(digishield, false)) < 0 {
		nActualTimespan.Set(minActualTimespan(digishield, false))
		log.Debug("CalcDifficulty Minimum Timespan set")
	} else if nActualTimespan.Cmp(maxActualTimespan(digishield, false)) > 0 {
		nActualTimespan.Set(maxActualTimespan(digishield, false))
		log.Debug("CalcDifficulty Maximum Timespan set")
	}

	log.Debug(fmt.Sprintf("CalcDifficulty nActualTimespan = %v final\n", nActualTimespan))

	// Retarget
	x.Mul(parentDiff, averagingWindowTimespan(digishield))
	log.Debug(fmt.Sprintf("CalcDifficulty parentDiff * AveragingWindowTimespan: %v", x))

	x.Div(x, nActualTimespan)
	log.Debug(fmt.Sprintf("CalcDifficulty x / nActualTimespan: %v", x))

	if x.Cmp(params.MinimumDifficulty) < 0 {
		x.Set(params.MinimumDifficulty)
	}

	return x
}

func CalcDifficultyFlux(chain consensus.ChainHeaderReader, time, parentTime, parentNumber, parentDiff *big.Int, parent *types.Header) *big.Int {
	x := new(big.Int)
	nFirstBlock := new(big.Int)
	nFirstBlock.Sub(parentNumber, fluxConfig.AveragingWindow)

	// Check we have enough blocks
	if parentNumber.Cmp(fluxConfig.AveragingWindow) < 1 {
		log.Debug(fmt.Sprintf("CalcDifficulty: parentNumber(%+x) < fluxConfig.AveragingWindow(%+x)", parentNumber, fluxConfig.AveragingWindow))
		x.Set(parentDiff)
		return x
	}

	diffTime := new(big.Int)
	diffTime.Sub(time, parentTime)

	nLastBlockTime := chain.CalcPastMedianTime(parentNumber.Uint64(), parent)
	nFirstBlockTime := chain.CalcPastMedianTime(nFirstBlock.Uint64(), parent)
	nActualTimespan := new(big.Int)
	nActualTimespan.Sub(nLastBlockTime, nFirstBlockTime)

	y := new(big.Int)
	y.Sub(nActualTimespan, averagingWindowTimespan(fluxConfig))
	y.Div(y, big.NewInt(4))
	nActualTimespan.Add(y, averagingWindowTimespan(fluxConfig))

	if nActualTimespan.Cmp(minActualTimespan(fluxConfig, false)) < 0 {
		doubleBig88 := new(big.Int)
		doubleBig88.Mul(big88, big.NewInt(2))
		if diffTime.Cmp(doubleBig88) > 0 {
			nActualTimespan.Set(minActualTimespan(fluxConfig, true))
		} else {
			nActualTimespan.Set(minActualTimespan(fluxConfig, false))
		}
	} else if nActualTimespan.Cmp(maxActualTimespan(fluxConfig, false)) > 0 {
		halfBig88 := new(big.Int)
		halfBig88.Div(big88, big.NewInt(2))
		if diffTime.Cmp(halfBig88) < 0 {
			nActualTimespan.Set(maxActualTimespan(fluxConfig, true))
		} else {
			nActualTimespan.Set(maxActualTimespan(fluxConfig, false))
		}
	}

	x.Mul(parentDiff, averagingWindowTimespan(fluxConfig))
	x.Div(x, nActualTimespan)

	if x.Cmp(params.MinimumDifficulty) < 0 {
		x.Set(params.MinimumDifficulty)
	}

	return x
}
