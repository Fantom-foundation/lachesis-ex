package main

import (
	"github.com/naoina/toml/ast"

	"github.com/Fantom-foundation/lachesis-ex/utils/migration"
	"github.com/Fantom-foundation/lachesis-ex/utils/toml"
)

func (c *config) migrate(t *ast.Table) (changed bool, err error) {
	data := toml.NewTomlHelper(t)
	migrations := c.migrations(data)
	versions := toml.NewIDStore(data, migrations.IdChain())

	before := versions.GetID()

	err = migrations.Exec(versions)
	if err != nil && err != toml.ErrorParamNotExists {
		return
	}

	after := versions.GetID()
	changed = before != after

	return
}

func (c *config) migrations(data *toml.Helper) *migration.Migration {
	return migration.
		Begin("lachesis-config")
}
