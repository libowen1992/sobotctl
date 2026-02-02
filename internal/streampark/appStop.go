package streampark

import (
	"github.com/pkg/errors"
	"sobotctl/global"
)

func (a *AppOps) Stop(appId string) error {
	client, err := NewStreamClient(global.StreamparkS.User,
		global.StreamparkS.Password, global.StreamparkS.Adder)

	if err != nil {
		return errors.Wrap(err, "初始化client失败")
	}
	if err := client.stopApp(appId); err != nil {
		return errors.Wrap(err, "stop服务失败")
	}

	return nil
}
