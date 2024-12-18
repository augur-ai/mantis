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

package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/opentofu/opentofu/internal/hof/cmd/hof/flags"
	"github.com/opentofu/opentofu/internal/hof/lib/runtime"
)

func log(R *runtime.Runtime, dflags flags.DatamodelPflagpole) error {
	// find max label width after indentation for column alignment
	max := findMaxLabelLen(R, dflags)

	cflags := flags.Datamodel__LogFlags

	if cflags.ByValue {
		for _, dm := range R.Datamodels {
			if err := dm.PrintLogByValue(os.Stdout, max, dflags, cflags); err != nil {
				return err
			}
		}
	} else {

		tm := make(map[string]string)
		for _, dm := range R.Datamodels {

			// snapshots (map[path][]snapshot)
			SL := dm.GetSnapshotList()

			// build timestamp -> message map
			for _, sl := range SL {
				for _, s := range sl {
					tm[s.Timestamp] = s.Message
				}
			}
		}

		// uniq timestamp strings
		ts := make([]string, 0, len(tm))
		for t, _ := range tm {
			ts = append(ts, t)
		}

		// sort and reverse so most recent first
		sort.Strings(ts)
		for i, j := 0, len(ts)-1; i < j; i, j = i+1, j-1 {
			ts[j], ts[i] = ts[i], ts[j]
		}

		for _, t := range ts {
			msg := tm[t]
			fmt.Fprintf(os.Stdout, "%s: %q\n", t, msg)
			for _, dm := range R.Datamodels {
				if err := dm.PrintLog(os.Stdout, max, t, dflags, cflags); err != nil {
					return err
				}
			}
			fmt.Fprintf(os.Stdout, "\n")
		}

	}
	return nil
}


