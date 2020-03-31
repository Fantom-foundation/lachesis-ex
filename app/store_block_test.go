package app

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Fantom-foundation/lachesis-ex/inter"
	"github.com/Fantom-foundation/lachesis-ex/inter/idx"
	"github.com/Fantom-foundation/lachesis-ex/logger"
)

func TestStoreGetBlock(t *testing.T) {
	logger.SetTestMode(t)

	expect := fakeBlock()
	store := cachedStore()
	store.SetBlock(expect)

	got := store.GetBlock(expect.Index)
	assert.EqualValues(t, expect, got)
}

func BenchmarkStoreGetBlock(b *testing.B) {
	logger.SetTestMode(b)

	b.Run("cache on", func(b *testing.B) {
		benchStoreGetBlock(b, cachedStore())
	})
	b.Run("cache off", func(b *testing.B) {
		benchStoreGetBlock(b, nonCachedStore())
	})
}

func benchStoreGetBlock(b *testing.B, store *Store) {
	block := fakeBlock()

	store.SetBlock(block)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if store.GetBlock(block.Index) == nil {
			b.Fatal("invalid result")
		}
	}
}

func BenchmarkStoreSetBlock(b *testing.B) {
	logger.SetTestMode(b)

	b.Run("cache on", func(b *testing.B) {
		benchStoreSetBlock(b, cachedStore())
	})
	b.Run("cache off", func(b *testing.B) {
		benchStoreSetBlock(b, nonCachedStore())
	})
}

func benchStoreSetBlock(b *testing.B, store *Store) {
	block := fakeBlock()

	for i := 0; i < b.N; i++ {
		store.SetBlock(block)
	}
}

func fakeBlock() *BlockInfo {
	return &BlockInfo{
		Index: idx.Block(1),
		Time:  inter.Timestamp(rand.Int63()),
	}
}
