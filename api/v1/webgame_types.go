package v1

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// WebGameSpec defines the desired state of WebGame
type WebGameSpec struct {
	DisplayName string `json:"displayName"`
	GameType    string `json:"gameType"`
	// +kubebuilder:default:=localhost
	Domain string `json:"domain"`
	// +kubebuilder:default:=/
	IndexPage    string             `json:"indexPage"`
	IngressClass string             `json:"ingressClass"`
	ServerPort   intstr.IntOrString `json:"serverPort"`
	Replicas     *int32             `json:"replicas"`
	Image        string             `json:"image"`
	// +kubebuilder:validation:Optional
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`
}

// WebGameStatus defines the observed state of WebGame
type WebGameStatus struct {
	DeploymentStatus appsv1.DeploymentStatus `json:"deploymentStatus,omitempty"`
	GameAddress      string                  `json:"gameAddress,omitempty"`
	ClusterIP        string                  `json:"clusterIP,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=wg
// +kubebuilder:printcolumn:name="DisplayName",type="string",JSONPath=".spec.displayName"
// +kubebuilder:printcolumn:name="GameType",type="string",JSONPath=".spec.gameType"
// +kubebuilder:printcolumn:name="ServerPort",type="string",JSONPath=".spec.serverPort"
// +kubebuilder:printcolumn:name="Replicas",type="integer",JSONPath=".spec.replicas"
// +kubebuilder:printcolumn:name="Available",type="integer",JSONPath=".status.deploymentStatus.availableReplicas"
// +kubebuilder:printcolumn:name="Ready",type="integer",JSONPath=".status.deploymentStatus.readyReplicas"
// +kubebuilder:printcolumn:name="Updated",type="integer",JSONPath=".status.deploymentStatus.updatedReplicas"
// +kubebuilder:printcolumn:name="Observed",type="integer",JSONPath=".status.deploymentStatus.observedGeneration"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"

// WebGame is the Schema for the webgames API
type WebGame struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WebGameSpec   `json:"spec,omitempty"`
	Status WebGameStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// WebGameList contains a list of WebGame
type WebGameList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WebGame `json:"items"`
}

func init() {
	SchemeBuilder.Register(&WebGame{}, &WebGameList{})
}
