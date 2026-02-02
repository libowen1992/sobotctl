package streampark

import (
	"database/sql"
	"errors"
	"sobotctl/global"
	"sobotctl/pkg/convert"
	"sobotctl/pkg/tableRender"
	"strings"
)

func (a *AppOps) List(filter string) error {
	db, err := NewMySQL(global.StreamparkS)
	if err != nil {
		return err
	}
	defer db.Close()

	dbData := make([]*TFlinkApp, 0)
	sqlStr := "select id,job_name,jar,main_class,app_id,state,build,start_time,duration,args from t_flink_app"
	if err := db.Select(&dbData, sqlStr); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			global.Logger.Warn("empty")
			return nil
		}
		global.Logger.Error("查询有误", err.Error())
		return nil
	}

	headers := []string{"id", "job_name", "state", "app_id", "jar", "main_class", "start_time"}
	renderData := make([][]string, 0)

	for _, v := range dbData {
		if strings.Contains(v.JobName, filter) {
			var t []string
			t = append(t, convert.Int64ToStr(v.Id),
				v.JobName, v.StateHuman(), v.AppId, v.Jar,
				v.MainClass, convert.TimeToDateTimeFormat(v.StartTime))
			renderData = append(renderData, t)
		}
	}
	tableRender.Render(headers, renderData)

	return nil
}

func (a *AppOps) ListCheck() (apps []*TFlinkApp, err error) {
	db, err := NewMySQL(global.StreamparkS)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	dbData := make([]*TFlinkApp, 0)
	sqlStr := "select id,job_name,jar,main_class,app_id,state,build,start_time,duration,args from t_flink_app"
	if err := db.Select(&dbData, sqlStr); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("app列表为空")
		}
		return nil, err
	}

	return dbData, nil
}
