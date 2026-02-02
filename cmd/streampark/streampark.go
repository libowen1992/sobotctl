package streampark

import (
	"github.com/spf13/cobra"
	"sobotctl/cmd/streampark/appManager"
	"sobotctl/cmd/streampark/variableManager"
)

func NewStreamParkCmd() *cobra.Command {
	var streamParkCmd = &cobra.Command{
		Use:   "streampark",
		Short: "streampark管理",
	}
	streamParkCmd.AddCommand(variableManager.NewVariableCmd())
	streamParkCmd.AddCommand(appManager.NewAppM())
	return streamParkCmd
}
