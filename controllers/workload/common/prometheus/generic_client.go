package prometheus

import (
	poclientset "github.com/prometheus-operator/prometheus-operator/pkg/client/versioned"
	kubeclientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type GenericClientset struct {
	KubeClient kubeclientset.Interface
	PoClient   poclientset.Interface
}

// NewForConfig creates a new Clientset for the given config.
func newForConfig(c *rest.Config) (*GenericClientset, error) {
	kubeClient, err := kubeclientset.NewForConfig(c)
	if err != nil {
		return nil, err
	}
	kruiseClient, err := poclientset.NewForConfig(c)
	if err != nil {
		return nil, err
	}
	return &GenericClientset{
		KubeClient: kubeClient,
		PoClient:   kruiseClient,
	}, nil
}
