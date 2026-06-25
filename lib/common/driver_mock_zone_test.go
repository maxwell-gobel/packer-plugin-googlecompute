// Copyright IBM Corp. 2013, 2026
// SPDX-License-Identifier: MPL-2.0

package common

import (
	"errors"
	"testing"
)

func TestDriverMock_GetInstanceZone(t *testing.T) {
	d := &DriverMock{
		GetInstanceZoneResult: "us-east1-c",
	}

	zone, err := d.GetInstanceZone("packer-123")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if zone != "us-east1-c" {
		t.Fatalf("expected zone us-east1-c, got %q", zone)
	}
	if d.GetInstanceZoneName != "packer-123" {
		t.Fatalf("expected recorded name packer-123, got %q", d.GetInstanceZoneName)
	}
}

func TestDriverMock_GetInstanceZone_error(t *testing.T) {
	wantErr := errors.New("boom")
	d := &DriverMock{GetInstanceZoneErr: wantErr}
	_, err := d.GetInstanceZone("packer-1")
	if err != wantErr {
		t.Fatalf("expected returned error %v, got %v", wantErr, err)
	}
}

func TestDriverMock_RunInstanceInRegion(t *testing.T) {
	d := &DriverMock{}
	c := &InstanceConfig{Name: "packer-x", Region: "us-east1"}

	_, err := d.RunInstanceInRegion(c)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if d.RunInstanceInRegionConfig == nil || d.RunInstanceInRegionConfig.Region != "us-east1" {
		t.Fatalf("expected recorded config with region us-east1")
	}
}
