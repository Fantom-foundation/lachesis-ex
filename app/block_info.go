package app

import (
	"github.com/Fantom-foundation/lachesis-ex/inter"
	"github.com/Fantom-foundation/lachesis-ex/inter/idx"
)

// BlockInfo includes only necessary information about inter.Block
type BlockInfo struct {
	Index idx.Block
	Time  inter.Timestamp
}

func blockInfo(b *inter.Block) *BlockInfo {
	return &BlockInfo{
		Index: b.Index,
		Time:  b.Time,
	}
}
