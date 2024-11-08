/*
/*
 * Copyright (c) 2024 Augur AI, Inc.
 * This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0. 
 * If a copy of the MPL was not distributed with this file, you can obtain one at https://mozilla.org/MPL/2.0/.
 *
 
 * Copyright (c) 2024 Augur AI, Inc.
 *
 * This file is licensed under the Augur AI Proprietary License.
 *
 * Attribution:
 * This work is based on code from https://github.com/hofstadter-io/hof, licensed under the Apache License 2.0.
 */

package types

// More TBD
type Appdir struct {
	Accounts   map[string]interface{}
	Workspaces map[string]interface{}
	Contexts   map[string]interface{}

	Clouds    map[string]interface{}
	Environs  map[string]interface{}
	Resources map[string]interface{}
}
