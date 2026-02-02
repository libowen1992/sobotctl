package variableManager

import (
	"github.com/spf13/cobra"
	"sobotctl/global"
	"sobotctl/internal/streampark"
)

func NewVariableUpdate() *cobra.Command {
	var name string
	var value string
	var VariableUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "修改变量值",
		Run: func(cmd *cobra.Command, args []string) {
			if err := streampark.NewVariableOps().Update(name, value); err != nil {
				global.Logger.Error("更新失败", err)
				return
			}
			global.Logger.Info("更新成功")
		},
	}

	VariableUpdateCmd.Flags().StringVarP(&name, "name", "n", "", "变量名")
	VariableUpdateCmd.MarkFlagRequired("name")
	VariableUpdateCmd.Flags().StringVarP(&value, "value", "v", "", "变量值")
	VariableUpdateCmd.MarkFlagRequired("value")

	return VariableUpdateCmd
}
