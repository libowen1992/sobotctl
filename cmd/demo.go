package cmd

import "github.com/spf13/cobra"

func demo() *cobra.Command {
	var filter string
	action := "demo"
	desc := "demo配置列表"
	var Cmd = &cobra.Command{
		Use:   action,
		Short: desc,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	Cmd.Flags().StringVarP(&filter, "filter", "f", "", "模糊查询")
	return Cmd
}
