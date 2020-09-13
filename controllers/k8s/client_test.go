package k8s_test

import (
	"context"
	"errors"

	"github.com/ghostbaby/zookeeper-operator/controllers/k8s"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ctxKey struct{}

var (
	userProvidedContextKey        = ctxKey{}
	errUsingUserProvidedContext   = errors.New("using user-provided context")
	errUsingDefaultTimeoutContext = errors.New("using default timeout context")
)

var (
	tests []struct {
		name string
		call func(c k8s.Client) error
	}
	ctx = context.Background()
)

var _ = Describe("Client", func() {
	BeforeEach(func() {
		tests = []struct {
			name string
			call func(c k8s.Client) error
		}{
			{
				name: "Get",
				call: func(c k8s.Client) error {
					return c.Get(types.NamespacedName{}, nil)
				},
			},
			{
				name: "List",
				call: func(c k8s.Client) error {
					return c.List(nil, nil)
				},
			},
			{
				name: "Create",
				call: func(c k8s.Client) error {
					return c.Create(nil)
				},
			},
			{
				name: "Update",
				call: func(c k8s.Client) error {
					return c.Update(nil)
				},
			},
			{
				name: "Patch",
				call: func(c k8s.Client) error {
					return c.Patch(nil, nil, nil)
				},
			},
		}
	})

	Describe("Wrapper k8s Client", func() {
		Context("Wrapper k8s Client", func() {
			It("should pass ", func() {
				for _, tt := range tests {
					// setup the Client with a timeout
					c := k8s.WrapClient(ctx, mockedClient{})

					// pass a custom context with the call
					ctx := context.WithValue(context.Background(), userProvidedContextKey, userProvidedContextKey)
					err := tt.call(c.WithContext(ctx))
					// make sure this custom context was used and not the timeout one
					Expect(err).To(Equal(errUsingUserProvidedContext))
				}
			})
		})

	})
})

// mockedClient's only purpose is to perform checks against the context
// passed in from the surrounding Client
type mockedClient struct{}

func (m mockedClient) checkCtx(ctx context.Context) error {
	if ctx == nil {
		return errors.New("using no context")
	}
	if ctx.Value(userProvidedContextKey) == userProvidedContextKey {
		return errUsingUserProvidedContext
	}
	// should be the init timeout context
	<-ctx.Done()
	return errUsingDefaultTimeoutContext
}

func (m mockedClient) Get(ctx context.Context, key client.ObjectKey, obj runtime.Object) error {
	return m.checkCtx(ctx)
}

func (m mockedClient) List(ctx context.Context, list runtime.Object, opts ...client.ListOption) error {
	return m.checkCtx(ctx)
}

func (m mockedClient) Create(ctx context.Context, obj runtime.Object, opts ...client.CreateOption) error {
	return m.checkCtx(ctx)
}

func (m mockedClient) Delete(ctx context.Context, obj runtime.Object, opts ...client.DeleteOption) error {
	return m.checkCtx(ctx)
}

func (m mockedClient) Update(ctx context.Context, obj runtime.Object, opts ...client.UpdateOption) error {
	return m.checkCtx(ctx)
}

func (m mockedClient) Patch(ctx context.Context, obj runtime.Object, patch client.Patch, opts ...client.PatchOption) error {
	return m.checkCtx(ctx)
}

func (m mockedClient) DeleteAllOf(ctx context.Context, obj runtime.Object, opts ...client.DeleteAllOfOption) error {
	return m.checkCtx(ctx)
}

func (m mockedClient) Status() client.StatusWriter {
	return mockedStatusWriter{c: m}
}

type mockedStatusWriter struct {
	c mockedClient
}

func (m mockedStatusWriter) Update(ctx context.Context, obj runtime.Object, opts ...client.UpdateOption) error {
	return m.c.checkCtx(ctx)
}

func (m mockedStatusWriter) Patch(ctx context.Context, obj runtime.Object, patch client.Patch, opts ...client.PatchOption) error {
	return m.c.checkCtx(ctx)
}
