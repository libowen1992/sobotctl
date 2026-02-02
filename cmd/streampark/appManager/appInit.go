package appManager

import (
	"github.com/spf13/cobra"
	"sobotctl/global"
	"sobotctl/internal/streampark"
	"sobotctl/pkg/userInput"
	"strings"
)

func NewAppInitCMD() *cobra.Command {
	guide := "此命令为初始化命令，会重新部署所有streampark服务，适合初次部署使用，确认输入y，推出输入n"
	var jarDir string
	var AppInit = &cobra.Command{
		Use:   "init",
		Short: "初始化，适合初次部署时使用，部署所有的服务",
		Run: func(cmd *cobra.Command, args []string) {
			choice, err := userInput.UserString(guide)
			if err != nil {
				global.Logger.Error("输入有误")
				return
			}
			if strings.ToLower(choice) == "y" {
				streampark.NewAppOps().Init(jarDir)
			}
		},
	}

	AppInit.Flags().StringVarP(&jarDir, "jarDir", "d", "", "jar包文件夹")
	AppInit.MarkFlagRequired("jarDir")

	return AppInit
}
