// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2023 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package utils

import "github.com/opentofu/opentofu/internal/command/workdir"

func WorkingDir(originalDir string, overrideDataDir string) *workdir.Dir {
	ret := workdir.NewDir(".") // caller should already have used os.Chdir in "-chdir=..." mode
	ret.OverrideOriginalWorkingDir(originalDir)
	if overrideDataDir != "" {
		ret.OverrideDataDir(overrideDataDir)
	}
	return ret
}
