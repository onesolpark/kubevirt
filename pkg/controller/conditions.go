/*
 * This file is part of the KubeVirt project
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Copyright 2019 Red Hat, Inc.
 *
 */
package controller

import (
	k8sv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "kubevirt.io/client-go/api/v1"
)

type VirtualMachineConditionManager struct {
}

func (d *VirtualMachineConditionManager) CheckFailure(vmi *v1.VirtualMachineInstance, syncErr error, reason string) (changed bool) {
	if syncErr != nil && !d.HasCondition(vmi, v1.VirtualMachineInstanceSynchronized) {
		vmi.Status.Conditions = append(vmi.Status.Conditions, v1.VirtualMachineInstanceCondition{
			Type:               v1.VirtualMachineInstanceSynchronized,
			Reason:             reason,
			Message:            syncErr.Error(),
			LastTransitionTime: metav1.Now(),
			Status:             k8sv1.ConditionFalse,
		})
		return true
	} else if syncErr == nil && d.HasCondition(vmi, v1.VirtualMachineInstanceSynchronized) {
		d.RemoveCondition(vmi, v1.VirtualMachineInstanceSynchronized)
		return true
	}
	return false
}

func (d *VirtualMachineConditionManager) HasCondition(vmi *v1.VirtualMachineInstance, cond v1.VirtualMachineInstanceConditionType) bool {
	if vmi == nil {
		return false
	}
	for _, c := range vmi.Status.Conditions {
		if c.Type == cond {
			return true
		}
	}
	return false
}

func (d *VirtualMachineConditionManager) HasConditionWithStatus(vmi *v1.VirtualMachineInstance, cond v1.VirtualMachineInstanceConditionType, status k8sv1.ConditionStatus) bool {
	for _, c := range vmi.Status.Conditions {
		if c.Type == cond {
			if c.Status == status {
				return true
			}
			return false
		}
	}
	return false
}

func (d *VirtualMachineConditionManager) RemoveCondition(vmi *v1.VirtualMachineInstance, cond v1.VirtualMachineInstanceConditionType) {
	var conds []v1.VirtualMachineInstanceCondition
	for _, c := range vmi.Status.Conditions {
		if c.Type == cond {
			continue
		}
		conds = append(conds, c)
	}
	vmi.Status.Conditions = conds
}

// AddPodCondition add pod condition to the VM.
func (d *VirtualMachineConditionManager) AddPodCondition(vmi *v1.VirtualMachineInstance, cond *k8sv1.PodCondition) {
	if !d.HasCondition(vmi, v1.VirtualMachineInstanceConditionType(cond.Type)) {
		vmi.Status.Conditions = append(vmi.Status.Conditions, v1.VirtualMachineInstanceCondition{
			LastProbeTime:      cond.LastProbeTime,
			LastTransitionTime: cond.LastTransitionTime,
			Message:            cond.Message,
			Reason:             cond.Reason,
			Status:             cond.Status,
			Type:               v1.VirtualMachineInstanceConditionType(cond.Type),
		})
	}
}

func (d *VirtualMachineConditionManager) PodHasCondition(pod *k8sv1.Pod, conditionType k8sv1.PodConditionType, status k8sv1.ConditionStatus) bool {
	for _, cond := range pod.Status.Conditions {
		if cond.Type == conditionType {
			if cond.Status == status {
				return true
			} else {
				return false
			}
		}
	}
	return false
}

func (d *VirtualMachineConditionManager) GetPodConditionWithStatus(pod *k8sv1.Pod, conditionType k8sv1.PodConditionType, status k8sv1.ConditionStatus) *k8sv1.PodCondition {
	for _, cond := range pod.Status.Conditions {
		if cond.Type == conditionType {
			if cond.Status == status {
				return &cond
			} else {
				return nil
			}
		}
	}
	return nil
}

func (d *VirtualMachineConditionManager) GetPodCondition(pod *k8sv1.Pod, conditionType k8sv1.PodConditionType) *k8sv1.PodCondition {
	for _, cond := range pod.Status.Conditions {
		if cond.Type == conditionType {
			return &cond
		}
	}
	return nil
}

func NewVirtualMachineInstanceConditionManager() *VirtualMachineConditionManager {
	return &VirtualMachineConditionManager{}
}

type VirtualMachineInstanceMigrationConditionManager struct {
}

func (d *VirtualMachineInstanceMigrationConditionManager) HasCondition(migration *v1.VirtualMachineInstanceMigration, cond v1.VirtualMachineInstanceMigrationConditionType) bool {
	for _, c := range migration.Status.Conditions {
		if c.Type == cond {
			return true
		}
	}
	return false
}

func (d *VirtualMachineInstanceMigrationConditionManager) HasConditionWithStatus(migration *v1.VirtualMachineInstanceMigration, cond v1.VirtualMachineInstanceMigrationConditionType, status k8sv1.ConditionStatus) bool {
	for _, c := range migration.Status.Conditions {
		if c.Type == cond {
			if c.Status == status {
				return true
			}
			return false
		}
	}
	return false
}

func (d *VirtualMachineInstanceMigrationConditionManager) RemoveCondition(migration *v1.VirtualMachineInstanceMigration, cond v1.VirtualMachineInstanceMigrationConditionType) {
	var conds []v1.VirtualMachineInstanceMigrationCondition
	for _, c := range migration.Status.Conditions {
		if c.Type == cond {
			continue
		}
		conds = append(conds, c)
	}
	migration.Status.Conditions = conds
}
func NewVirtualMachineInstanceMigrationConditionManager() *VirtualMachineInstanceMigrationConditionManager {
	return &VirtualMachineInstanceMigrationConditionManager{}
}
