package main

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/rlp"

	"github.com/Fantom-foundation/lachesis-ex/gossip"
	"github.com/Fantom-foundation/lachesis-ex/inter/idx"
	"github.com/Fantom-foundation/lachesis-ex/kvdb"
	"github.com/Fantom-foundation/lachesis-ex/kvdb/table"
)

func checkPacks(p kvdb.DbProducer) {
	db := p.OpenDb("gossip-main")
	defer db.Close()

	t := table.New(db, []byte("p"))

	it := t.NewIterator()
	defer it.Release()

	for it.Next() {
		buf := it.Key()
		w := it.Value()

		if strings.HasPrefix(string(buf), "serverPool") {
			fmt.Printf("skip %s key\n", string(buf))
			continue
		}

		var info gossip.PackInfo
		err := rlp.DecodeBytes(w, &info)
		if err != nil {
			fmt.Printf(">>> %s\n ", string(buf))
			continue
		}

		epoch := idx.BytesToEpoch(buf[0:4])
		pack := idx.BytesToEpoch(buf[4:8])
		fmt.Printf("%d:%d %+v\n", epoch, pack, info)
	}
}
