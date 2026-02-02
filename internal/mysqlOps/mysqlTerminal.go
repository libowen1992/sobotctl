package mysqlOps

import (
	"fmt"
	"os"
	"os/exec"
	"sobotctl/global"
)

func (mo *MysqlOps) Terminal(dbName string) error {
	dbInfo, err := global.MySQLSetting.GetInfo(dbName)
	if err != nil {
		return err
	}
	connect := fmt.Sprintf("mysql -h %s -u%s  -P%d -p'%s'",
		dbInfo.IP, dbInfo.User,
		dbInfo.Port, dbInfo.Pass)
	cmd := exec.Command("bash", "-c", connect)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
