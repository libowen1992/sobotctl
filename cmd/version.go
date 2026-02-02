package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var ctlVersion = "1.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(ctlVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
