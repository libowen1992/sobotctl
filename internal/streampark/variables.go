package streampark

import (
	"database/sql"
	"errors"
	"sobotctl/global"
	"sobotctl/pkg/convert"
	"sobotctl/pkg/tableRender"
	"strings"
	"time"
)

type TVariable struct {
	Id            int64     `db:"id" json:"id"`
	VariableCode  string    `db:"variable_code" json:"variable_code"`   //Variable code is used for parameter names passed to the program or as placeholders
	VariableValue string    `db:"variable_value" json:"variable_value"` //The specific value corresponding to the variable
	CreateTime    time.Time `db:"create_time" json:"create_time"`       //create time
	ModifyTime    time.Time `db:"modify_time" json:"modify_time"`       //modify time
}

type TVariables []*TVariable

type variableOps struct {
}

func NewVariableOps() *variableOps {
	return &variableOps{}
}

func (v *variableOps) List(name string) error {
	db, err := NewMySQL(global.StreamparkS)
	if err != nil {
		return err
	}
	defer db.Close()

	dbData := make([]*TVariable, 0)
	sqlStr := "select id,variable_code,variable_value,create_time,modify_time from t_variable"
	if err := db.Select(&dbData, sqlStr); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			global.Logger.Warn("没有变量")
			return nil
		}
		global.Logger.Error("查询有误", err.Error())
		return nil
	}

	headers := []string{"id", "变量名", "变量值", "最近修改时间"}
	renderData := make([][]string, 0)

	for _, v := range dbData {
		if strings.Contains(v.VariableCode, name) {
			var t []string
			t = append(t, convert.Int64ToStr(v.Id), strings.TrimSpace(v.VariableCode), strings.TrimSpace(v.VariableValue), convert.TimeToDateTimeFormat(v.ModifyTime))
			renderData = append(renderData, t)
		}
	}
	tableRender.Render(headers, renderData)

	return nil
}

//func (v *variableOps) List(name string) error {
//	client, err := NewStreamClient(global.StreamparkS.User,
//		global.StreamparkS.Password, global.StreamparkS.Adder)
//
//	if err != nil {
//		return err
//	}
//
//	data, err := client.variableList()
//	if err != nil {
//		return err
//	}
//
//	headers := []string{"变量名", "变量值", "最近修改时间"}
//	renderData := make([][]string, 0)
//
//	for _, v := range data.Data.Records {
//		if strings.Contains(v.VariableCode, name) {
//			var t []string
//			t = append(t, strings.TrimSpace(v.VariableCode), strings.TrimSpace(v.VariableValue), v.ModifyTime)
//			renderData = append(renderData, t)
//		}
//	}
//	tableRender.Render(headers, renderData)
//	return nil
//}
