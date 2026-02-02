package cmd

import (
	"os"
	"sobotctl/cmd/minioOps"
	"sobotctl/cmd/streampark"

	"github.com/spf13/cobra"
	"sobotctl/cmd/kafka"
)

var rootCmd = &cobra.Command{
	Use:   "sobotctl",
	Short: "运维管理工具",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

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
