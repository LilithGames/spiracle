package v1

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
	v1 "github.com/LilithGames/spiracle/api/v1"

	"github.com/LilithGames/spiracle/client/v1/scheme"
)

type SpiracleV1Interface interface {
	RESTClient() rest.Interface
	RoomIngresses(namespace string) RoomIngressInterface
}

type SpiracleV1Client struct {
	restClient rest.Interface
}

func NewForConfig(c *rest.Config) (*SpiracleV1Client, error) {
	config := *c
	config.GroupVersion = &v1.GroupVersion
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()
	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &SpiracleV1Client{restClient: client}, nil
}
func (it *SpiracleV1Client) RESTClient() rest.Interface {
	return it.restClient
}

func (it *SpiracleV1Client) RoomIngresses(namespace string) RoomIngressInterface {
	return &roomIngressClient{client: it.restClient, ns: namespace}
}

type RoomIngressInterface interface {
	Create(ctx context.Context, ring *v1.RoomIngress, opts metav1.CreateOptions) (*v1.RoomIngress, error)
	Update(ctx context.Context, ring *v1.RoomIngress, opts metav1.UpdateOptions) (*v1.RoomIngress, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.RoomIngress, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.RoomIngressList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.RoomIngress, err error)
}

type roomIngressClient struct {
	client rest.Interface
	ns     string
}

func (it *roomIngressClient) Create(ctx context.Context, ring *v1.RoomIngress, opts metav1.CreateOptions) (result *v1.RoomIngress, err error) {
	result = &v1.RoomIngress{}
	err = it.client.Post().
		Namespace(it.ns).
		Resource("roomingresses").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(ring).
		Do(ctx).
		Into(result)
	return

}
func (it *roomIngressClient) Update(ctx context.Context, ring *v1.RoomIngress, opts metav1.UpdateOptions) (result *v1.RoomIngress, err error) {
	result = &v1.RoomIngress{}
	err = it.client.Put().
		Namespace(it.ns).
		Resource("roomingresses").
		Name(ring.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(ring).
		Do(ctx).
		Into(result)
	return
}
func (it *roomIngressClient) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return it.client.Delete().
		Namespace(it.ns).
		Resource("roomingresses").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}
func (it *roomIngressClient) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return it.client.Delete().
		Namespace(it.ns).
		Resource("roomingresses").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}
func (it *roomIngressClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (result *v1.RoomIngress, err error) {
	result = &v1.RoomIngress{}
	err = it.client.Get().
		Namespace(it.ns).
		Resource("roomingresses").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}
func (it *roomIngressClient) List(ctx context.Context, opts metav1.ListOptions) (result *v1.RoomIngressList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.RoomIngressList{}
	err = it.client.Get().
		Namespace(it.ns).
		Resource("roomingresses").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}
func (it *roomIngressClient) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return it.client.Get().
		Namespace(it.ns).
		Resource("roomingresses").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}
func (it *roomIngressClient) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.RoomIngress, err error) {
	result = &v1.RoomIngress{}
	err = it.client.Patch(pt).
		Namespace(it.ns).
		Resource("roomingresses").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
