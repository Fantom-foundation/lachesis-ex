package poset

import (
	"github.com/Fantom-foundation/lachesis-ex/kvdb"
	"github.com/Fantom-foundation/lachesis-ex/utils/migration"
)

func (s *Store) migrate() {
	versions := kvdb.NewIDStore(s.table.Version)
	err := s.migrations().Exec(versions)
	if err != nil {
		s.Log.Crit("poset store migrations", "err", err)
	}
}

func (s *Store) migrations() *migration.Migration {
	return migration.Begin("lachesis-poset-store")
}
