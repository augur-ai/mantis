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

package components

import (
	"fmt"

	"github.com/opentofu/opentofu/internal/hof/lib/tui/components/widget"

	// cue widgets
	"github.com/opentofu/opentofu/internal/hof/lib/tui/components/cue/browser"
	// "github.com/opentofu/opentofu/internal/hof/lib/tui/components/cue/flower"
	"github.com/opentofu/opentofu/internal/hof/lib/tui/components/cue/playground"
)

var _registry = map[string]func (input map[string]any) (widget.Widget, error){
	// common widgets
	"widget/Box": (widget.NewBox()).Decode,
	"widget/TextView": (widget.NewTextView()).Decode,

	// cue widgets
	"cue/browser": (&browser.Browser{}).Decode,
	// "cue/flower": (&flower.Flower{}).Decode,
	"cue/playground": (&playground.Playground{}).Decode,
}

func DecodeWidget(input map[string]any) (widget.Widget, error) {

	typename, ok := input["typename"]
	if !ok {
		return nil, fmt.Errorf("input to DecodeWidget did not contain 'typename'")
	}

	decoder, ok := _registry[typename.(string)]
	if !ok {
		return nil, fmt.Errorf("unknown 'typename': %q in DecodeWidget", typename)
	}

	return decoder(input)
}
