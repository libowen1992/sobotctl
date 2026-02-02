package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path"
	"sobotctl/global"
	"sobotctl/internal/mysqlOps"
	"strings"
)

func NewMysqlManager() *cobra.Command {
	var mysqlManagerCmd = &cobra.Command{
		Use:   "mysql",
		Short: "mysql管理工具",
	}

	mysqlManagerCmd.AddCommand(NewMysqlTerminal())
	mysqlManagerCmd.AddCommand(NewMysqlInitDb())
	mysqlManagerCmd.AddCommand(NewMysqlCheckCmd())
	return mysqlManagerCmd
}

func NewMysqlTerminal() *cobra.Command {
	var mysqlTerminalCmd = &cobra.Command{
		Use:   "terminal",
		Short: "mysql终端",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.Errorf("请输入要连接的dbname %v", global.MySQLSetting.ListDBName())
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			_, err := exec.LookPath("mysql")
			if err != nil {
				global.Logger.Error("请先安装mysql客户端")
				return
			}
			dbName := strings.TrimSpace(args[0])
			if err := mysqlOps.NewMysqlOps().Terminal(dbName); err != nil {
				global.Logger.Error(err)
			}
		},
	}

	return mysqlTerminalCmd
}

func NewMysqlInitDb() *cobra.Command {
	var argus mysqlOps.MysqlInitArgus
	var mysqlInitDBCmd = &cobra.Command{
		Use:   "init",
		Short: "初始化数据库，适合初始化客户的数据库，其他场景不要使用",
		Args: func(cmd *cobra.Command, args []string) error {
			if !path.IsAbs(argus.DBFile) {
				return errors.New("初始化文件必须为绝对路径")
			}
			if _, err := global.MySQLSetting.GetInfo(argus.DB); err != nil {
				return errors.Errorf("%s实例找不到，请配置在配置文件里面，可用的数据库实例 %v", argus.DB, global.MySQLSetting.ListDBName())
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			_, err := exec.LookPath("mysql")
			if err != nil {
				global.Logger.Error("请先安装mysql客户端")
				return
			}
			if err := mysqlOps.NewMysqlOps().InitDb(argus); err != nil {
				global.Logger.Error(err)
				os.Exit(1)
			}
		},
	}

	mysqlInitDBCmd.Flags().StringVarP(&argus.DBFile, "file", "f", "", "初始化文件")
	mysqlInitDBCmd.MarkFlagRequired("file")

	mysqlInitDBCmd.Flags().StringVarP(&argus.DB, "db", "d", "", "数据库，如 sobot_db")
	mysqlInitDBCmd.MarkFlagRequired("db")

	return mysqlInitDBCmd
}

func NewMysqlCheckCmd() *cobra.Command {
	action := "check"
	desc := "mysql检查状态"
	var Cmd = &cobra.Command{
		Use:   action,
		Short: desc,
		Run: func(cmd *cobra.Command, args []string) {
			mysqlOps.NewMysqlOps().Check()
		},
	}
	return Cmd
}
