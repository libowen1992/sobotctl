package redisManage

import (
	"github.com/pkg/errors"
	"sobotctl/pkg/tableRender"
)

func (ro *RedisOps) Config() error {
	client, err := NewGoRedis(0)
	if err != nil {
		return errors.WithStack(err)
	}
	defer client.Close()
	configs, err := client.GetConfig()
	if err != nil {
		return errors.WithStack(err)
	}

	headers := []string{"配置", "值"}
	data := make([][]string, 0, len(configs))
	for _, v := range configs {
		data = append(data, []string{v.Name, v.Value})
	}
	tableRender.Render(headers, data)
	return nil
}
