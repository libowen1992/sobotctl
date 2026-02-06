package cmd

import (
	"github.com/spf13/cobra"
	"sobotctl/internal/elasticOps"
)


//定义子命令，无操作命令
func NewElastic() *cobra.Command {
	action := "elastic"
	desc := "elastic工具"
	var Cmd = &cobra.Command{
		Use:   action,
		Short: desc,
	}
	Cmd.AddCommand(NewElasticCheckCmd())
	return Cmd
}

func NewElasticCheckCmd() *cobra.Command {
	action := "check"
	desc := "elastic检查状态"
	var Cmd = &cobra.Command{
		Use:   action,
		Short: desc,
		Run: func(cmd *cobra.Command, args []string) {
			elasticOps.NewElasticOps().Check()
		},
	}
	return Cmd
}
