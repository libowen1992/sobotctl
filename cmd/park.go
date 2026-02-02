package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path"
	"sobotctl/global"
	"strings"
)

const (
	toolYamlPath = "ops/tool/"
)

var (
	tools = []string{"kuboard", "apisix", "kong", "minio", "nacos", "powerjob", "streampark"}
)

func NewParkCmd() *cobra.Command {
	action := "park"
	desc := "切换工具域名用途"
	var Cmd = &cobra.Command{
		Use:   action,
		Short: desc,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 || !toolExist(args[0]) {
				fmt.Println(toolHelpStr())
				return
			}
			yamlFIle := fmt.Sprintf("%s.yml", strings.TrimSpace(args[0]))
			yamlFullPath := path.Join(global.K8sS.YamlDir, toolYamlPath, yamlFIle)
			commandStr := fmt.Sprintf("kubectl apply -f %s", yamlFullPath)
			command := exec.Command("bash", "-c", commandStr)
			command.Stdin = os.Stdin
			command.Stdout = os.Stdout
			command.Stderr = os.Stderr
			command.Run()
		},
	}
	return Cmd
}

func toolHelpStr() string {
	t := strings.Join(tools, "|")
	return fmt.Sprintf("Usage: sobotctl park |%v|", t)
}

func toolExist(tool string) bool {
	for _, t := range tools {
		if tool == t {
			return true
		}
	}
	return false
}
