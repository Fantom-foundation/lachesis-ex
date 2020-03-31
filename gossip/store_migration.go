package gossip

import (
	"github.com/Fantom-foundation/lachesis-ex/kvdb"
	"github.com/Fantom-foundation/lachesis-ex/kvdb/table"
	"github.com/Fantom-foundation/lachesis-ex/utils/migration"
)

func (s *Store) migrate() {
	versions := kvdb.NewIDStore(s.table.Version)
	err := s.migrations().Exec(versions)
	if err != nil {
		s.Log.Crit("gossip store migrations", "err", err)
	}
}

func (s *Store) migrations() *migration.Migration {
	return migration.
		Begin("lachesis-gossip-store").
		Next("service db",
			func() error {
				dst := table.New(s.serviceDb, []byte("Z")) // service.Peers

				old1 := table.New(s.mainDb, []byte("p")) // table.PackInfos
				err := kvdb.Move(old1, dst, []byte("serverPool"))
				if err != nil {
					return err
				}

				old2 := table.New(s.mainDb, []byte("Z"))
				err = kvdb.Move(old2, dst, nil)
				if err != nil {
					return err
				}

				s.Log.Warn("service db migration has been applied")
				return nil
			})
}
