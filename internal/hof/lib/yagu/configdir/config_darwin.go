/*
 * Augur AI Proprietary
 * Copyright (c) 2024 Augur AI, Inc.
 *
 * This file is licensed under the Augur AI Proprietary License.
 *
 * Attribution:
 * This work is based on code from https://github.com/hofstadter-io/hof, licensed under the Apache License 2.0.
 */

package configdir

import "os"

var (
	systemConfig []string
	localConfig  string
	localCache   string
)

func findPaths() {
	systemConfig = []string{"/Library/Application Support"}
	localConfig = os.Getenv("HOME") + "/Library/Application Support"
	localCache = os.Getenv("HOME") + "/Library/Caches"
}