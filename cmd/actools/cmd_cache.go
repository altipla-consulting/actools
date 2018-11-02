package main

import (
	"github.com/spf13/cobra"
)

func init() {
	CmdRoot.AddCommand(CmdCache)
}

var CmdCache = &cobra.Command{
	Use:   "cache",
	Short: "Manage the tools local cache.",
}
