package k8s_ops

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sOps struct {
	Kubeconfig *string
	ClientSet  *kubernetes.Clientset
}

func NewK8sOps(kubeconfig *string) (*K8sOps, error) {
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return nil, err
	}
	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &K8sOps{
		Kubeconfig: kubeconfig,
		ClientSet:  clientset,
	}, nil
}
