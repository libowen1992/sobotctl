package global

import (
	"go.uber.org/zap"
	"sobotctl/setting"
)

var (
	HostSetting   *setting.Hosts
	MySQLSetting  *setting.Mysql
	RedisSetting  *setting.Redis
	AppSetting    setting.Apps
	KafkaSetting  *setting.Kafka
	StreamparkS   *setting.StreamPark
	DorisS        *setting.Doris
	ElasticS      *setting.Elastic
	ZookeeperS    *setting.Zookeeper
	HadoopS       *setting.Hadoop
	NacosS        *setting.Nacos
	K8sS          *setting.K8s
	HarborS       *setting.Harbor
	TCPPortCheckS *setting.TCPPortCheck
	Logger        *zap.SugaredLogger
)
