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

package runtime

func (ts *Script) CmdStdin(neg int, args []string) {
	if neg != 0 {
		ts.Fatalf("unsupported: !? stdin")
	}
	if len(args) != 1 {
		ts.Fatalf("usage: stdin filename")
	}
	ts.stdin = ts.ReadFile(args[0])
}

// stdout checks that the last go command standard output matches a regexp.
func (ts *Script) CmdStdout(neg int, args []string) {
	scriptMatch(ts, neg, args, ts.stdout, "stdout")
}

// stderr checks that the last go command standard output matches a regexp.
func (ts *Script) CmdStderr(neg int, args []string) {
	scriptMatch(ts, neg, args, ts.stderr, "stderr")
}
