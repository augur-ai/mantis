/*
 * Augur AI Proprietary
 * Copyright (c) 2024 Augur AI, Inc.
 *
 * This file is licensed under the Augur AI Proprietary License.
 *
 * Attribution:
 * This work is based on code from https://github.com/hofstadter-io/hof, licensed under the Apache License 2.0.
 */

// Copyright 2022 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (js && wasm) || plan9

package mmap

import (
	"io"
	"os"
)

// mmapFile on other systems doesn't mmap the file. It just reads everything.
func mmapFile(f *os.File) (Data, error) {
	b, err := io.ReadAll(f)
	if err != nil {
		return Data{}, err
	}
	return Data{f, b}, nil
}
