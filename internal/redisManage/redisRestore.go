package redisManage

import (
	"fmt"
	"os"
	"os/exec"
	"sobotctl/global"
	"text/template"
)

const (
	redisRestoreConfigTemplate = "./tools/redis/restore.tmpl"
	redisShakeBin              = "./tools/redis/redis-shake"
	redisRestoreConfig         = "./tmp/redis_restore.toml"
)

type RedisRestore struct {
	RdbFile   string
	RedisAddr string
	RedisPort string
	RedisPass string
}

func (r *RedisRestore) Restore() error {
	// 生成模版
	file, err := template.ParseFiles(redisRestoreConfigTemplate)
	if err != nil {
		return err
	}

	redisC, err := os.OpenFile(redisRestoreConfig, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer redisC.Close()

	err = file.Execute(redisC, r)
	if err != nil {
		return err
	}

	// 执行命令
	//（实时）
	cmdStr := fmt.Sprintf("%s %s", redisShakeBin, redisRestoreConfig)
	cmd := exec.Command("/bin/bash", "-c", cmdStr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = cmd.Stdout
	global.Logger.Debug(cmd.String())
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
