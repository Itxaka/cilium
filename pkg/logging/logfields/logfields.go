// Copyright 2017-2018 Authors of Cilium
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

// Package logfields defines common logging fields which are used across packages
package logfields

const (

	// LogSubsys is the field denoting the subsystem when logging
	LogSubsys = "subsys"

	// Signal is the field to print os signals on exit etc.
	Signal = "signal"

	// Node is a host machine in the cluster, running cilium
	Node = "node"

	// NodeName is a human readable name for the node
	NodeName = "nodeName"

	// EndpointID is the numeric endpoint identifier
	EndpointID = "endpointID"

	// EndpointState is the current endpoint state
	EndpointState = "endpointState"

	// EventUUID is an event unique identifier
	EventUUID = "eventID"

	// ContainerID is the container identifier
	ContainerID = "containerID"

	// IdentityLabels are the labels relevant for the security identity
	IdentityLabels = "identityLabels"

	// InfoLabels are the labels relevant for the security identity
	InfoLabels = "infoLabels"

	// Labels are any label, they may not be relevant to the security identity.
	Labels = "labels"

	// Controller is the name of the controller to log it.
	Controller = "controller"

	// Identity is the identifier of a security identity
	Identity = "identity"

	// OldIdentity is a previously used security identity
	OldIdentity = "oldIdentity"

	// PolicyRevision is the revision of the policy in the repository or of
	// the object in question
	PolicyRevision = "policyRevision"

	// DatapathPolicyRevision is the policy revision currently running in
	// the datapath
	DatapathPolicyRevision = "datapathPolicyRevision"

	// DesiredPolicyRevision is the latest policy revision as evaluated for
	// an endpoint. It is the desired policy revision to be implemented
	// into the datapath.
	DesiredPolicyRevision = "desiredPolicyRevision"

	// PolicyID is the identifier of a L3, L4 or L7 Policy. Ideally the .NumericIdentity
	PolicyID = "policyID"

	// AddedPolicyID is the .NumericIdentity, or set or them
	AddedPolicyID = "policyID.Added"

	// DeletedPolicyID is the .NumericIdentity, or set or them
	DeletedPolicyID = "policyID.Deleted"

	// L3PolicyID is the identifier of a L3 Policy
	L3PolicyID = "policyID.L3"

	// L4PolicyID is the identifier of a L4 Policy
	L4PolicyID = "PolicyID.L4"

	// DNSName is a FQDN or not fully qualified name intended for DNS lookups
	DNSName = "dnsName"

	// DNSRequestID is the DNS request id used by dns-proxy
	DNSRequestID = "DNSRequestID"

	// IPAddr is an IPV4 or IPv6 address
	IPAddr = "ipAddr"

	// IPMask is an IPV4 or IPv6 address mask
	IPMask = "ipMask"

	// IPv4 is an IPv4 address
	IPv4 = "ipv4"

	// IPv6 is an IPv6 address
	IPv6 = "ipv6"

	// BuildDuration is the time elapsed to build a BPF program
	BuildDuration = "buildDuration"

	// BPFCompilationTime is the time elapsed to build a BPF endpoint program
	BPFCompilationTime = "BPFCompilationTime"

	// EndpointRegenerationTime is the time elapsed to generate an endpoint
	EndpointRegenerationTime = "endpointRegenerationTime"

	// StartTime is the start time of an event
	StartTime = "startTime"

	// EndTime is the end time of an event
	EndTime = "endTime"

	// Duration is the duration of a measured operation
	Duration = "duration"

	// V4HealthIP is an address used to contact the cilium-health endpoint
	V4HealthIP = "v4healthIP.IPv4"

	// V6HealthIP is an address used to contact the cilium-health endpoint
	V6HealthIP = "v6healthIP.IPv6"

	// V4CiliumHostIP is an address used for the cilium_host interface.
	V4CiliumHostIP = "v4CiliumHostIP.IPv4"

	// V6CiliumHostIP is an address used for the cilium_host interface.
	V6CiliumHostIP = "v6CiliumHostIP.IPv6"

	// L3n4Addr is a L3 (IP) + L4 (port and protocol) address object.
	L3n4Addr = "l3n4Addr"

	// L3n4AddrID is the allocated ID for a L3n4Addr object
	L3n4AddrID = "l3n4AddrID"

	// Port is a L4 port
	Port = "port"

	// Family is the L3 protocol family
	Family = "family"

	// Protocol is the L4 protocol
	Protocol = "protocol"

	// V4Prefix is a IPv4 subnet/CIDR prefix
	V4Prefix = "v4Prefix"

	// V6Prefix is a IPv6 subnet/CIDR prefix
	V6Prefix = "v6Prefix"

	// Interface is an interface id/name on the system
	Interface = "interface"

	// Ipvlan is a ipvlan object or ID
	Ipvlan = "ipvlan"

	// Veth is a veth object or ID
	Veth = "veth"

	// VethPair is a tuple of Veth that are paired
	VethPair = "vethPair"

	// NetNSName is a name of a network namespace
	NetNSName = "netNSName"

	// Hash is a hash of something
	Hash = "hash"

	// ServiceName is the orchestration framework name for a service
	ServiceName = "serviceName"

	// ClusterName is the name of the cluster
	ClusterName = "clusterName"

	// ServiceID is the orchestration unique ID of a service
	ServiceID = "serviceID"

	// ServiceIP is the IP of the service
	ServiceIP = "serviceIP"

	// ServiceKey is the key of the service in a BPF map
	ServiceKey = "svcKey"

	// ServiceValue is the value of the service in a BPF map
	ServiceValue = "svcVal"

	// ServiceType is the type of the service
	ServiceType = "svcType"

	// BackendIDs is the map of backend IDs (lbmap) indexed by backend address
	BackendIDs = "backendIDs"

	// BackendID is the ID of the backend
	BackendID = "backendID"

	// Backends is the list of the service backends
	Backends = "backends"

	// BackendName is the name of the backend
	BackendName = "backendName"

	// SlaveSlot is the slot number in a service BPF map
	SlaveSlot = "slaveSlot"

	// CiliumNetworkPolicy is a cilium specific NetworkPolicy
	CiliumNetworkPolicy = "ciliumNetworkPolicy"

	// CiliumNetworkPolicyName is the name of a CiliumNetworkPolicy
	CiliumNetworkPolicyName = "ciliumNetworkPolicyName"

	// BPFMapKey is a key from a BPF map
	BPFMapKey = "bpfMapKey"

	// BPFMapValue is a value from a BPF map
	BPFMapValue = "bpfMapValue"

	// XDPDevice is the device name
	XDPDevice = "xdpDevice"

	// Device is the device name
	Device = "device"

	// IpvlanMasterDevice is the ipvlan master device name
	IpvlanMasterDevice = "ipvlanMasterDevice"

	// DatapathMode is the datapath mode name
	DatapathMode = "datapathMode"

	// Tunnel is the tunnel name
	Tunnel = "tunnel"

	// EndpointLabelSelector is a selector for Endpoints by label
	EndpointLabelSelector = "EndpointLabelSelector"

	// EndpointSelector is a selector for Endpoints
	EndpointSelector = "EndpointSelector"

	// Path is a filesystem path. It can be a file or directory.
	// Note: pkg/proxy/accesslog points to this variable so be careful when
	// changing the value
	Path = "file-path"

	// Line is a line number within a file
	Line = "line"

	// Object is used when "%+v" printing Go objects for debug or error handling.
	// It is often paired with logfields.Repr to render the object.
	Object = "obj"

	// Request is a request object received by us, reported in debug or error.
	// It is often paired with logfields.Repr to render the object.
	Request = "req"

	// Params are the parameters of a request, reported in debug or error.
	Params = "params"

	// Response is a response object received by us, reported in debug or error.
	// It is often paired with logfields.Repr to render the object.
	Response = "resp"

	// Route is a L2 or L3 Linux route
	Route = "route"

	// RetryUUID is an UUID identical for all retries of a set
	RetryUUID = "retryUUID"

	// Envoy xDS-protocol-specific

	// XDSStreamID is the ID of an xDS request stream.
	XDSStreamID = "xdsStreamID"

	// XDSVersionInfo is the version info of an xDS resource.
	XDSVersionInfo = "xdsVersionInfo"

	// XDSTypeURL is the URL that identifies an xDS resource type.
	XDSTypeURL = "xdsTypeURL"

	// XDSNonce is a nonce sent in xDS requests and responses.
	XDSNonce = "xdsNonce"

	// XDSCanary is a boolean indicating whether a response is a dry run.
	XDSCanary = "xdsCanary"

	// XDSResourceName is the name of an xDS resource.
	XDSResourceName = "xdsResourceName"

	// XDSClientNode is the ID of an XDS client, e.g. an Envoy node.
	XDSClientNode = "xdsClientNode"

	// XDSResource is an xDS resource message.
	XDSResource = "xdsResource"

	// K8s-specific

	// K8sNodeID is the k8s ID of a K8sNode
	K8sNodeID = "k8sNodeID"

	// K8sPodName is the name of a k8s pod
	K8sPodName = "k8sPodName"

	// K8sSvcName is the name of a K8s service
	K8sSvcName = "k8sSvcName"

	// K8sSvcType is the k8s service type (e.g. NodePort, Loadbalancer etc.)
	K8sSvcType = "k8sSvcType"

	// K8sEndpointName is the k8s name for a k8s Endpoint (not a cilium Endpoint)
	K8sEndpointName = "k8sEndpointName"

	// K8sNamespace is the namespace something belongs to
	K8sNamespace = "k8sNamespace"

	// K8sIdentityAnnotation is a k8s non-identifying annotations on k8s objects
	K8sIdentityAnnotation = "k8sIdentityAnnotation"

	// K8sNetworkPolicy is a k8s NetworkPolicy object (not a CiliumNetworkObject, above).
	K8sNetworkPolicy = "k8sNetworkPolicy"

	// K8sNetworkPolicyName is the name of a K8sPolicyObject
	K8sNetworkPolicyName = "k8sNetworkPolicyName"

	// K8sIngress is a k8s Ingress service object
	K8sIngress = "k8sIngress"

	// K8sIngressName is the name of a K8sIngress
	K8sIngressName = "k8sIngressName"

	// K8sAPIVersion is the version of the k8s API an object has
	K8sAPIVersion = "k8sApiVersion"

	// Attempt is the attempt number if an operation is attempted multiple times
	Attempt = "attempt"

	// TrafficDirection represents the directionality of traffic with respect
	// to an endpoint.
	TrafficDirection = "trafficDirection"

	// Modification represents a type of state change operation (insert, delete,
	// upsert, etc.).
	Modification = "modification"

	// BPFMapName is the name of a BPF map.
	BPFMapName = "bpfMapName"

	// BPFHeaderHash is the hash of the BPF header.
	BPFHeaderfileHash = "bpfHeaderfileHash"

	// BPFMapPath is the path of a BPF map in the filesystem.
	BPFMapPath = "bpfMapPath"

	// BPFMapFD is the file descriptor for a BPF map.
	BPFMapFD = "bpfMapFileDescriptor"

	// ThreadID is the Envoy thread ID.
	ThreadID = "threadID"

	// Reason is a human readable string describing why an operation was
	// performed
	Reason = "reason"

	// Debug is a boolean value for whether debug is set or not.
	Debug = "debug"

	// PID is an integer value for the process identifier of a process.
	PID = "pid"

	// PIDFile is a string value for the path to a file containing a PID.
	PIDFile = "pidfile"

	// Probe is the name of a status probe.
	Probe = "probe"

	// Key is the identity of the encryption key
	Key = "key"

	// SysParamName is the name of the kernel parameter (sysctl)
	SysParamName = "sysParamName"

	// SysParamValue is the value of the kernel parameter (sysctl)
	SysParamValue = "sysParamValue"
)
