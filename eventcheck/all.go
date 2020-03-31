package eventcheck

import (
	"github.com/Fantom-foundation/lachesis-ex/eventcheck/basiccheck"
	"github.com/Fantom-foundation/lachesis-ex/eventcheck/epochcheck"
	"github.com/Fantom-foundation/lachesis-ex/eventcheck/gaspowercheck"
	"github.com/Fantom-foundation/lachesis-ex/eventcheck/heavycheck"
	"github.com/Fantom-foundation/lachesis-ex/eventcheck/parentscheck"
	"github.com/Fantom-foundation/lachesis-ex/inter"
)

// Checkers is collection of all the checkers
type Checkers struct {
	Basiccheck    *basiccheck.Checker
	Epochcheck    *epochcheck.Checker
	Parentscheck  *parentscheck.Checker
	Gaspowercheck *gaspowercheck.Checker
	Heavycheck    *heavycheck.Checker
}

// Validate runs all the checks except Poset-related. intended only for tests
func (v *Checkers) Validate(e *inter.Event, parents []*inter.EventHeaderData) error {
	if err := v.Basiccheck.Validate(e); err != nil {
		return err
	}
	if err := v.Epochcheck.Validate(e); err != nil {
		return err
	}
	if err := v.Parentscheck.Validate(e, parents); err != nil {
		return err
	}
	var selfParent *inter.EventHeaderData
	if e.SelfParent() != nil {
		selfParent = parents[0]
	}
	if err := v.Gaspowercheck.Validate(e, selfParent); err != nil {
		return err
	}
	if err := v.Heavycheck.Validate(e); err != nil {
		return err
	}
	return nil
}
