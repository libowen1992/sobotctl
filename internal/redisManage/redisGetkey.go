package redisManage

import (
	"fmt"
	"github.com/pkg/errors"
	"sobotctl/global"
)

func (ro *RedisOps) GetKey(db int, key string) error {
	client, err := NewGoRedis(db)
	if err != nil {
		return errors.WithStack(err)
	}
	defer client.Close()

	keyT, err := client.keyType(key)
	if err != nil {
		return errors.WithStack(err)
	}
	global.Logger.Infof("key类型: %s", keyT)
	var v interface{}
	switch keyT {
	case keyTypeString:
		v, err = client.GetStringValue(key)
		if err != nil {
			return errors.WithStack(err)
		}
		//fmt.Println(v)
	case keyTypeHash:
		v, err = client.GetHashValue(key)
		if err != nil {
			return errors.WithStack(err)
		}
	case keyTypeList:
		v, err = client.GetListValue(key)
		if err != nil {
			return errors.WithStack(err)
		}
	case keyTypeSet:
		v, err = client.GetSetValue(key)
		if err != nil {
			return errors.WithStack(err)
		}
		//fmt.Println(v)
	default:
		global.Logger.Warn("没有找到key")
	}
	fmt.Println(v)

	return nil
}
