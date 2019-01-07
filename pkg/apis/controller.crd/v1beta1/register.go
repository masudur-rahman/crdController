package v1beta1

import (
	"github.com/masudur-rahman/crdController/pkg/apis/controller.crd"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var SchemeGroupVersion =schema.GroupVersion{Group:controllercrd.GroupName, Version: "v1beta1"}

func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

var (
	SchemeBuilder 	= runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme 	= SchemeBuilder.AddToScheme
)
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&CustomDeployment{},
		&CustomDeploymentList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
