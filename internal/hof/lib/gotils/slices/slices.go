/*
 * Augur AI Proprietary
 * Copyright (c) 2024 Augur AI, Inc.
 *
 * This file is licensed under the Augur AI Proprietary License.
 *
 * Attribution:
 * This work is based on code from https://github.com/hofstadter-io/hof, licensed under the Apache License 2.0.
 */

// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// TODO: Replace with slices package when it lands in standard library.

package slices

// Clip removes unused capacity from the slice, returning s[:len(s):len(s)].
func Clip[S ~[]E, E any](s S) S {
	return s[:len(s):len(s)]
}
