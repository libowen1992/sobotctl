package variableManager

import (
	"github.com/spf13/cobra"
)

func NewVariableCmd() *cobra.Command {
	var variableCmd = &cobra.Command{
		Use:   "variable",
		Short: "变量管理",
	}
	variableCmd.AddCommand(NewVariableList())
	variableCmd.AddCommand(NewVariableInit())
	variableCmd.AddCommand(NewVariableUpdate())
	return variableCmd
}
