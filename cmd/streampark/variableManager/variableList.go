package variableManager

import (
	"github.com/spf13/cobra"
	"sobotctl/internal/streampark"
)

func NewVariableList() *cobra.Command {
	var filter string
	var VariableListCmd = &cobra.Command{
		Use:   "list",
		Short: "变量列表",
		Run: func(cmd *cobra.Command, args []string) {
			if err := streampark.NewVariableOps().List(filter); err != nil {
				panic(err)
			}
		},
	}

	VariableListCmd.Flags().StringVarP(&filter, "filter", "f", "", "过滤器，模糊匹配")

	return VariableListCmd
}
