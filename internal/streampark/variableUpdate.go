package streampark

import "sobotctl/global"

func (v *variableOps) Update(name, value string) error {
	db, err := NewMySQL(global.StreamparkS)
	if err != nil {
		return err
	}
	defer db.Close()

	SqlStr := "update t_variable set variable_value = ? where variable_code = ?"
	if _, err := db.Exec(SqlStr, value, name); err != nil {
		return err
	}

	return nil
}
