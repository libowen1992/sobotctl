package cmd

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os/exec"
	"path"
	"sobotctl/global"
	"sobotctl/internal/redisManage"
)

func NewRedisManager() *cobra.Command {
	var RedisManagerCmd = &cobra.Command{
		Use:   "redis",
		Short: "redis管理",
	}
	RedisManagerCmd.AddCommand(newRedisRestore())
	RedisManagerCmd.AddCommand(NewRedisTerminal())
	RedisManagerCmd.AddCommand(NewRedisSlowLog())
	RedisManagerCmd.AddCommand(NewRedisKey())
	RedisManagerCmd.AddCommand(NewRedisConfig())
	RedisManagerCmd.AddCommand(NewRedisCheckCmd())
	return RedisManagerCmd
}

func newRedisRestore() *cobra.Command {
	var rs = &redisManage.RedisRestore{}
	//var
	var redisSyncCmd = &cobra.Command{
		Use:   "restore",
		Short: "同步rdb文件数据到redis实例",
		Args: func(cmd *cobra.Command, args []string) error {
			if !path.IsAbs(rs.RdbFile) {
				return errors.New("rdb文件必须为绝对路径")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := rs.Restore(); err != nil {
				global.Logger.Error(err)
			}
		},
	}
	redisSyncCmd.Flags().StringVarP(&rs.RdbFile, "rdb", "s", "", "rdb文件的绝对路径路径")
	redisSyncCmd.MarkFlagRequired("rdb")

	redisSyncCmd.Flags().StringVarP(&rs.RedisAddr, "redisAddr", "d", "", "redis地址")
	redisSyncCmd.MarkFlagRequired("redisAddr")

	redisSyncCmd.Flags().StringVarP(&rs.RedisPort, "redisPort", "p", "6379", "redis端口，默认6379")

	redisSyncCmd.Flags().StringVarP(&rs.RedisPass, "redisPass", "a", "", "redis密码")
	redisSyncCmd.MarkFlagRequired("redisPass")

	return redisSyncCmd
}

func NewRedisTerminal() *cobra.Command {
	var redisTerminalCMD = &cobra.Command{
		Use:   "terminal",
		Short: "终端",
		Run: func(cmd *cobra.Command, args []string) {
			_, err := exec.LookPath("redis-cli")
			if err != nil {
				global.Logger.Error("请先安装redis客户端 redis-cli")
				return
			}
			if err := redisManage.NewRedisOps().Terminal(); err != nil {
				global.Logger.Error(err)
			}
		},
	}
	return redisTerminalCMD
}

func NewRedisSlowLog() *cobra.Command {
	var redisSlowLogCMD = &cobra.Command{
		Use:   "slowlog",
		Short: "慢日志",
		Run: func(cmd *cobra.Command, args []string) {
			if err := redisManage.NewRedisOps().SlowLog(); err != nil {
				global.Logger.Error(err)
			}
		},
	}

	return redisSlowLogCMD
}

func NewRedisConfig() *cobra.Command {
	var redisConfigCMD = &cobra.Command{
		Use:   "config",
		Short: "查看配置",
		Run: func(cmd *cobra.Command, args []string) {
			if err := redisManage.NewRedisOps().Config(); err != nil {
				global.Logger.Error(err)
			}
		},
	}

	return redisConfigCMD
}

func NewRedisKey() *cobra.Command {
	var redisKeyCMD = &cobra.Command{
		Use:   "key",
		Short: "key相关操作",
	}
	redisKeyCMD.AddCommand(NewRedisSearchKey())
	redisKeyCMD.AddCommand(NewRedisGetKey())
	return redisKeyCMD
}

func NewRedisSearchKey() *cobra.Command {
	var db int
	var redisSearchKeyCMD = &cobra.Command{
		Use:   "search",
		Short: "模糊查询key",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("必须给定要查询的模糊key")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			key := args[0]
			key = fmt.Sprintf("*%s*", key)
			if err := redisManage.NewRedisOps().SearchKey(db, key); err != nil {
				global.Logger.Error(err)
			}
		},
	}
	redisSearchKeyCMD.Flags().IntVarP(&db, "db", "d", 0, "库号，如0-15")
	return redisSearchKeyCMD
}

func NewRedisGetKey() *cobra.Command {
	var db int
	var redisGetKeyCMD = &cobra.Command{
		Use:   "value",
		Short: "查询key值",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("必须给定要查询的key")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			key := args[0]
			if err := redisManage.NewRedisOps().GetKey(db, key); err != nil {
				global.Logger.Error(err)
			}
		},
	}
	redisGetKeyCMD.Flags().IntVarP(&db, "db", "d", 0, "库号，如0-15")
	return redisGetKeyCMD
}

func NewRedisCheckCmd() *cobra.Command {
	action := "check"
	desc := "检查redis状态"
	var Cmd = &cobra.Command{
		Use:   action,
		Short: desc,
		Run: func(cmd *cobra.Command, args []string) {
			redisManage.NewRedisOps().Status()
		},
	}
	return Cmd
}
