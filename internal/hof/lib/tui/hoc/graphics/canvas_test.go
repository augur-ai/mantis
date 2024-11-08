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

// Copyright 2017 Zack Guo <zack.y.guo@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

// +build ignore

package vermui

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestCanvasSet(t *testing.T) {
	c := NewCanvas()
	c.Set(0, 0)
	c.Set(0, 1)
	c.Set(0, 2)
	c.Set(0, 3)
	c.Set(1, 3)
	c.Set(2, 3)
	c.Set(3, 3)
	c.Set(4, 3)
	c.Set(5, 3)
	spew.Dump(c)
}

func TestCanvasUnset(t *testing.T) {
	c := NewCanvas()
	c.Set(0, 0)
	c.Set(0, 1)
	c.Set(0, 2)
	c.Unset(0, 2)
	spew.Dump(c)
	c.Unset(0, 3)
	spew.Dump(c)
}

func TestCanvasBuffer(t *testing.T) {
	c := NewCanvas()
	c.Set(0, 0)
	c.Set(0, 1)
	c.Set(0, 2)
	c.Set(0, 3)
	c.Set(1, 3)
	c.Set(2, 3)
	c.Set(3, 3)
	c.Set(4, 3)
	c.Set(5, 3)
	c.Set(6, 3)
	c.Set(7, 2)
	c.Set(8, 1)
	c.Set(9, 0)
	bufs := c.Buffer()
	spew.Dump(bufs)
}
