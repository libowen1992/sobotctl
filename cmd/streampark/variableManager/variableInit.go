package variableManager

import (
	"github.com/spf13/cobra"
	"sobotctl/global"
	"sobotctl/internal/streampark"
)

func NewVariableInit() *cobra.Command {
	var variableInit streampark.VariableInit

	var variableInitCMD = &cobra.Command{
		Use:   "init",
		Short: "初始化,主要用于第一次部署后，批量修改变量",
		Run: func(cmd *cobra.Command, args []string) {
			if err := streampark.NewVariableOps().Init(variableInit); err != nil {
				global.Logger.Error(err)
				return
			}
			global.Logger.Info("Success!!")
		},
	}

	variableInitCMD.Flags().StringVarP(&variableInit.DbHost, "dbHost", "", "", "数据库地址,如 10.0.0.42")
	variableInitCMD.Flags().StringVarP(&variableInit.DbPort, "dbPort", "", "", "数据库端口, 如 3306")
	variableInitCMD.Flags().StringVarP(&variableInit.DbName, "dbName", "", "", "数据库名,主库名，如 sobot_db")
	variableInitCMD.Flags().StringVarP(&variableInit.DbUser, "dbUser", "", "", "数据库用户")
	variableInitCMD.Flags().StringVarP(&variableInit.DbPass, "dbPass", "", "", "数据库密码")
	variableInitCMD.Flags().StringVarP(&variableInit.IcalldbName, "icalldbName", "", "", "智能外呼数据库名称")
	variableInitCMD.Flags().StringVarP(&variableInit.IcalldbUser, "icalldbUser", "", "", "智能外呼数据库用户")
	variableInitCMD.Flags().StringVarP(&variableInit.IcalldbPass, "icalldbPass", "", "", "智能外呼数据库密码")

	variableInitCMD.Flags().StringVarP(&variableInit.RedisHost, "redisHost", "", "", "redis地址")
	variableInitCMD.Flags().StringVarP(&variableInit.RedisPort, "redisPort", "", "", "redis端口")
	variableInitCMD.Flags().StringVarP(&variableInit.RedisPass, "redisPass", "", "", "redis密码")

	variableInitCMD.Flags().StringVarP(&variableInit.DorisNodes, "dorisNodes", "", "", "doris fe nodes,如10.0.0.42:18030")
	variableInitCMD.Flags().StringVarP(&variableInit.DorisUser, "dorisUser", "", "", "doris用户，如 sobot_mysql_user")
	variableInitCMD.Flags().StringVarP(&variableInit.DorisPass, "dorisPass", "", "", "doris密码 如 xxxxx")

	variableInitCMD.Flags().StringVarP(&variableInit.KafkaAddrs, "kafkaAddrs", "", "", "kafka的地址，如 10.0.0.42:909,10.0.0.13:9092")

	variableInitCMD.Flags().StringVarP(&variableInit.HdfsUri, "hdfsUri", "", "", "hdfs地址，如 hdfs://10.0.0.42:9000")

	variableInitCMD.Flags().StringVarP(&variableInit.ElasticUrl, "elasticUrl", "", "", "elastic地址，如 10.0.0.42")
	variableInitCMD.Flags().StringVarP(&variableInit.ElasticPort, "elasticPort", "", "", "elastic端口，如 9200")

	return variableInitCMD
}
