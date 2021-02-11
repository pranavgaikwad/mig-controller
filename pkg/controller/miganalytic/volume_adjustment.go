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
	"fmt"
	"regexp"
	"strconv"
	"strings"

	migapi "github.com/konveyor/mig-controller/pkg/apis/migration/v1alpha1"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// DFCommandExecutor defines an executor responsible for running DF
type DFCommandExecutor interface {
	Execute(map[string][]MigAnalyticPersistentVolumeDetails) ([]PersistentVolumeUsageData, error)
}

// PersistentVolumeAdjuster defines volume adjustment context
type PersistentVolumeAdjuster struct {
	Owner      *migapi.MigAnalytic
	Client     client.Client
	DFExecutor DFCommandExecutor
}

// DFBaseUnit defines supported block sizes for df command
type DFBaseUnit string

// Base units used for df command
const (
	DecimalSIGiga = DFBaseUnit("GB")
	DecimalSIMega = DFBaseUnit("MB")
	BinarySIMega  = DFBaseUnit("M")
	BinarySIGiga  = DFBaseUnit("G")
)

// DFCommand represent a df command
type DFCommand struct {
	// stdout from df
	StdOut string
	// stderr from df
	StdErr string
	// exit code from df
	ExitCode int
	// Base unit used for df
	BlockSize DFBaseUnit
	// BaseLocation defines path where volumes can be found
	BaseLocation string
}

// DFDiskPath defines format of expected path of the volume present on Pod
const DFDiskPath = "%s/%s/volumes/*/%s"

// PersistentVolumeUsageData defines structured output of df per PV
type PersistentVolumeUsageData struct {
	Name            string
	Namespace       string
	UsagePercentage int
	TotalSize       resource.Quantity
	IsError         bool
}

// convertDFQuantityToKubernetesResource converts a quantity present in df output to k8s.Resource
func (cmd *DFCommand) convertDFQuantityToKubernetesResource(quantity string) (resource.Quantity, error) {
	var parsedQuantity resource.Quantity
	unitMatcher, _ := regexp.Compile(
		fmt.Sprintf(
			"(\\d+)%s", cmd.BlockSize))
	matched := unitMatcher.FindStringSubmatch(quantity)
	if len(matched) != 2 {
		return parsedQuantity, errors.Errorf("Invalid quantity or block size unit")
	}
	switch cmd.BlockSize {
	case DecimalSIGiga, DecimalSIMega:
		quantity = strings.ReplaceAll(quantity, "B", "")
		break
	case BinarySIGiga, BinarySIMega:
		quantity = fmt.Sprintf("%si", quantity)
		break
	}
	return resource.ParseQuantity(quantity)
}

// GetPVUsage finds output for a particular pv in stdout of DFCommand
// only works on commands created by DFCommand
func (cmd *DFCommand) GetPVUsage(volName string, podUID types.UID) (pv PersistentVolumeUsageData) {
	var err error
	stdOutLines, stdErrLines := strings.Split(cmd.StdOut, "\n"), strings.Split(cmd.StdErr, "\n")
	lineMatcher, _ := regexp.Compile(
		fmt.Sprintf(
			strings.Replace(DFDiskPath, "*", ".*", 1),
			cmd.BaseLocation,
			podUID,
			volName))
	percentageMatcher, _ := regexp.Compile("(\\d+)%")
	for _, line := range stdOutLines {
		if lineMatcher.MatchString(line) {
			cols := strings.Fields(line)
			if len(cols) != 6 {
				pv.IsError = true
				return
			}
			pv.TotalSize, err = cmd.convertDFQuantityToKubernetesResource(cols[1])
			pv.IsError = (err != nil)
			matched := percentageMatcher.FindStringSubmatch(cols[4])
			if len(matched) > 1 {
				pv.UsagePercentage, err = strconv.Atoi(matched[1])
				pv.IsError = (err != nil)
			}
			return
		}
	}
	for _, line := range stdErrLines {
		if lineMatcher.MatchString(line) {
			pv.IsError = true
			return
		}
	}
	return
}

// PrepareDFCommand given a list of volumes, creates a bulk df command for all volumes
func (cmd *DFCommand) PrepareDFCommand(pvcs []MigAnalyticPersistentVolumeDetails) []string {
	command := []string{
		"/bin/bash",
		"-c",
	}
	// location = baseLocation + podUID + volumeName
	volPaths := []string{}
	for _, pvc := range pvcs {
		volPaths = append(volPaths,
			fmt.Sprintf(
				DFDiskPath,
				cmd.BaseLocation,
				pvc.PodUID,
				pvc.VolumeName))
	}
	return append(command, fmt.Sprintf("df -B%s %s", cmd.BlockSize, strings.Join(volPaths, " ")))
}

// Run runs df_executor, uses df output to calculate proposed volume, updates owner.status with proposed sizes
func (pva *PersistentVolumeAdjuster) Run(pvNodeMap map[string][]MigAnalyticPersistentVolumeDetails) error {
	_, err := pva.DFExecutor.Execute(pvNodeMap)
	if err != nil {
		return err
	}
	return nil
}
