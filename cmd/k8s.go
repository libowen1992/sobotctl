package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"
	"time"

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
	Cmd.AddCommand(K8sServiceManager())
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

func K8sServiceManager() *cobra.Command {
	action := "svc"
	desc := "service管理"
	var Cmd = &cobra.Command{
		Use:   action,
		Short: desc,
	}
	Cmd.AddCommand(K8sServiceList())
	return Cmd
}

func K8sServiceList() *cobra.Command {
	var kubeconfig string
	var namespace string
	action := "check"
	desc := "service检查"
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
			svcList, err := clientset.CoreV1().Services(ns).List(ctx, metav1.ListOptions{})
			if err != nil {
				fmt.Fprintf(os.Stderr, "获取 Service 列表失败: %v\n", err)
				return
			}

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
			fmt.Fprintln(w, "NAME\tTYPE\tCLUSTER-IP\tEXTERNAL-IP\tPORT(S)\tAGE")
			now := time.Now()
			for _, svc := range svcList.Items {
				// 拼接端口信息
				portStr := ""
				for i, port := range svc.Spec.Ports {
					if i > 0 {
						portStr += ","
					}
					portStr += fmt.Sprintf("%d/%s", port.Port, port.Protocol)
					if port.NodePort > 0 {
						portStr += fmt.Sprintf(":%d", port.NodePort)
					}
				}
				// 计算 age
				age := now.Sub(svc.CreationTimestamp.Time)
				ageStr := formatDuration(age)
				// 获取 external IP
				externalIP := "<none>"
				if len(svc.Status.LoadBalancer.Ingress) > 0 {
					externalIP = svc.Status.LoadBalancer.Ingress[0].IP
					if externalIP == "" {
						externalIP = svc.Status.LoadBalancer.Ingress[0].Hostname
					}
				}
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", svc.Name, svc.Spec.Type, svc.Spec.ClusterIP, externalIP, portStr, ageStr)
			}
			w.Flush()
		},
	}
	Cmd.Flags().StringVarP(&kubeconfig, "kubeconfig", "c", filepath.Join(homedir.HomeDir(), ".kube", "config"), "kubeconfig绝对路径")
	Cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "命名空间，留空表示所有命名空间")
	return Cmd
}

// formatDuration 将 time.Duration 格式化为人类可读的字符串（如 109d, 33d, 3h2m 等）
func formatDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	if d.Hours() >= 24 {
		days := int(d.Hours()) / 24
		return fmt.Sprintf("%dd", days)
	}
	if d.Hours() >= 1 {
		return fmt.Sprintf("%dh%dm", int(d.Hours()), int(d.Minutes())%60)
	}
	return fmt.Sprintf("%dm", int(d.Minutes()))
}
