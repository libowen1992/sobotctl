package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func K8sManager() *cobra.Command {
	action := "k8s"
	desc := "k8s管理"
	var Cmd = &cobra.Command{
		Use:   action,
		Short: desc,
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
	Cmd.AddCommand(K8sPodManager())
	return Cmd
}

func K8sPodManager() *cobra.Command {
	action := "pod"
	desc := "pod管理"
	var Cmd = &cobra.Command{
		Use:   action,
		Short: desc,
	}
	Cmd.AddCommand(K8sPodList())
	return Cmd
}

func K8sPodList() *cobra.Command {
	var kubeconfig string
	action := "check"
	desc := "pod检查"
	var Cmd = &cobra.Command{
		Use:   action,
		Short: desc,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	Cmd.Flags().StringVarP(&kubeconfig, "kubeconfig", "c", filepath.Join(homedir.HomeDir(), ".kube", "config"), "kubeconfig绝对路径")
	return Cmd
}
