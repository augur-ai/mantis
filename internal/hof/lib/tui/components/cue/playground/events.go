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

package playground

import (
	"github.com/gdamore/tcell/v2"

	"github.com/opentofu/opentofu/internal/hof/lib/tui"
	"github.com/opentofu/opentofu/internal/hof/lib/tui/components/cue/helpers"
	"github.com/opentofu/opentofu/internal/hof/lib/tui/tview"
)

func (C *Playground) setupKeybinds() {
	// events (hotkeys)
	C.SetInputCapture(func(evt *tcell.EventKey) *tcell.EventKey {
		switch evt.Key() {
		case tcell.KeyRune:
			if (evt.Modifiers() & tcell.ModAlt) == tcell.ModAlt {
				switch evt.Rune() {
				case 'f':
					flexDir := C.GetDirection()
					if flexDir == tview.FlexRow {
						C.SetDirection(tview.FlexColumn)
					} else {
						C.SetDirection(tview.FlexRow)
					}

				case 's':
					C.ToggleShowScope()
					C.Rebuild()

				case 'S':
					C.UseScope(!C.useScope)
					if C.useScope {
						C.scope.Rebuild()
					}
					C.Rebuild()

				case 'E':
					C.mode = ModeEval
					C.Rebuild()
				case 'W':
					C.mode = ModeFlow
					C.Rebuild()

				case 'R':
					C.scope.Rebuild()

					switch C.editCfg.Source {
					case helpers.EvalText, helpers.EvalFile, helpers.EvalHttp, helpers.EvalBash:
						txt, err := C.editCfg.GetText()
						if err != nil {
							tui.Log("error", err)
							txt = err.Error()
						}
						C.edit.SetText(txt, false)
					}
					C.Rebuild()
					tui.SetFocus(C)

				default: 
					return evt
				}

				return nil
			}

			return evt

		default:
			return evt
		}
	})	
}
