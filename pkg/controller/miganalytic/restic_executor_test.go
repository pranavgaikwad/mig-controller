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
	migapi "github.com/konveyor/mig-controller/pkg/apis/migration/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func getResticPod(ns string, name string, nodeName string) corev1.Pod {
	return corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: corev1.PodSpec{
			NodeName: nodeName,
		},
	}
}

var resticPod1 = getResticPod("openshift-migration", "restic-00", "node1.migration.internal")
var resticPod2 = getResticPod("openshift-migration", "restic-01", "node2.migration.internal")

var testAnalytic = migapi.MigAnalytic{
	ObjectMeta: metav1.ObjectMeta{
		Namespace: "openshift-migration",
	},
}
