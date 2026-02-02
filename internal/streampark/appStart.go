package streampark

import "sobotctl/global"

func (a *AppOps) Start(appId string) error {
	client, err := NewStreamClient(global.StreamparkS.User,
		global.StreamparkS.Password, global.StreamparkS.Adder)

	if err != nil {
		return err
	}
	if err := client.startApp(appId); err != nil {
		return err
	}

	return nil
}
