// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cli

import (
	"testing"
)

func TestMockCommand_implements(t *testing.T) {
	var _ Command = new(MockCommand)
}
