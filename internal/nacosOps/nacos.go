package nacosOps

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"sobotctl/global"
	"sobotctl/pkg/mysql"
	"sobotctl/pkg/tableRender"
	"strings"
)

type NacosOps struct {
}

func NewNacosOps() *NacosOps {
	return &NacosOps{}
}

type NacosConfigInfo struct {
	Id       int64  `db:"id" json:"id" form:"id"`                //id
	DataId   string `db:"data_id" json:"data_id" form:"data_id"` //data_id
	GroupId  string `db:"group_id" json:"group_id" form:"group_id"`
	Content  string `db:"content" json:"content" form:"content"`       //content
	TenantId string `db:"tenant_id" json:"tenant_id" form:"tenant_id"` //租户字段
	Type     string `db:"type" json:"type" form:"type"`
}
//链接nacos数据库
func (no *NacosOps) ConfigList(filter string) error {
	db, err := mysql.NewSqx(global.NacosS.Dbip,
		global.NacosS.Dbuser,
		global.NacosS.Dbpass,
		global.NacosS.DbName,
		"utf8", global.NacosS.Dbport, 3, 1)
	if err != nil {
		return errors.WithStack(err)
	}
	defer db.Close()
///查询nacos表中数据
	dbData := make([]*NacosConfigInfo, 0) //定义空切片
	sqlStr := "select id,data_id,group_id,content,tenant_id,type from config_info"
	if err := db.Select(&dbData, sqlStr); err != nil {   //查询的数据写入空切片中
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		global.Logger.Error("查询有误", err.Error())
		return nil
	}

	headers := []string{"配置", "组", "租户", "类型"} //请求头书写方式
	renderData := make([][]string, 0)  //二维字符串切片

	for _, v := range dbData {   //for循环dbData切片
		if strings.Contains(v.DataId, filter) {
			var t []string  //切片
			t = append(t, v.DataId, v.GroupId, v.TenantId, v.Type)  //切片追加
			renderData = append(renderData, t)
		}
	}
	tableRender.Render(headers, renderData)  //输出标准设置

	return nil
}

func (no *NacosOps) ConfigDetail(config string) error {
	db, err := mysql.NewSqx(global.NacosS.Dbip,
		global.NacosS.Dbuser,
		global.NacosS.Dbpass,
		global.NacosS.DbName,
		"utf8", global.NacosS.Dbport, 3, 1)
	if err != nil {
		return errors.WithStack(err)
	}
	defer db.Close()  //执行关闭数据库的动作

	data := NacosConfigInfo{}  //空结构体
	sqlStr := "select id,data_id,group_id,content,tenant_id,type from config_info where data_id = ?"
	if err := db.Get(&data, sqlStr, config); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			global.Logger.Warn("没查到相关配置")
			return nil
		}
		return errors.WithStack(err)
	}

	fmt.Println(data.Content)

	return nil
}

func (no *NacosOps) ConfigUpdate(file string) error {
	dataId := path.Base(file)
	content, err := os.ReadFile(file)
	if err != nil {
		return errors.WithStack(err)
	}

	db, err := mysql.NewSqx(global.NacosS.Dbip,
		global.NacosS.Dbuser,
		global.NacosS.Dbpass,
		global.NacosS.DbName,
		"utf8", global.NacosS.Dbport, 3, 1)
	if err != nil {
		return errors.WithStack(err)
	}
	defer db.Close()

	sqlStr := "update config_info set content = ? where data_id = ?"
	_, err = db.Exec(sqlStr, string(content), dataId)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

type NacosImportArgus struct {
	FIle      string
	NameSpace string
}

type NacosImportResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		SuccCount int `json:"succCount"`
		SkipCount int `json:"skipCount"`
	} `json:"data"`
}

func (no *NacosOps) ConfigImport(argus NacosImportArgus) error {
	nacosImportUrl := fmt.Sprintf("%s/nacos/v1/cs/configs?import=true&namespace=%s&tenant=%s",
		global.NacosS.Host, argus.NameSpace, argus.NameSpace)

	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)

	// file
	fileName := path.Base(argus.FIle)
	fileWriter, err := bodyWriter.CreateFormFile("file", fileName)
	if err != nil {
		return errors.Cause(err)
	}
	file, err := os.Open(argus.FIle)
	if err != nil {
		return errors.Cause(err)
	}
	defer file.Close()
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return errors.Cause(err)
	}

	// other form data
	bodyWriter.WriteField("policy", "OVERWRITE")
	if err != nil {
		return errors.Cause(err)
	}

	// request
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	resp, err := http.Post(nacosImportUrl, contentType, bodyBuffer)
	if err != nil {
		return errors.Cause(err)
	}
	defer resp.Body.Close()

	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Cause(err)
	}

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("请求失败 msg: %s", resp_body)
	}

	return nil
}

func (no *NacosOps) Status() error {

	return nil
}
