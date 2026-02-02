package appManager

import (
	"github.com/spf13/cobra"
	"sobotctl/global"
	"sobotctl/internal/streampark"
)

func NewAppUpdateCMD() *cobra.Command {
	var appid string
	var jar string
	var AppUpdate = &cobra.Command{
		Use:   "update",
		Short: "更新app",
		Run: func(cmd *cobra.Command, args []string) {
			if err := streampark.NewAppOps().Update(appid, jar); err != nil {
				global.Logger.Error(err)
				return
			}
			global.Logger.Info("update stop ok!")
		},
	}

	AppUpdate.Flags().StringVarP(&appid, "appid", "", "", "appid")
	AppUpdate.MarkFlagRequired("appid")
	AppUpdate.Flags().StringVarP(&jar, "jar", "", "", "jar包")
	AppUpdate.MarkFlagRequired("jar")

	return AppUpdate
}
