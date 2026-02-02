package appManager

import (
	"github.com/spf13/cobra"
	"sobotctl/global"
	"sobotctl/internal/streampark"
)

func NewAppList() *cobra.Command {
	var filter string
	var APPListCmd = &cobra.Command{
		Use:   "list",
		Short: "app列表",
		Run: func(cmd *cobra.Command, args []string) {
			if err := streampark.NewAppOps().List(filter); err != nil {
				global.Logger.Error(err)
			}
		},
	}

	APPListCmd.Flags().StringVarP(&filter, "filter", "f", "", "过滤器，模糊匹配")

	return APPListCmd
}
