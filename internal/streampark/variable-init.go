package streampark

import (
	"sobotctl/global"
)

type VariableInit struct {
	DbPass      string
	DbUser      string
	DbName      string
	DbPort      string
	DbHost      string
	IcalldbName string
	IcalldbUser string
	IcalldbPass string
	RedisHost   string
	RedisPort   string
	RedisPass   string
	DorisNodes  string
	DorisUser   string
	DorisPass   string
	KafkaAddrs  string
	HdfsUri     string
	ElasticUrl  string
	ElasticPort string
}

func (v *variableOps) Init(s VariableInit) error {
	db, err := NewMySQL(global.StreamparkS)
	if err != nil {
		return err
	}
	defer db.Close()

	if len(s.DbHost) != 0 {
		SqlStr := "update t_variable set variable_value = ? where variable_code  in" +
			" ('report.jdbc.host','callv6.jdbc.host','basic.jdbc.host','icall_new.jdbc.host','boss.jdbc.host')"
		if _, err := db.Exec(SqlStr, s.DbHost); err != nil {
			return err
		}
	}

	if len(s.DbPort) != 0 {
		SqlStr := "update t_variable set variable_value = ? where variable_code  in" +
			" ('report.jdbc.port','callv6.jdbc.port','basic.jdbc.port','icall_new.jdbc.port','boss.jdbc.port')"
		if _, err := db.Exec(SqlStr, s.DbPort); err != nil {
			return err
		}
	}

	if len(s.DbName) != 0 {
		SqlStr := "update t_variable set variable_value = ? where variable_code  in" +
			" ('report.jdbc.dbname','callv6.jdbc.dbname','basic.jdbc.dbName','boss.jdbc.dbname')"
		if _, err := db.Exec(SqlStr, s.DbName); err != nil {
			return err
		}
	}

	if len(s.DbUser) != 0 {
		SqlStr := "update t_variable set variable_value = ? where variable_code  in" +
			" ('report.jdbc.username','callv6.jdbc.username','basic.jdbc.username','boss.jdbc.username')"
		if _, err := db.Exec(SqlStr, s.DbUser); err != nil {
			return err
		}
	}

	if len(s.DbPass) != 0 {
		SqlStr := "update t_variable set variable_value = ? where variable_code  in" +
			" ('report.jdbc.password','callv6.jdbc.password','basic.jdbc.password','boss.jdbc.password')"
		if _, err := db.Exec(SqlStr, s.DbPass); err != nil {
			return err
		}
	}

	if len(s.IcalldbName) != 0 {
		SqlStr := "update t_variable set variable_value = ? where variable_code  in" +
			" ('icall_new.jdbc.dbname')"
		if _, err := db.Exec(SqlStr, s.IcalldbName); err != nil {
			return err
		}
	}

	if len(s.IcalldbUser) != 0 {
		SqlStr := "update t_variable set variable_value = ? where variable_code  in" +
			" ('icall_new.jdbc.username')"
		if _, err := db.Exec(SqlStr, s.IcalldbUser); err != nil {
			return err
		}
	}

	if len(s.IcalldbPass) != 0 {
		SqlStr := "update t_variable set variable_value = ? where variable_code  in" +
			" ('icall_new.jdbc.password')"
		if _, err := db.Exec(SqlStr, s.IcalldbPass); err != nil {
			return err
		}
	}

	if len(s.DorisNodes) != 0 {
		SqlStr := "update t_variable set variable_value = ? where variable_code  in" +
			" ('doris.fenodes')"
		if _, err := db.Exec(SqlStr, s.DorisNodes); err != nil {
			return err
		}
	}

	if len(s.DorisUser) != 0 {
		SqlStr := "update t_variable set variable_value = ? where variable_code  in" +
			" ('doris.username')"
		if _, err := db.Exec(SqlStr, s.DorisUser); err != nil {
			return err
		}
	}

	if len(s.DorisPass) != 0 {
		SqlStr := "update t_variable set variable_value = ? where variable_code  in" +
			" ('doris.password')"
		if _, err := db.Exec(SqlStr, s.DorisPass); err != nil {
			return err
		}
	}

	if len(s.RedisHost) != 0 {
		SqlStr := "update t_variable set variable_value = ? where variable_code  in" +
			" ('basic.redis.host','call.redis.host','chat.redis.host','icall.redis.host','robot.redis.host')"
		if _, err := db.Exec(SqlStr, s.RedisHost); err != nil {
			return err
		}
	}

	if len(s.RedisPort) != 0 {
		SqlStr := "update t_variable set variable_value = ? where variable_code  in" +
			" ('basic.redis.port','call.redis.port','chat.redis.port','icall.redis.port','robot.redis.port')"
		if _, err := db.Exec(SqlStr, s.RedisPort); err != nil {
			return err
		}
	}

	if len(s.RedisPass) != 0 {
		SqlStr := "update t_variable set variable_value = ? where variable_code  in" +
			" ('basic.redis.password','call.redis.password','chat.redis.password','icall.redis.password','robot.redis.password')"
		if _, err := db.Exec(SqlStr, s.RedisPass); err != nil {
			return err
		}
	}

	if len(s.ElasticUrl) != 0 {
		SqlStr := "update t_variable set variable_value = ? where variable_code  in" +
			" ('es6.hosts','es7.hosts')"
		if _, err := db.Exec(SqlStr, s.ElasticUrl); err != nil {
			return err
		}
	}

	if len(s.ElasticPort) != 0 {
		SqlStr := "update t_variable set variable_value = ? where variable_code  in" +
			" ('es6.port','es7.port')"
		if _, err := db.Exec(SqlStr, s.ElasticPort); err != nil {
			return err
		}
	}

	if len(s.HdfsUri) != 0 {
		SqlStr := "update t_variable set variable_value = ? where variable_code  in" +
			" ('hdfs.uri')"
		if _, err := db.Exec(SqlStr, s.HdfsUri); err != nil {
			return err
		}
	}

	if len(s.KafkaAddrs) != 0 {
		SqlStr := "update t_variable set variable_value = ? where variable_code  in" +
			" ('call-bootstrapServers','common-bootstrapServers')"
		if _, err := db.Exec(SqlStr, s.KafkaAddrs); err != nil {
			return err
		}
	}

	return nil
}
