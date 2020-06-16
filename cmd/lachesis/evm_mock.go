package main

import (
	"gopkg.in/urfave/cli.v1"
)

var (
	// VmFlag kind of virtual machine.
	VmFlag = cli.StringFlag{
		Name:  "vm",
		Usage: "Kind of virtual machine to use [evm, mock, fvm]",
		Value: "evm",
	}
)
