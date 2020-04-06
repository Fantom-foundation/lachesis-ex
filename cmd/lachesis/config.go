package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/params"
	"gopkg.in/urfave/cli.v1"

	"github.com/Fantom-foundation/lachesis-ex/app"
	"github.com/Fantom-foundation/lachesis-ex/evmcore"
	"github.com/Fantom-foundation/lachesis-ex/gossip"
	"github.com/Fantom-foundation/lachesis-ex/gossip/gasprice"
	"github.com/Fantom-foundation/lachesis-ex/lachesis"
	"github.com/Fantom-foundation/lachesis-ex/utils/toml"
)

var (
	dumpConfigCommand = cli.Command{
		Action:      utils.MigrateFlags(dumpConfig),
		Name:        "dumpconfig",
		Usage:       "Show configuration values",
		ArgsUsage:   "",
		Flags:       append(nodeFlags, testFlags...),
		Category:    "MISCELLANEOUS COMMANDS",
		Description: `The dumpconfig command shows configuration values.`,
	}

	configFileFlag = cli.StringFlag{
		Name:  "config",
		Usage: "TOML configuration file",
	}

	// GpoDefaultFlag defines a starting gas price for the oracle (GPO)
	GpoDefaultFlag = utils.BigFlag{
		Name:  "gpofloor",
		Usage: "The default suggested gas price",
		Value: big.NewInt(params.GWei),
	}

	// DataDirFlag defines directory to store Lachesis state and user's wallets
	DataDirFlag = utils.DirectoryFlag{
		Name:  "datadir",
		Usage: "Data directory for the databases and keystore",
		Value: utils.DirectoryString(DefaultDataDir()),
	}
)

type config struct {
	Version  string
	Node     node.Config
	Lachesis gossip.Config
	App      app.Config
}

func (cfg *config) Load(file string) error {
	data, err := toml.ParseFile(file)
	if err != nil {
		return err
	}

	changed, err := cfg.migrate(data)
	if err != nil {
		return err
	}

	err = toml.Settings.UnmarshalTable(data, cfg)
	if err != nil {
		err = errors.New(file + ", " + err.Error())
		return err
	}

	if !changed {
		return nil
	}

	_, fname := filepath.Split(file)
	updated, err := ioutil.TempFile("", fname)
	if err != nil {
		return err
	}

	err = toml.WriteTo(data, updated)
	updated.Close()
	if err != nil {
		return err
	}

	err = copyFile(file, file+".bak")
	if err != nil {
		return err
	}
	err = copyFile(updated.Name(), file)
	if err != nil {
		return err
	}

	return nil
}

func defaultLachesisConfig(ctx *cli.Context) lachesis.Config {
	var cfg lachesis.Config

	switch {
	case ctx.GlobalIsSet(FakeNetFlag.Name):
		_, accs, err := parseFakeGen(ctx.GlobalString(FakeNetFlag.Name))
		if err != nil {
			log.Crit("Invalid flag", "flag", FakeNetFlag.Name, "err", err)
		}
		cfg = lachesis.FakeNetConfig(accs)
	case ctx.GlobalBool(utils.TestnetFlag.Name):
		cfg = lachesis.TestNetConfig()
	default:
		cfg = lachesis.MainNetConfig()
	}

	return cfg
}

func setDataDir(ctx *cli.Context, cfg *node.Config) {
	defaultDataDir := DefaultDataDir()

	switch {
	case ctx.GlobalIsSet(utils.DataDirFlag.Name):
		cfg.DataDir = ctx.GlobalString(utils.DataDirFlag.Name)
	case ctx.GlobalIsSet(FakeNetFlag.Name):
		_, accs, err := parseFakeGen(ctx.GlobalString(FakeNetFlag.Name))
		if err != nil {
			log.Crit("Invalid flag", "flag", FakeNetFlag.Name, "err", err)
		}
		cfg.DataDir = filepath.Join(defaultDataDir, fmt.Sprintf("fakenet-%d", len(accs.Accounts)))
	case ctx.GlobalBool(utils.TestnetFlag.Name):
		cfg.DataDir = filepath.Join(defaultDataDir, "testnet")
	default:
		cfg.DataDir = defaultDataDir
	}
}

func setGPO(ctx *cli.Context, cfg *gasprice.Config) {
	if ctx.GlobalIsSet(utils.GpoBlocksFlag.Name) {
		cfg.Blocks = ctx.GlobalInt(utils.GpoBlocksFlag.Name)
	}
	if ctx.GlobalIsSet(utils.GpoPercentileFlag.Name) {
		cfg.Percentile = ctx.GlobalInt(utils.GpoPercentileFlag.Name)
	}
	if ctx.GlobalIsSet(GpoDefaultFlag.Name) {
		cfg.Default = utils.GlobalBig(ctx, GpoDefaultFlag.Name)
	}
}

func setTxPool(ctx *cli.Context, cfg *evmcore.TxPoolConfig) {
	if ctx.GlobalIsSet(utils.TxPoolLocalsFlag.Name) {
		locals := strings.Split(ctx.GlobalString(utils.TxPoolLocalsFlag.Name), ",")
		for _, account := range locals {
			if trimmed := strings.TrimSpace(account); !common.IsHexAddress(trimmed) {
				utils.Fatalf("Invalid account in --txpool.locals: %s", trimmed)
			} else {
				cfg.Locals = append(cfg.Locals, common.HexToAddress(account))
			}
		}
	}
	if ctx.GlobalIsSet(utils.TxPoolNoLocalsFlag.Name) {
		cfg.NoLocals = ctx.GlobalBool(utils.TxPoolNoLocalsFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolJournalFlag.Name) {
		cfg.Journal = ctx.GlobalString(utils.TxPoolJournalFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolRejournalFlag.Name) {
		cfg.Rejournal = ctx.GlobalDuration(utils.TxPoolRejournalFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolPriceLimitFlag.Name) {
		cfg.PriceLimit = ctx.GlobalUint64(utils.TxPoolPriceLimitFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolPriceBumpFlag.Name) {
		cfg.PriceBump = ctx.GlobalUint64(utils.TxPoolPriceBumpFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolAccountSlotsFlag.Name) {
		cfg.AccountSlots = ctx.GlobalUint64(utils.TxPoolAccountSlotsFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolGlobalSlotsFlag.Name) {
		cfg.GlobalSlots = ctx.GlobalUint64(utils.TxPoolGlobalSlotsFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolAccountQueueFlag.Name) {
		cfg.AccountQueue = ctx.GlobalUint64(utils.TxPoolAccountQueueFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolGlobalQueueFlag.Name) {
		cfg.GlobalQueue = ctx.GlobalUint64(utils.TxPoolGlobalQueueFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolLifetimeFlag.Name) {
		cfg.Lifetime = ctx.GlobalDuration(utils.TxPoolLifetimeFlag.Name)
	}
}

func gossipConfigWithFlags(ctx *cli.Context, src gossip.Config) gossip.Config {
	cfg := src

	// Avoid conflicting network flags
	utils.CheckExclusive(ctx, FakeNetFlag, utils.DeveloperFlag, utils.TestnetFlag)
	utils.CheckExclusive(ctx, FakeNetFlag, utils.DeveloperFlag, utils.ExternalSignerFlag) // Can't use both ephemeral unlocked and external signer

	setGPO(ctx, &cfg.GPO)
	setTxPool(ctx, &cfg.TxPool)

	if ctx.GlobalIsSet(utils.NetworkIdFlag.Name) {
		cfg.Net.NetworkID = ctx.GlobalUint64(utils.NetworkIdFlag.Name)
	}
	// TODO cache config
	//if ctx.GlobalIsSet(utils.CacheFlag.Name) || ctx.GlobalIsSet(utils.CacheDatabaseFlag.Name) {
	//	cfg.DatabaseCache = ctx.GlobalInt(utils.CacheFlag.Name) * ctx.GlobalInt(utils.CacheDatabaseFlag.Name) / 100
	//}
	//if ctx.GlobalIsSet(utils.CacheFlag.Name) || ctx.GlobalIsSet(CacheTrieFlag.Name) {
	//	cfg.TrieCleanCache = ctx.GlobalInt(utils.CacheFlag.Name) * ctx.GlobalInt(CacheTrieFlag.Name) / 100
	//}
	//if ctx.GlobalIsSet(utils.CacheFlag.Name) || ctx.GlobalIsSet(CacheGCFlag.Name) {
	//	cfg.TrieDirtyCache = ctx.GlobalInt(utils.CacheFlag.Name) * ctx.GlobalInt(CacheGCFlag.Name) / 100
	//}

	if ctx.GlobalIsSet(utils.VMEnableDebugFlag.Name) {
		cfg.EnablePreimageRecording = ctx.GlobalBool(utils.VMEnableDebugFlag.Name)
	}

	if ctx.GlobalIsSet(utils.EWASMInterpreterFlag.Name) {
		cfg.EWASMInterpreter = ctx.GlobalString(utils.EWASMInterpreterFlag.Name)
	}

	if ctx.GlobalIsSet(utils.EVMInterpreterFlag.Name) {
		cfg.EVMInterpreter = ctx.GlobalString(utils.EVMInterpreterFlag.Name)
	}
	if ctx.GlobalIsSet(utils.RPCGlobalGasCap.Name) {
		cfg.RPCGasCap = new(big.Int).SetUint64(ctx.GlobalUint64(utils.RPCGlobalGasCap.Name))
	}

	return cfg
}

func nodeConfigWithFlags(ctx *cli.Context, cfg node.Config) node.Config {
	utils.SetNodeConfig(ctx, &cfg)
	setDataDir(ctx, &cfg)
	return cfg
}

func makeAllConfigs(ctx *cli.Context) config {
	// Defaults (low priority)
	net := defaultLachesisConfig(ctx)
	cfg := config{
		Lachesis: gossip.DefaultConfig(net),
		Node:     defaultNodeConfig(),
		App:      app.DefaultConfig(net),
	}

	// zero SelfForkProtection interval
	cfg.Lachesis.Emitter.EmitIntervals.SelfForkProtection = 0

	// Load config file (medium priority)
	if file := ctx.GlobalString(configFileFlag.Name); file != "" {
		if err := cfg.Load(file); err != nil {
			utils.Fatalf("%v", err)
		}
	}

	// Apply flags (high priority)
	cfg.Lachesis = gossipConfigWithFlags(ctx, cfg.Lachesis)
	cfg.Node = nodeConfigWithFlags(ctx, cfg.Node)

	return cfg
}

func defaultNodeConfig() node.Config {
	cfg := NodeDefaultConfig
	cfg.Name = clientIdentifier
	cfg.Version = params.VersionWithCommit(gitCommit, gitDate)
	cfg.HTTPModules = append(cfg.HTTPModules, "eth", "ftm", "sfc", "web3")
	cfg.WSModules = append(cfg.WSModules, "eth", "ftm", "sfc", "web3")
	cfg.IPCPath = "lachesis.ipc"
	cfg.DataDir = DefaultDataDir()
	return cfg
}

// dumpConfig is the dumpconfig command.
func dumpConfig(ctx *cli.Context) error {
	cfg := makeAllConfigs(ctx)
	comment := ""

	out, err := toml.Settings.Marshal(&cfg)
	if err != nil {
		return err
	}

	dump := os.Stdout
	if ctx.NArg() > 0 {
		dump, err = os.OpenFile(ctx.Args().Get(0), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer dump.Close()
	}
	dump.WriteString(comment)
	dump.Write(out)

	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
