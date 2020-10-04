package prometheus

import (
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

var (
	genericClient *GenericClientset
)

func NewRegistry(mgr manager.Manager) (*GenericClientset, error) {
	var err error
	genericClient, err = newForConfig(mgr.GetConfig())
	if err != nil {
		return nil, err
	}
	return genericClient, nil
}

func GetGenericClient() GenericClientset {
	return *genericClient
}
