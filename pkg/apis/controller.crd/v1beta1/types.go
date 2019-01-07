package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Foo is a specification for Foo Resource
type Foo struct {
	metav1.TypeMeta		`json:",inline"`
	metav1.ObjectMeta	`json:"metadata,omitempty"`

	Spec 	FooSpec		`json:"spec"`
	Status 	FooStatus	`json:"status"`
}


// Spec for Foo resource
type FooSpec struct {
	DeploymentName	string	`json:"deploymentName"`
	Replicas		*int32	`json:"replicas"`
}

// Status of Foo resource
type FooStatus struct {
	AvailableReplicas int32	`json:"availableReplicas"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// List of Foo type resource
type FooList struct {
	metav1.TypeMeta	`json:",inline"`
	metav1.ListMeta	`json:"metadata"`

	Items []Foo		`json:"items"`
}