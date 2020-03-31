package integration

import (
	"time"

	"github.com/ethereum/go-ethereum/p2p/simulations/adapters"

	"github.com/Fantom-foundation/lachesis-ex/app"
	"github.com/Fantom-foundation/lachesis-ex/gossip"
	"github.com/Fantom-foundation/lachesis-ex/lachesis"
)

// NewIntegration creates gossip service for the integration test
func NewIntegration(ctx *adapters.ServiceContext, network lachesis.Config) *gossip.Service {
	gossipCfg := gossip.DefaultConfig(network)
	appCfg := app.DefaultConfig(network)

	engine, adb, gdb := MakeEngine(ctx.Config.DataDir, &gossipCfg, &appCfg)
	abci := MakeABCI(network, adb)

	coinbase := SetAccountKey(
		ctx.NodeContext.AccountManager,
		ctx.Config.PrivateKey,
		"fakepassword",
	)

	gossipCfg.Emitter.Validator = coinbase.Address
	gossipCfg.Emitter.EmitIntervals.Max = 3 * time.Second
	gossipCfg.Emitter.EmitIntervals.SelfForkProtection = 0

	svc, err := gossip.NewService(ctx.NodeContext, &gossipCfg, gdb, engine, abci)
	if err != nil {
		panic(err)
	}

	return svc
}
