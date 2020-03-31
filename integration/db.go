package integration

import (
	"github.com/Fantom-foundation/lachesis-ex/kvdb"
	"github.com/Fantom-foundation/lachesis-ex/kvdb/leveldb"
	"github.com/Fantom-foundation/lachesis-ex/kvdb/memorydb"
)

func dbProducer(dbdir string) kvdb.DbProducer {
	if dbdir == "inmemory" || dbdir == "" {
		return memorydb.NewProducer("")
	}

	return leveldb.NewProducer(dbdir)
}
