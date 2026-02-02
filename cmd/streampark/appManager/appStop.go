package appManager

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"sobotctl/global"
	"sobotctl/internal/streampark"
)

func NewAppStopCMD() *cobra.Command {
	var AppStop = &cobra.Command{
		Use:   "stop",
		Short: "停止app",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New(" sobotctl streampark app stop  appid")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := streampark.NewAppOps().Stop(args[0]); err != nil {
				global.Logger.Error(err)
			}
			global.Logger.Info("app stop ok!")
		},
	}

	return AppStop
}
