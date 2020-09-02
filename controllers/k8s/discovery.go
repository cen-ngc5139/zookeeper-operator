package k8s

import (
	"context"

	openapi "github.com/googleapis/gnostic/openapiv2"
	"k8s.io/client-go/discovery"
)

// WrapClient returns a Client that performs requests within DefaultTimeout.
func WrapDiscoveryClient(ctx context.Context, client discovery.DiscoveryClient) DisClient {
	return &ClusterDiscoveryClient{
		crClient: client,
		ctx:      ctx,
	}
}

// Client wraps a discovery client to use a
// default context with a timeout if no context is passed.
type DisClient interface {
	// WithContext returns a client configured to use the provided context on
	// subsequent requests, instead of one created from the preconfigured timeout.
	WithContext(ctx context.Context) DisClient

	OpenAPISchema() (*openapi.Document, error)
}

type ClusterDiscoveryClient struct {
	crClient discovery.DiscoveryClient
	ctx      context.Context
}

// WithContext returns a client configured to use the provided context on
// subsequent requests, instead of one created from the preconfigured timeout.
func (w *ClusterDiscoveryClient) WithContext(ctx context.Context) DisClient {
	w.ctx = ctx
	return w
}

func (w *ClusterDiscoveryClient) OpenAPISchema() (*openapi.Document, error) {
	return w.crClient.OpenAPISchema()
}
