package appManager

import "github.com/spf13/cobra"

func NewAppM() *cobra.Command {
	var AppMCMD = &cobra.Command{
		Use:   "app",
		Short: "app管理",
	}
	AppMCMD.AddCommand(NewAppList())
	AppMCMD.AddCommand(NewAppStopCMD())
	AppMCMD.AddCommand(NewAppStartCMD())
	AppMCMD.AddCommand(NewAppUpdateCMD())
	AppMCMD.AddCommand(NewAppInitCMD())
	AppMCMD.AddCommand(NewAppLog())
	return AppMCMD
}
