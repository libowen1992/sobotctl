package topicManager

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"path"
	"sobotctl/global"
	"sobotctl/internal/kafka"
)

func NewTopicCmd() *cobra.Command {
	var TopicCmd = &cobra.Command{
		Use:   "topic",
		Short: "topic管理",
		Long:  `topic管理`,
	}
	TopicCmd.AddCommand(NewTopicCreate())
	TopicCmd.AddCommand(NewTopicList())
	TopicCmd.AddCommand(NewTopicDescribeCmd())
	return TopicCmd
}

func NewTopicCreate() *cobra.Command {
	var tcf kafka.TopicCreateFlags

	var createCmd = &cobra.Command{
		Use:   "create",
		Short: "创建topic",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(tcf.TopicFile) != 0 && !path.IsAbs(tcf.TopicFile) {
				return errors.New("主题文件必须给绝对路径")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := kafka.NewTopicOps().CreateTopic(tcf); err != nil {
				global.Logger.Error(err)
				os.Exit(1)
			}
		},
	}

	createCmd.Flags().StringVarP(&tcf.TopicName, "topic", "t", "", "主题名称")
	createCmd.Flags().StringVarP(&tcf.TopicFile, "file", "f", "", "主题文件,一行一个主题")
	createCmd.MarkFlagsMutuallyExclusive("topic", "file")
	createCmd.Flags().Int32VarP(&tcf.TopicPartitions, "partitions", "p", 3, "主题分区")
	createCmd.Flags().Int16VarP(&tcf.TopicReplicationFactor, "replication-factor", "r", 1, "主题副本")

	return createCmd
}

func NewTopicList() *cobra.Command {
	var filterFlags kafka.TopicListFilter
	var TopicListCmd = &cobra.Command{
		Use:   "list",
		Short: "topic列表",
		RunE: func(cmd *cobra.Command, args []string) error {
			return kafka.NewTopicOps().List(filterFlags)
		},
	}

	TopicListCmd.Flags().StringVarP(&filterFlags.Filter, "filter", "f", "", "过滤器，模糊匹配")
	return TopicListCmd
}

func NewTopicDescribeCmd() *cobra.Command {
	var filterFlags kafka.TopicDescribeFilter

	var topicDescribe = &cobra.Command{
		Use:   "desc",
		Short: "主题详情",
		RunE: func(cmd *cobra.Command, args []string) error {
			return kafka.NewTopicOps().Describe(filterFlags)
		},
	}
	topicDescribe.Flags().StringVarP(&filterFlags.Topic, "topic", "t", "", "topic名称")

	return topicDescribe
}
