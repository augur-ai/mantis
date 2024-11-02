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

package ast

import (
	"fmt"
	"strings"
	"time"
)

type Phase struct {
	NodeBase

	// how many '%' signs
	level int
	// the string after
	title string

	// TODO, will be same as Cue attrs: @label(...)
	// attrs []Attr

	begTime time.Time
	endTime time.Time

	parent *Phase

	// subnodes / expressions, etx
	nodes []Node
}

func (P *Parser) parsePhase() error {
	Ph := P.phase
	N := P.node
	S := P.script

	// grab current line
	line := stripTrailingWhitespace(S.Lines[N.BegLine()])
	spc := strings.Index(line, " ")
	if spc == -1 {
		return fmt.Errorf("Phase missing title in %s:%d", S.Path, N.BegLine())
	}

	lvl, title := line[0:spc], line[spc+1:]
	title = cleanLine(title)

	ph := &Phase{
		NodeBase: P.node.CloneNodeBase(),
		level:    len(lvl),
		title:    title,
	}

	// no phase yet
	if Ph == nil {
		// add phase to script
		S.AddPhase(ph)

	} else if ph.level < Ph.level {
		// new is sub of current

		// set parent to current phase
		ph.parent = Ph

		// update current  phase
		Ph.SetEndLine(P.lineno)

	} else if ph.level == Ph.level {
		// new is same as current

		// set same parent
		ph.parent = Ph.parent

	} else if ph.level > Ph.level {
		// new is bigger, need to walk

		// walk up parent chain, closing phases
		p := Ph
		for ; p != nil && p.level < ph.level; p = p.parent {
			p.SetEndLine(ph.DocLine() - 1)
		}

		// if p is not nil, we found our sibling
		if p != nil {
			p.SetEndLine(ph.DocLine() - 1)
			// set same parent
			ph.parent = p.parent
		}
	}

	// add to script if == first level encountered
	f := S.Phases[0]
	if f.level == ph.level && f.title != ph.title {
		S.AddPhase(ph)
	}

	// add new phase to parent
	if ph.parent != nil {
		ph.parent.AppendNode(ph)
	}

	// update parser values
	P.phase = ph
	P.node = nil

	return nil
}

func (Ph *Phase) AppendNode(n Node) {
	Ph.nodes = append(Ph.nodes, n)
}
