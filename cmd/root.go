package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
    Use:   "aqua-sec-cloud-inventory",
    Short: "A CLI for the Aqua Security Cloud Resource Inventory Management System",
    Long:  `A CLI tool to manage migrations, seeding, and running the server for the Aqua Security Cloud Resource Inventory Management System.`,
}

// Execute executes the root command.
func Execute() {
    if err := RootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}