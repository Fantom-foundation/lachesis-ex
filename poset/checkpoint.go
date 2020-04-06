package poset

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/Fantom-foundation/lachesis-ex/hash"
	"github.com/Fantom-foundation/lachesis-ex/inter"
	"github.com/Fantom-foundation/lachesis-ex/inter/idx"
	"github.com/Fantom-foundation/lachesis-ex/poset/election"
	"github.com/Fantom-foundation/lachesis-ex/vector"
)

// Checkpoint is for persistent storing.
type Checkpoint struct {
	// fields can change only after a frame is decided
	LastDecidedFrame idx.Frame
	LastBlockN       idx.Block
	LastAtropos      hash.Event
	AppHash          common.Hash
}

/*
 * Poset's methods:
 */

// State saves Checkpoint.
func (p *Poset) saveCheckpoint() {
	p.store.SetCheckpoint(p.Checkpoint)
}

// Bootstrap restores poset's state from store.
func (p *Poset) Bootstrap(callback inter.ConsensusCallbacks) {
	if p.Checkpoint != nil {
		return
	}
	// block handler must be set before p.handleElection
	p.callback = callback

	// restore Checkpoint
	p.Checkpoint = p.store.GetCheckpoint()
	if p.Checkpoint == nil {
		p.Log.Crit("Apply genesis for store first")
	}

	// restore current epoch
	p.loadEpoch()
	p.vecClock = vector.NewIndex(p.dag.VectorClockConfig, p.Validators, p.store.epochTable.VectorIndex, func(id hash.Event) *inter.EventHeaderData {
		return p.input.GetEventHeader(p.EpochN, id)
	})
	p.election = election.New(p.Validators, p.LastDecidedFrame+1, p.vecClock.Cause, p.store.GetFrameRoots)

	// events reprocessing
	p.handleElection(nil)
}
