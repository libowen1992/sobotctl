package kafka

import (
	"sobotctl/pkg/convert"
	"sobotctl/pkg/tableRender"
	"strings"
)

var TopicListHeader = []string{"主题", "分区数", "副本数"}

type TopicListFilter struct {
	Filter string
}

func (t *topicOps) List(filter TopicListFilter) error {
	var data = make([][]string, 0, 0)
	kfkops, err := NewKfkOpsAdmin()
	if err != nil {
		return err
	}
	defer kfkops.CloseAdmin()

	// List topics
	topics, err := kfkops.Admin.ListTopics()
	if err != nil {
		return err
	}

	data = make([][]string, 0, 0)
	for topic, detail := range topics {
		if strings.Contains(topic, filter.Filter) {
			item := make([]string, 0, 0)
			item = append(item, topic)
			item = append(item, convert.Int32ToStr(detail.NumPartitions))
			item = append(item, convert.Int16ToStr(detail.ReplicationFactor))
			data = append(data, item)
		}
	}

	tableRender.Render(TopicListHeader, data)

	return nil
}
