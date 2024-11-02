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

package test

import (
	"time"

	"cuelang.org/go/cue"

	"github.com/opentofu/opentofu/internal/hof/lib/cuetils"
)

// Known testers in the hof system
var knownTesters = []string{
	"bash",
	"exec",
	"tsuite",
	"api",
	"hls",
}

// A suite is a collection of testers
type Suite struct {
	// Name of the Suite (Cue field)
	Name string

	// Runtime the Cue values were built from
	CTX *cue.Context

	// Cue Value for the Suite
	Value cue.Value

	// Extracted testers, including any selectors
	Tests []Tester

	// pass/fail/skip stats
	Stats Stats

	// Total Suite runtime (to account for gaps between tests)
	Runtime time.Duration

	// Errors encountered during the testing
	Errors []error
}

func getValueTestSuites(ctx *cue.Context, val cue.Value, labels []string) ([]Suite, error) {
	vals, err := cuetils.GetByAttrKeys(val, "test", append(labels, "suite"), nil)
	suites := []Suite{}
	for _, v := range vals {
		suites = append(suites, Suite{Name: v.Key, CTX: ctx, Value: v.Val})
	}
	return suites, err
}

// A tester has configuration for running a set of tests
type Tester struct {
	// Name of the Tester (Cue field)
	Name string

	// Type of the Tester (@test(key[0]))
	Type string

	// Runtime the Cue values were built from
	CTX *cue.Context

	// Cue Value for the Tester
	Value cue.Value

	// Execution output
	Output string
	//Stdout string
	//Stderr string

	// pass/fail/skip stats
	Stats Stats

	// Errors encountered during the testing
	Errors []error
}

func getValueTestSuiteTesters(ctx *cue.Context, val cue.Value, labels []string) ([]Tester, error) {
	vals, err := cuetils.GetByAttrKeys(val, "test", labels, []string{})
	testers := []Tester{}
	for _, v := range vals {
		a := v.Val.Attribute("test")
		typ, err := a.String(0)
		if err != nil {
			return testers, err
		}
		testers = append(testers, Tester{Name: v.Key, Type: typ, CTX: ctx, Value: v.Val})
	}
	return testers, err
}

type Stats struct {
	Pass int
	Fail int
	Skip int

	Start time.Time
	End   time.Time
	Time  time.Duration
}

func (S *Stats) add(s Stats) {
	S.Pass += s.Pass
	S.Fail += s.Fail
	S.Skip += s.Skip
	S.Time += s.Time
}
