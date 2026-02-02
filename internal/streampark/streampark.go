package streampark

import (
	"context"
	"fmt"
	"github.com/carlmjohnson/requests"
	"github.com/pkg/errors"
	"net/url"
)

const (
	PATH_SIGN           = "passport/signin"
	PATH_VARIABLE_LIST  = "variable/page"
	PATH_APP_STOP       = "flink/app/cancel"
	PATH_APP_START      = "flink/app/start"
	PATH_APP_savepoints = "flink/savepoint/history"
	TEAMID              = "100000"
)

type StreamParkClient struct {
	Token    *string
	TeamId   string
	User     string
	Password string
	Url      string
}

func NewStreamClient(user, pass, addr string) (*StreamParkClient, error) {
	var loginRes LoginResp
	req := make(url.Values)
	req["username"] = []string{user}
	req["password"] = []string{pass}
	req["loginType"] = []string{"PASSWORD"}
	reqUrl := fmt.Sprintf("%s/%s", addr, PATH_SIGN)

	ctx := context.Background()
	err := requests.
		URL(reqUrl).BodyForm(req).ToJSON(&loginRes).Fetch(ctx)
	if err != nil {
		return nil, err
	}

	return &StreamParkClient{
		Token:    &loginRes.Data.Token,
		User:     user,
		Password: pass,
		Url:      addr,
	}, nil
}

func (s *StreamParkClient) variableList() (data *VariableListResp, err error) {
	if s.Token == nil {
		return nil, errors.New("请先初始化client")
	}

	var resp VariableListResp
	req := make(url.Values)
	req["pageNum"] = []string{"1"}
	req["pageSize"] = []string{"100"}
	req["teamId"] = []string{"100000"}
	reqUrl := fmt.Sprintf("%s/%s", s.Url, PATH_VARIABLE_LIST)

	if err := requests.URL(reqUrl).BodyForm(req).Header("Authorization", *s.Token).ToJSON(&resp).Fetch(context.Background()); err != nil {
		return nil, errors.New(err.Error())
	}
	return &resp, nil
}

func (s *StreamParkClient) stopApp(appId string) error {
	if s.Token == nil {
		return errors.New("请先初始化client")
	}
	req := make(url.Values)
	req["id"] = []string{appId}
	req["savePointed"] = []string{"false"}
	req["drain"] = []string{"false"}
	req["teamId"] = []string{TEAMID}
	reqUrl := fmt.Sprintf("%s/%s", s.Url, PATH_APP_STOP)
	if err := requests.URL(reqUrl).BodyForm(req).Header("Authorization", *s.Token).Fetch(context.Background()); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (s *StreamParkClient) startApp(appId string) error {
	if s.Token == nil {
		return errors.New("请先初始化client")
	}

	savepointsResp, err := s.AppSavePointHistory(appId)
	if err != nil {
		return errors.Wrapf(err, "获取app: %s savepoint历史失败", appId)
	}

	req := make(url.Values)
	req["id"] = []string{appId}
	req["teamId"] = []string{TEAMID}
	req["allowNonRestored"] = []string{"false"}
	if len(savepointsResp.Data.Records) == 0 {
		req["savePointed"] = []string{"false"}
	} else {
		req["savePointed"] = []string{"true"}
		req["savePoint"] = []string{savepointsResp.Data.Records[0].Path}
	}
	reqUrl := fmt.Sprintf("%s/%s", s.Url, PATH_APP_START)

	if err := requests.URL(reqUrl).BodyForm(req).Header("Authorization", *s.Token).Fetch(context.Background()); err != nil {
		return errors.New(err.Error())
	}
	return nil
}

func (s *StreamParkClient) AppSavePointHistory(appId string) (data *AppSavePointResp, err error) {
	if s.Token == nil {
		return nil, errors.New("请先初始化client")
	}
	resp := AppSavePointResp{}
	req := make(url.Values)
	req["appId"] = []string{appId}
	req["pageNum"] = []string{"1"}
	req["pageSize"] = []string{"10"}
	req["teamId"] = []string{TEAMID}
	reqUrl := fmt.Sprintf("%s/%s", s.Url, PATH_APP_savepoints)

	if err := requests.URL(reqUrl).BodyForm(req).Header("Authorization", *s.Token).ToJSON(&resp).Fetch(context.Background()); err != nil {
		return nil, errors.New(err.Error())
	}
	return &resp, nil
}

type LoginResp struct {
	Data LoginUserInfo `json:"data"`
}

type LoginUserInfo struct {
	Token  string `json:"token"`
	Expire string `json:"expire"`
}

type VariableListResp struct {
	Data VariableListRespData `json:"data"`
}

type VariableListRespData struct {
	Total   string         `json:"total"`
	Records []VariableInfo `json:"records"`
}

type VariableInfo struct {
	VariableCode  string `json:"variableCode"`
	VariableValue string `json:"variableValue"`
	ModifyTime    string `json:"modifyTime"`
}

type AppSavePointResp struct {
	Code string `json:"code"`
	Data struct {
		Records []struct {
			ID          string `json:"id"`
			AppID       string `json:"appId"`
			ChkID       string `json:"chkId"`
			Latest      bool   `json:"latest"`
			Type        int    `json:"type"`
			Path        string `json:"path"`
			TriggerTime string `json:"triggerTime"`
			CreateTime  string `json:"createTime"`
		} `json:"records"`
		Total   string `json:"total"`
		Size    string `json:"size"`
		Current string `json:"current"`
		Orders  []struct {
			Column string `json:"column"`
			Asc    bool   `json:"asc"`
		} `json:"orders"`
		OptimizeCountSql bool        `json:"optimizeCountSql"`
		SearchCount      bool        `json:"searchCount"`
		MaxLimit         interface{} `json:"maxLimit"`
		CountID          interface{} `json:"countId"`
		Pages            string      `json:"pages"`
	} `json:"data"`
	Status string `json:"status"`
}
