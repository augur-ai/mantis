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

package hof

import (
	"fmt"
)

func (hn *Node[T]) Print() {
	if hn.Hof.Datamodel.Root {
		hn.PrintDatamodel("  ", "")
	}

	if hn.Hof.Gen.Root {
		hn.PrintGen("  ", "")
	}

	if hn.Hof.Flow.Root {
		hn.PrintFlow("  ", "")
	}
}

func (hn *Node[T]) PrintDatamodel(spacer, indent string) {
	s := fmt.Sprintf("%s%s:", indent, hn.Hof.Path)
	dm := hn.Hof.Datamodel
	found := false
	if dm.Root {
		s += " @datamodel()"
		found = true
	}
	if dm.Node {
		s += " @node()"
		found = true
	}
	if dm.History {
		s += " @history()"
		found = true
	}
	if dm.Ordered {
		s += " @ordered()"
		found = true
	}
	if dm.Cue {
		s += " @cue()"
		found = true
	}
	
	if !found {
		return
	}
	// print line
	fmt.Println(s)
	// pretty.Println(hn.Hof)

	// recurse
	for _, c := range hn.Children {
		c.PrintDatamodel(spacer, indent+spacer)
	}
}

func (hn *Node[T]) PrintGen(spacer, indent string) {
	fmt.Printf("%s%s: @gen(%s)\n", indent, hn.Hof.Path, hn.Hof.Gen.Name)
	for _, c := range hn.Children {
		c.PrintDatamodel(spacer, indent+spacer)
	}
}

func (hn *Node[T]) PrintFlow(spacer, indent string) {
	f := hn.Hof.Flow
	s := fmt.Sprintf("%s%s:", indent, hn.Hof.Path)

	if f.Root {
		s += fmt.Sprintf(" @flow(%s)", f.Name)
	}
	if f.Task != "" {
		s += fmt.Sprintf(" @task(%s)", f.Task)
	}
	// print line
	fmt.Println(s)
	// pretty.Println(hn.Hof)

	// recurse
	for _, c := range hn.Children {
		c.PrintFlow(spacer, indent+spacer)
	}
}
