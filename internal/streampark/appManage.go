package streampark

import "time"

type AppOps struct {
}

func NewAppOps() *AppOps {
	return &AppOps{}
}

type TFlinkApp struct {
	Id        int64     `db:"id"`
	JobName   string    `db:"job_name"`
	Jar       string    `db:"jar"`
	MainClass string    `db:"main_class"`
	Args      string    `db:"args"`
	AppId     string    `db:"app_id"`
	State     int       `db:"state"`
	Build     int8      `db:"build"`
	StartTime time.Time `db:"start_time"`
	Duration  int64     `db:"duration"`
}

func (t *TFlinkApp) StateHuman() string {
	switch t.State {
	case 5:
		return "RUNNING"
	case 7:
		return "FAILED"
	case 9:
		return "CANCELED"
	default:
		return "UNKNOWN"
	}
}

func UnixMilliToTime(milli int64) time.Time {
	return time.Unix(milli/1000, (milli%1000)*(1000*1000))
}

func (t *TFlinkApp) Uptime() string {
	//d := time.Unix(t.Duration/1000, 0)
	return ""
}
