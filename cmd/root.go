package cmd

import (
	"os"
	"sobotctl/cmd/minioOps"
	"sobotctl/cmd/streampark"

	"github.com/spf13/cobra"
	"sobotctl/cmd/kafka"
)
//轻松创建基于子命令的 CLI：如app server、app fetch等、自动添加-h,–help等帮助性Flag、自动生成命令和Flag的帮助信息
//定义结构体指针--根命令
var rootCmd = &cobra.Command{
	Use:   "sobotctl",  //命令名称，用户输入的命令名，运行方式：sobotctl [子命令] [参数]
	Short: "运维管理工具", //简单描述
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

//通过Execute调用的函数，首先执行init初始化，且只执行一次，再次调用的话，会直接执行Execute
//go中包的init()函数会在import时执行，通过AddCommand(NewHostManage)将NewHostManage添加到结构体Command 成员变量commands中
//初始化命令数，添加子命令
func init() {
	rootCmd.AddCommand(NewHostManage())
	rootCmd.AddCommand(kafka.NewKafkaCmd())
	rootCmd.AddCommand(streampark.NewStreamParkCmd())
	rootCmd.AddCommand(minioOps.NewMinioOpsCmd())
	rootCmd.AddCommand(NewRedisManager())
	rootCmd.AddCommand(NewMysqlManager())
	rootCmd.AddCommand(NewNacosCmd())
	rootCmd.AddCommand(NewElastic())
	rootCmd.AddCommand(NewCheckCmd())
	rootCmd.AddCommand(NewParkCmd())
	rootCmd.AddCommand(NewHarborManager())
	rootCmd.AddCommand(K8sManager())
}
