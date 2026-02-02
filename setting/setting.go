package setting

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Hosts struct {
	User            string   `mapstructure:"user"`
	SshType         string   `mapstructure:"ssh_type"`
	SshPass         string   `mapstructure:"ssh_pass"`
	SshKey          string   `mapstructure:"ssh_key"`
	HostctlPath     string   `mapstructure:"hostctl_path"`
	Port            int      `mapstructure:"port"`
	DiskCheckPoints []string `mapstructure:"disk_check_points"`
	NtpServer       string   `mapstructure:"ntp_server"`
	IPS             []string `mapstructure:"ips"`
}
type Redis struct {
	Cluster     bool     `mapstructure:"cluster"`
	ClusterAddr []string `mapstructure:"cluster_addr"`
	IP          string   `mapstructure:"ip"`
	Port        int      `mapstructure:"port"`
	Pass        string   `mapstructure:"pass"`
}

type Mysql []MySQLInfo

func (m Mysql) ListDBName() []string {
	var names []string
	for _, db := range m {
		names = append(names, db.DBName)
	}
	return names
}

func (m Mysql) GetInfo(db string) (*MySQLInfo, error) {
	for _, info := range m {
		if db == info.DBName {
			return &info, nil
		}
	}
	return nil, errors.Errorf("没有这个db名称")
}

type MySQLInfo struct {
	IP     string `mapstructure:"ip"`
	DBName string `mapstructure:"dbName"`
	Port   int    `mapstructure:"port"`
	User   string `mapstructure:"user"`
	Pass   string `mapstructure:"pass"`
}

type Elastic struct {
	Auth    bool     `mapstructure:"auth"`
	User    string   `mapstructure:"user"`
	Pass    string   `mapstructure:"pass"`
	Address []string `mapstructure:"address"`
}

type Zookeeper struct {
	Address []string `mapstructure:"address"`
}

type App struct {
	Name  string   `mapstructure:"name"`
	Nodes []string `mapstructure:"nodes"`
}

type Kafka struct {
	Adders []string `mapstructure:"adders"`
}

type Apps []*App

type Hadoop struct {
	NameNode string   `mapstructure:"namenode"`
	DataNode []string `mapstructure:"datanode"`
	User     string   `mapstructure:"user"`
}

type StreamPark struct {
	Adder             string          `mapstructure:"adder"`
	User              string          `mapstructure:"user"`
	Password          string          `mapstructure:"password"`
	DBHost            string          `mapstructure:"db_host"`
	DBPort            int             `mapstructure:"db_port"`
	DBName            string          `mapstructure:"db_name"`
	DBUser            string          `mapstructure:"db_user"`
	DBPass            string          `mapstructure:"db_pass"`
	JarHdfsBasePath   string          `mapstructure:"jar_hdfs_base_path"`
	JarHdfsBackupPath string          `mapstructure:"jar_hdfs_backup_path"`
	Apps              []StreamParkApp `mapstructure:"apps"`
}

type Doris struct {
	IP       string `mapstructure:"ip"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbName"`
}

type StreamParkApp struct {
	AppId  string `mapstructure:"app_id"`
	AppJar string `mapstructure:"app_jar"`
}

type Nacos struct {
	Host   string `mapstructure:"host"`
	Dbip   string `mapstructure:"dbip"`
	DbName string `mapstructure:"dbName"`
	Dbport int    `mapstructure:"dbport"`
	Dbuser string `mapstructure:"dbuser"`
	Dbpass string `mapstructure:"dbpass"`
}

type K8s struct {
	YamlDir string `mapstructure:"yamlDir"`
}

type Harbor struct {
	URL      string `mapstructure:"url"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type PortCheckItem struct {
	Name  string `mapstructure:"name"`
	Value string `mapstructure:"value"`
}

type TCPPortCheck []*PortCheckItem

type setting struct {
	vp *viper.Viper
}

func New() *setting {
	vp := viper.New()
	return &setting{vp: vp}
}

func (s setting) Init(file string) error {
	s.vp.SetConfigFile(file)
	s.vp.AddConfigPath(".")
	if err := s.vp.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

func (s setting) SetSection(k string, v interface{}) error {
	if err := s.vp.UnmarshalKey(k, v); err != nil {
		return err
	}
	return nil
}
