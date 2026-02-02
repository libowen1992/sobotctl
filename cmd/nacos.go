package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"path"
	"sobotctl/global"
	"sobotctl/internal/nacosOps"
)

func NewNacosCmd() *cobra.Command {
	var NacosCmd = &cobra.Command{
		Use:   "nacos",
		Short: "nacos管理工具",
	}
	NacosCmd.AddCommand(newNacosConfigCmd())
	return NacosCmd
}

func newNacosConfigCmd() *cobra.Command {
	var NacosConfigCmd = &cobra.Command{
		Use:   "config",
		Short: "nacos配置管理",
	}
	NacosConfigCmd.AddCommand(newNacosConfigListCmd())
	NacosConfigCmd.AddCommand(newNacosConfigDetailCmd())
	NacosConfigCmd.AddCommand(newNacosConfigUpdateCmd())
	NacosConfigCmd.AddCommand(newNacosConfigImportCmd())
	return NacosConfigCmd
}

func newNacosConfigListCmd() *cobra.Command {
	var filter string
	action := "list"
	desc := "nacos配置列表"
	var Cmd = &cobra.Command{
		Use:   action,
		Short: desc,
		Run: func(cmd *cobra.Command, args []string) {
			if err := nacosOps.NewNacosOps().ConfigList(filter); err != nil {
				global.Logger.Error(err)
			}
		},
	}
	Cmd.Flags().StringVarP(&filter, "filter", "f", "", "模糊查询")
	return Cmd
}

func newNacosConfigDetailCmd() *cobra.Command {
	var config string
	action := "detail"
	desc := "nacos配置详情"
	var Cmd = &cobra.Command{
		Use:   action,
		Short: desc,
		Run: func(cmd *cobra.Command, args []string) {
			if err := nacosOps.NewNacosOps().ConfigDetail(config); err != nil {
				global.Logger.Error(err)
			}
		},
	}
	Cmd.Flags().StringVarP(&config, "config", "c", "", "配置文件名称")
	Cmd.MarkFlagRequired("config")

	return Cmd
}

func newNacosConfigUpdateCmd() *cobra.Command {
	var file string
	action := "update"
	desc := "nacos更新配置"
	var Cmd = &cobra.Command{
		Use:   action,
		Short: desc,
		Args: func(cmd *cobra.Command, args []string) error {
			if !path.IsAbs(file) {
				return errors.New("文件必须为绝对路径")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := nacosOps.NewNacosOps().ConfigUpdate(file); err != nil {
				global.Logger.Error(err)
			}
		},
	}

	Cmd.Flags().StringVarP(&file, "file", "f", "", "配置文件内容")
	Cmd.MarkFlagRequired("file")

	return Cmd
}

func newNacosConfigImportCmd() *cobra.Command {
	var argus nacosOps.NacosImportArgus
	action := "import"
	desc := "nacos导入配置，注意是.zip的压缩包"
	var Cmd = &cobra.Command{
		Use:   action,
		Short: desc,
		Args: func(cmd *cobra.Command, args []string) error {
			if !path.IsAbs(argus.FIle) {
				return errors.New("文件必须为绝对路径")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := nacosOps.NewNacosOps().ConfigImport(argus); err != nil {
				global.Logger.Error(err)
				os.Exit(1)
			}
		},
	}
	Cmd.Flags().StringVarP(&argus.NameSpace, "namespace", "n", "qa", "命名空间，默认qa")
	Cmd.Flags().StringVarP(&argus.FIle, "file", "f", "", "压缩包路径，注意全路径")
	Cmd.MarkFlagRequired("file")
	return Cmd
}
