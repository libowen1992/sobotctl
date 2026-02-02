package mysqlOps

import (
	"fmt"
	"os"
	"os/exec"
	"sobotctl/global"
)

type MysqlInitArgus struct {
	DB     string
	DBFile string
}

func (mo *MysqlOps) InitDb(argus MysqlInitArgus) error {
	mInfo, _ := global.MySQLSetting.GetInfo(argus.DB)
	cmdStr := fmt.Sprintf("mysql -h %s  -P %d  -B %s -u %s -p%s < %s",
		mInfo.IP, mInfo.Port, argus.DB, mInfo.User, mInfo.Pass, argus.DBFile)
	cmd := exec.Command("bash", "-c", cmdStr)
	global.Logger.Info(cmd.String())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
