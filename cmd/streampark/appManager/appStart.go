package appManager

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"sobotctl/global"
	"sobotctl/internal/streampark"
)

func NewAppStartCMD() *cobra.Command {
	var AppStop = &cobra.Command{
		Use:   "start",
		Short: "启动app",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New(" sobotctl streampark app start  appid")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := streampark.NewAppOps().Start(args[0]); err != nil {
				global.Logger.Error(err)
			}
			global.Logger.Info("app start ok!")
		},
	}

	return AppStop
}
