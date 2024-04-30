/*
 * Augur AI Proprietary
 * Copyright (c) 2024 Augur AI, Inc.
 *
 * This file is licensed under the Augur AI Proprietary License.
 *
 * Attribution:
 * This work is based on code from https://github.com/hofstadter-io/hof, licensed under the Apache License 2.0.
 */

package types

import (
	"cuelang.org/go/cue"
)

type Config struct {
	Modelsets map[string]Modelset
	Stores    map[string]Store

	CueValue cue.Value
}
