package kafka

import (
	"bufio"
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"os"
	"sobotctl/global"
	"strings"
)

var (
	FlagInvalidMsg     = "至少指定一个主题名称或者一个主题文件"
	SingleTopicCreate  = 1
	FileTopicCreate    = 2
	UnknownTopicCreate = 3
)

type TopicCreateFlags struct {
	TopicName              string // 主题
	TopicFile              string // 主题文件
	TopicPartitions        int32  // 分区数
	TopicReplicationFactor int16  // 副本数
}

func (tcf TopicCreateFlags) Check() (ok bool, method int, errMsg string) {
	if len(tcf.TopicName) != 0 {
		return true, SingleTopicCreate, ""
	} else if len(tcf.TopicFile) != 0 {
		return true, FileTopicCreate, ""
	} else {
		return false, UnknownTopicCreate, FlagInvalidMsg
	}
}

func (t *topicOps) CreateTopic(flags TopicCreateFlags) error {
	ok, method, errMsg := flags.Check()
	if !ok {
		return errors.New(errMsg)
	}

	switch method {
	case SingleTopicCreate:
		if err := t.createSingleTopic(flags); err != nil {
			return err
		}
	case FileTopicCreate:
		global.Logger.Infof("开始创建topic,文件为%s", flags.TopicFile)
		if err := t.createTopicFromFile(flags); err != nil {
			return err
		}
	}
	return nil
}

func (t *topicOps) createSingleTopic(flags TopicCreateFlags) error {
	kfkops, err := NewKfkOpsAdmin()
	if err != nil {
		return errors.WithStack(err)
	}
	defer kfkops.CloseAdmin()

	topicDetail := &sarama.TopicDetail{
		NumPartitions:     flags.TopicPartitions,
		ReplicationFactor: flags.TopicReplicationFactor,
	}
	if err := kfkops.Admin.CreateTopic(flags.TopicName, topicDetail, false); err != nil {
		if errors.Is(err, sarama.ErrTopicAlreadyExists) {
			global.Logger.Warn("topic已经存在了: ", flags.TopicName)
			return nil
		}
		return errors.WithStack(err)
	}
	global.Logger.Infof("%s: success create", flags.TopicName)

	return nil
}

func (t *topicOps) createTopicFromFile(flags TopicCreateFlags) error {
	fh, err := os.Open(flags.TopicFile)
	if err != nil {
		return errors.WithStack(err)
	}
	defer fh.Close()

	kfkops, err := NewKfkOpsAdmin()
	if err != nil {
		return errors.WithStack(err)
	}
	defer kfkops.CloseAdmin()
	topicDetail := &sarama.TopicDetail{
		NumPartitions:     flags.TopicPartitions,
		ReplicationFactor: flags.TopicReplicationFactor,
	}

	fs := bufio.NewScanner(fh)
	fs.Split(bufio.ScanLines)
	for fs.Scan() {
		topic := strings.TrimSpace(fs.Text())
		global.Logger.Infof("%s: 开始创建topic", topic)
		if err := kfkops.Admin.CreateTopic(topic, topicDetail, false); err != nil {
			if errors.Is(err, sarama.ErrTopicAlreadyExists) {
				global.Logger.Warnf("%s: topic已经存在了", topic)
				continue
			}
			return errors.WithStack(err)
		}
		global.Logger.Infof("%s: 创建成功", topic)
	}

	return nil
}
