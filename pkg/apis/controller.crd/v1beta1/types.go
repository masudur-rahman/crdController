package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CustomDeployment is a specification for CustomDeployment Resource
type CustomDeployment struct {
	metav1.TypeMeta		`json:",inline"`
	metav1.ObjectMeta	`json:"metadata,omitempty"`

	Spec 	CustomDeploymentSpec	`json:"spec"`
	Status 	CustomDeploymentStatus	`json:"status"`
}


// Spec for CustomDeployment resource
type CustomDeploymentSpec struct {
	DeploymentName	string	`json:"deploymentName"`
	Replicas		*int32	`json:"replicas"`
}

// Status of CustomDeployment resource
type CustomDeploymentStatus struct {
	AvailableReplicas int32	`json:"availableReplicas"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// List of CustomDeployment type resource
type CustomDeploymentList struct {
	metav1.TypeMeta	`json:",inline"`
	metav1.ListMeta	`json:"metadata"`

	Items []CustomDeployment `json:"items"`
}