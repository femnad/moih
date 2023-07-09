package main

import (
	"github.com/alexflint/go-arg"
	"github.com/femnad/moih/cmd"
	"log"
	"os"
)

const (
	version = "0.4.2"
)

type GetCmd struct {
	cmd.PrivateKeyInfo
}

type PutCmd struct {
	cmd.PrivateKeyInfo
}

type UpdateCmd struct {
	cmd.UpdateCfg
}

type Base struct{}

func (Base) Version() string {
	return version
}

var args struct {
	Base
	Get    *GetCmd    `arg:"subcommand:get" help:"get a key from GCP Cloud Storage"`
	Put    *PutCmd    `arg:"subcommand:put" help:"put a key into GCP Cloud Storage"`
	Update *UpdateCmd `arg:"subcommand:update" help:"update a key in GitHub"`
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
		err := cmd.Put(args.Put.PrivateKeyInfo)
		mustSucceed(err)

	case args.Get != nil:
		err := cmd.Get(args.Get.PrivateKeyInfo)
		mustSucceed(err)

	case args.Update != nil:
		err := cmd.Update(args.Update.UpdateCfg)
		mustSucceed(err)

	case true:
		p.WriteHelp(os.Stdout)
	}
}
