// xvoke is the entry point of the xvoke command. It hands the arguments to internal/cli.
package main

import (
	"os"

	"github.com/progrhyme/xvoke/internal/cli"
	"github.com/progrhyme/xvoke/internal/version"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, cli.Deps{Version: version.Get()}))
}
