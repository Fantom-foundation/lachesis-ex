package app

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/Fantom-foundation/lachesis-ex/evmcore"
	"github.com/Fantom-foundation/lachesis-ex/hash"
	"github.com/Fantom-foundation/lachesis-ex/inter"
	"github.com/Fantom-foundation/lachesis-ex/inter/idx"
)

// TxPosition includes block and event.
type TxPosition struct {
	Block       idx.Block
	Event       hash.Event
	Creator     idx.StakerID
	EventOffset uint32
	BlockOffset uint32
}

// BlocksMissed is information about missed blocks from a staker.
type BlocksMissed struct {
	Num    idx.Block
	Period inter.Timestamp
}

// updateOriginationScores calculates the origination scores
func (a *App) updateOriginationScores(
	epoch idx.Epoch,
	evmBlock *evmcore.EvmBlock,
	receipts types.Receipts,
	txPositions map[common.Hash]TxPosition,
) {
	for i, tx := range evmBlock.Transactions {
		txEventPos := txPositions[receipts[i].TxHash]
		txFee := new(big.Int).Mul(new(big.Int).SetUint64(receipts[i].GasUsed), tx.GasPrice())
		a.store.AddDirtyOriginationScore(txEventPos.Creator, txFee)
	}

	if a.ctx.sealEpoch {
		a.store.DelAllActiveOriginationScores()
		a.store.MoveDirtyOriginationScoresToActive()
	}
}

// updateValidationScores calculates the validation scores
func (a *App) updateValidationScores(
	epoch idx.Epoch,
	block *inter.Block,
	blockParticipated map[idx.StakerID]bool,
) {
	blockTimeDiff := block.Time - a.blockTime(block.Index-1)

	// Calc validation scores
	for _, it := range a.store.GetActiveSfcStakers() {
		// validators only
		if !a.store.HasEpochValidator(epoch, it.StakerID) {
			continue
		}

		// Check if validator has confirmed events by this Atropos
		missedBlock := !blockParticipated[it.StakerID]

		// If have no confirmed events by this Atropos - just add missed blocks for validator
		if missedBlock {
			a.store.IncBlocksMissed(it.StakerID, blockTimeDiff)
			continue
		}

		missedNum := a.store.GetBlocksMissed(it.StakerID).Num
		if missedNum > a.config.Net.Economy.BlockMissedLatency {
			missedNum = a.config.Net.Economy.BlockMissedLatency
		}

		// Add score for previous blocks, but no more than FrameLatency prev blocks
		a.store.AddDirtyValidationScore(it.StakerID, new(big.Int).SetUint64(uint64(blockTimeDiff)))
		for i := idx.Block(1); i <= missedNum && i < block.Index; i++ {
			curBlockTime := a.blockTime(block.Index - i)
			prevBlockTime := a.blockTime(block.Index - i - 1)
			timeDiff := curBlockTime - prevBlockTime
			a.store.AddDirtyValidationScore(it.StakerID, new(big.Int).SetUint64(uint64(timeDiff)))
		}
		a.store.ResetBlocksMissed(it.StakerID)
	}

	if !a.ctx.sealEpoch {
		return
	}

	if a.config.EpochDowntimeIndex {
		a.store.NewDowntimeSnapshotEpoch(epoch)
	}
	if a.config.EpochActiveValidationScoreIndex {
		a.store.NewScoreSnapshotEpoch(epoch)
	}

	a.store.DelAllActiveValidationScores()
	a.store.MoveDirtyValidationScoresToActive()
}
