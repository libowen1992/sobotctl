package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Shopify/sarama"
	_ "github.com/go-sql-driver/mysql"
	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
	goredis "github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"net"
	"path/filepath"
	"sobotctl/global"
	"sobotctl/internal/hostManage"
	"sobotctl/internal/redisManage"
	"sobotctl/internal/streampark"
	"sobotctl/pkg/mysql"
	"sobotctl/pkg/tableRender"
	"sobotctl/setting"
	"strings"
	"time"
)

func NewCheckCmd() *cobra.Command {
	action := "check"
	desc := "检查服务"
	var Cmd = &cobra.Command{
		Use:   action,
		Short: desc,
		Run: func(cmd *cobra.Command, args []string) {
			taskStartT := "%s: 检测开始"
			taskEndT := "%s: 检测结束"
			headers := []string{"资源", "问题", "详细信息"}
			problems := make([][]string, 0, 0)

			// 主机检测
			global.Logger.Info("---------------------------------------")
			global.Logger.Infof(taskStartT, "主机")
			if ps, err := hostCheck(); err != nil {
				global.Logger.Error("主机检测失败，请检查主机配置是否正确")
			} else {
				problems = append(problems, ps...)
			}
			global.Logger.Infof(taskEndT, "主机")

			// redis检测
			global.Logger.Info("---------------------------------------")
			global.Logger.Infof(taskStartT, "Redis")
			ps := RedisCheck()
			if len(ps) != 0 {
				problems = append(problems, ps...)
			}
			global.Logger.Infof(taskEndT, "Redis")

			// mysql检测
			global.Logger.Info("---------------------------------------")
			global.Logger.Infof(taskStartT, "MySQL")
			ps = MysqlCheck()
			if len(ps) != 0 {
				problems = append(problems, ps...)
			}
			global.Logger.Infof(taskEndT, "MySQL")

			// elastic检测
			global.Logger.Info("---------------------------------------")
			global.Logger.Infof(taskStartT, "Elasticsearch")
			ps = ElasticCheck()
			if len(ps) != 0 {
				problems = append(problems, ps...)
			}
			global.Logger.Infof(taskEndT, "Elasticsearch")

			// zookeeper检测
			global.Logger.Info("---------------------------------------")
			global.Logger.Infof(taskStartT, "Zookeeper")
			ps = ZookeeperCheck()
			if len(ps) != 0 {
				problems = append(problems, ps...)
			}
			global.Logger.Infof(taskEndT, "Zookeeper")

			// kafka检测
			global.Logger.Info("---------------------------------------")
			global.Logger.Infof(taskStartT, "Kafka")
			ps = KafkaCheck()
			if len(ps) != 0 {
				problems = append(problems, ps...)
			}
			global.Logger.Infof(taskEndT, "Kafka")

			// hadoop检测
			global.Logger.Info("---------------------------------------")
			global.Logger.Infof(taskStartT, "Hadoop")
			ps = HadoopCheck()
			if len(ps) != 0 {
				problems = append(problems, ps...)
			}
			global.Logger.Infof(taskEndT, "Hadoop")

			// hbase检测
			global.Logger.Info("---------------------------------------")
			global.Logger.Infof(taskStartT, "Hbase")
			ps = HbaseCheck()
			if len(ps) != 0 {
				problems = append(problems, ps...)
			}
			global.Logger.Infof(taskEndT, "Hbase")

			// streampark
			global.Logger.Info("---------------------------------------")
			global.Logger.Infof(taskStartT, "Streampark")
			ps = StreamParkCheck()
			if len(ps) != 0 {
				problems = append(problems, ps...)
			}
			global.Logger.Infof(taskEndT, "Streampark")

			// Doris
			global.Logger.Info("---------------------------------------")
			global.Logger.Infof(taskStartT, "Doris")
			ps = DorisCheck()
			if len(ps) != 0 {
				problems = append(problems, ps...)
			}
			global.Logger.Infof(taskEndT, "Doris")

			// Doris
			global.Logger.Info("---------------------------------------")
			global.Logger.Infof(taskStartT, "端口检测")
			ps = PortCheck()
			if len(ps) != 0 {
				problems = append(problems, ps...)
			}
			global.Logger.Infof(taskEndT, "端口检测")

			// Doris
			global.Logger.Info("---------------------------------------")
			global.Logger.Infof(taskStartT, "K8s检测")
			ps = K8SCheck()
			if len(ps) != 0 {
				problems = append(problems, ps...)
			}
			global.Logger.Infof(taskEndT, "K8s检测")

			// 输出整体检查结果
			global.Logger.Info("---------------------------------------")
			if len(problems) == 0 {
				global.Logger.Info("系统运行正常")
				return
			}
			global.Logger.Error("系统有运行异常情况，如下表")
			tableRender.Render(headers, problems)
		},
	}
	return Cmd
}

func hostCheck() (ps [][]string, err error) {
	problems := make([][]string, 0, 0)  //二维字符串切片
	data, err := hostManage.NewHostOps().Check()
	if err != nil {
		global.Logger.Error(err)
		return nil, err
	}
	for _, item := range data {
		// 如果负载大于cpu核心
		global.Logger.Debug(fmt.Sprintf("检测负载(负载是否大于cpu cores): %s", item.IP))
		if item.Load > float64(item.CPUCores) {
			resource := "主机"
			problemTitle := "cpu负载异常"
			problemDesc := fmt.Sprintf("主机: %s, cpu核心: %d, 当前负载:%f", item.HostName, item.CPUCores, item.Load)
			problems = append(problems, []string{resource, problemTitle, problemDesc})
		}
		// 如果cpu使用率大于80%
		global.Logger.Debug(fmt.Sprintf("检测cpu使用率(使用率是否大于80%%): %s", item.IP))
		if item.CPUPercent > float64(80) {
			resource := "主机"
			problemTitle := "cpu使用率异常"
			problemDesc := fmt.Sprintf("主机: %s, cpu使用率%f", item.HostName, item.CPUPercent)
			problems = append(problems, []string{resource, problemTitle, problemDesc})
		}

		// 如果内存使用率大于80%
		global.Logger.Debug(fmt.Sprintf("检测内存使用率(使用率是否大于80%%): %s", item.IP))
		if item.MemPercent > float64(80) {
			resource := "主机"
			problemTitle := "内存使用率异常"
			problemDesc := fmt.Sprintf("主机: %s, 内存使用率%f", item.HostName, item.MemPercent)
			problems = append(problems, []string{resource, problemTitle, problemDesc})
		}

		// 如果磁盘使用率大于50%
		global.Logger.Debug(fmt.Sprintf("检测磁盘使用率(使用率是否大于50%%): %s", item.IP))
		for _, disk := range item.Disk {
			// block
			if disk.UsedPercent > float64(50) {
				resource := "主机"
				problemTitle := "磁盘使用率（block）异常"
				problemDesc := fmt.Sprintf("主机: %s, 挂载点: %s, 磁盘block使用率%f", item.HostName, disk.Path, disk.UsedPercent)
				problems = append(problems, []string{resource, problemTitle, problemDesc})
			}
			// inode
			if disk.InodesUsedPercent > float64(50) {
				resource := "主机"
				problemTitle := "磁盘使用率（inode）异常"
				problemDesc := fmt.Sprintf("主机: %s, 挂载点: %s, 磁盘inode使用率%f", item.HostName, disk.Path, disk.InodesUsedPercent)
				problems = append(problems, []string{resource, problemTitle, problemDesc})
			}
		}

		// 如果链接不上时间同步服务器，或者ntpserver和本地时间相差1s
		global.Logger.Debug(fmt.Sprintf("检测时间是否同步: %s", item.IP))
		if !item.NTPState {
			resource := "主机"
			problemTitle := "时间同步异常"
			problemDesc := fmt.Sprintf("主机: %s, 时间同步异常", item.HostName)
			problems = append(problems, []string{resource, problemTitle, problemDesc})
		}

	}
	// 渲染
	return problems, nil
}

type RedisCheckOps struct {
	Cluster       bool
	Client        *goredis.Client
	ClusterClient *goredis.ClusterClient
}

func RedisCheck() (problems [][]string) {
	problems = make([][]string, 0, 0)
	// redis客户端实例化
	redisClient := &RedisCheckOps{}
	if global.RedisSetting.Cluster {
		redisClient.Cluster = true   //集群判断
		global.Logger.Debug(fmt.Sprintf("检测redis联通性"))
		client, err := redisManage.SetUpGoRedisCluster(global.RedisSetting.ClusterAddr, global.RedisSetting.Pass)   //检测命令
		if err != nil {
			resource := "redis"
			problemTitle := "redis链接失败"
			problemDesc := fmt.Sprintf("redis节点: %s, 错误: %v", strings.Join(global.RedisSetting.ClusterAddr, ","), err)
			problems = append(problems, []string{resource, problemTitle, problemDesc})
		}
		redisClient.ClusterClient = client

		global.Logger.Debug(fmt.Sprintf("检测redis数据是否做了初始化"))
		dbSize, err := redisClient.ClusterClient.DBSize(context.Background()).Result()
		if err != nil {
			resource := "redis"
			problemTitle := "redis获取keys失败"
			problemDesc := fmt.Sprintf("redis节点: %s, 错误: %v", strings.Join(global.RedisSetting.ClusterAddr, ","), err)
			problems = append(problems, []string{resource, problemTitle, problemDesc})
		}

		var keyErrNumber int64 = 1000
		if dbSize < keyErrNumber {
			resource := "redis"
			problemTitle := "redis初始化异常"
			problemDesc := fmt.Sprintf("redis节点: %s, 错误: key数量为%d少于%d", strings.Join(global.RedisSetting.ClusterAddr, ","), dbSize, keyErrNumber)
			problems = append(problems, []string{resource, problemTitle, problemDesc})
		}

	} else {
		redisClient.Cluster = false
		global.Logger.Debug(fmt.Sprintf("检测redis联通性"))
		client, err := redisManage.SetUpGoRedis(0)
		if err != nil {
			resource := "redis"
			problemTitle := "redis链接失败"
			problemDesc := fmt.Sprintf("redis节点: %s, 错误: %v", global.RedisSetting.IP, err)
			problems = append(problems, []string{resource, problemTitle, problemDesc})
		}
		redisClient.Client = client

		global.Logger.Debug(fmt.Sprintf("检测redis数据是否做了初始化"))
		dbSize, err := redisClient.Client.DBSize(context.Background()).Result()
		if err != nil {
			resource := "redis"
			problemTitle := "redis获取keys失败"
			problemDesc := fmt.Sprintf("redis节点: %s, 错误: %v", global.RedisSetting.IP, err)
			problems = append(problems, []string{resource, problemTitle, problemDesc})
		}

		var keyErrNumber int64 = 1000
		if dbSize < keyErrNumber {
			resource := "redis"
			problemTitle := "redis初始化异常"
			problemDesc := fmt.Sprintf("redis节点: %s, 错误: key数量为%d少于%d", global.RedisSetting.IP, dbSize, keyErrNumber)
			problems = append(problems, []string{resource, problemTitle, problemDesc})
		}
	}

	return problems
}

func MysqlCheck() (problems [][]string) {
	problems = make([][]string, 0, 0)
	for _, dbInfo := range *global.MySQLSetting {
		p := MysqlSingleDBCheck(dbInfo)
		if p != nil {
			problems = append(problems, p...)
		}
	}
	return problems
}

func MysqlSingleDBCheck(m setting.MySQLInfo) (p [][]string) {
	p = make([][]string, 0, 0)
	global.Logger.Debug(fmt.Sprintf("检测mysql联通性-%s", m.DBName))  //取数据库名称
	db, err := mysql.NewSqx(m.IP, m.User, m.Pass, "information_schema", "utf8", m.Port, 2, 1)
	if err != nil {
		resource := "mysql"
		problemTitle := "mysql链接异常"
		problemDesc := fmt.Sprintf("mysql节点: %s, 库: %s 错误: %v", m.IP, m.DBName, err)
		p = append(p, []string{resource, problemTitle, problemDesc})
		return p
	}
	defer db.Close()

	global.Logger.Debug(fmt.Sprintf("检测mysql数据是否已经初始化-%s", m.DBName))
	var count int
	sqlStr := `SELECT COUNT(*) FROM tables WHERE table_schema = ?`
	if err := db.Get(&count, sqlStr, m.DBName); err != nil {
		resource := "mysql"
		problemTitle := "mysql链接异常"
		problemDesc := fmt.Sprintf("mysql节点: %s, 库: %s 错误: %v", m.IP, m.DBName, err)
		p = append(p, []string{resource, problemTitle, problemDesc})
		return p
	}

	if count == 0 {
		resource := "mysql"
		problemTitle := "mysql初始化失败"
		problemDesc := fmt.Sprintf("mysql节点: %s, 库: %s表数量0", m.IP, m.DBName)
		p = append(p, []string{resource, problemTitle, problemDesc})
	}

	return p
}

func ElasticCheck() (problems [][]string) {
	problems = make([][]string, 0, 0)
	global.Logger.Debug(fmt.Sprintf("检测elastic联通性"))
	for _, addrc := range global.ElasticS.Address {
		addr := strings.Split(addrc, "//")[1]
		ok, err := CheckTcpPort(addr)
		if !ok {
			resource := "elastic"
			problemTitle := "端口异常"
			problemDesc := fmt.Sprintf("节点: %s, 错误: %v", addrc, err)
			problems = append(problems, []string{resource, problemTitle, problemDesc})
		}
	}
	var client *elastic.Client
	var err error
	if global.ElasticS.Auth {
		client, err = elastic.NewSimpleClient(
			elastic.SetURL(global.ElasticS.Address...), elastic.SetBasicAuth(global.ElasticS.User, global.ElasticS.Pass))
	} else {
		client, err = elastic.NewSimpleClient(
			elastic.SetURL(global.ElasticS.Address...))
	}

	global.Logger.Debug(fmt.Sprintf("检测elastic状态是否健康"))
	if err != nil {
		resource := "elastic"
		problemTitle := "集群连接异常"
		problemDesc := fmt.Sprintf("错误: %v", err)
		problems = append(problems, []string{resource, problemTitle, problemDesc})
		return problems
	}

	chr, err := client.ClusterHealth().Do(context.Background())
	if err != nil {
		resource := "elastic"
		problemTitle := "集群连接异常"
		problemDesc := fmt.Sprintf("错误: %v", err)
		problems = append(problems, []string{resource, problemTitle, problemDesc})
		return problems
	} else {
		if chr.Status != "green" {
			resource := "elastic"
			problemTitle := "集群状态异常"
			problemDesc := fmt.Sprintf("状态: %s", chr.Status)
			problems = append(problems, []string{resource, problemTitle, problemDesc})
		}
	}

	global.Logger.Debug(fmt.Sprintf("检测elastic数据是否已经做了初始化"))
	cir, err := client.CatIndices().Do(context.Background())
	if err != nil {
		resource := "elastic"
		problemTitle := "集群连接异常"
		problemDesc := fmt.Sprintf("错误: %v", err)
		problems = append(problems, []string{resource, problemTitle, problemDesc})
		return problems
	} else {
		if len(cir) < 20 {
			resource := "elastic"
			problemTitle := "集群初始化异常"
			problemDesc := fmt.Sprintf("索引个数异常，%d", len(cir))
			problems = append(problems, []string{resource, problemTitle, problemDesc})
		}
	}
	return problems
}

func ZookeeperCheck() (problems [][]string) {
	global.Logger.Debug(fmt.Sprintf("检测zookeeper联通性"))
	problems = make([][]string, 0, 0)
	for _, addr := range global.ZookeeperS.Address {
		ok, err := CheckTcpPort(addr)
		if !ok {
			resource := "zookeeper"
			problemTitle := "端口异常"
			problemDesc := fmt.Sprintf("节点: %s, 错误: %v", addr, err)
			problems = append(problems, []string{resource, problemTitle, problemDesc})
		}
	}

	//conn, _, err := zk.Connect(global.ZookeeperS.Address, time.Second*5)
	//if err != nil {
	//	resource := "zookeeper"
	//	problemTitle := "集群连接异常"
	//	problemDesc := fmt.Sprintf("错误: %v", err)
	//	problems = append(problems, []string{resource, problemTitle, problemDesc})
	//	return problems
	//}
	//defer conn.Close()
	//
	//data, _, err := conn.Get("/zookeeper/stats")
	//fmt.Println(string(data))

	return problems
}

type KfkOps struct {
	Admin  sarama.ClusterAdmin
	Client sarama.Client
}

func NewKfkOps() (*KfkOps, error) {
	if len(global.KafkaSetting.Adders) == 0 {
		return nil, errors.New("创建kafka实例失败, 请在配置文件里面配置kafka实例")
	}
	config := sarama.NewConfig()
	admin, err := sarama.NewClusterAdmin(global.KafkaSetting.Adders, config)
	if err != nil {
		return nil, err
	}
	client, err := sarama.NewClient(global.KafkaSetting.Adders, config)
	if err != nil {
		return nil, err
	}

	return &KfkOps{Admin: admin, Client: client}, nil
}

func (k *KfkOps) Close() {
	if k.Admin != nil {
		k.Admin.Close()
	}
	if k.Client != nil {
		k.Client.Closed()
	}
}

func KafkaCheck() (problems [][]string) {
	global.Logger.Debug(fmt.Sprintf("检测kafka联通性"))
	problems = make([][]string, 0, 0)
	for _, addr := range global.KafkaSetting.Adders {
		ok, err := CheckTcpPort(addr)
		if !ok {
			resource := "kafka"
			problemTitle := "端口异常"
			problemDesc := fmt.Sprintf("节点: %s, 错误: %v", addr, err)
			problems = append(problems, []string{resource, problemTitle, problemDesc})
		}
	}

	global.Logger.Debug("检测kafka集群是否健康")
	kafkaOps, err := NewKfkOps()
	if err != nil {
		resource := "kafka"
		problemTitle := "集群初始化失败"
		problemDesc := fmt.Sprintf("错误: %v", err)
		problems = append(problems, []string{resource, problemTitle, problemDesc})
		return
	}
	defer kafkaOps.Close()

	if len(global.KafkaSetting.Adders) != len(kafkaOps.Client.Brokers()) {
		resource := "kafka"
		problemTitle := "集群节点不健康"
		problemDesc := fmt.Sprintf("期望节点个数: %d, 实际节点个数%d", len(global.KafkaSetting.Adders), len(kafkaOps.Client.Brokers()))
		problems = append(problems, []string{resource, problemTitle, problemDesc})
	}

	global.Logger.Debug("检测kafka集群数据是否已做初始化")
	topics, err := kafkaOps.Admin.ListTopics()
	if err != nil {
		resource := "kafka"
		problemTitle := "获取topic失败"
		problemDesc := fmt.Sprintf("%v", err)
		problems = append(problems, []string{resource, problemTitle, problemDesc})
		return
	}

	if len(topics) < 10 {
		resource := "kafka"
		problemTitle := "数据初始化失败"
		problemDesc := fmt.Sprintf("topic数: %d", len(topics))
		problems = append(problems, []string{resource, problemTitle, problemDesc})
		return
	}

	return problems
}

var (
	DataNodePort = 9866
	HMasterPort  = 60010
	HRegionPort  = 16020
)

func HadoopCheck() (problems [][]string) {
	problems = make([][]string, 0, 0)
	global.Logger.Debug(fmt.Sprintf("检测hadoop namenode 联通性"))
	ok, err := CheckTcpPort(global.HadoopS.NameNode)
	if !ok {
		resource := "hadoop"
		problemTitle := "namenode端口异常"
		problemDesc := fmt.Sprintf("节点: %s, 错误: %v", global.HadoopS.NameNode, err)
		problems = append(problems, []string{resource, problemTitle, problemDesc})
	}

	global.Logger.Debug(fmt.Sprintf("检测hadoop datanode 联通性"))
	for _, ip := range global.HadoopS.DataNode {
		addr := fmt.Sprintf("%s:%d", ip, DataNodePort)
		ok, err := CheckTcpPort(addr)
		if !ok {
			resource := "hadoop"
			problemTitle := "datanode端口异常"
			problemDesc := fmt.Sprintf("节点: %s, 错误: %v", addr, err)
			problems = append(problems, []string{resource, problemTitle, problemDesc})
		}
	}
	return problems
}

func HbaseCheck() (problems [][]string) {
	problems = make([][]string, 0, 0)
	global.Logger.Debug(fmt.Sprintf("检测hbase master 联通性"))
	ip := strings.Split(global.HadoopS.NameNode, ":")[0]
	addr := fmt.Sprintf("%s:%d", ip, HMasterPort)
	ok, err := CheckTcpPort(addr)
	if !ok {
		resource := "hbase"
		problemTitle := "HMaster端口异常"
		problemDesc := fmt.Sprintf("节点: %s, 错误: %v", addr, err)
		problems = append(problems, []string{resource, problemTitle, problemDesc})
	}

	global.Logger.Debug(fmt.Sprintf("检测hbase region server 联通性"))
	for _, ip := range global.HadoopS.DataNode {
		addr := fmt.Sprintf("%s:%d", ip, HRegionPort)
		ok, err := CheckTcpPort(addr)
		if !ok {
			resource := "hbase"
			problemTitle := "Region Server 端口异常"
			problemDesc := fmt.Sprintf("节点: %s, 错误: %v", addr, err)
			problems = append(problems, []string{resource, problemTitle, problemDesc})
		}
	}
	return problems
}

var (
	StreamparkAppStateOk = "RUNNING"
)

func StreamParkCheck() (problems [][]string) {
	problems = make([][]string, 0, 0)
	global.Logger.Debug(fmt.Sprintf("检测streampark联通性"))
	addr := strings.Split(global.StreamparkS.Adder, "//")[1]
	ok, err := CheckTcpPort(addr)
	if !ok {
		resource := "streampark"
		problemTitle := "端口异常"
		problemDesc := fmt.Sprintf("节点: %s, 错误: %v", addr, err)
		problems = append(problems, []string{resource, problemTitle, problemDesc})
	}

	apps, err := streampark.NewAppOps().ListCheck()
	if err != nil {
		resource := "streampark"
		problemTitle := "获取app列表异常"
		problemDesc := fmt.Sprintf("错误: %v", err)
		problems = append(problems, []string{resource, problemTitle, problemDesc})
		return
	}

	for _, app := range apps {
		if app.StateHuman() != StreamparkAppStateOk {
			resource := "streampark"
			problemTitle := "app运行异常"
			problemDesc := fmt.Sprintf("app: %s 状态: %v", app.JobName, app.StateHuman())
			problems = append(problems, []string{resource, problemTitle, problemDesc})
		}
	}

	return problems
}

var (
	DorisFePort = 9030
	DorisBePort = 9050
)

func DorisCheck() (problems [][]string) {
	problems = make([][]string, 0, 0)
	global.Logger.Debug(fmt.Sprintf("检测联通性"))
	addr := fmt.Sprintf("%s:%d", global.DorisS.IP, DorisFePort)
	ok, err := CheckTcpPort(addr)
	if !ok {
		resource := "doris"
		problemTitle := "FE端口异常"
		problemDesc := fmt.Sprintf("节点: %s, 错误: %v", addr, err)
		problems = append(problems, []string{resource, problemTitle, problemDesc})
		return
	}

	addr = fmt.Sprintf("%s:%d", global.DorisS.IP, DorisBePort)
	ok, err = CheckTcpPort(addr)
	if !ok {
		resource := "doris"
		problemTitle := "BE端口异常"
		problemDesc := fmt.Sprintf("节点: %s, 错误: %v", addr, err)
		problems = append(problems, []string{resource, problemTitle, problemDesc})
		return
	}

	global.Logger.Debug(fmt.Sprintf("检测数据"))
	db, err := NewSqlDb(global.DorisS.IP, global.DorisS.User, global.DorisS.Password, "information_schema", global.DorisS.Port)
	if err != nil {
		resource := "doris"
		problemTitle := "库连接异常"
		problemDesc := fmt.Sprintf("节点: %s, 库: %s 错误: %v", fmt.Sprintf("%s:%d", global.DorisS.IP, DorisFePort), global.DorisS.DBName, err)
		problems = append(problems, []string{resource, problemTitle, problemDesc})
		return problems
	}
	defer db.Close()

	var count int
	sqlStr := fmt.Sprintf("SELECT COUNT(*) FROM tables WHERE table_schema = '%s'", global.DorisS.DBName)
	if err := db.QueryRow(sqlStr).Scan(&count); err != nil {
		resource := "doris"
		problemTitle := "库连接异常"
		problemDesc := fmt.Sprintf("节点: %s, 库: %s 错误: %v", fmt.Sprintf("%s:%d", global.DorisS.IP, DorisFePort), global.DorisS.DBName, err)
		problems = append(problems, []string{resource, problemTitle, problemDesc})
		return problems
	}

	if count == 0 {
		resource := "doris"
		problemTitle := "初始化失败"
		problemDesc := fmt.Sprintf("节点: %s, 库: %s表数量0", fmt.Sprintf("%s:%d", global.DorisS.IP, DorisFePort), global.DorisS.DBName)
		problems = append(problems, []string{resource, problemTitle, problemDesc})
	}

	return problems
}

func PortCheck() (problems [][]string) {
	problems = make([][]string, 0, 0)
	for _, item := range *global.TCPPortCheckS {
		global.Logger.Debug(fmt.Sprintf("探测%s: %s", item.Name, item.Value))
		ok, err := CheckTcpPort(item.Value)
		if !ok {
			resource := "端口检测"
			problemTitle := "检测失败"
			problemDesc := fmt.Sprintf("探测%s: %s, 错误: %v", item.Name, item.Value, err)
			problems = append(problems, []string{resource, problemTitle, problemDesc})
		}
	}

	return problems
}

func K8SCheck() (problems [][]string) {
	problems = make([][]string, 0, 0)
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		resource := "k8s"
		problemTitle := "config创建失败"
		problemDesc := fmt.Sprintf("错误: %v", err)
		problems = append(problems, []string{resource, problemTitle, problemDesc})
		return
	}
	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		resource := "k8s"
		problemTitle := "client创建失败"
		problemDesc := fmt.Sprintf("错误: %v", err)
		problems = append(problems, []string{resource, problemTitle, problemDesc})
		return
	}

	global.Logger.Debug("检测node")
	nodeList, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		resource := "k8s"
		problemTitle := "获取节点列表失败"
		problemDesc := fmt.Sprintf("错误: %v", err)
		problems = append(problems, []string{resource, problemTitle, problemDesc})
	}
	for _, node := range nodeList.Items {
		if node.Status.NodeInfo.KubeletVersion == "" {
			resource := "k8s"
			problemTitle := "节点状态异常"
			problemDesc := fmt.Sprintf("节点: %s", node.Name)
			problems = append(problems, []string{resource, problemTitle, problemDesc})
		}
	}

	global.Logger.Debug("检测pod")
	podList, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		resource := "k8s"
		problemTitle := "获取pod列表失败"
		problemDesc := fmt.Sprintf("错误: %v", err)
		problems = append(problems, []string{resource, problemTitle, problemDesc})
		return
	}
	for _, pod := range podList.Items {
		switch pod.Status.Phase {
		case v1.PodFailed, v1.PodReasonUnschedulable, v1.PodPending:
			podErrMsg := fmt.Sprintf("%s: %s", pod.Status.Reason, pod.Status.Message)
			resource := "k8s"
			problemTitle := "pod运行异常"
			problemDesc := fmt.Sprintf("pod: %s, pod状态: %s, pod详情: %s", pod.GetName(), string(pod.Status.Phase), podErrMsg)
			problems = append(problems, []string{resource, problemTitle, problemDesc})
		default:
			continue
		}

	}
	return
}

func CheckTcpPort(addr string) (ok bool, err error) {
	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		return false, err
	}
	defer conn.Close()
	return true, nil
}

func NewSqlDb(host, user, password, dbname string, port int) (db *sql.DB, err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		user, password, host, port, dbname)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("校验失败,err", err)
		return nil, err
	}
	return db, nil
}
