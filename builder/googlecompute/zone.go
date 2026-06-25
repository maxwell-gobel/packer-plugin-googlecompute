// Copyright IBM Corp. 2013, 2026
// SPDX-License-Identifier: MPL-2.0

package googlecompute

import "github.com/hashicorp/packer-plugin-sdk/multistep"

// stateZone returns the resolved build zone from the state bag. The zone is
// seeded by the builder and, in bulk mode, overwritten with the
// automatically-selected zone by StepCreateInstance.
func stateZone(state multistep.StateBag) string {
	if z, ok := state.GetOk("zone"); ok {
		return z.(string)
	}
	return ""
}
