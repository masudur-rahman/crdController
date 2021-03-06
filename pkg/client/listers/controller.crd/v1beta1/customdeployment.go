/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by lister-gen. DO NOT EDIT.

package v1beta1

import (
	v1beta1 "github.com/masudur-rahman/crdController/pkg/apis/controller.crd/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// CustomDeploymentLister helps list CustomDeployments.
type CustomDeploymentLister interface {
	// List lists all CustomDeployments in the indexer.
	List(selector labels.Selector) (ret []*v1beta1.CustomDeployment, err error)
	// CustomDeployments returns an object that can list and get CustomDeployments.
	CustomDeployments(namespace string) CustomDeploymentNamespaceLister
	CustomDeploymentListerExpansion
}

// customDeploymentLister implements the CustomDeploymentLister interface.
type customDeploymentLister struct {
	indexer cache.Indexer
}

// NewCustomDeploymentLister returns a new CustomDeploymentLister.
func NewCustomDeploymentLister(indexer cache.Indexer) CustomDeploymentLister {
	return &customDeploymentLister{indexer: indexer}
}

// List lists all CustomDeployments in the indexer.
func (s *customDeploymentLister) List(selector labels.Selector) (ret []*v1beta1.CustomDeployment, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.CustomDeployment))
	})
	return ret, err
}

// CustomDeployments returns an object that can list and get CustomDeployments.
func (s *customDeploymentLister) CustomDeployments(namespace string) CustomDeploymentNamespaceLister {
	return customDeploymentNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// CustomDeploymentNamespaceLister helps list and get CustomDeployments.
type CustomDeploymentNamespaceLister interface {
	// List lists all CustomDeployments in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1beta1.CustomDeployment, err error)
	// Get retrieves the CustomDeployment from the indexer for a given namespace and name.
	Get(name string) (*v1beta1.CustomDeployment, error)
	CustomDeploymentNamespaceListerExpansion
}

// customDeploymentNamespaceLister implements the CustomDeploymentNamespaceLister
// interface.
type customDeploymentNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all CustomDeployments in the indexer for a given namespace.
func (s customDeploymentNamespaceLister) List(selector labels.Selector) (ret []*v1beta1.CustomDeployment, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.CustomDeployment))
	})
	return ret, err
}

// Get retrieves the CustomDeployment from the indexer for a given namespace and name.
func (s customDeploymentNamespaceLister) Get(name string) (*v1beta1.CustomDeployment, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1beta1.Resource("customdeployment"), name)
	}
	return obj.(*v1beta1.CustomDeployment), nil
}
