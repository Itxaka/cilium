// Copyright 2018-2019 Authors of Cilium
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package k8s

import (
	"net"
	"reflect"
	"strings"

	"github.com/cilium/cilium/pkg/annotation"
	"github.com/cilium/cilium/pkg/comparator"
	cilium_v2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
	"github.com/cilium/cilium/pkg/k8s/types"
	"github.com/cilium/cilium/pkg/logging/logfields"

	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/client-go/tools/cache"
)

func CopyObjToV1NetworkPolicy(obj interface{}) *types.NetworkPolicy {
	k8sNP, ok := obj.(*types.NetworkPolicy)
	if !ok {
		log.WithField(logfields.Object, logfields.Repr(obj)).
			Warn("Ignoring invalid k8s v1 NetworkPolicy")
		return nil
	}
	return k8sNP.DeepCopy()
}

func CopyObjToV1Services(obj interface{}) *types.Service {
	svc, ok := obj.(*types.Service)
	if !ok {
		log.WithField(logfields.Object, logfields.Repr(obj)).
			Warn("Ignoring invalid k8s v1 Service")
		return nil
	}
	return svc.DeepCopy()
}

func CopyObjToV1Endpoints(obj interface{}) *types.Endpoints {
	ep, ok := obj.(*types.Endpoints)
	if !ok {
		log.WithField(logfields.Object, logfields.Repr(obj)).
			Warn("Ignoring invalid k8s v1 Endpoints")
		return nil
	}
	return ep.DeepCopy()
}

func CopyObjToV1beta1Ingress(obj interface{}) *types.Ingress {
	ing, ok := obj.(*types.Ingress)
	if !ok {
		log.WithField(logfields.Object, logfields.Repr(obj)).
			Warn("Ignoring invalid k8s v1beta1 Ingress")
		return nil
	}
	return ing.DeepCopy()
}

func CopyObjToV2CNP(obj interface{}) *types.SlimCNP {
	cnp, ok := obj.(*types.SlimCNP)
	if !ok {
		log.WithField(logfields.Object, logfields.Repr(obj)).
			Warn("Ignoring invalid k8s v2 CiliumNetworkPolicy")
		return nil
	}
	return cnp.DeepCopy()
}

func CopyObjToV1Pod(obj interface{}) *types.Pod {
	pod, ok := obj.(*types.Pod)
	if !ok {
		log.WithField(logfields.Object, logfields.Repr(obj)).
			Warn("Ignoring invalid k8s v1 Pod")
		return nil
	}
	return pod.DeepCopy()
}

func CopyObjToV1Node(obj interface{}) *types.Node {
	node, ok := obj.(*types.Node)
	if !ok {
		log.WithField(logfields.Object, logfields.Repr(obj)).
			Warn("Ignoring invalid k8s v1 Node")
		return nil
	}
	return node.DeepCopy()
}

func CopyObjToV1Namespace(obj interface{}) *types.Namespace {
	ns, ok := obj.(*types.Namespace)
	if !ok {
		log.WithField(logfields.Object, logfields.Repr(obj)).
			Warn("Ignoring invalid k8s v1 Namespace")
		return nil
	}
	return ns.DeepCopy()
}

func EqualV1NetworkPolicy(np1, np2 *types.NetworkPolicy) bool {
	// As Cilium uses all of the Spec from a NP it's not probably not worth
	// it to create a dedicated deep equal	 function to compare both network
	// policies.
	return np1.Name == np2.Name &&
		np1.Namespace == np2.Namespace &&
		reflect.DeepEqual(np1.Spec, np2.Spec)
}

func EqualV1Services(svc1, svc2 *types.Service) bool {
	// Service annotations are used to mark services as global, shared, etc.
	if !comparator.MapStringEquals(svc1.GetAnnotations(), svc2.GetAnnotations()) {
		return false
	}

	clusterIP := net.ParseIP(svc1.Spec.ClusterIP)
	headless := false
	if strings.ToLower(svc1.Spec.ClusterIP) == "none" {
		headless = true
	}
	si1 := NewService(clusterIP, headless, svc1.Labels, svc1.Spec.Selector)

	clusterIP = net.ParseIP(svc2.Spec.ClusterIP)
	headless = false
	if strings.ToLower(svc2.Spec.ClusterIP) == "none" {
		headless = true
	}
	si2 := NewService(clusterIP, headless, svc2.Labels, svc2.Spec.Selector)

	// Please write all the equalness logic inside the K8sServiceInfo.Equals()
	// method.
	return si1.DeepEquals(si2)
}

func EqualV1Endpoints(ep1, ep2 *types.Endpoints) bool {
	// We only care about the Name, Namespace and Subsets of a particular
	// endpoint.
	return ep1.Name == ep2.Name &&
		ep1.Namespace == ep2.Namespace &&
		reflect.DeepEqual(ep1.Subsets, ep2.Subsets)
}

func EqualV1beta1Ingress(ing1, ing2 *types.Ingress) bool {
	if ing1.Name != ing2.Name || ing1.Namespace != ing2.Namespace {
		return false
	}
	switch {
	case (ing1.Spec.Backend == nil) != (ing2.Spec.Backend == nil):
		return false
	case (ing1.Spec.Backend == nil) && (ing2.Spec.Backend == nil):
		return true
	}

	return ing1.Spec.Backend.ServicePort.IntVal ==
		ing2.Spec.Backend.ServicePort.IntVal &&
		ing1.Spec.Backend.ServicePort.StrVal ==
			ing2.Spec.Backend.ServicePort.StrVal
}

func EqualV2CNP(cnp1, cnp2 *types.SlimCNP) bool {
	if !(cnp1.Name == cnp2.Name && cnp1.Namespace == cnp2.Namespace) {
		return false
	}

	// Ignore v1.LastAppliedConfigAnnotation annotation
	lastAppliedCfgAnnotation1, ok1 := cnp1.GetAnnotations()[v1.LastAppliedConfigAnnotation]
	lastAppliedCfgAnnotation2, ok2 := cnp2.GetAnnotations()[v1.LastAppliedConfigAnnotation]
	defer func() {
		if ok1 && cnp1.GetAnnotations() != nil {
			cnp1.GetAnnotations()[v1.LastAppliedConfigAnnotation] = lastAppliedCfgAnnotation1
		}
		if ok2 && cnp2.GetAnnotations() != nil {
			cnp2.GetAnnotations()[v1.LastAppliedConfigAnnotation] = lastAppliedCfgAnnotation2
		}
	}()
	delete(cnp1.GetAnnotations(), v1.LastAppliedConfigAnnotation)
	delete(cnp2.GetAnnotations(), v1.LastAppliedConfigAnnotation)

	return comparator.MapStringEquals(cnp1.GetAnnotations(), cnp2.GetAnnotations()) &&
		reflect.DeepEqual(cnp1.Spec, cnp2.Spec) &&
		reflect.DeepEqual(cnp1.Specs, cnp2.Specs)
}

// AnnotationsEqual returns whether the annotation with any key in
// relevantAnnotations is equal in anno1 and anno2.
func AnnotationsEqual(relevantAnnotations []string, anno1, anno2 map[string]string) bool {
	for _, an := range relevantAnnotations {
		if anno1[an] != anno2[an] {
			return false
		}
	}
	return true
}

func EqualV1Pod(pod1, pod2 *types.Pod) bool {
	// We only care about the HostIP, the PodIP and the labels of the pods.
	if pod1.StatusPodIP != pod2.StatusPodIP ||
		pod1.StatusHostIP != pod2.StatusHostIP {
		return false
	}

	if !AnnotationsEqual([]string{annotation.ProxyVisibility}, pod1.GetAnnotations(), pod2.GetAnnotations()) {
		return false
	}

	oldPodLabels := pod1.GetLabels()
	newPodLabels := pod2.GetLabels()
	return comparator.MapStringEquals(oldPodLabels, newPodLabels)
}

func EqualV1Node(node1, node2 *types.Node) bool {
	if node1.GetObjectMeta().GetName() != node2.GetObjectMeta().GetName() {
		return false
	}

	anno1 := node1.GetAnnotations()
	anno2 := node2.GetAnnotations()
	annotationsWeCareAbout := []string{
		annotation.CiliumHostIP,
		annotation.CiliumHostIPv6,
		annotation.V4HealthName,
		annotation.V6HealthName,
	}
	for _, an := range annotationsWeCareAbout {
		if anno1[an] != anno2[an] {
			return false
		}
	}
	return true
}

func EqualV1Namespace(ns1, ns2 *types.Namespace) bool {
	// we only care about namespace labels.
	return ns1.Name == ns2.Name &&
		comparator.MapStringEquals(ns1.GetLabels(), ns2.GetLabels())
}

// ConvertToNetworkPolicy converts a *networkingv1.NetworkPolicy into a
// *types.NetworkPolicy or a cache.DeletedFinalStateUnknown into
// a cache.DeletedFinalStateUnknown with a *types.NetworkPolicy in its Obj.
// If the given obj can't be cast into either *networkingv1.NetworkPolicy
// nor cache.DeletedFinalStateUnknown, the original obj is returned.
func ConvertToNetworkPolicy(obj interface{}) interface{} {
	// TODO check which fields we really need
	switch concreteObj := obj.(type) {
	case *networkingv1.NetworkPolicy:
		return &types.NetworkPolicy{
			NetworkPolicy: concreteObj,
		}
	case cache.DeletedFinalStateUnknown:
		netPol, ok := concreteObj.Obj.(*networkingv1.NetworkPolicy)
		if !ok {
			return obj
		}
		return cache.DeletedFinalStateUnknown{
			Key: concreteObj.Key,
			Obj: &types.NetworkPolicy{
				NetworkPolicy: netPol,
			},
		}
	default:
		return obj
	}
}

// ConvertToK8sService converts a *v1.Service into a
// *types.Service or a cache.DeletedFinalStateUnknown into
// a cache.DeletedFinalStateUnknown with a *types.Service in its Obj.
// If the given obj can't be cast into either *v1.Service
// nor cache.DeletedFinalStateUnknown, the original obj is returned.
func ConvertToK8sService(obj interface{}) interface{} {
	// TODO check which fields we really need
	switch concreteObj := obj.(type) {
	case *v1.Service:
		return &types.Service{
			Service: concreteObj,
		}
	case cache.DeletedFinalStateUnknown:
		svc, ok := concreteObj.Obj.(*v1.Service)
		if !ok {
			return obj
		}
		return cache.DeletedFinalStateUnknown{
			Key: concreteObj.Key,
			Obj: &types.Service{
				Service: svc,
			},
		}
	default:
		return obj
	}
}

// ConvertToK8sEndpoints converts a *v1.Endpoints into a
// *types.Endpoints or a cache.DeletedFinalStateUnknown into
// a cache.DeletedFinalStateUnknown with a *types.Endpoints in its Obj.
// If the given obj can't be cast into either *v1.Endpoints
// nor cache.DeletedFinalStateUnknown, the original obj is returned.
func ConvertToK8sEndpoints(obj interface{}) interface{} {
	// TODO check which fields we really need
	switch concreteObj := obj.(type) {
	case *v1.Endpoints:
		return &types.Endpoints{
			Endpoints: concreteObj,
		}
	case cache.DeletedFinalStateUnknown:
		eps, ok := concreteObj.Obj.(*v1.Endpoints)
		if !ok {
			return obj
		}
		return cache.DeletedFinalStateUnknown{
			Key: concreteObj.Key,
			Obj: &types.Endpoints{
				Endpoints: eps,
			},
		}
	default:
		return obj
	}
}

// ConvertToIngress converts a *v1beta1.Ingress into a
// *types.Ingress or a cache.DeletedFinalStateUnknown into
// a cache.DeletedFinalStateUnknown with a *types.Ingress in its Obj.
// If the given obj can't be cast into either *v1beta1.Ingress
// nor cache.DeletedFinalStateUnknown, the original obj is returned.
func ConvertToIngress(obj interface{}) interface{} {
	// TODO check which fields we really need
	switch concreteObj := obj.(type) {
	case *v1beta1.Ingress:
		return &types.Ingress{
			Ingress: concreteObj,
		}
	case cache.DeletedFinalStateUnknown:
		ingrss, ok := concreteObj.Obj.(*v1beta1.Ingress)
		if !ok {
			return obj
		}
		return cache.DeletedFinalStateUnknown{
			Key: concreteObj.Key,
			Obj: &types.Ingress{
				Ingress: ingrss,
			},
		}
	default:
		return obj
	}
}

// ConvertToCNPWithStatus converts a *cilium_v2.CiliumNetworkPolicy into a
// *types.SlimCNP or a cache.DeletedFinalStateUnknown into
// a cache.DeletedFinalStateUnknown with a *types.SlimCNP in its Obj.
// If the given obj can't be cast into either *cilium_v2.CiliumNetworkPolicy
// nor cache.DeletedFinalStateUnknown, the original obj is returned.
func ConvertToCNPWithStatus(obj interface{}) interface{} {
	switch concreteObj := obj.(type) {
	case *cilium_v2.CiliumNetworkPolicy:
		return &types.SlimCNP{
			CiliumNetworkPolicy: concreteObj,
		}
	case cache.DeletedFinalStateUnknown:
		cnp, ok := concreteObj.Obj.(*cilium_v2.CiliumNetworkPolicy)
		if !ok {
			return obj
		}
		return cache.DeletedFinalStateUnknown{
			Key: concreteObj.Key,
			Obj: &types.SlimCNP{
				CiliumNetworkPolicy: cnp,
			},
		}
	default:
		return obj
	}
}

// ConvertToCNP converts a *cilium_v2.CiliumNetworkPolicy into a
// *types.SlimCNP without the Status field of the given CNP, or a
// cache.DeletedFinalStateUnknown into a cache.DeletedFinalStateUnknown with a
// *types.SlimCNP, also without the Status field of the given CNP, in its Obj.
// If the given obj can't be cast into either *cilium_v2.CiliumNetworkPolicy
// nor cache.DeletedFinalStateUnknown, the original obj is returned.
// WARNING calling this function will set *all* fields of the given CNP as
// empty.
func ConvertToCNP(obj interface{}) interface{} {
	switch concreteObj := obj.(type) {
	case *cilium_v2.CiliumNetworkPolicy:
		cnp := &types.SlimCNP{
			CiliumNetworkPolicy: &cilium_v2.CiliumNetworkPolicy{
				TypeMeta:   concreteObj.TypeMeta,
				ObjectMeta: concreteObj.ObjectMeta,
				Spec:       concreteObj.Spec,
				Specs:      concreteObj.Specs,
			},
		}
		*concreteObj = cilium_v2.CiliumNetworkPolicy{}
		return cnp
	case cache.DeletedFinalStateUnknown:
		cnp, ok := concreteObj.Obj.(*cilium_v2.CiliumNetworkPolicy)
		if !ok {
			return obj
		}
		dfsu := cache.DeletedFinalStateUnknown{
			Key: concreteObj.Key,
			Obj: &types.SlimCNP{
				CiliumNetworkPolicy: &cilium_v2.CiliumNetworkPolicy{
					TypeMeta:   cnp.TypeMeta,
					ObjectMeta: cnp.ObjectMeta,
					Spec:       cnp.Spec,
					Specs:      cnp.Specs,
				},
			},
		}
		*cnp = cilium_v2.CiliumNetworkPolicy{}
		return dfsu
	default:
		return obj
	}
}

// ConvertToPod converts a *v1.Pod into a
// *types.Pod or a cache.DeletedFinalStateUnknown into
// a cache.DeletedFinalStateUnknown with a *types.Pod in its Obj.
// If the given obj can't be cast into either *v1.Pod
// nor cache.DeletedFinalStateUnknown, the original obj is returned.
// WARNING calling this function will set *all* fields of the given Pod as
// empty.
func ConvertToPod(obj interface{}) interface{} {
	switch concreteObj := obj.(type) {
	case *v1.Pod:
		p := &types.Pod{
			TypeMeta:        concreteObj.TypeMeta,
			ObjectMeta:      concreteObj.ObjectMeta,
			StatusPodIP:     concreteObj.Status.PodIP,
			StatusHostIP:    concreteObj.Status.HostIP,
			SpecHostNetwork: concreteObj.Spec.HostNetwork,
		}
		*concreteObj = v1.Pod{}
		return p
	case cache.DeletedFinalStateUnknown:
		pod, ok := concreteObj.Obj.(*v1.Pod)
		if !ok {
			return obj
		}
		dfsu := cache.DeletedFinalStateUnknown{
			Key: concreteObj.Key,
			Obj: &types.Pod{
				TypeMeta:        pod.TypeMeta,
				ObjectMeta:      pod.ObjectMeta,
				StatusPodIP:     pod.Status.PodIP,
				StatusHostIP:    pod.Status.HostIP,
				SpecHostNetwork: pod.Spec.HostNetwork,
			},
		}
		*pod = v1.Pod{}
		return dfsu
	default:
		return obj
	}
}

// ConvertToNode converts a *v1.Node into a
// *types.Node or a cache.DeletedFinalStateUnknown into
// a cache.DeletedFinalStateUnknown with a *types.Node in its Obj.
// If the given obj can't be cast into either *v1.Node
// nor cache.DeletedFinalStateUnknown, the original obj is returned.
// WARNING calling this function will set *all* fields of the given Node as
// empty.
func ConvertToNode(obj interface{}) interface{} {
	switch concreteObj := obj.(type) {
	case *v1.Node:
		p := &types.Node{
			TypeMeta:        concreteObj.TypeMeta,
			ObjectMeta:      concreteObj.ObjectMeta,
			StatusAddresses: concreteObj.Status.Addresses,
			SpecPodCIDR:     concreteObj.Spec.PodCIDR,
			SpecPodCIDRs:    concreteObj.Spec.PodCIDRs,
		}
		*concreteObj = v1.Node{}
		return p
	case cache.DeletedFinalStateUnknown:
		node, ok := concreteObj.Obj.(*v1.Node)
		if !ok {
			return obj
		}
		dfsu := cache.DeletedFinalStateUnknown{
			Key: concreteObj.Key,
			Obj: &types.Node{
				TypeMeta:        node.TypeMeta,
				ObjectMeta:      node.ObjectMeta,
				StatusAddresses: node.Status.Addresses,
				SpecPodCIDR:     node.Spec.PodCIDR,
				SpecPodCIDRs:    node.Spec.PodCIDRs,
			},
		}
		*node = v1.Node{}
		return dfsu
	default:
		return obj
	}
}

// ConvertToNamespace converts a *v1.Namespace into a
// *types.Namespace or a cache.DeletedFinalStateUnknown into
// a cache.DeletedFinalStateUnknown with a *types.Namespace in its Obj.
// If the given obj can't be cast into either *v1.Namespace
// nor cache.DeletedFinalStateUnknown, the original obj is returned.
// WARNING calling this function will set *all* fields of the given Namespace as
// empty.
func ConvertToNamespace(obj interface{}) interface{} {
	switch concreteObj := obj.(type) {
	case *v1.Namespace:
		p := &types.Namespace{
			TypeMeta:   concreteObj.TypeMeta,
			ObjectMeta: concreteObj.ObjectMeta,
		}
		*concreteObj = v1.Namespace{}
		return p
	case cache.DeletedFinalStateUnknown:
		namespace, ok := concreteObj.Obj.(*v1.Namespace)
		if !ok {
			return obj
		}
		dfsu := cache.DeletedFinalStateUnknown{
			Key: concreteObj.Key,
			Obj: &types.Namespace{
				TypeMeta:   namespace.TypeMeta,
				ObjectMeta: namespace.ObjectMeta,
			},
		}
		*namespace = v1.Namespace{}
		return dfsu
	default:
		return obj
	}
}

// ConvertToCiliumNode converts a *cilium_v2.CiliumNode into a
// *cilium_v2.CiliumNode or a cache.DeletedFinalStateUnknown into
// a cache.DeletedFinalStateUnknown with a *cilium_v2.CiliumNode in its Obj.
// If the given obj can't be cast into either *cilium_v2.CiliumNode
// nor cache.DeletedFinalStateUnknown, the original obj is returned.
func ConvertToCiliumNode(obj interface{}) interface{} {
	// TODO create a slim type of the CiliumNode
	switch concreteObj := obj.(type) {
	case *cilium_v2.CiliumNode:
		return concreteObj
	case cache.DeletedFinalStateUnknown:
		ciliumNode, ok := concreteObj.Obj.(*cilium_v2.CiliumNode)
		if !ok {
			return obj
		}
		return cache.DeletedFinalStateUnknown{
			Key: concreteObj.Key,
			Obj: ciliumNode,
		}
	default:
		return obj
	}
}

// CopyObjToCiliumNode attempts to cast object to a CiliumNode object and
// returns a deep copy if the castin succeeds. Otherwise, nil is returned.
func CopyObjToCiliumNode(obj interface{}) *cilium_v2.CiliumNode {
	cn, ok := obj.(*cilium_v2.CiliumNode)
	if !ok {
		log.WithField(logfields.Object, logfields.Repr(obj)).
			Warn("Ignoring invalid CiliumNode")
		return nil
	}
	return cn.DeepCopy()
}

// ConvertToCiliumEndpoint converts a *cilium_v2.CiliumEndpoint into a
// *types.CiliumEndpoint or a cache.DeletedFinalStateUnknown into a
// cache.DeletedFinalStateUnknown with a *types.CiliumEndpoint in its Obj.
// If the given obj can't be cast into either *cilium_v2.CiliumEndpoint nor
// cache.DeletedFinalStateUnknown, the original obj is returned.
func ConvertToCiliumEndpoint(obj interface{}) interface{} {
	switch concreteObj := obj.(type) {
	case *cilium_v2.CiliumEndpoint:
		p := &types.CiliumEndpoint{
			TypeMeta:   concreteObj.TypeMeta,
			ObjectMeta: concreteObj.ObjectMeta,
			Encryption: concreteObj.Status.Encryption.DeepCopy(),
			Identity:   concreteObj.Status.Identity.DeepCopy(),
			Networking: concreteObj.Status.Networking.DeepCopy(),
		}
		*concreteObj = cilium_v2.CiliumEndpoint{}
		return p
	case cache.DeletedFinalStateUnknown:
		ciliumEndpoint, ok := concreteObj.Obj.(*cilium_v2.CiliumEndpoint)
		if !ok {
			return obj
		}
		dfsu := cache.DeletedFinalStateUnknown{
			Key: concreteObj.Key,
			Obj: &types.CiliumEndpoint{
				TypeMeta:   ciliumEndpoint.TypeMeta,
				ObjectMeta: ciliumEndpoint.ObjectMeta,
				Encryption: ciliumEndpoint.Status.Encryption.DeepCopy(),
				Identity:   ciliumEndpoint.Status.Identity.DeepCopy(),
				Networking: ciliumEndpoint.Status.Networking.DeepCopy(),
			},
		}
		*ciliumEndpoint = cilium_v2.CiliumEndpoint{}
		return dfsu
	default:
		return obj
	}
}

// CopyObjToCiliumEndpoint attempts to cast object to a CiliumEndpoint object
// and returns a deep copy if the castin succeeds. Otherwise, nil is returned.
func CopyObjToCiliumEndpoint(obj interface{}) *types.CiliumEndpoint {
	ce, ok := obj.(*types.CiliumEndpoint)
	if !ok {
		log.WithField(logfields.Object, logfields.Repr(obj)).
			Warn("Ignoring invalid CiliumEndpoint")
		return nil
	}
	return ce.DeepCopy()
}
