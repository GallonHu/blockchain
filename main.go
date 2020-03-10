package main

import (
	"blockchain/cli"
	"blockchain/core"
)

func main() {
	bc := core.NewBlockchain()
	defer bc.Db.Close()

	cli := cli.CLI{bc}
	cli.Run()
}
