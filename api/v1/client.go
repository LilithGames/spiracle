package v1

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	types "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
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
	config.ContentConfig.GroupVersion = &GroupVersion
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	config.UserAgent = rest.DefaultKubernetesUserAgent()

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
	Create(ctx context.Context, ring *RoomIngress, opts metav1.CreateOptions) (*RoomIngress, error)
	Update(ctx context.Context, ring *RoomIngress, opts metav1.UpdateOptions) (*RoomIngress, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*RoomIngress, error)
	List(ctx context.Context, opts metav1.ListOptions) (*RoomIngressList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *RoomIngress, err error)
}

type roomIngressClient struct {
	client rest.Interface
	ns     string
}

func (it *roomIngressClient) Create(ctx context.Context, ring *RoomIngress, opts metav1.CreateOptions) (result *RoomIngress, err error) {
	result = &RoomIngress{}
	err = it.client.Post().
		Namespace(it.ns).
		Resource("roomingresses").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(ring).
		Do(ctx).
		Into(result)
	return

}
func (it *roomIngressClient) Update(ctx context.Context, ring *RoomIngress, opts metav1.UpdateOptions) (result *RoomIngress, err error) {
	result = &RoomIngress{}
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
func (it *roomIngressClient) Get(ctx context.Context, name string, opts metav1.GetOptions) (result *RoomIngress, err error) {
	result = &RoomIngress{}
	err = it.client.Get().
		Namespace(it.ns).
		Resource("roomingresses").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}
func (it *roomIngressClient) List(ctx context.Context, opts metav1.ListOptions) (result *RoomIngressList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &RoomIngressList{}
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
func (it *roomIngressClient) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *RoomIngress, err error) {
	result = &RoomIngress{}
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
