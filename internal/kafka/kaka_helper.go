package kafka

import (
	"errors"
	"github.com/Shopify/sarama"
	"sobotctl/global"
)

type KfkOps struct {
	Admin  sarama.ClusterAdmin
	Client sarama.Client
}

func NewKfkOpsAdmin() (*KfkOps, error) {
	if len(global.KafkaSetting.Adders) == 0 {
		return nil, errors.New("创建kafka实例失败, 请在配置文件里面配置kafka实例")
	}
	// Create Kafka admin client
	config := sarama.NewConfig()
	admin, err := sarama.NewClusterAdmin(global.KafkaSetting.Adders, config)
	if err != nil {
		return nil, err
	}
	return &KfkOps{Admin: admin}, nil
}

func NewKfkClient() (*KfkOps, error) {
	if len(global.KafkaSetting.Adders) == 0 {
		return nil, errors.New("创建kafka实例失败, 请在配置文件里面配置kafka实例")
	}
	config := sarama.NewConfig()
	client, err := sarama.NewClient(global.KafkaSetting.Adders, config)
	if err != nil {
		return nil, err
	}
	return &KfkOps{Client: client}, nil
}

func (k *KfkOps) CloseAdmin() {
	k.Admin.Close()
}

func (k *KfkOps) CloseClient() {
	k.Client.Closed()
}
