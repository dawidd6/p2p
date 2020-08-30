package main

import (
	"github.com/dawidd6/p2p/pkg/version"
	"github.com/spf13/cobra"
)

var cmdRoot = &cobra.Command{
	Use:                        "p2p-tracker",
	Short:                      "CLI interface to P2P tracker daemon.",
	Version:                    version.Version,
	RunE:                       run,
	SilenceErrors:              false,
	SilenceUsage:               false,
}

func run(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}

func main() {
	err := cmdRoot.Execute()
	if err != nil {
		panic(err)
	}
}
