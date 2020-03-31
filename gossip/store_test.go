package gossip

import (
	"time"

	"github.com/Fantom-foundation/lachesis-ex/kvdb"
	"github.com/Fantom-foundation/lachesis-ex/kvdb/flushable"
	"github.com/Fantom-foundation/lachesis-ex/kvdb/leveldb"
	"github.com/Fantom-foundation/lachesis-ex/kvdb/memorydb"
)

func cachedStore() *Store {
	mems := memorydb.NewProducer("", withDelay)
	dbs := flushable.NewSyncedPool(mems)
	cfg := LiteStoreConfig()

	return NewStore(dbs, cfg)
}

func nonCachedStore() *Store {
	mems := memorydb.NewProducer("", withDelay)
	dbs := flushable.NewSyncedPool(mems)
	cfg := StoreConfig{}

	return NewStore(dbs, cfg)
}

func realStore(dir string) *Store {
	disk := leveldb.NewProducer(dir)
	dbs := flushable.NewSyncedPool(disk)
	cfg := LiteStoreConfig()

	return NewStore(dbs, cfg)
}

func withDelay(db kvdb.KeyValueStore) kvdb.KeyValueStore {
	mem, ok := db.(*memorydb.Database)
	if ok {
		mem.SetDelay(time.Millisecond)

	}

	return db
}
