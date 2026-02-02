package streampark

import (
	"github.com/jmoiron/sqlx"
	"sobotctl/pkg/mysql"
	"sobotctl/setting"
)

var (
	charset = "utf8mb4"
)

func NewMySQL(c *setting.StreamPark) (db *sqlx.DB, err error) {
	return mysql.NewSqx(c.DBHost, c.DBUser,
		c.DBPass, c.DBName, charset, c.DBPort, 5, 2)
}
