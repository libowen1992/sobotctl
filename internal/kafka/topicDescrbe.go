package kafka

var TopicDescribeHeader = []string{"主题", "分区数", "副本数"}

type TopicDescribeFilter struct {
	Topic string
}

func (t *topicOps) Describe(filter TopicDescribeFilter) error {
	//kfkops, err := NewKfkOpsAdmin()
	//if err != nil {
	//	return err
	//}
	//defer kfkops.CloseAdmin()
	//
	//topicDetails, err := kfkops.Admin.DescribeTopics([]string{filter.Topic})
	//if err != nil {
	//	return err
	//}
	//topicDetail := topicDetails[0]
	//if len(topicDetail.Partitions) == 0 {
	//	return errors.New("topic不存在")
	//}
	//var partionsNumbers []int32
	//for _, p := range topicDetail.Partitions {
	//	partionsNumbers = append(partionsNumbers, p.ID)
	//}
	//
	//client, err := NewKfkClient()
	//if err != nil {
	//	return err
	//}
	//defer client.CloseClient()
	//client.Client.GetOffset()

	return nil
}
