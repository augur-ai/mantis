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

package eval

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"

	"github.com/opentofu/opentofu/internal/hof/lib/tui"
	"github.com/opentofu/opentofu/internal/hof/lib/tui/components/panel"
	"github.com/opentofu/opentofu/internal/hof/lib/tui/components/widget"
	"github.com/opentofu/opentofu/internal/hof/lib/tui/tview"
)

type Eval struct {
	*panel.Panel

	// border display
	showPanel, showOther bool

	// default overide to all panels
	// would it be better as a widget creator? (after refactor 1)
	// or a function that can take a widget creator with a default ItemBase++
	_creator panel.ItemCreator
}

func NewEval() *Eval {
	M := &Eval{
		showPanel: true,
		showOther: true,
	}
	M.Panel = panel.New(nil, M.creator)
	M.Panel.SetBorderColor(tcell.Color42).SetBorder(true)

	// add dummy primitive
	M.AddItem(tview.NewBox(), 0, 1, true)


	// do layout setup here
	M.SetName("eval")
	M.SetBorder(true)
	M.SetDirection(tview.FlexColumn)

	return M
}

func (M *Eval) Mount(context map[string]any) error {

	// this will mount the core element and all children
	M.Flex.Mount(context)
	// tui.Log("trace", "Eval.Mount")

	// probably want to do some self mount first?
	M.setupEventHandlers()

	// and then refresh?
	// should this even happen, or just a rebuild
	M.CommandCallback(context)

	return nil
}

func (M *Eval) Unmount() error {
	// remove keybinds
	M.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey { return event })

	// handle border display
	tui.RemoveWidgetHandler(M.Panel, "/sys/key/A-P")
	tui.RemoveWidgetHandler(M.Panel, "/sys/key/A-O")

	// this is where we can do some unloading, depending on the application
	M.Flex.Unmount()

	return nil
}

// todo, add more functions so that we can separate new command messages from refresh?

func (M *Eval) showError(err error) error {
	txt := widget.NewTextView()
	fmt.Fprint(txt, err)

	I := panel.NewBaseItem(M.Panel)
	I.SetWidget(txt)

	M.Panel.AddItem(I, 0, 1, true)

	return err
}



func (M *Eval) Focus(delegate func(p tview.Primitive)) {
	// tui.Log("warn", "Eval.Focus")
	delegate(M.Panel)
	// M.Panel.Focus(delegate)
}

func (M *Eval) getPanelByPath(path string) (*panel.Panel, error) {
	if path == "" {
		return M.Panel, nil
	}
	parts := strings.Split(path, ".")

	// set at our panel
	curr := M.Panel

	for _, part := range parts {
		p := curr.GetItemByName(part)
		if p == nil {
			p = curr.GetItemById(part)
			if p == nil {
				return nil, fmt.Errorf("unable to find node %q in %q", part, path)
			}
		}
		switch t := p.(type) {
		case *panel.Panel:
			curr = t	
		}
	}

	return curr, nil

	return nil, fmt.Errorf("did not find item at path %q", path)
}

func (M *Eval) getItemByPath(path string) (panel.PanelItem, error) {
	parts := strings.Split(path, ".")

	// set at our panel
	curr := M.Panel

	for _, part := range parts {
		p := curr.GetItemByName(part)
		if p == nil {
			p = curr.GetItemById(part)
			if p == nil {
				return nil, fmt.Errorf("unable to find node %q in %q", part, path)
			}
		}
		switch t := p.(type) {
		case *panel.Panel:
			curr = t	
		case panel.PanelItem:
			return t, nil
		}
	}

	return nil, fmt.Errorf("did not find item at path %q", path)
}
