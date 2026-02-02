package kafka

import (
	"github.com/spf13/cobra"
	"sobotctl/cmd/kafka/groupManager"
	"sobotctl/cmd/kafka/topicManager"
)

func NewKafkaCmd() *cobra.Command {
	var kafkaCmd = &cobra.Command{
		Use:   "kafka",
		Short: "kafka管理工具",
	}
	kafkaCmd.AddCommand(topicManager.NewTopicCmd())
	kafkaCmd.AddCommand(groupManager.NewGroupCmd())
	return kafkaCmd
}
