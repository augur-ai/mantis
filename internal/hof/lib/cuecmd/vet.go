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

package cuecmd

import (
	"fmt"
	"os"
	"time"

	"cuelang.org/go/cue"

	"github.com/opentofu/opentofu/internal/hof/cmd/hof/flags"
	"github.com/opentofu/opentofu/internal/hof/lib/cuetils"
	"github.com/opentofu/opentofu/internal/hof/lib/runtime"
)

func Vet(args []string, rflags flags.RootPflagpole, cflags flags.VetFlagpole) error {

	start := time.Now()
	R, err := runtime.New(args, rflags)

	defer func() {
		if R.Flags.Stats {
			fmt.Println(R.Stats)
			end := time.Now()
			fmt.Printf("\nTotal Elapsed Time: %s\n\n", end.Sub(start))
		}
	}()

	if err != nil {
		return err
	}

	wantErrors := rflags.IngoreErrors || rflags.AllErrors

	// this is a bit hacky (?), but we use this so we can vet the orphaned files rather than add them to the value
	R.DontPlaceOrphanedFiles = true
	err = R.Load()
	if err != nil {
		return cuetils.ExpandCueError(err)
	}

	// build options
	opts := []cue.Option{
		cue.Docs(cflags.Comments),
		cue.Attributes(cflags.Attributes),
		cue.Definitions(cflags.Definitions),
		cue.Optional(cflags.Optional),
		cue.ErrorsAsValues(wantErrors),
	}

	// these two have to be done specially
	// because there are three options [true, false, missing]
	if cflags.Concrete {
		opts = append(opts, cue.Concrete(true))
	}
	if cflags.Hidden {
		opts = append(opts, cue.Hidden(true))
	}

	exs := cflags.Expression
	if len(exs) == 0 {
		exs = []string{""}
	}

	hadError := false
	handleErr := func(ex string, err error) {
		if err == nil {
			return
		}
		err = cuetils.ExpandCueError(err)
		hadError = true
		if len(exs) > 1 {
			fmt.Fprintln(os.Stderr, "//", ex)
		}
		fmt.Fprint(os.Stderr, err)
	}

	// TODO, how do we think about the cross-product of { files } x { -e } x { -l }
	// maybe -l doesn't make sense here? (or only files that can be placed)

	// setup our bi and other stuff
	bi := R.BuildInstances[0]
	if R.Flags.Verbosity > 1 {
		fmt.Println("ID:", bi.ID(), bi.PkgName, bi.Module)
	}
	pkg := bi.PkgName
	if bi.Module == "" {
		pkg = bi.ID()
	}

	// vet the orphaned files
	hadOrphan := false // MORE HACKS FOR INCUESISTENCY, .txt files are now showing up here
	if len(bi.OrphanedFiles) > 0 {
		for i, f := range bi.OrphanedFiles {
			F, err := R.LoadOrphanedFile(f, pkg, bi.Root, bi.Dir, i, len(bi.OrphanedFiles))
			if err != nil {
				handleErr("during load", err)
				continue
			}
			// probably a filetype CUE does not understand
			if F == nil {
				if R.Flags.Verbosity > 1 {
					fmt.Printf("nil file for %s\n", f.Filename)
				}
				continue
			}
			hadOrphan = true
			fv := R.CueContext.BuildFile(F)

			// vet the value with each expression
			for _, ex := range exs {

				v := getValByEx(ex, pkg, R.Value)	
				if !v.Exists() {
					handleErr(ex, v.Err())
					continue
				}

				v = v.Unify(fv)
			
				// we want to ensure concrete when validating data (orphaned files)
				opts = append(opts, cue.Concrete(true))
				err := v.Validate(append(opts, )...)
				handleErr(ex, err)
			}

		}
	}

	// ugh, more hacks because inCUEsistency...
	if !hadOrphan {
		// vet the root value at each expression
		// often this will default to [""] which is just the whole value
		for _, ex := range exs {

			v := getValByEx(ex, pkg, R.Value)
			if !v.Exists() {
				handleErr(ex, v.Err())
				continue
			}
		
			err := v.Validate(append(opts, )...)
			handleErr(ex, err)
		}
	}



	if hadError {
		// messages already printed, we want an empty message
		return fmt.Errorf("")
	}

	return nil
}
