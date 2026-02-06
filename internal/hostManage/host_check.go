package hostManage

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"sobotctl/global"
	"strings"
)

type DiskUsageStat struct {
	Path              string  `json:"path"`
	Fstype            string  `json:"fstype"`
	Total             uint64  `json:"total"`
	Free              uint64  `json:"free"`
	Used              uint64  `json:"used"`
	UsedPercent       float64 `json:"usedPercent"`
	InodesTotal       uint64  `json:"inodesTotal"`
	InodesUsed        uint64  `json:"inodesUsed"`
	InodesFree        uint64  `json:"inodesFree"`
	InodesUsedPercent float64 `json:"inodesUsedPercent"`
}

type HostInfo struct {
	HostName   string           `json:"host_name"`
	IP         string           `json:"ip"`
	OS         string           `json:"os"`
	Load       float64          `json:"load"`
	CPUCores   int              `json:"cpu_cores"`
	CPUPercent float64          `json:"cpu_percent"`
	MemTotal   string           `json:"mem_total"`
	MemPercent float64          `json:"mem_percent"`
	NTPState   bool             `json:"ntp_state"`
	Disk       []*DiskUsageStat `json:"disk"`
}

func (ho *HostOps) Check() (data []HostInfo, err error) {
	if len(global.HostSetting.IPS) == 0 {
		return nil, nil
	}

	data = make([]HostInfo, 0)  //创建一个空切片
	for _, item := range global.HostSetting.IPS {   //ip集群
		h := &Host{
			User:    global.HostSetting.User,
			SshType: global.HostSetting.SshType,
			SshPass: global.HostSetting.SshPass,
			SshKey:  global.HostSetting.SshKey,
			Port:    global.HostSetting.Port,
			IP:      item,
		}
		global.Logger.Debug(fmt.Sprintf("获取主机信息: %s", h.IP)) //日志输出信息
		info, err := ho.GetHostInfo(h)
		if err != nil {
			global.Logger.Errorf("获取主机信息发生错误, ip:%s, err: %v", h.IP, err)
		}
		data = append(data, info)  //切片追加
	}

	return
}

func (ho *HostOps) GetHostInfo(h *Host) (info HostInfo, err error) {
	if h.Client == nil {
		if err := h.NewSSHClient(); err != nil {
			return HostInfo{}, err
		}
	}
	defer h.CloseClient()

	diskPoints := strings.Join(global.HostSetting.DiskCheckPoints, ",")
	command := fmt.Sprintf("%s check -d %s -n %s", global.HostSetting.HostctlPath, diskPoints, global.HostSetting.NtpServer)
	stdout, stderr, err := h.RunSingleCommand(command)
	if err != nil {
		return HostInfo{}, err
	}
	if len(stderr) != 0 {
		err = errors.New(fmt.Sprintf("检测时间同步失败 %s", string(stderr)))
	}

	info = HostInfo{}
	if err := json.Unmarshal(stdout, &info); err != nil {
		return HostInfo{}, err
	}
	info.IP = h.IP

	return info, err
}
