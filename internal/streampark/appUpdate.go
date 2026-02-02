package streampark

import (
	"fmt"
	"github.com/pkg/errors"
	"path"
	"path/filepath"
	"sobotctl/global"
	"sobotctl/pkg/sobothdfs"
	"sobotctl/setting"
	"time"
)

func (a *AppOps) VerifyApp(appid string, jar string) (*setting.StreamParkApp, bool) {
	file := filepath.Base(jar)
	for _, app := range global.StreamparkS.Apps {
		if app.AppId == appid && app.AppJar == file {
			return &app, true
		}
	}
	return nil, false
}

func (a *AppOps) Update(appid, jar string) error {
	app, ok := a.VerifyApp(appid, jar)
	if !ok {
		a.List("")
		//global.Logger.Error(, zap.Any("apps", global.StreamparkS.Apps))
		return errors.New("请检查appid是否在配置文件里面,或者appid和jar包对象关系是否是对的")
	}

	client, err := NewStreamClient(global.StreamparkS.User,
		global.StreamparkS.Password, global.StreamparkS.Adder)

	if err != nil {
		return err
	}

	global.Logger.Infof("关闭服务：%s", appid)
	if err := client.stopApp(appid); err != nil {
		global.Logger.Error("关闭服务failed")
		return err
	}

	global.Logger.Info("备份jar包")
	hdfsClient, err := sobothdfs.NewHdfsClient(global.HadoopS.NameNode, global.HadoopS.User)
	if err != nil {
		return err
	}
	defer hdfsClient.Close()

	// 备份
	jarPathDir := path.Join(global.StreamparkS.JarHdfsBasePath, appid)
	jarBackupPathDir := path.Join(global.StreamparkS.JarHdfsBackupPath, appid)
	ok, _ = hdfsClient.FileExist(jarPathDir)
	if !ok {
		if err := hdfsClient.Client.MkdirAll(jarPathDir, 0755); err != nil {
			return err
		}
	}
	ok, _ = hdfsClient.FileExist(jarBackupPathDir)
	if !ok {
		if err := hdfsClient.Client.MkdirAll(jarBackupPathDir, 0755); err != nil {
			return err
		}
	}
	jarPath := path.Join(jarPathDir, app.AppJar)
	jarBackupPath := path.Join(global.StreamparkS.JarHdfsBackupPath, appid, fmt.Sprintf("%s_%d", app.AppJar, time.Now().Unix()))
	ok, _ = hdfsClient.FileExist(jarPath)
	if ok {
		if err := hdfsClient.Client.Rename(jarPath, jarBackupPath); err != nil {
			return err
		}
	}
	global.Logger.Info("上传jar包")
	if err := hdfsClient.Client.CopyToRemote(jar, jarPath); err != nil {
		return err
	}

	global.Logger.Info("启动服务")
	if err := client.startApp(appid); err != nil {
		global.Logger.Error("启动服务failed")
		return err
	}
	return nil
}

func (a *AppOps) Init(jarDir string) {
	for _, Sapp := range global.StreamparkS.Apps {
		if err := a.Update(Sapp.AppId, path.Join(jarDir, Sapp.AppJar)); err != nil {
			global.Logger.Error(err)
		}
	}
	a.List("")
}
