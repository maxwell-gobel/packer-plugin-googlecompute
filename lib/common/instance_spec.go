// Copyright IBM Corp. 2013, 2026
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"log"

	compute "google.golang.org/api/compute/v1"
)

// The builders below assemble instance-spec sub-objects that are identical
// between the classic single-zone path (RunInstance) and the bulk
// auto-zone path (RunInstanceInRegion). Keeping them in one place removes
// duplication and makes the assembly logic unit-testable. Pieces that differ
// between the two paths (machine type, disks, static-address region, and the
// API call) intentionally remain inline in each method.

func buildMetadataItems(c *InstanceConfig) []*compute.MetadataItems {
	metadata := make([]*compute.MetadataItems, 0, len(c.Metadata))
	for k, v := range c.Metadata {
		vCopy := v
		metadata = append(metadata, &compute.MetadataItems{Key: k, Value: &vCopy})
	}
	return metadata
}

func buildGuestAccelerators(c *InstanceConfig) []*compute.AcceleratorConfig {
	if c.AcceleratorCount <= 0 {
		return nil
	}
	return []*compute.AcceleratorConfig{
		{
			AcceleratorCount: c.AcceleratorCount,
			AcceleratorType:  c.AcceleratorType,
		},
	}
}

func buildServiceAccount(c *InstanceConfig) *compute.ServiceAccount {
	serviceAccount := &compute.ServiceAccount{}
	if !c.DisableDefaultServiceAccount {
		serviceAccount.Email = "default"
		serviceAccount.Scopes = c.Scopes
	}
	if c.ServiceAccountEmail != "" {
		serviceAccount.Email = c.ServiceAccountEmail
		serviceAccount.Scopes = c.Scopes
	}
	return serviceAccount
}

func buildScheduling(c *InstanceConfig) *compute.Scheduling {
	scheduling := &compute.Scheduling{
		OnHostMaintenance: c.OnHostMaintenance,
		Preemptible:       c.Preemptible,
	}
	if c.MaxRunDurationInSeconds > 0 {
		log.Printf("[DEBUG] setting max run duration to %d seconds", c.MaxRunDurationInSeconds)
		scheduling.MaxRunDuration = &compute.Duration{Seconds: c.MaxRunDurationInSeconds}
		log.Printf("[DEBUG] setting instance termination action to %s", c.InstanceTerminationAction)
		scheduling.InstanceTerminationAction = c.InstanceTerminationAction
	}
	if len(c.NodeAffinities) > 0 {
		scheduling.NodeAffinities = make([]*compute.SchedulingNodeAffinity, 0, len(c.NodeAffinities))
		for _, na := range c.NodeAffinities {
			scheduling.NodeAffinities = append(scheduling.NodeAffinities, na.ComputeType())
		}
	}
	return scheduling
}

func buildShieldedInstanceConfig(c *InstanceConfig) *compute.ShieldedInstanceConfig {
	if !c.EnableSecureBoot && !c.EnableVtpm && !c.EnableIntegrityMonitoring {
		return nil
	}
	return &compute.ShieldedInstanceConfig{
		EnableSecureBoot:          c.EnableSecureBoot,
		EnableVtpm:                c.EnableVtpm,
		EnableIntegrityMonitoring: c.EnableIntegrityMonitoring,
	}
}
