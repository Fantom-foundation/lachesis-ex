package app

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"

	"github.com/Fantom-foundation/lachesis-ex/evmcore"
	"github.com/Fantom-foundation/lachesis-ex/inter/sfctype"
	"github.com/Fantom-foundation/lachesis-ex/lachesis"
)

// ApplyGenesis writes initial state.
func (s *Store) ApplyGenesis(net *lachesis.Config) (block *evmcore.EvmBlock, isNew bool, err error) {
	s.migrate()

	stored := s.getGenesisState()

	if stored != nil {
		block, err = calcGenesisBlock(net)
		if err != nil {
			return
		}

		if block.Root != *stored {
			err = fmt.Errorf("database contains incompatible state hash (have %s, new %s)", *stored, block.Root)
		}

		return
	}

	// if we'here, then it's first time genesis is applied
	isNew = true
	s.SetBlock(&BlockInfo{
		Index: 0,
		Time:  net.Genesis.Time,
	})
	block, err = s.applyGenesis(net)
	if err != nil {
		return
	}
	s.setGenesisState(block.Root)
	s.FlushState()
	return
}

// calcGenesisBlock calcs hash of genesis state.
func calcGenesisBlock(net *lachesis.Config) (*evmcore.EvmBlock, error) {
	s := NewMemStore()
	defer s.Close()

	s.Log.SetHandler(log.DiscardHandler())

	return s.applyGenesis(net)
}

func (s *Store) applyGenesis(net *lachesis.Config) (evmBlock *evmcore.EvmBlock, err error) {
	evmBlock, err = evmcore.ApplyGenesis(s.table.Evm, net)
	if err != nil {
		return
	}

	// calc total pre-minted supply
	totalSupply := big.NewInt(0)
	for _, account := range net.Genesis.Alloc.Accounts {
		totalSupply.Add(totalSupply, account.Balance)
	}
	s.SetTotalSupply(totalSupply)

	validatorsArr := []sfctype.SfcStakerAndID{}
	for _, validator := range net.Genesis.Alloc.Validators {
		staker := &sfctype.SfcStaker{
			Address:      validator.Address,
			CreatedEpoch: 0,
			CreatedTime:  net.Genesis.Time,
			StakeAmount:  validator.Stake,
			DelegatedMe:  big.NewInt(0),
		}
		s.SetSfcStaker(validator.ID, staker)
		validatorsArr = append(validatorsArr, sfctype.SfcStakerAndID{
			StakerID: validator.ID,
			Staker:   staker,
		})
	}
	s.SetEpochValidators(1, validatorsArr)
	s.SetLastVoting(0, net.Genesis.Time)

	return
}

func (s *Store) setGenesisState(root common.Hash) {
	key := []byte("genesis")

	if err := s.table.Genesis.Put(key, root.Bytes()); err != nil {
		s.Log.Crit("Failed to put key-value", "err", err)
	}
}

func (s *Store) getGenesisState() *common.Hash {
	key := []byte("genesis")

	buf, err := s.table.Genesis.Get(key)
	if err != nil {
		s.Log.Crit("Failed to get key-value", "err", err)
	}
	if buf == nil {
		return nil
	}

	root := common.BytesToHash(buf)
	return &root
}
