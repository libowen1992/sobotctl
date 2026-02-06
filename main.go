/*
Copyright © 2023 wangys@sobot.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"github.com/pkg/errors"
	"sobotctl/cmd"
	"sobotctl/global"
	"sobotctl/pkg/logger"
	"sobotctl/setting"
)

func main() {
	// init logger  日志输出格式
	global.Logger = logger.NewStdoutConsole()

	// init setting
	if err := SetupSetting(); err != nil {
		global.Logger.Fatal(err)
	}
	cmd.Execute()   // Cobra 命令行框架
}

//main() 执行流程
//↓
//rootCmd.Execute()
//↓
//├── 解析命令行参数
//├── 匹配命令/子命令
//├── 执行对应的 Run 函数
//└── 返回错误（如果有）

//定义函数---把配置文件读取到结构体上
func SetupSetting() error {
	st := setting.New()  //返回v
	if err := st.Init("./config.yml"); err != nil {
		return errors.Wrap(err, "读取配置文件失败")
	}
	if err := st.SetSection("hosts", &global.HostSetting); err != nil {   //通过键值对的方式进行数据传参
		return err
	}
	if err := st.SetSection("redis", &global.RedisSetting); err != nil {
		return err
	}
	if err := st.SetSection("mysql", &global.MySQLSetting); err != nil {
		return err
	}
	if err := st.SetSection("apps", &global.AppSetting); err != nil {
		return err
	}
	if err := st.SetSection("kafka", &global.KafkaSetting); err != nil {
		return err
	}
	if err := st.SetSection("streampark", &global.StreamparkS); err != nil {
		return err
	}
	if err := st.SetSection("doris", &global.DorisS); err != nil {
		return err
	}
	if err := st.SetSection("hadoop", &global.HadoopS); err != nil {
		return err
	}
	if err := st.SetSection("elastic", &global.ElasticS); err != nil {
		return err
	}
	if err := st.SetSection("zookeeper", &global.ZookeeperS); err != nil {
		return err
	}
	if err := st.SetSection("nacos", &global.NacosS); err != nil {
		return err
	}
	if err := st.SetSection("k8s", &global.K8sS); err != nil {
		return err
	}
	if err := st.SetSection("harbor", &global.HarborS); err != nil {
		return err
	}
	if err := st.SetSection("TCPPortCheck", &global.TCPPortCheckS); err != nil {
		return err
	}
	return nil
}
