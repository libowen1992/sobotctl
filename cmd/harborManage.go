package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	harborOps "sobotctl/internal/harborOps"
	"sobotctl/pkg/convert"
	"sobotctl/pkg/tableRender"
)

func NewHarborManager() *cobra.Command {
	action := "harbor"
	desc := "harbor工具"
	var Cmd = &cobra.Command{
		Use:   action,
		Short: desc,
	}
	Cmd.AddCommand(NewHarborProjectCmd())
	return Cmd
}

func NewHarborProjectCmd() *cobra.Command {
	action := "project"
	desc := "project管理"
	var Cmd = &cobra.Command{
		Use:   action,
		Short: desc,
	}
	Cmd.AddCommand(NewHarborProjectListCmd())
	Cmd.AddCommand(NewHarborProjectCreateCmd())
	return Cmd
}

func NewHarborProjectListCmd() *cobra.Command {
	action := "list"
	desc := "project列表"
	var Cmd = &cobra.Command{
		Use:   action,
		Short: desc,
		Run: func(cmd *cobra.Command, args []string) {
			harbor, err := harborOps.NewHarborOps()
			if err != nil {
				panic(err)
			}
			data, err := harbor.List()
			if err != nil {
				panic(err)
			}
			renderData := make([][]string, 0)
			headers := []string{"项目名称", "访问级别", "镜像仓库数", "创建时间"}
			for _, item := range data {
				s := make([]string, 0)
				access := "公开"
				if item.Public == "false" {
					access = "私有"
				}
				s = append(s, item.Name, access, convert.Int64ToStr(item.RepoCount), item.CreationTime.String())
				renderData = append(renderData, s)
			}
			tableRender.Render(headers, renderData)

		},
	}
	return Cmd
}

func NewHarborProjectCreateCmd() *cobra.Command {
	action := "create"
	desc := "创建project"
	var Cmd = &cobra.Command{
		Use:   action,
		Short: desc,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.Errorf("请至少提供一个项目名称")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			harbor, err := harborOps.NewHarborOps()
			if err != nil {
				panic(err)
			}
			err = harbor.Create(args)
			if err != nil {
				panic(err)
			}
		},
	}
	return Cmd
}
