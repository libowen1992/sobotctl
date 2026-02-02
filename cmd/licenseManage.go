package cmd

import (
	"github.com/spf13/cobra"
)

func LicenseManage() *cobra.Command {
	action := "license"
	desc := "license管理"
	var Cmd = &cobra.Command{
		Use:   action,
		Short: desc,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	return Cmd
}
