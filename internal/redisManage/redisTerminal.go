package redisManage

import (
	"fmt"
	"os"
	"os/exec"
	"sobotctl/global"
)

func (ro *RedisOps) Terminal() error {
	connect := fmt.Sprintf("redis-cli --raw -a '%s' -h %s -p %d",
		global.RedisSetting.Pass, global.RedisSetting.IP, global.RedisSetting.Port)
	cmd := exec.Command("bash", "-c", connect)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
