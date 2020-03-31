package app

import (
	"github.com/Fantom-foundation/lachesis-ex/kvdb"
	"github.com/Fantom-foundation/lachesis-ex/kvdb/flushable"
	"github.com/Fantom-foundation/lachesis-ex/utils/migration"
)

func (s *Store) migrate() {
	versions := kvdb.NewIDStore(s.table.Version)
	err := s.migrations(s.dbs).Exec(versions)
	if err != nil {
		s.Log.Crit("app store migrations", "err", err)
	}
}

func (s *Store) migrations(dbs *flushable.SyncedPool) *migration.Migration {
	return migration.
		Begin("lachesis-app-store")
}
