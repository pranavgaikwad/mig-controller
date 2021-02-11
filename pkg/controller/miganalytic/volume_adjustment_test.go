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
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/types"
)

func getSamplePersistentVolume(name string, ns string, podUID string, volName string) MigAnalyticPersistentVolumeDetails {
	return MigAnalyticPersistentVolumeDetails{
		Name:       name,
		Namespace:  ns,
		PodUID:     types.UID(podUID),
		VolumeName: volName,
	}
}

var pvd1 = getSamplePersistentVolume("pvc-1", "test-ns", "0000", "pvc-1-vol")
var pvd2 = getSamplePersistentVolume("pvc-2", "test-ns", "0000", "pvc-2-vol")

func TestDFCommand_PrepareDFCommand(t *testing.T) {
	type fields struct {
		StdOut       string
		StdErr       string
		ExitCode     int
		BlockSize    DFBaseUnit
		BaseLocation string
	}
	type args struct {
		pvcs []MigAnalyticPersistentVolumeDetails
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name: "given an empty volume, should return command without args",
			args: args{
				pvcs: []MigAnalyticPersistentVolumeDetails{},
			},
			fields: fields{
				BaseLocation: "/host_pods",
				BlockSize:    DecimalSIMega,
			},
			want: []string{"/bin/bash", "-c", "df -BMB "},
		},
		{
			name: "given two volumes, should return command with two vols as args",
			args: args{
				pvcs: []MigAnalyticPersistentVolumeDetails{pvd1, pvd2},
			},
			fields: fields{
				BaseLocation: "/host_pods",
				BlockSize:    DecimalSIMega,
			},
			want: []string{"/bin/bash", "-c", "df -BMB /host_pods/0000/volumes/*/pvc-1-vol /host_pods/0000/volumes/*/pvc-2-vol"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &DFCommand{
				StdOut:       tt.fields.StdOut,
				StdErr:       tt.fields.StdErr,
				ExitCode:     tt.fields.ExitCode,
				BlockSize:    tt.fields.BlockSize,
				BaseLocation: tt.fields.BaseLocation,
			}
			if got := cmd.PrepareDFCommand(tt.args.pvcs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DFCommand.PrepareDFCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

var testDFStderr = `
df: /host_pods/280f7572-6590-11eb-b436-0a916cc7c396/volumes/kubernetes.io~empty-dir/tmp-volum: No such file or directory
df: /l: No such file or directory
`
var testDFStdout = `
Filesystem     1M-blocks  Used Available Use% Mounted on
tmpfs              7942M    1M     7942M   1% /host_pods/280f7572-6590-11eb-b436-0a916cc7c396/volumes/kubernetes.io~secret/sock-shop-token-pn8n9
shm                  64M    0M       64M   0% /dev/shm
tmpfs              7942M    1M     7942M   1% /credentials
/dev/xvda2        51188M 7613M    43576M  15% /host_pods
`

func TestDFCommand_GetPVUsage(t *testing.T) {
	type fields struct {
		StdOut       string
		StdErr       string
		ExitCode     int
		BlockSize    DFBaseUnit
		BaseLocation string
	}
	type args struct {
		volName string
		podUID  types.UID
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		wantPv PersistentVolumeUsageData
	}{
		{
			name: "given a volume that we know exists, correct pv usage info should be returned",
			fields: fields{
				BlockSize:    BinarySIMega,
				BaseLocation: "/host_pods",
				StdOut:       testDFStdout,
				StdErr:       testDFStderr,
			},
			args: args{
				volName: "sock-shop-token-pn8n9",
				podUID:  types.UID("280f7572-6590-11eb-b436-0a916cc7c396"),
			},
			wantPv: PersistentVolumeUsageData{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &DFCommand{
				StdOut:       tt.fields.StdOut,
				StdErr:       tt.fields.StdErr,
				ExitCode:     tt.fields.ExitCode,
				BlockSize:    tt.fields.BlockSize,
				BaseLocation: tt.fields.BaseLocation,
			}
			if gotPv := cmd.GetPVUsage(tt.args.volName, tt.args.podUID); !reflect.DeepEqual(gotPv, tt.wantPv) {
				t.Errorf("DFCommand.GetPVUsage() = %v, want %v", gotPv, tt.wantPv)
			}
		})
	}
}

type testQuantity struct {
	dfQuantity  string
	k8sQuantity string
}

var testQuantities = []testQuantity{
	{
		dfQuantity:  "1090MB",
		k8sQuantity: "1090M",
	},
	{
		dfQuantity:  "1090M",
		k8sQuantity: "1090Mi",
	},
	{
		dfQuantity:  "1GB",
		k8sQuantity: "1G",
	},
	{
		dfQuantity:  "1000G",
		k8sQuantity: "1000Gi",
	},
}

func TestDFCommand_convertDFQuantityToKubernetesResource(t *testing.T) {
	type fields struct {
		StdOut       string
		StdErr       string
		ExitCode     int
		BlockSize    DFBaseUnit
		BaseLocation string
	}
	type args struct {
		quantity string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "given a valid decimal SI quantity, should return valid k8s resource.Quantity",
			fields: fields{
				BlockSize: DecimalSIMega,
			},
			args: args{
				quantity: testQuantities[0].dfQuantity,
			},
			want:    testQuantities[0].k8sQuantity,
			wantErr: false,
		},
		{
			name: "given an invalid pair of SI quantity and block size unit, should return error",
			fields: fields{
				BlockSize: DecimalSIGiga,
			},
			args: args{
				quantity: testQuantities[0].dfQuantity,
			},
			want:    "0",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &DFCommand{
				StdOut:       tt.fields.StdOut,
				StdErr:       tt.fields.StdErr,
				ExitCode:     tt.fields.ExitCode,
				BlockSize:    tt.fields.BlockSize,
				BaseLocation: tt.fields.BaseLocation,
			}
			got, err := cmd.convertDFQuantityToKubernetesResource(tt.args.quantity)
			if (err != nil) != tt.wantErr {
				t.Errorf("DFCommand.convertDFQuantityToKubernetesResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.String(), tt.want) {
				t.Errorf("DFCommand.convertDFQuantityToKubernetesResource() = %v, want %v", got.String(), tt.want)
			}
		})
	}
}
