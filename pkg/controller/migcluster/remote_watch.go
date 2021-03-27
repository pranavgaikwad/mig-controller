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

package migcluster

import (
	"strconv"

	liberr "github.com/konveyor/controller/pkg/error"
	migapi "github.com/konveyor/mig-controller/pkg/apis/migration/v1alpha1"
	"github.com/konveyor/mig-controller/pkg/remote"
	kapi "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// StartRemoteWatch will configure a new RemoteWatcher manager + controller to monitor Velero
// events on a remote cluster. A GenericEvent channel will be configured to funnel events from
// the RemoteWatcher controller to the MigCluster controller.
func StartRemoteWatch(r *ReconcileMigCluster, config remote.ManagerConfig) error {
	remoteWatchMap := remote.GetWatchMap()

	mgr, err := manager.New(config.RemoteRestConfig, manager.Options{Scheme: config.Scheme})
	if err != nil {
		return liberr.Wrap(err)
	}

	sigStopChan := make(chan struct{})
	log.Info("[rWatch] Starting manager")
	go mgr.Start(sigStopChan)

	// Indexes
	indexer := mgr.GetFieldIndexer()

	// Plan
	err = indexer.IndexField(
		&migapi.MigPlan{},
		migapi.ClosedIndexField,
		func(rawObj runtime.Object) []string {
			p, cast := rawObj.(*migapi.MigPlan)
			if !cast {
				return nil
			}
			return []string{
				strconv.FormatBool(p.Spec.Closed),
			}
		})
	if err != nil {
		return err
	}
	// Pod
	err = indexer.IndexField(
		&kapi.Pod{},
		"status.phase",
		func(rawObj runtime.Object) []string {
			p, cast := rawObj.(*kapi.Pod)
			if !cast {
				return nil
			}
			return []string{
				string(p.Status.Phase),
			}
		})
	if err != nil {
		return err
	}

	log.Info("[rWatch] Manager started")
	// TODO: provide a way to dynamically change where events are being forwarded to (multiple controllers)
	// Create remoteWatchCluster tracking obj and attach reference to parent object so we don't create extra
	remoteWatchCluster := &remote.WatchCluster{RemoteManager: mgr, StopChannel: sigStopChan}

	// MigClusters have a 1:1 association with a RemoteWatchCluster, so we will store the mapping
	// to avoid creating duplicate remote managers in the future.
	remoteWatchMap.Set(config.ParentNsName, remoteWatchCluster)

	return nil
}

// StopRemoteWatch will close a remote watch's stop channel
// and delete it from the remote watch map.
func StopRemoteWatch(nsName types.NamespacedName) {
	remoteWatchMap := remote.GetWatchMap()
	remoteWatchCluster := remoteWatchMap.Get(nsName)
	if remoteWatchCluster != nil {
		close(remoteWatchCluster.StopChannel)
		remoteWatchMap.Delete(nsName)
	}
}
