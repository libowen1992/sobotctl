package groupManager

import (
	"github.com/spf13/cobra"
)

func NewGroupCmd() *cobra.Command {
	var GroupCmd = &cobra.Command{
		Use:   "group",
		Short: "消费组管理",
		Long:  `消费组管理`,
	}
	GroupCmd.AddCommand(NewGroupList())
	return GroupCmd
}
