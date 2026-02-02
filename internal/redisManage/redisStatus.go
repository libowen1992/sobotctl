package redisManage

import (
	"context"
	"sobotctl/global"
	"time"
)

func (ro *RedisOps) Status() {
	okMsg := "redis is ok!"
	errMsg := "redis not ok, please check redis !"

	db, err := SetUpGoRedis(0)
	if err != nil {
		global.Logger.Error(errMsg, err)
		return
	}
	defer db.Close()

	if err = db.Set(context.Background(), "check", "check", time.Second).Err(); err != nil {
		global.Logger.Error(errMsg, err)
		return
	}
	global.Logger.Info(okMsg)
}
