package k8s

import (
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	k8sruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/resource"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewKubeClient() (client.Client, *rest.Config, error) {

	scheme := k8sruntime.NewScheme()
	_ = clientgoscheme.AddToScheme(scheme)

	kubeConfigFlags := genericclioptions.NewConfigFlags(true).WithDeprecatedPasswordFlag()
	f := NewFactory(kubeConfigFlags)
	restConf, err := f.ToRESTConfig()
	if err != nil {
		fmt.Println("get kubeconfig err", err)
		return nil, nil, err
	}
	client, err := client.New(restConf, client.Options{Scheme: scheme})
	if err != nil {
		fmt.Println("create client from kubeconfig err", err)
		return nil, nil, err
	}
	return client, restConf, nil
}

type Factory interface {
	genericclioptions.RESTClientGetter
	NewBuilder() *resource.Builder
}

type factoryImpl struct {
	clientGetter genericclioptions.RESTClientGetter
}

func (f *factoryImpl) ToRESTConfig() (*rest.Config, error) {
	return f.clientGetter.ToRESTConfig()
}

func (f *factoryImpl) ToRESTMapper() (meta.RESTMapper, error) {
	return f.clientGetter.ToRESTMapper()
}

func (f *factoryImpl) ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error) {
	return f.clientGetter.ToDiscoveryClient()
}

func (f *factoryImpl) ToRawKubeConfigLoader() clientcmd.ClientConfig {
	return f.clientGetter.ToRawKubeConfigLoader()
}

// NewBuilder returns a new resource builder for structured api objects.
func (f *factoryImpl) NewBuilder() *resource.Builder {
	return resource.NewBuilder(f.clientGetter)
}

func NewFactory(clientGetter genericclioptions.RESTClientGetter) Factory {
	if clientGetter == nil {
		panic("attempt to instantiate client_access_factory with nil clientGetter")
	}

	f := &factoryImpl{
		clientGetter: clientGetter,
	}

	return f
}
