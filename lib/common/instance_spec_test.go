// Copyright IBM Corp. 2013, 2026
// SPDX-License-Identifier: MPL-2.0

package common

import "testing"

func TestBuildServiceAccount(t *testing.T) {
	// Explicit email + scopes.
	sa := buildServiceAccount(&InstanceConfig{
		ServiceAccountEmail: "me@example.com",
		Scopes:              []string{"scope-a"},
	})
	if sa.Email != "me@example.com" {
		t.Fatalf("expected explicit email, got %q", sa.Email)
	}
	if len(sa.Scopes) != 1 || sa.Scopes[0] != "scope-a" {
		t.Fatalf("expected scope-a, got %v", sa.Scopes)
	}

	// Default service account.
	sa = buildServiceAccount(&InstanceConfig{Scopes: []string{"scope-b"}})
	if sa.Email != "default" {
		t.Fatalf("expected default email, got %q", sa.Email)
	}

	// Disabled default service account, no explicit email -> empty.
	sa = buildServiceAccount(&InstanceConfig{DisableDefaultServiceAccount: true})
	if sa.Email != "" {
		t.Fatalf("expected empty email, got %q", sa.Email)
	}
}

func TestBuildGuestAccelerators(t *testing.T) {
	if got := buildGuestAccelerators(&InstanceConfig{}); got != nil {
		t.Fatalf("expected nil when count is 0, got %v", got)
	}
	got := buildGuestAccelerators(&InstanceConfig{AcceleratorCount: 2, AcceleratorType: "nvidia-tesla-t4"})
	if len(got) != 1 || got[0].AcceleratorCount != 2 || got[0].AcceleratorType != "nvidia-tesla-t4" {
		t.Fatalf("unexpected accelerators: %v", got)
	}
}

func TestBuildShieldedInstanceConfig(t *testing.T) {
	if got := buildShieldedInstanceConfig(&InstanceConfig{}); got != nil {
		t.Fatalf("expected nil when no shielded option set, got %v", got)
	}
	got := buildShieldedInstanceConfig(&InstanceConfig{EnableSecureBoot: true})
	if got == nil || !got.EnableSecureBoot {
		t.Fatalf("expected secure boot enabled config, got %v", got)
	}
}

func TestBuildScheduling(t *testing.T) {
	s := buildScheduling(&InstanceConfig{
		OnHostMaintenance:         "MIGRATE",
		Preemptible:               true,
		MaxRunDurationInSeconds:   3600,
		InstanceTerminationAction: "STOP",
	})
	if s.OnHostMaintenance != "MIGRATE" || !s.Preemptible {
		t.Fatalf("unexpected base scheduling: %#v", s)
	}
	if s.MaxRunDuration == nil || s.MaxRunDuration.Seconds != 3600 {
		t.Fatalf("expected MaxRunDuration 3600, got %#v", s.MaxRunDuration)
	}
	if s.InstanceTerminationAction != "STOP" {
		t.Fatalf("expected termination action STOP, got %q", s.InstanceTerminationAction)
	}

	// NodeAffinities branch: each NodeAffinity is mapped via ComputeType().
	withAffinities := buildScheduling(&InstanceConfig{
		NodeAffinities: []NodeAffinity{
			{Key: "compute.googleapis.com/node-group-name", Operator: "IN", Values: []string{"sole-tenant-group"}},
		},
	})
	if len(withAffinities.NodeAffinities) != 1 {
		t.Fatalf("expected 1 node affinity, got %d", len(withAffinities.NodeAffinities))
	}
	na := withAffinities.NodeAffinities[0]
	if na.Key != "compute.googleapis.com/node-group-name" {
		t.Fatalf("unexpected node affinity key: %q", na.Key)
	}
	if na.Operator != "IN" {
		t.Fatalf("unexpected node affinity operator: %q", na.Operator)
	}
	if len(na.Values) != 1 || na.Values[0] != "sole-tenant-group" {
		t.Fatalf("unexpected node affinity values: %v", na.Values)
	}

	// No affinities -> nil slice.
	if got := buildScheduling(&InstanceConfig{}).NodeAffinities; got != nil {
		t.Fatalf("expected nil node affinities when none set, got %v", got)
	}
}

func TestBuildMetadataItems(t *testing.T) {
	items := buildMetadataItems(&InstanceConfig{Metadata: map[string]string{"k": "v"}})
	if len(items) != 1 || items[0].Key != "k" || items[0].Value == nil || *items[0].Value != "v" {
		t.Fatalf("unexpected metadata items: %#v", items)
	}
}
