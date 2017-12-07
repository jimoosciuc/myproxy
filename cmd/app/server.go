package app

import (
	"github.com/spf13/cobra"
	"fmt"
)

var RootCmd = &cobra.Command{
	Use: "open-falcon",
	RunE: func(c *cobra.Command, args []string) error {
		if versionFlag {
			fmt.Printf("Open-Falcon version %s, build %s\n", Version, GitCommit)
			return nil
		}
		return c.Usage()
	},
}