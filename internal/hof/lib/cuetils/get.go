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

package cuetils

import (
	"fmt"
	"regexp"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/errors"
)

type KeyVal struct {
	Key string
	Val cue.Value
}

// GetByAttrAndKeys extracts fields from a value by attribute and key names
// if all or any is empty, the condition check is skipped and all values will pass
// so to get all values with an attribute, with no concern for the contents, use:
//
//   GetByAttrAndKeys(val, "myattr", []string{}, []string{})
//
// TODO, change this to return map[string]cue.Value
func GetByAttrKeys(val cue.Value, attr string, all, any []string) ([]KeyVal, error) {
	// Todo, rewrite this to use structural

	// fmt.Println("GET:", name, attr, all, any)

	rets := []KeyVal{}

	S, err := val.Struct()
	if err != nil {
		es := errors.Errors(err)
		for _, e := range es {
			fmt.Println(e)
		}
		return rets, fmt.Errorf("Error loading cue code")
	}

	// Loop through all top level fields
	iter := S.Fields(
		cue.Attributes(true),
		cue.Definitions(true),
	)
	for iter.Next() {

		label := iter.Label()
		value := iter.Value()
		attrs := value.Attributes(cue.ValueAttr)

		// fmt.Println("  -", label, attrs)

		// find top-level with gen attr
		hasattr := false
		for _, A := range attrs {
			// does it have an "@<attr>(...)"
			if A.Name() == attr {

				vals := AttrToMap(A)

				// must match all
				if len(all) > 0 {
					match := true

					// loop over the all list
					for _, l := range all {
						R := regexp.MustCompile(l)
						// loop over the field attt key names
						found := false
						for v, _ := range vals {
							m := R.MatchString(v)
							if m {
								found = true
								break
							}
						}
						// break one more time if we have failed
						if !found {
							match = false
							break
						}
					}

					// did we not match all?
					if !match {
						continue
					}
				}

				// match one of any
				if len(any) > 0 {
					match := false

					// loop over the any list
					for _, l := range any {
						R := regexp.MustCompile(l)
						// loop over the field attt key names
						for v, _ := range vals {
							m := R.MatchString(v)
							if m {
								match = true
								break
							}
						}

						// break again if we have matched
						if match {
							break
						}
					}

					// did we not match any?
					if !match {
						continue
					}
				}

				// fmt.Println("  ...Has", label, A.Name())
				// passed, we should include
				hasattr = true
				break
			}
		}

		// fmt.Println("  ...Attr", label, attr, hasattr)
		// ok, we're back outside the attrs look now, did we match on it?
		// if no, let's try the next field
		if !hasattr {
			continue
		}

		// add it and move on!
		rets = append(rets, KeyVal{Key: label, Val: value})
	}

	return rets, nil
}
