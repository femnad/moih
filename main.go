package main

import (
	"github.com/alexflint/go-arg"
	"github.com/femnad/moih/cmd"
	"log"
	"os"
)

const (
	version = "0.5.0"
)

type Base struct{}

func (Base) Version() string {
	return version
}

var args struct {
	Base
	Get    *cmd.KeyCfg    `arg:"subcommand:get" help:"get a key from GCP Cloud Storage"`
	Put    *cmd.KeyCfg    `arg:"subcommand:put" help:"put a key into GCP Cloud Storage"`
	Update *cmd.UpdateCfg `arg:"subcommand:update" help:"update a key in GitHub"`
}

func mustSucceed(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	p := arg.MustParse(&args)

	switch {
	case args.Put != nil:
		err := cmd.Put(*args.Put)
		mustSucceed(err)

	case args.Get != nil:
		err := cmd.Get(*args.Get)
		mustSucceed(err)

	case args.Update != nil:
		err := cmd.Update(*args.Update)
		mustSucceed(err)

	case true:
		p.WriteHelp(os.Stdout)
	}
}
