package appManager

import (
	"fmt"
	"github.com/spf13/cobra"
)

func NewAppLog() *cobra.Command {
	var APPLogCmd = &cobra.Command{
		Use:   "log",
		Short: "log查看日志",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("功能开发中，登陆hadoop01服务器，yarn logs -applicationId xxx， appid 通过app list 查看")
		},
	}

	return APPLogCmd
}
