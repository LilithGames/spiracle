package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//+kubebuilder:rbac:groups=projectdavinci.com,resources=roomingresses,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=projectdavinci.com,resources=roomingresses/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=projectdavinci.com,resources=roomingresses/finalizers,verbs=update
//+kubebuilder:webhook:path=/mutate,mutating=true,failurePolicy=fail,groups=projectdavinci.com,resources=roomingresses,verbs=create;update;patch;delete,versions=v1,name=roomingress-webhook.projectdavinci.com,sideEffects=None,admissionReviewVersions=v1

func init() {
	SchemeBuilder.Register(&RoomIngress{}, &RoomIngressList{})
}

type RoomIngressPlayer struct {
	//+kubebuilder:validation:MinLength=1
	Id string `json:"id"`
	//+kubebuilder:validation:Minimum=0
	//+kubebuilder:validation:Maximum=4294967295
	Token int64 `json:"token"`
}

type RoomIngressRoom struct {
	//+kubebuilder:validation:MinLength=1
	Id string `json:"id,omitempty"`
	//+kubebuilder:validation:MinLength=1
	Server string `json:"server,omitempty"`
	//+kubebuilder:validation:MinLength=1
	Upstream string              `json:"upstream,omitempty"`
	Players  []RoomIngressPlayer `json:"players"`
}

type RoomIngressSpec struct {
	//+kubebuilder:validation:MinItems=1
	Rooms []RoomIngressRoom `json:"rooms,omitempty"`
}

// +kubebuilder:validation:Enum=Success;Pending;Failure;Expired;Retry
type PlayerStatus string

const PlayerStatusSuccess PlayerStatus = "Success"
const PlayerStatusPending PlayerStatus = "Pending"
const PlayerStatusFailure PlayerStatus = "Failure"
const PlayerStatusExpired PlayerStatus = "Expired"
const PlayerStatusRetry PlayerStatus = "Retry"

type RoomIngressPlayerStatus struct {
	Id string `json:"id"`
	//+kubebuilder:validation:Minimum=0
	//+kubebuilder:validation:Maximum=4294967295
	Token     int64        `json:"token"`
	Externals []string     `json:"externals,omitempty"`
	Timestamp *metav1.Time  `json:"timestamp,omitempty"`
	Expire    *metav1.Time  `json:"expire,omitempty"`
	Status    PlayerStatus `json:"status"`
	Detail    string       `json:"detail"`
}

type RoomIngressRoomStatus struct {
	Id       string                    `json:"id,omitempty"`
	Server   string                    `json:"server,omitempty"`
	Upstream string                    `json:"upstream,omitempty"`
	Players  []RoomIngressPlayerStatus `json:"players,omitempty"`
}

type RoomIngressStatus struct {
	Rooms []RoomIngressRoomStatus `json:"rooms,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Server",type=string,JSONPath=`.spec.rooms[0].server`
//+kubebuilder:printcolumn:name="Room",type=string,JSONPath=`.spec.rooms[0].id`
//+kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`
//+kubebuilder:resource:shortName="ring"
type RoomIngress struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RoomIngressSpec   `json:"spec,omitempty"`
	Status RoomIngressStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
type RoomIngressList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RoomIngress `json:"items"`
}
