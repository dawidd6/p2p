package main

import "github.com/spf13/cobra"

var cmdRoot = &cobra.Command{
	Use:                        "p2p",
	Example:                    "",
	Version:                    "0.1.0",
	PersistentPreRun:           func(*cobra.Command, []string) { panic("not implemented") },
	PersistentPreRunE:          func(*cobra.Command, []string) error { panic("not implemented") },
	PreRun:                     func(*cobra.Command, []string) { panic("not implemented") },
	PreRunE:                    func(*cobra.Command, []string) error { panic("not implemented") },
	Run:                        func(*cobra.Command, []string) { panic("not implemented") },
	RunE:                       func(*cobra.Command, []string) error { panic("not implemented") },
	PostRun:                    func(*cobra.Command, []string) { panic("not implemented") },
	PostRunE:                   func(*cobra.Command, []string) error { panic("not implemented") },
	PersistentPostRun:          func(*cobra.Command, []string) { panic("not implemented") },
	PersistentPostRunE:         func(*cobra.Command, []string) error { panic("not implemented") },
	SilenceErrors:              false,
	SilenceUsage:               false,
	DisableFlagParsing:         false,
	DisableAutoGenTag:          false,
	DisableFlagsInUseLine:      false,
	DisableSuggestions:         false,
	SuggestionsMinimumDistance: 0,
	TraverseChildren:           false,
}

var cmdDaemon = &cobra.Command{
	Use:                        "daemon",
	Example:                    "",
	PersistentPreRun:           func(*cobra.Command, []string) { panic("not implemented") },
	PersistentPreRunE:          func(*cobra.Command, []string) error { panic("not implemented") },
	PreRun:                     func(*cobra.Command, []string) { panic("not implemented") },
	PreRunE:                    func(*cobra.Command, []string) error { panic("not implemented") },
	Run:                        func(*cobra.Command, []string) { panic("not implemented") },
	RunE:                       func(*cobra.Command, []string) error { panic("not implemented") },
	PostRun:                    func(*cobra.Command, []string) { panic("not implemented") },
	PostRunE:                   func(*cobra.Command, []string) error { panic("not implemented") },
	PersistentPostRun:          func(*cobra.Command, []string) { panic("not implemented") },
	PersistentPostRunE:         func(*cobra.Command, []string) error { panic("not implemented") },
	SilenceErrors:              false,
	SilenceUsage:               false,
	DisableFlagParsing:         false,
	DisableAutoGenTag:          false,
	DisableFlagsInUseLine:      false,
	DisableSuggestions:         false,
	SuggestionsMinimumDistance: 0,
	TraverseChildren:           false,
}

func init() {
	cmdRoot.AddCommand(
		cmdDaemon,
	)
}

func main() {
	err := cmdRoot.Execute()
	if err != nil {
		panic(err)
	}
}
