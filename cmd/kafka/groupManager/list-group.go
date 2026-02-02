package groupManager

import (
	"github.com/spf13/cobra"
	"sobotctl/internal/kafka"
)

func NewGroupList() *cobra.Command {
	var filterFlags kafka.GroupListFilter
	var GroupListCmd = &cobra.Command{
		Use:   "list",
		Short: "Group列表",
		RunE: func(cmd *cobra.Command, args []string) error {
			return kafka.NewGroupOps().List(filterFlags)
		},
	}

	GroupListCmd.Flags().StringVarP(&filterFlags.Filter, "filter", "f", "", "过滤器，模糊匹配")
	return GroupListCmd
}
