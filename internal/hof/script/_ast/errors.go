/*
 * Augur AI Proprietary
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
)

type ScriptError struct {
	Message string
	Node    Node
	Err     error
}

func NewScriptError(msg string, node Node, err error) *ScriptError {
	return &ScriptError{msg, node, err}
}

func (e *ScriptError) Error() string {
	src := ""
	n := e.Node
	s := n.Script()
	lines := s.Lines[n.DocLine():n.EndLine()]
	if len(lines) == 1 {
		src = lines[0]
	} else {
		src = strings.Join(lines, "\n  > ")
	}
	msg := `%s:%d:%d:%d: %s %v
	  > %s
	`
	return fmt.Sprintf(msg, s.Path, n.BegLine()+1, n.EndLine()+1, e.Message, e.Err, src)
}
