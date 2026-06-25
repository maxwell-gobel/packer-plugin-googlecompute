// Copyright IBM Corp. 2013, 2026
// SPDX-License-Identifier: MPL-2.0

package googlecompute

import (
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/stretchr/testify/assert"
)

func TestStateZone(t *testing.T) {
	state := new(multistep.BasicStateBag)
	assert.Equal(t, "", stateZone(state))
	state.Put("zone", "us-east1-b")
	assert.Equal(t, "us-east1-b", stateZone(state))
}
