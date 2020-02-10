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

// Code generated by informer-gen. DO NOT EDIT.

package v1

import (
	time "time"

	cephrookiov1 "github.com/rook/rook/pkg/apis/ceph.rook.io/v1"
	versioned "github.com/rook/rook/pkg/client/clientset/versioned"
	internalinterfaces "github.com/rook/rook/pkg/client/informers/externalversions/internalinterfaces"
	v1 "github.com/rook/rook/pkg/client/listers/ceph.rook.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// CephObjectStoreUserInformer provides access to a shared informer and lister for
// CephObjectStoreUsers.
type CephObjectStoreUserInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1.CephObjectStoreUserLister
}

type cephObjectStoreUserInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewCephObjectStoreUserInformer constructs a new informer for CephObjectStoreUser type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewCephObjectStoreUserInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredCephObjectStoreUserInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredCephObjectStoreUserInformer constructs a new informer for CephObjectStoreUser type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredCephObjectStoreUserInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CephV1().CephObjectStoreUsers(namespace).List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CephV1().CephObjectStoreUsers(namespace).Watch(options)
			},
		},
		&cephrookiov1.CephObjectStoreUser{},
		resyncPeriod,
		indexers,
	)
}

func (f *cephObjectStoreUserInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredCephObjectStoreUserInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *cephObjectStoreUserInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&cephrookiov1.CephObjectStoreUser{}, f.defaultInformer)
}

func (f *cephObjectStoreUserInformer) Lister() v1.CephObjectStoreUserLister {
	return v1.NewCephObjectStoreUserLister(f.Informer().GetIndexer())
}
