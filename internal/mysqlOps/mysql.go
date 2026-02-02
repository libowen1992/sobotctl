package mysqlOps

import (
	"sobotctl/global"
	"sobotctl/pkg/mysql"
	"sobotctl/setting"
)

type MysqlOps struct {
}

func NewMysqlOps() *MysqlOps {
	return &MysqlOps{}
}

func (mo *MysqlOps) Check() {

	for _, dbInfo := range *global.MySQLSetting {
		mo.mysqlStatus(dbInfo)
	}
}

func (mo *MysqlOps) mysqlStatus(dbInfo setting.MySQLInfo) {
	errMsg := "mysql not ok, %s"
	db, err := mysql.NewSqx(dbInfo.IP, dbInfo.User, dbInfo.Pass, dbInfo.DBName, "utf8", dbInfo.Port, 2, 1)
	if err != nil {
		global.Logger.Errorf(errMsg, err.Error())
		return
	}
	defer db.Close()
	global.Logger.Infof("mysql: db %s ok !", dbInfo.DBName)
}

//func ()
