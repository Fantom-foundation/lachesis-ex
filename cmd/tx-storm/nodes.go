package main

import (
	"fmt"
	"github.com/Fantom-foundation/lachesis-ex/logger"
	"github.com/Fantom-foundation/lachesis-ex/utils"
)

var (
	max float64
	maxAvg float64
)

// Nodes pool.
type Nodes struct {
	tps    chan float64
	conns  []*Sender
	blocks chan Block
	done   chan struct{}
	logger.Instance
}

func NewNodes(cfg *Config, input <-chan *Transaction) *Nodes {
	n := &Nodes{
		tps:      make(chan float64, 1),
		blocks:   make(chan Block, 1),
		done:     make(chan struct{}),
		Instance: logger.MakeInstance(),
	}
	was := make(map[string]struct{}, len(cfg.URLs))
	for _, url := range cfg.URLs {
		_, double := was[url]
		n.add(url, !double)
		was[url] = struct{}{}
	}

	n.notifyTPS(0)
	go n.background(input)
	go n.measureTPS()
	return n
}

func (n *Nodes) Count() int {
	return len(n.conns)
}

func (n *Nodes) TPS() <-chan float64 {
	return n.tps
}

func (n *Nodes) notifyTPS(tps float64) {
	select {
	case n.tps <- tps:
		break
	default:
	}
}

func (n *Nodes) measureTPS() {
	var (
		lastBlock Block
		avgbuff   = utils.NewAvgBuff(100)
	)
	for b := range n.blocks {
		if lastBlock.Number != nil && b.Number.Cmp(lastBlock.Number) < 1 || b.TxCount == 0 {
			continue
		}

		txCountGotMeter.Inc(int64(b.TxCount))

		dur := b.Timestamp.Sub(lastBlock.Timestamp).Seconds()
		if dur == 0 {
			continue
		}
		tps := float64(b.TxCount) / dur
		// Protection for first block (current time - 0)
		if dur < 1000.0 {
			avgbuff.Push(float64(b.TxCount), dur)
		}
		if tps > max {
			max = tps
		}

		txTpsMeter.Update(int64(tps))

		lastBlock = b
		avg := avgbuff.Avg()
		if avg > maxAvg {
			maxAvg = avg
		}

		n.notifyTPS(avg)
		n.Log.Info("TPS", "block", b.Number, "value", tps, "avg", avg, "max", max, "maxAvg", maxAvg)
	}
}

func (n *Nodes) add(url string, is1st bool) {
	c := NewSender(url, is1st)
	c.SetName(fmt.Sprintf("Node-%d", len(n.conns)))
	n.conns = append(n.conns, c)

	go func() {
		defer n.stop()
		for b := range c.Blocks() {
			n.blocks <- b
		}
	}()
}

func (n *Nodes) stop() {
	// TODO: mutex
	close(n.blocks)
}

func (n *Nodes) background(input <-chan *Transaction) {
	if len(n.conns) < 1 {
		panic("no connections")
	}

	i := 0
	for tx := range input {
		c := n.conns[i]
		c.Send(tx)
		i = (i + 1) % len(n.conns)
	}

	for _, c := range n.conns {
		c.Close()
	}
}

