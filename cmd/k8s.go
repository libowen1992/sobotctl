package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func K8sManager() *cobra.Command {
	action := "k8s"
	desc := "k8s管理"
	var Cmd = &cobra.Command{
		Use:   action,
		Short: desc,
		Run:   func(cmd *cobra.Command, args []string) {},
	}
	Cmd.AddCommand(K8sNamespaceManager())
	Cmd.AddCommand(K8sPodManager())
	return Cmd
}

func K8sNamespaceManager() *cobra.Command {
	action := "namespace"
	desc := "命名空间管理"
	var Cmd = &cobra.Command{
		Use:   action,
		Short: desc,
	}
	Cmd.AddCommand(K8sNamespaceList())
	return Cmd
}

func K8sNamespaceList() *cobra.Command {
	var kubeconfig string
	action := "list"
	desc := "列出命名空间"
	var Cmd = &cobra.Command{
		Use:   action,
		Short: desc,
		Run: func(cmd *cobra.Command, args []string) {
			cfgPath := kubeconfig
			if cfgPath == "" {
				cfgPath = filepath.Join(homedir.HomeDir(), ".kube", "config")
			}
			config, err := clientcmd.BuildConfigFromFlags("", cfgPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "构建 kubeconfig 失败: %v\n", err)
				return
			}
			clientset, err := kubernetes.NewForConfig(config)
			if err != nil {
				fmt.Fprintf(os.Stderr, "创建 Kubernetes 客户端失败: %v\n", err)
				return
			}
			ctx := context.Background()
			nsList, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
			if err != nil {
				fmt.Fprintf(os.Stderr, "获取命名空间列表失败: %v\n", err)
				return
			}
			for _, ns := range nsList.Items {
				fmt.Println(ns.Name)
			}
		},
	}
	Cmd.Flags().StringVarP(&kubeconfig, "kubeconfig", "c", filepath.Join(homedir.HomeDir(), ".kube", "config"), "kubeconfig绝对路径")
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
	var namespace string
	action := "check"
	desc := "pod检查"
	var Cmd = &cobra.Command{
		Use:   action,
		Short: desc,
		Run: func(cmd *cobra.Command, args []string) {
			cfgPath := kubeconfig
			if cfgPath == "" {
				cfgPath = filepath.Join(homedir.HomeDir(), ".kube", "config")
			}
			config, err := clientcmd.BuildConfigFromFlags("", cfgPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "构建 kubeconfig 失败: %v\n", err)
				return
			}
			clientset, err := kubernetes.NewForConfig(config)
			if err != nil {
				fmt.Fprintf(os.Stderr, "创建 Kubernetes 客户端失败: %v\n", err)
				return
			}
			ctx := context.Background()
			ns := namespace
			// 如果没有指定命名空间，使用 "" 表示所有命名空间
			podList, err := clientset.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{})
			if err != nil {
				fmt.Fprintf(os.Stderr, "获取 Pod 列表失败: %v\n", err)
				return
			}
			for _, p := range podList.Items {
				fmt.Printf("%s\t%s\t%s\t%s\n", p.Namespace, p.Name, p.Status.Phase, p.Spec.NodeName)
			}
		},
	}
	Cmd.Flags().StringVarP(&kubeconfig, "kubeconfig", "c", filepath.Join(homedir.HomeDir(), ".kube", "config"), "kubeconfig绝对路径")
	Cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "命名空间，留空表示所有命名空间")
	return Cmd
}
