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

package form

import (
	"github.com/gdamore/tcell/v2"
	"github.com/opentofu/opentofu/internal/hof/lib/tui/tview"
)

type Flex struct {
	*tview.Flex

	name    string
	items   []FormItem
	buttons []FormButton
}

func NewFlex(name string) *Flex {
	F := &Flex{
		Flex:    tview.NewFlex(),
		name:    name,
		items:   []FormItem{},
		buttons: []FormButton{},
	}

	return F
}

func (F *Flex) Name() string {
	return F.name
}

func (F *Flex) SetFinishedFunction(handler func(key tcell.Key)) {}

func (F *Flex) GetValues() (values map[string]interface{}) {
	values = make(map[string]interface{})
	for _, item := range F.items {
		vals := item.GetValues()
		for field, value := range vals {
			values[field] = value
		}
	}
	return values
}

func (F *Flex) SetValues(values map[string]interface{}) {
	for _, item := range F.items {
		item.SetValues(values)
	}
}

func (F *Flex) GetItem(name string) FormItem {
	for _, item := range F.GetItems() {
		if item.Name() == name {
			return item
		}
	}

	return nil
}

func (F *Flex) GetItems() []FormItem {
	items := []FormItem{}
	for _, item := range F.items {
		switch typ := item.(type) {
		case *Flex:
			itms := typ.GetItems()
			items = append(items, itms...)
		default:
			items = append(items, item)
		}
	}

	return items
}

func (F *Flex) AddItem(item FormItem, fixedSize, proportion int) {
	F.items = append(F.items, item)

	F.Flex.AddItem(item, fixedSize, proportion, true)
}

func (F *Flex) GetButton(name string) FormButton {
	for _, item := range F.GetButtons() {
		if item.Name() == name {
			return item
		}
	}
	return nil
}

func (F *Flex) GetButtons() []FormButton {
	buttons := []FormButton{}
	for _, button := range F.buttons {
		buttons = append(buttons, button)
	}
	for _, item := range F.items {
		switch typ := item.(type) {
		case *Flex:
			btns := typ.GetButtons()
			buttons = append(buttons, btns...)
		}
	}
	return buttons
}

func (F *Flex) AddButton(button FormButton, fixedSize, proportion int) {
	F.buttons = append(F.buttons, button)
	F.Flex.AddItem(button, fixedSize, proportion, true)
}
