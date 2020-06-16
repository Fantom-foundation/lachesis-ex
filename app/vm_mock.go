package app

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/Fantom-foundation/lachesis-ex/evmcore"
)

type VmMock struct {
}

func (vm *VmMock) Process(
	block *evmcore.EvmBlock, statedb *state.StateDB, cfg vm.Config, strict bool,
) (
	receipts types.Receipts, logs []*types.Log, gasUsed uint64, totalFee *big.Int, skipped []uint, err error,
) {
	totalFee = big.NewInt(0)
	return
}
