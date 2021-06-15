package v1

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	labels "k8s.io/apimachinery/pkg/labels"
	watch "k8s.io/apimachinery/pkg/watch"
	types "k8s.io/apimachinery/pkg/types"
	v1 "github.com/LilithGames/spiracle/api/v1"
)

var ringsResource = schema.GroupVersionResource{Group: "projectdavinci.com", Version: "v1", Resource: "roomingresses"}
var ringsKind = schema.GroupVersionKind{Group: "projectdavinci.com", Version: "v1", Kind: "RoomIngress"}

type FakeSpiracleV1 struct {
	*testing.Fake
}

func (it *FakeSpiracleV1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
func (it *FakeSpiracleV1) RoomIngresses(namespace string) RoomIngressInterface {
	return &FakeRoomIngress{it, namespace}
}

type FakeRoomIngress struct {
	Fake *FakeSpiracleV1
	ns string
}
func (it *FakeRoomIngress) Create(ctx context.Context, ring *v1.RoomIngress, opts metav1.CreateOptions) (*v1.RoomIngress, error) {
	obj, err := it.Fake.Invokes(testing.NewCreateAction(ringsResource, it.ns, ring), &v1.RoomIngress{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1.RoomIngress), nil
}
func (it *FakeRoomIngress) Update(ctx context.Context, ring *v1.RoomIngress, opts metav1.UpdateOptions) (*v1.RoomIngress, error) {
	obj, err := it.Fake.Invokes(testing.NewUpdateAction(ringsResource, it.ns, ring), &v1.RoomIngress{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1.RoomIngress), nil
}
func (it *FakeRoomIngress) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	_, err := it.Fake.Invokes(testing.NewDeleteAction(ringsResource, it.ns, name), &v1.RoomIngress{})
	return err
}
func (it *FakeRoomIngress) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	_, err := it.Fake.Invokes(testing.NewDeleteCollectionAction(ringsResource, it.ns, listOpts), &v1.RoomIngress{})
	return err
}
func (it *FakeRoomIngress) Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.RoomIngress, error) {
	obj, err := it.Fake.Invokes(testing.NewGetAction(ringsResource, it.ns, name), &v1.RoomIngress{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1.RoomIngress), nil
}
func (it *FakeRoomIngress) List(ctx context.Context, opts metav1.ListOptions) (*v1.RoomIngressList, error) {
	obj, err := it.Fake.Invokes(testing.NewListAction(ringsResource, ringsKind, it.ns, opts), &v1.RoomIngressList{})
	if obj == nil {
		return nil, err
	}
	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1.RoomIngressList{ListMeta: obj.(*v1.RoomIngressList).ListMeta}
	for _, item := range obj.(*v1.RoomIngressList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err

}
func (it *FakeRoomIngress) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	return it.Fake.InvokesWatch(testing.NewWatchAction(ringsResource, it.ns, opts))
}
func (it *FakeRoomIngress) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.RoomIngress, err error) {
	obj, err := it.Fake.Invokes(testing.NewPatchSubresourceAction(ringsResource, it.ns, name, pt, data, subresources...), &v1.RoomIngress{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1.RoomIngress), err
}
