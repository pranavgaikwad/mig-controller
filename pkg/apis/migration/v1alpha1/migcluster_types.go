/*
Copyright 2019 Red Hat Inc.

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

package v1alpha1

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	liberr "github.com/konveyor/controller/pkg/error"
	pvdr "github.com/konveyor/mig-controller/pkg/cloudprovider"
	"github.com/konveyor/mig-controller/pkg/compat"
	"github.com/konveyor/mig-controller/pkg/remote"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	ocapi "github.com/openshift/api/apps/v1"
	imgapi "github.com/openshift/api/image/v1"
	"github.com/openshift/library-go/pkg/image/reference"
	velero "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	kapi "k8s.io/api/core/v1"
	storageapi "k8s.io/api/storage/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// SA secret keys.
const (
	SaToken = "saToken"
)

// migration-cluster-config configmap
const (
	ClusterConfigMapName  = "migration-cluster-config"
	RegistryImageKey      = "REGISTRY_IMAGE"
	StagePodImageKey      = "STAGE_IMAGE"
	RsyncTransferImageKey = "RSYNC_TRANSFER_IMAGE"
	ClusterSubdomainKey   = "CLUSTER_SUBDOMAIN"
	OperatorVersionKey    = "OPERATOR_VERSION"
)

// MigClusterSpec defines the desired state of MigCluster
type MigClusterSpec struct {
	// Specifies if the cluster is host (where the controller is installed) or not. This is a required field.
	IsHostCluster bool `json:"isHostCluster"`

	// Stores the url of the remote cluster. The field is only required for the source cluster object.
	URL string `json:"url,omitempty"`

	ServiceAccountSecretRef *kapi.ObjectReference `json:"serviceAccountSecretRef,omitempty"`

	// If the migcluster needs SSL verification for connections a user can supply a custom CA bundle. This field is required only when spec.Insecure is set false
	CABundle []byte `json:"caBundle,omitempty"`

	// For azure clusters -- it's the resource group that in-cluster volumes use.
	AzureResourceGroup string `json:"azureResourceGroup,omitempty"`

	// If set false, user will need to provide CA bundle for TLS connection to the remote cluster.
	Insecure bool `json:"insecure,omitempty"`

	// An override setting to tell the controller that the source cluster restic needs to be restarted after stage pod creation.
	RestartRestic *bool `json:"restartRestic,omitempty"`

	// If set True, forces the controller to run a full suite of validations on migcluster.
	Refresh bool `json:"refresh,omitempty"`

	// Stores the path of registry route when using direct migration.
	ExposedRegistryPath string `json:"exposedRegistryPath,omitempty"`
}

// MigClusterStatus defines the observed state of MigCluster
type MigClusterStatus struct {
	Conditions      `json:","`
	ObservedDigest  string `json:"observedDigest,omitempty"`
	RegistryPath    string `json:"registryPath,omitempty"`
	OperatorVersion string `json:"operatorVersion,omitempty"`
}

var clientMap compatClientCache

type compatClientCache struct {
	cMap  map[types.UID]compat.Client
	mutex sync.RWMutex
}

func (cm *compatClientCache) init() {
	if cm.cMap == nil {
		cm.cMap = make(map[types.UID]compat.Client)
	}
}

func (cm *compatClientCache) Get(key types.UID) (compat.Client, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	val, found := cm.cMap[key]
	return val, found
}

func (cm *compatClientCache) Set(key types.UID, val compat.Client) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	cm.cMap[key] = val
}

func init() {
	clientMap = compatClientCache{}
	clientMap.init()
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MigCluster is the Schema for the migclusters API
// +k8s:openapi-gen=true
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="URL",type=string,JSONPath=".spec.url"
// +kubebuilder:printcolumn:name="Host",type=boolean,JSONPath=".spec.isHostCluster"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type MigCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MigClusterSpec   `json:"spec,omitempty"`
	Status MigClusterStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MigClusterList contains a list of MigCluster
type MigClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MigCluster `json:"items"`
}

// StorageClass is an available storage class in the cluster
// Name - the storage class name
// Provisioner - the dynamic provisioner for the storage class
// Default - whether or not this storage class is the default
// AccessModes - access modes supported by the dynamic provisioner
type StorageClass struct {
	Name        string                            `json:"name,omitempty"`
	Provisioner string                            `json:"provisioner,omitempty"`
	Default     bool                              `json:"default,omitempty"`
	AccessModes []kapi.PersistentVolumeAccessMode `json:"accessModes,omitempty" protobuf:"bytes,1,rep,name=accessModes,casttype=PersistentVolumeAccessMode"`
}

func init() {
	SchemeBuilder.Register(&MigCluster{}, &MigClusterList{})
}

// Get the service account secret.
// Returns `nil` when the reference cannot be resolved.
func (m *MigCluster) GetServiceAccountSecret(client k8sclient.Client) (*kapi.Secret, error) {
	return GetSecret(client, m.Spec.ServiceAccountSecretRef)
}

// // GetClient get a local or remote client using a MigCluster and an existing client
// func (m *MigCluster) GetClient(c k8sclient.Client) (compat.Client, error) {
// 	restConfig, err := m.BuildRestConfig(c)
// 	if err != nil {
// 		return nil, err
// 	}
// 	client, err := compat.NewClient(restConfig)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return client, nil
// }

// GetClient get a local or remote client using a MigCluster and an existing client
// GetClient get a local or remote client using a MigCluster and an existing client
func (m *MigCluster) GetClient(c k8sclient.Client) (compat.Client, error) {
	client, ok := clientMap.Get(m.UID)
	if ok {
		return client, nil
	}

	// Building a compat client requires both restConfig and a k8s client
	var rClient *k8sclient.Client

	// Retrieve cached client if it exists
	if m.Spec.IsHostCluster {
		rClient = &c
	} else {
		rwm := remote.GetWatchMap()
		remoteCluster := rwm.Get(types.NamespacedName{Namespace: m.Namespace, Name: m.Name})
		if remoteCluster != nil {
			cachedClient := remoteCluster.RemoteManager.GetClient()
			rClient = &cachedClient
		}
	}

	restConfig, err := m.BuildRestConfig(c)
	if err != nil {
		return nil, err
	}

	rwStarted := false
	// Build client without cache if remote watch hasn't been started yet
	if rClient == nil {
		rwStarted = true
		uncachedClient, err := k8sclient.New(
			restConfig,
			k8sclient.Options{
				Scheme: scheme.Scheme,
			})
		if err != nil {
			return nil, err
		}
		rClient = &uncachedClient
	}

	compatClient, err := compat.NewClient(restConfig, rClient)
	if err != nil {
		return nil, err
	}
	if rwStarted {
		clientMap.Set(m.UID, compatClient)
	}
	return compatClient, nil
}

func (m *MigCluster) GetClusterConfigMap(c k8sclient.Client) (*corev1.ConfigMap, error) {
	clusterConfig := &corev1.ConfigMap{}
	clusterConfigRef := types.NamespacedName{Name: ClusterConfigMapName, Namespace: VeleroNamespace}
	err := c.Get(context.TODO(), clusterConfigRef, clusterConfig)
	if err != nil {
		return nil, liberr.Wrap(err)
	}
	if clusterConfig.Data == nil {
		return nil, liberr.Wrap(errors.New("failed to find data in cluster configmap"))
	}
	return clusterConfig, nil
}

// GetRegistryImage gets a MigCluster specific registry image from ConfigMap
func (m *MigCluster) GetRegistryImage(c k8sclient.Client) (string, error) {
	clusterConfig, err := m.GetClusterConfigMap(c)
	if err != nil {
		return "", liberr.Wrap(err)
	}
	registryImage, ok := clusterConfig.Data[RegistryImageKey]
	if !ok {
		return "", liberr.Wrap(errors.Errorf("configmap key not found: %v", RegistryImageKey))
	}
	return registryImage, nil
}

// GetRsyncTransferImage gets a MigCluster specific rsync transfer image from ConfigMap
func (m *MigCluster) GetRsyncTransferImage(c k8sclient.Client) (string, error) {
	client, err := m.GetClient(c)
	if err != nil {
		return "", err
	}
	clusterConfig, err := m.GetClusterConfigMap(client)
	if err != nil {
		return "", liberr.Wrap(err)
	}
	rsyncImage, ok := clusterConfig.Data[RsyncTransferImageKey]
	if !ok {
		return "", liberr.Wrap(errors.Errorf("configmap key not found: %v", RsyncTransferImageKey))
	}
	return rsyncImage, nil
}

// GetClusterSubdomain gets a MigCluster specific subdomain value to be used for DVM routes
func (m *MigCluster) GetClusterSubdomain(c k8sclient.Client) (string, error) {
	client, err := m.GetClient(c)
	if err != nil {
		return "", err
	}
	clusterConfig, err := m.GetClusterConfigMap(client)
	if err != nil {
		return "", liberr.Wrap(err)
	}
	clusterSubdomain, ok := clusterConfig.Data[ClusterSubdomainKey]
	if !ok || clusterSubdomain == "" {
		return "", liberr.Wrap(errors.Errorf("configmap key not found: %v", ClusterSubdomainKey))
	}
	return clusterSubdomain, nil
}

// GetOperatorVersion retrieves the operator version from the respective controllers ConfigMap
func (m *MigCluster) GetOperatorVersion(c k8sclient.Client) (string, error) {
	clusterConfig, err := m.GetClusterConfigMap(c)
	if err != nil {
		return "", liberr.Wrap(err)
	}

	operatorVersion, ok := clusterConfig.Data[OperatorVersionKey]
	if !ok {
		return "", liberr.Wrap(errors.Errorf("configmap key %v not found in configmap %v/%v for migcluster %v/%v",
			OperatorVersionKey, clusterConfig.Namespace, clusterConfig.Name, m.Namespace, m.Name))
	}

	return operatorVersion, nil
}

// Test the connection settings by building a client.
func (m *MigCluster) TestConnection(c k8sclient.Client, timeout time.Duration) error {
	if m.Spec.IsHostCluster {
		return nil
	}
	restConfig, err := m.BuildRestConfig(c)
	if err != nil {
		return err
	}
	restConfig.Timeout = timeout
	_, err = k8sclient.New(restConfig, k8sclient.Options{Scheme: scheme.Scheme})
	if err != nil {
		return err
	}

	return nil
}

// Build a REST configuration.
func (m *MigCluster) BuildRestConfig(c k8sclient.Client) (*rest.Config, error) {
	if m.Spec.IsHostCluster {
		return config.GetConfig()
	}
	secret, err := GetSecret(c, m.Spec.ServiceAccountSecretRef)
	if err != nil {
		return nil, err
	}
	if secret == nil {
		return nil, errors.Errorf("Service Account Secret not found for %v", m.Name)
	}
	var tlsClientConfig rest.TLSClientConfig
	if m.Spec.Insecure {
		tlsClientConfig = rest.TLSClientConfig{Insecure: true}
	} else {
		tlsClientConfig = rest.TLSClientConfig{Insecure: false, CAData: m.Spec.CABundle}
	}
	restConfig := &rest.Config{
		Host:            m.Spec.URL,
		BearerToken:     string(secret.Data[SaToken]),
		TLSClientConfig: tlsClientConfig,
		Burst:           1000,
		QPS:             100,
	}

	return restConfig, nil
}

// Test whether OPERATOR_VERSION in configmap on MigCluster matches status.OperatorVersion
func (m *MigCluster) OperatorVersionMatchesConfigmap(c k8sclient.Client) (bool, error) {
	clusterClient, err := m.GetClient(c)
	if err != nil {
		return false, liberr.Wrap(err)
	}
	operatorVersion, err := m.GetOperatorVersion(clusterClient)
	if operatorVersion == m.Status.OperatorVersion {
		return true, nil
	}
	return false, nil
}

// Delete resources on the cluster by label.
func (m *MigCluster) DeleteResources(client k8sclient.Client, labels map[string]string) error {
	client, err := m.GetClient(client)
	if err != nil {
		return err
	}
	if labels == nil {
		labels = map[string]string{PartOfLabel: Application}
	}

	options := k8sclient.MatchingLabels(labels)

	// Deployment
	dList := appv1.DeploymentList{}
	err = client.List(context.TODO(), options, &dList)
	if err != nil {
		return err
	}
	for _, r := range dList.Items {
		err = client.Delete(context.TODO(), &r, k8sclient.PropagationPolicy(metav1.DeletePropagationForeground))
		if err != nil && !k8serror.IsNotFound(err) {
			return err
		}
	}

	// DeploymentConfig
	dcList := ocapi.DeploymentConfigList{}
	err = client.List(context.TODO(), options, &dcList)
	if err != nil {
		return err
	}
	for _, r := range dcList.Items {
		err = client.Delete(context.TODO(), &r)
		if err != nil && !k8serror.IsNotFound(err) {
			return err
		}
	}

	// Service
	svList := kapi.ServiceList{}
	err = client.List(context.TODO(), options, &svList)
	if err != nil {
		return err
	}
	for _, r := range svList.Items {
		err = client.Delete(context.TODO(), &r)
		if err != nil && !k8serror.IsNotFound(err) {
			return err
		}
	}

	// Pod
	pList := kapi.PodList{}
	err = client.List(context.TODO(), options, &pList)
	if err != nil {
		return err
	}
	for _, r := range pList.Items {
		err = client.Delete(context.TODO(), &r)
		if err != nil && !k8serror.IsNotFound(err) {
			return err
		}
	}

	// Secret
	sList := kapi.SecretList{}
	err = client.List(context.TODO(), options, &sList)
	if err != nil {
		return err
	}
	for _, r := range sList.Items {
		err = client.Delete(context.TODO(), &r)
		if err != nil && !k8serror.IsNotFound(err) {
			return err
		}
	}

	// ImageStream
	iList := imgapi.ImageStreamList{}
	err = client.List(context.TODO(), options, &iList)
	if err != nil {
		return err
	}
	for _, r := range iList.Items {
		err = client.Delete(context.TODO(), &r)
		if err != nil && !k8serror.IsNotFound(err) {
			return err
		}
	}

	// Backup
	bList := velero.BackupList{}
	err = client.List(context.TODO(), options, &bList)
	if err != nil {
		return err
	}
	for _, r := range bList.Items {
		err = client.Delete(context.TODO(), &r)
		if err != nil && !k8serror.IsNotFound(err) {
			return err
		}
	}

	// Restore
	rList := velero.RestoreList{}
	err = client.List(context.TODO(), options, &rList)
	if err != nil {
		return err
	}
	for _, r := range rList.Items {
		err = client.Delete(context.TODO(), &r)
		if err != nil && !k8serror.IsNotFound(err) {
			return err
		}
	}

	// BSL
	bslList := velero.BackupStorageLocationList{}
	err = client.List(context.TODO(), options, &bslList)
	if err != nil {
		return err
	}
	for _, r := range bslList.Items {
		err = client.Delete(context.TODO(), &r)
		if err != nil && !k8serror.IsNotFound(err) {
			return err
		}
	}

	// VSL
	vslList := velero.VolumeSnapshotLocationList{}
	err = client.List(context.TODO(), options, &vslList)
	if err != nil {
		return err
	}
	for _, r := range vslList.Items {
		err = client.Delete(context.TODO(), &r)
		if err != nil && !k8serror.IsNotFound(err) {
			return err
		}
	}

	return nil
}

// Get the list StorageClasses in the format expected by PV discovery
func (r *MigCluster) GetStorageClasses(client k8sclient.Client) ([]StorageClass, error) {
	kubeStorageClasses, err := r.GetKubeStorageClasses(client)
	if err != nil {
		return nil, err
	}
	// Transform kube storage classes into format used in PV discovery
	var storageClasses []StorageClass
	for _, clusterStorageClass := range kubeStorageClasses {
		storageClass := StorageClass{
			Name:        clusterStorageClass.Name,
			Provisioner: clusterStorageClass.Provisioner,
			AccessModes: r.accessModesForProvisioner(clusterStorageClass.Provisioner),
		}
		if clusterStorageClass.Annotations != nil {
			storageClass.Default, _ = strconv.ParseBool(clusterStorageClass.Annotations["storageclass.kubernetes.io/is-default-class"])
		}
		storageClasses = append(storageClasses, storageClass)
	}
	return storageClasses, nil
}

// Gets the list of supported access modes for a provisioner
// TODO: allow the in-file mapping to be overridden by a configmap
func (r *MigCluster) accessModesForProvisioner(provisioner string) []kapi.PersistentVolumeAccessMode {
	for _, pModes := range accessModeList {
		if pModes.MatchBySuffix {
			if strings.HasSuffix(provisioner, pModes.Provisioner) {
				return pModes.AccessModes
			}
		} else if pModes.MatchByPrefix {
			if strings.HasPrefix(provisioner, pModes.Provisioner) {
				return pModes.AccessModes
			}
		} else {
			if pModes.Provisioner == provisioner {
				return pModes.AccessModes
			}
		}
	}

	// default value
	return []kapi.PersistentVolumeAccessMode{kapi.ReadWriteOnce}
}

type provisionerAccessModes struct {
	Provisioner   string
	MatchBySuffix bool
	MatchByPrefix bool
	AccessModes   []kapi.PersistentVolumeAccessMode
}

// Since the StorageClass API doesn't provide this information, the support list has been
// compiled from Kubernetes API docs. Most of the below comes from:
// https://kubernetes.io/docs/concepts/storage/persistent-volumes/#access-modes
// https://kubernetes.io/docs/concepts/storage/storage-classes/#provisioner
var accessModeList = []provisionerAccessModes{
	provisionerAccessModes{
		Provisioner: "kubernetes.io/aws-ebs",
		AccessModes: []kapi.PersistentVolumeAccessMode{kapi.ReadWriteOnce},
	},
	provisionerAccessModes{
		Provisioner: "kubernetes.io/azure-file",
		AccessModes: []kapi.PersistentVolumeAccessMode{kapi.ReadWriteOnce, kapi.ReadOnlyMany, kapi.ReadWriteMany},
	},
	provisionerAccessModes{
		Provisioner: "kubernetes.io/azure-disk",
		AccessModes: []kapi.PersistentVolumeAccessMode{kapi.ReadWriteOnce},
	},
	provisionerAccessModes{
		Provisioner: "kubernetes.io/cinder",
		AccessModes: []kapi.PersistentVolumeAccessMode{kapi.ReadWriteOnce},
	},
	// FC : {kapi.ReadWriteOnce, kapi.ReadOnlyMany},
	// Flexvolume : {kapi.ReadWriteOnce, kapi.ReadOnlyMany}, RWX?
	// Flocker . : {kapi.ReadWriteOnce},
	provisionerAccessModes{
		Provisioner: "kubernetes.io/gce-pd",
		AccessModes: []kapi.PersistentVolumeAccessMode{kapi.ReadWriteOnce, kapi.ReadOnlyMany},
	},
	provisionerAccessModes{
		Provisioner: "kubernetes.io/glusterfs",
		AccessModes: []kapi.PersistentVolumeAccessMode{kapi.ReadWriteOnce, kapi.ReadOnlyMany, kapi.ReadWriteMany},
	},
	provisionerAccessModes{
		Provisioner:   "gluster.org/glusterblock",
		MatchByPrefix: true,
		AccessModes:   []kapi.PersistentVolumeAccessMode{kapi.ReadWriteOnce, kapi.ReadOnlyMany},
	}, // verify glusterblock ROX
	// ISCSI : {kapi.ReadWriteOnce, kapi.ReadOnlyMany},
	provisionerAccessModes{
		Provisioner: "kubernetes.io/quobyte",
		AccessModes: []kapi.PersistentVolumeAccessMode{kapi.ReadWriteOnce, kapi.ReadOnlyMany, kapi.ReadWriteMany},
	},
	// NFS : {kapi.ReadWriteOnce, kapi.ReadOnlyMany, kapi.ReadWriteMany},
	provisionerAccessModes{
		Provisioner: "kubernetes.io/rbd",
		AccessModes: []kapi.PersistentVolumeAccessMode{kapi.ReadWriteOnce, kapi.ReadOnlyMany},
	},
	provisionerAccessModes{
		Provisioner: "kubernetes.io/vsphere-volume",
		AccessModes: []kapi.PersistentVolumeAccessMode{kapi.ReadWriteOnce},
	},
	provisionerAccessModes{
		Provisioner: "kubernetes.io/portworx-volume",
		AccessModes: []kapi.PersistentVolumeAccessMode{kapi.ReadWriteOnce, kapi.ReadWriteMany},
	},
	provisionerAccessModes{
		Provisioner: "kubernetes.io/scaleio",
		AccessModes: []kapi.PersistentVolumeAccessMode{kapi.ReadWriteOnce, kapi.ReadOnlyMany},
	},
	provisionerAccessModes{
		Provisioner: "kubernetes.io/storageos",
		AccessModes: []kapi.PersistentVolumeAccessMode{kapi.ReadWriteOnce},
	},
	// other CSI?
	// other OCP4?
	provisionerAccessModes{
		Provisioner:   "rbd.csi.ceph.com",
		MatchBySuffix: true,
		AccessModes:   []kapi.PersistentVolumeAccessMode{kapi.ReadWriteOnce},
	},
	provisionerAccessModes{
		Provisioner:   "cephfs.csi.ceph.com",
		MatchBySuffix: true,
		AccessModes:   []kapi.PersistentVolumeAccessMode{kapi.ReadWriteOnce, kapi.ReadOnlyMany, kapi.ReadWriteMany},
	},
	provisionerAccessModes{
		Provisioner: "netapp.io/trident",
		// Note: some backends won't support RWX
		AccessModes: []kapi.PersistentVolumeAccessMode{kapi.ReadWriteOnce, kapi.ReadOnlyMany, kapi.ReadWriteMany},
	},
}

// Get the list of k8s StorageClasses from the cluster.
func (r *MigCluster) GetKubeStorageClasses(client k8sclient.Client) ([]storageapi.StorageClass, error) {
	list := storageapi.StorageClassList{}
	err := client.List(
		context.TODO(),
		&k8sclient.ListOptions{},
		&list)
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}

func (r *MigCluster) UpdateProvider(provider pvdr.Provider) {
	switch provider.GetName() {
	case pvdr.Azure:
		p, cast := provider.(*pvdr.AzureProvider)
		if cast {
			p.ClusterResourceGroup = r.Spec.AzureResourceGroup
		}
	}
}

func (m *MigCluster) GetInternalRegistryPath(c k8sclient.Client) (string, error) {
	client, err := m.GetClient(c)
	if err != nil {
		return "", err
	}
	isList := imgapi.ImageStreamList{}
	err = client.List(
		context.TODO(),
		k8sclient.InNamespace("openshift"),
		&isList)
	if err == nil && len(isList.Items) > 0 {
		if value := isList.Items[0].Status.DockerImageRepository; len(value) > 0 {
			ref, err := reference.Parse(value)
			if err == nil {
				return ref.Registry, nil
			}
		}
	}
	if client.MajorVersion() != 1 {
		return "", errors.New(fmt.Sprintf("server version %v.%v not supported. Must be 1.x", client.MajorVersion(), client.MinorVersion()))
	}
	if client.MinorVersion() < 7 {
		return "", errors.New(fmt.Sprintf("Kubernetes version 1.%v not supported. Must be 1.7 or greater", client.MinorVersion()))
	} else if client.MinorVersion() <= 11 {
		registrySvc := kapi.Service{}
		err := client.Get(
			context.TODO(),
			k8sclient.ObjectKey{
				Namespace: "default",
				Name:      "docker-registry",
			},
			&registrySvc)
		if err != nil {
			// Return empty registry host but no error; registry not found
			return "", nil
		}
		internalRegistry := registrySvc.Spec.ClusterIP + ":" + strconv.Itoa(int(registrySvc.Spec.Ports[0].Port))
		return internalRegistry, nil
	} else {
		config := kapi.ConfigMap{}
		err := client.Get(
			context.TODO(),
			k8sclient.ObjectKey{
				Namespace: "openshift-apiserver",
				Name:      "config",
			},
			&config)
		if err != nil {
			return "", err
		}
		serverConfig := apiServerConfig{}
		err = json.Unmarshal([]byte(config.Data["config.yaml"]), &serverConfig)
		if err != nil {
			return "", err
		}
		internalRegistry := serverConfig.ImagePolicyConfig.InternalRegistryHostname
		if len(internalRegistry) == 0 {
			return "", nil
		}
		return internalRegistry, nil
	}
}

func (m *MigCluster) SetRegistryPath(c k8sclient.Client) error {
	newRegistryPath, err := m.GetRegistryPath(c)
	if err != nil {
		return err
	}
	m.Status.RegistryPath = newRegistryPath
	return nil
}

func (m *MigCluster) GetRegistryPath(c k8sclient.Client) (string, error) {
	if len(m.Spec.ExposedRegistryPath) > 0 {
		splitPath := strings.Split(m.Spec.ExposedRegistryPath, "//")
		if len(splitPath) == 2 {
			return splitPath[1], nil
		}
		return m.Spec.ExposedRegistryPath, nil
	} else if !m.Spec.IsHostCluster {
		// not host cluster and no path specified, return empty path
		return "", nil
	}
	return m.GetInternalRegistryPath(c)
}

// Pulls the operatorVersion from the migration-cluster-config configmap and
// loads it into MigCluster.Status.OperatorVersion
func (m *MigCluster) SetOperatorVersion(c k8sclient.Client) error {
	oldOperatorVersion := m.Status.OperatorVersion
	clusterClient, err := m.GetClient(c)
	if err != nil {
		return liberr.Wrap(err)
	}
	// Ignore error here. Missing configmap/key is already raised to user in validation,
	// we don't want reconcile to exit w/ error on MTC < 1.4.2. GetOperatorVersion will
	// return "" on error which is usable below.
	newOperatorVersion, _ := m.GetOperatorVersion(clusterClient)

	// When operator version changes, all other MigClusters will be updated
	// NOTE: In future if we need to support concurrent reconciles for MigCluster,
	//       we should remove this and write a watch with a predicate that would enqueue
	//       reconciles when other MigClusters change versions. This is simpler
	//       for now, but will break if MaxConcurrentReconciles is turned up,
	//       since other clusters would start reconciling before our change
	//       is written.
	if oldOperatorVersion != newOperatorVersion {
		clusterList, err := ListClusters(c)
		if err != nil {
			return liberr.Wrap(err)
		}
		// Update all other MigClusters to propagate version checks
		for _, cluster := range clusterList {
			if cluster.UID == m.UID {
				continue
			}
			if !cluster.Spec.Refresh {
				cluster.Spec.Refresh = true
				err := c.Update(context.Background(), &cluster)
				if err != nil {
					return liberr.Wrap(err)
				}
			}
		}
	}

	m.Status.OperatorVersion = newOperatorVersion

	return nil
}

type routingConfig struct {
	Subdomain string `json:"subdomain"`
}
type imagePolicyConfig struct {
	InternalRegistryHostname string `json:"internalRegistryHostname"`
}

// apiServerConfig stores configuration information about the current cluster
type apiServerConfig struct {
	ImagePolicyConfig imagePolicyConfig `json:"imagePolicyConfig"`
	RoutingConfig     routingConfig     `json:"routingConfig"`
}

// Get object reference for migcluster
func (r *MigCluster) GetObjectReference() *kapi.ObjectReference {
	return &kapi.ObjectReference{
		Name:      r.Name,
		Namespace: r.Namespace,
	}
}
