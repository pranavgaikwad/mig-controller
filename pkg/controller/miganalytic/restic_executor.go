/*
Copyright 2021 Red Hat Inc.

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

package miganalytic

import (
	"context"
	"fmt"

	"github.com/konveyor/mig-controller/pkg/compat"
	"github.com/konveyor/mig-controller/pkg/pods"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// ResticPodLabelKey is the key of the label used to discover Restic pod
	ResticPodLabelKey = "name"
	// ResticPodLabelValue is the value of the label used to discover Restic pod
	ResticPodLabelValue = "restic"
)

// ResticDFCommandExecutor uses Restic pod to run DF command
type ResticDFCommandExecutor struct {
	// Namespace is the ns in which df command pods are present
	Namespace string
	// Client to interact with restic pods
	Client compat.Client
	// ResticPodReferences is a local cache of Restic pods found per node
	ResticPodReferences map[string]*corev1.Pod
	// BaseUnit is the base unit this executor uses for conversions
	BaseUnit resource.Format
}

// DF runs df command given a podRef and a list of volumes
// any errors running the df command are suppressed here
// DFCommand.ExitCode field should be used to determine failure
func (r *ResticDFCommandExecutor) DF(podRef *corev1.Pod, persistentVolumes []MigAnalyticPersistentVolumeDetails) DFCommand {
	dfCmd := DFCommand{
		BaseLocation: "/host_pods",
		BlockSize:    DecimalSIMega,
		ExitCode:     0,
		StdOut:       "",
		StdErr:       "",
	}
	cmdString := dfCmd.PrepareDFCommand(persistentVolumes)
	restCfg := r.Client.RestConfig()
	podCommand := pods.PodCommand{
		Pod:     podRef,
		RestCfg: restCfg,
		Args:    cmdString,
	}
	err := podCommand.Run()
	if err != nil {
		log.Info(
			fmt.Sprintf("[PG] found error %v", podCommand.Err.String()))
		dfCmd.ExitCode = 1
	}
	dfCmd.StdErr = podCommand.Err.String()
	dfCmd.StdOut = podCommand.Out.String()
	return dfCmd
}

// finds a Restic pod running on a particular node
func (r *ResticDFCommandExecutor) getResticPodForNode(nodeName string) *corev1.Pod {
	if podRef, exists := r.ResticPodReferences[nodeName]; exists {
		return podRef
	}
	return nil
}

// loads restic pod information in-memory
func (r *ResticDFCommandExecutor) loadResticPodReferences() error {
	if r.ResticPodReferences == nil {
		r.ResticPodReferences = make(map[string]*corev1.Pod)
	}
	resticPodList := corev1.PodList{}
	labelSelector := client.InNamespace(r.Namespace).MatchingLabels(
		map[string]string{
			ResticPodLabelKey: ResticPodLabelValue})
	// NOTE: +1 List call
	err := r.Client.List(context.TODO(), labelSelector, &resticPodList)
	if err != nil {
		return err
	}
	for i := range resticPodList.Items {
		if resticPodList.Items[i].Spec.NodeName != "" {
			r.ResticPodReferences[resticPodList.Items[i].Spec.NodeName] = &resticPodList.Items[i]
		}
	}
	return nil
}

// Execute runs Df command for all known PVCs and updates MigAnalytic
func (r *ResticDFCommandExecutor) Execute(pvNodeMap map[string][]MigAnalyticPersistentVolumeDetails) ([]PersistentVolumeUsageData, error) {
	gatheredData := []PersistentVolumeUsageData{}
	err := r.loadResticPodReferences()
	if err != nil {
		return gatheredData, err
	}
	// TODO: run this concurrently
	for node, pvcs := range pvNodeMap {
		resticPodRef := r.getResticPodForNode(node)
		// if no Restic pod is found for this node, all PVCs on this node are skipped
		if resticPodRef == nil {
			continue
		}
		// run bulk Restic command for all pvcs present on this node
		cmdOutput := r.DF(resticPodRef, pvcs)
		// split bulk Restic output into output per PV
		for _, pvc := range pvcs {
			gatheredData = append(gatheredData,
				cmdOutput.GetPVUsage(pvc.VolumeName, pvc.PodUID))
		}

	}
	return gatheredData, nil
}
