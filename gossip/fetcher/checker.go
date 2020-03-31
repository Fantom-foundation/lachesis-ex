package fetcher

import (
	"github.com/Fantom-foundation/lachesis-ex/eventcheck/heavycheck"
	"github.com/Fantom-foundation/lachesis-ex/inter"
)

// Checker is an interface that represents abstract logic for a checker object
type Checker interface {
	Start()
	Stop()
	Overloaded() bool
	Enqueue(events inter.Events, onValidated heavycheck.OnValidatedFn) error
}
