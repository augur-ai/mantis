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

package os

import (
	"time"

	"cuelang.org/go/cue"

	hofcontext "github.com/opentofu/opentofu/internal/hof/flow/context"
)

type Sleep struct{}

func NewSleep(val cue.Value) (hofcontext.Runner, error) {
	return &Sleep{}, nil
}

func (T *Sleep) Run(ctx *hofcontext.Context) (interface{}, error) {

	v := ctx.Value

	var D time.Duration

	ferr := func() error {
		ctx.CUELock.Lock()
		defer func() {
			ctx.CUELock.Unlock()
		}()

		d := v.LookupPath(cue.ParsePath("duration"))
		if d.Err() != nil {
			return d.Err()
		} else if d.Exists() {
			ds, err := d.String()
			if err != nil {
				return err
			}
			D, err = time.ParseDuration(ds)
			if err != nil {
				return err
			}
		}
		return nil
	}()
	if ferr != nil {
		return nil, ferr
	}

	time.Sleep(D)

	var res interface{}
	func() {
		ctx.CUELock.Lock()
		defer func() {
			ctx.CUELock.Unlock()
		}()
		res = v.FillPath(cue.ParsePath("done"), true)
	}()

	return res, nil
}
