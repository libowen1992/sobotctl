package redisManage

import (
	"github.com/pkg/errors"
	"sobotctl/pkg/convert"
	"sobotctl/pkg/tableRender"
)

func (ro *RedisOps) SearchKey(db int, key string) error {
	client, err := NewGoRedis(db)
	if err != nil {
		return errors.WithStack(err)
	}
	defer client.Close()
	keys, err := client.SearchKey(key)
	if err != nil {
		return errors.WithStack(err)
	}
	headers := []string{"key", "类型", "失效时间"}
	data := make([][]string, 0, len(keys))
	for _, v := range keys {
		data = append(data, []string{v.Key, v.KeyType, convert.Int64ToStr(v.TTL)})
	}
	tableRender.Render(headers, data)
	return nil
}
