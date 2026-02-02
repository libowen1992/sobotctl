package cmd

import (
	"github.com/spf13/cobra"
	"sobotctl/global"
	"sobotctl/internal/hostManage"
	"sobotctl/pkg/convert"
	"sobotctl/pkg/tableRender"
)

func NewHostManage() *cobra.Command {
	action := "host"
	desc := "主机管理"
	var Cmd = &cobra.Command{
		Use:   action,
		Short: desc,
	}
	Cmd.AddCommand(NewHostCheck())
	Cmd.AddCommand(NewHostTerminal())
	return Cmd
}

func NewHostCheck() *cobra.Command {
	action := "check"
	desc := "主机状态"
	var Cmd = &cobra.Command{
		Use:   action,
		Short: desc,
		Run: func(cmd *cobra.Command, args []string) {
			renderData := make([][]string, 0)
			headers := []string{"主机名", "ip", "系统", "负载", "cpu核心", "cpu使用率", "内存总量", "内存使用率", "磁盘使用情况"}
			data, err := hostManage.NewHostOps().Check()
			if err != nil {
				global.Logger.Error(err)
				return
			}
			for _, item := range data {
				s := make([]string, 0)
				s = append(s, item.HostName, item.IP, item.OS,
					convert.Float64ToPercentString(item.Load),
					convert.IntToStr(item.CPUCores),
					convert.Float64ToPercentString(item.CPUPercent),
					item.MemTotal,
					convert.Float64ToPercentString(item.MemPercent),
					//item.Disk,
				)
				renderData = append(renderData, s)
			}
			tableRender.Render(headers, renderData)
		},
	}
	return Cmd
}

func NewHostTerminal() *cobra.Command {
	action := "terminal"
	desc := "主机shell"
	var Cmd = &cobra.Command{
		Use:   action,
		Short: desc,
	}
	return Cmd
}
