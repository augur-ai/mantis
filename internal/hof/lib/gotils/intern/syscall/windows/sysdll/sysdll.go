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

// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package sysdll is an internal leaf package that records and reports
// which Windows DLL names are used by Go itself. These DLLs are then
// only loaded from the System32 directory. See Issue 14959.
package sysdll

// IsSystemDLL reports whether the named dll key (a base name, like
// "foo.dll") is a system DLL which should only be loaded from the
// Windows SYSTEM32 directory.
//
// Filenames are case sensitive, but that doesn't matter because
// the case registered with Add is also the same case used with
// LoadDLL later.
//
// It has no associated mutex and should only be mutated serially
// (currently: during init), and not concurrent with DLL loading.
var IsSystemDLL = map[string]bool{}

// Add notes that dll is a system32 DLL which should only be loaded
// from the Windows SYSTEM32 directory. It returns its argument back,
// for ease of use in generated code.
func Add(dll string) string {
	IsSystemDLL[dll] = true
	return dll
}
