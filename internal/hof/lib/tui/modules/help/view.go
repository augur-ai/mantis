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

package help

import (
	"fmt"

	"github.com/opentofu/opentofu/internal/hof/lib/tui/modules/eval"
	"github.com/opentofu/opentofu/internal/hof/lib/tui/tview"
)

// Both a Module and a Layout and a Switcher.SubLayout
type Help struct {
	*tview.TextView
}

func NewHelp() *Help {
	view := tview.NewTextView()
	view.
		SetTitle("  Help  ").
		SetBorder(true).
		SetBorderPadding(1, 1, 2, 2)
	view.
		SetWrap(false).
		SetScrollable(true).
		SetDynamicColors(true).
		SetRegions(true)

	fmt.Fprintln(view, eval.EvalHelpText)

	h := &Help{
		TextView: view,
	}

	return h
}

func (H *Help) Id() string {
	return "help"
}

func (H *Help) Name() string {
	return "Help"
}

func (H *Help) CommandName() string {
	return "help"
}

func (H *Help) CommandUsage() string {
	return "help <topic> [sub-topics...]"
}

func (H *Help) CommandHelp() string {
	return "displays the home view"
}
func (H *Help) CommandCallback(context map[string]interface{}) {
	//helpPath := "/help"

	//args := []string{}
	//if a, ok := context["args"]; ok {
	//  args, _ = a.([]string)
	//}

	//if len(args) > 0 {
	//  H.Clear()
	//  fmt.Fprintln(H, "Help -", args, "\n\n")
	//  helpPath += "/" + strings.Join(args, "/")
	//} else {
	//  H.Clear()
	//  fmt.Fprintln(H, "Help - Main\n\n", helpContent)
	//}
}
