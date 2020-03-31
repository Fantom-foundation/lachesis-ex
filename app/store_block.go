package app

import (
	"github.com/Fantom-foundation/lachesis-ex/inter/idx"
)

// SetBlock stores chain block.
func (s *Store) SetBlock(b *BlockInfo) {
	s.set(s.table.Blocks, b.Index.Bytes(), b)

	// Add to LRU cache.
	if b != nil && s.cache.Blocks != nil {
		s.cache.Blocks.Add(b.Index, b)
	}
}

// GetBlock returns stored block.
func (s *Store) GetBlock(n idx.Block) *BlockInfo {
	// Get block from LRU cache first.
	if s.cache.Blocks != nil {
		if c, ok := s.cache.Blocks.Get(n); ok {
			if b, ok := c.(*BlockInfo); ok {
				return b
			}
		}
	}

	block, _ := s.get(s.table.Blocks, n.Bytes(), &BlockInfo{}).(*BlockInfo)

	// Add to LRU cache.
	if block != nil && s.cache.Blocks != nil {
		s.cache.Blocks.Add(n, block)
	}

	return block
}
