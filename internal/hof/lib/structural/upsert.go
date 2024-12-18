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

package structural

import (
	"fmt"

	"cuelang.org/go/cue"
)

func UpsertValue(up, val cue.Value, opts *Options) (cue.Value, error) {
	if opts == nil {
		opts = &Options{}
	}
	r, _ := upsertValue(up, val, opts)
	return r, nil
}

func upsertValue(up, val cue.Value, opts *Options) (cue.Value, bool) {
	switch val.IncompleteKind() {
	case cue.StructKind:
		return upsertStruct(up, val, opts)

	case cue.ListKind:
		return upsertList(up, val, opts)

	default:
		// should already have the same label by now
		// but maybe not if target is basic and up is not
		return up, true
	}
}

func upsertStruct(up, val cue.Value, opts *Options) (cue.Value, bool) {
	ctx := val.Context()
	result := newStruct(ctx)

	// first loop over val
	iter, _ := val.Fields(defaultWalkOptions...)
	for iter.Next() {
		s := iter.Selector()
		// HACK, this works around a bug in CUE
		// p := cue.MakePath(s)
		p := cue.ParsePath(fmt.Sprint(s))
		u := up.LookupPath(p)

		// check that field exists in from. Should we be checking f.Err()?
		if u.Exists() {
			r, ok := upsertValue(u, iter.Value(), opts)
			if ok {
				result = result.FillPath(p, r)
			}
		} else {
			// include if not in val
			result = result.FillPath(p, iter.Value())
		}
	}

	// add anything in ins that is not in val
	iter, _ = up.Fields(defaultWalkOptions...)
	for iter.Next() {
		s := iter.Selector()

		// HACK, this works around a bug in CUE
		// p := cue.MakePath(s)
		p := cue.ParsePath(fmt.Sprint(s))

		v := val.LookupPath(p)

		// check that field exists in from. Should we be checking f.Err()?
		if !v.Exists() {
			result = result.FillPath(p, iter.Value())
		}
	}

	return result, true
}

func upsertList(up, val cue.Value, opts *Options) (cue.Value, bool) {
	ctx := val.Context()

	ui, _ := up.List()
	vi, _ := val.List()

	result := []cue.Value{}
	for ui.Next() && vi.Next() {
		r, ok := upsertValue(ui.Value(), vi.Value(), opts)
		if ok {
			result = append(result, r)
		}
	}
	return ctx.NewList(result...), true

}
