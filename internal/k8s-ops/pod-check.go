package k8s_ops

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sobotctl/global"
	"sobotctl/pkg/tableRender"
)

func (k *K8sOps) Check() {
	podList, err := k.ClientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	headers := []string{
		"容器组",
		"状态",
		"ip",
		"宿主机",
		"宿主机ip",
		"命名空间",
		"错误信息",
		"创建时间",
	}
	data := make([][]string, 0, 0)
	for _, pod := range podList.Items {
		item := make([]string, 0, 0)
		podErrMsg := ""
		switch pod.Status.Phase {
		case v1.PodFailed, v1.PodReasonUnschedulable, v1.PodPending:
			podErrMsg = fmt.Sprintf("%s: %s", pod.Status.Reason, pod.Status.Message)
			item = append(item,
				pod.GetName(),
				string(pod.Status.Phase),
				pod.Status.PodIP,
				pod.Spec.NodeName,
				pod.Status.HostIP,
				pod.GetNamespace(),
				podErrMsg,
				pod.GetCreationTimestamp().String())
			data = append(data, item)
		default:
			continue
		}

	}
	if len(data) > 0 {
		global.Logger.Error("k8s pod有异常")
		tableRender.Render(headers, data)
	} else {
		global.Logger.Error("k8s pod 运行正常")
	}
}
