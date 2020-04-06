package gossip

import (
	"github.com/ethereum/go-ethereum/core/types"
)

const (
	// This is the target size for the packs of transactions sent by txsyncLoop.
	// A pack can get larger than this if a single transactions exceeds this size.
	txsyncPackSize = 100 * 1024
)

type txsync struct {
	p   *peer
	txs []*types.Transaction
}

// syncer is responsible for periodically synchronising with the network, both
// downloading hashes and events as well as handling the announcement handler.
func (pm *ProtocolManager) syncer() {
	// Start and ensure cleanup of sync mechanisms
	pm.fetcher.Start()
	defer pm.fetcher.Stop()
	defer pm.downloader.Terminate()

	for {
		select {
		case <-pm.newPeerCh:
		case <-pm.noMorePeers:
			return
		}
	}
}
