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

package browser

import (
	"strconv"

	"github.com/opentofu/opentofu/internal/hof/lib/tui/components/cue/helpers"
)

func handleClear(B *Browser, action string, args []string, context map[string]any) (bool, error) {

	for _, cfg := range B.sources {
		if cfg.WatchTime > 0 {
			cfg.WatchTime = 0
			cfg.StopWatch()
		}
	}

	B.sources = []*helpers.SourceConfig{}
	B.RebuildValue()
	B.Rebuild()

	return true, nil
}

func handleSet(B *Browser, action string, args []string, context map[string]any) (bool, error) {
	sc, err := helpers.CreateFrom(args, context)
	if err != nil {
		return true, err
	}
	B.sources = []*helpers.SourceConfig{sc}

	B.SetThinking(true)
	defer B.SetThinking(false)
	B.RebuildValue()
	B.Rebuild()

	return true, nil
}

func handleAdd(B *Browser, action string, args []string, context map[string]any) (bool, error) {
	sc, err := helpers.CreateFrom(args, context)
	if err != nil {
		return true, err
	}
	B.sources = append(B.sources, sc)

	B.SetThinking(true)
	defer B.SetThinking(false)
	B.RebuildValue()
	B.Rebuild()

	return true, nil
}

func handleWatchConfig(B *Browser, action string, args []string, context map[string]any) (bool, error) {

	do := func(cfg *helpers.SourceConfig, action string, args []string, context map[string]any) (bool, error) {
		h, err := cfg.UpdateFrom(action, args, context)
		if err != nil {
			return h, err
		}
		cfg.WatchFunc = func() {
			B.SetThinking(true)
			defer B.SetThinking(false)
			B.RebuildValue()
			B.Rebuild()
		}
		go cfg.Watch()
		return true, nil
	}

	return B.handleOneOrAll(do, action, args, context)
}

type doer func(cfg *helpers.SourceConfig, action string, args []string, context map[string]any) (bool, error)

func (B *Browser) handleOneOrAll(do doer, action string, args []string, context map[string]any) (bool, error) {
	var err error
	var handled bool

	// look for an index marker (2 args)
	idx := -1
	if len(args) > 1 {
		// there was an index setting
		idxStr := args[0]
		args = args[1:]
		idx, err = strconv.Atoi(idxStr)
	}
	if err != nil {
		return true, err
	}

	if idx > -1 {
		cfg := B.sources[idx]
		handled, err = do(cfg, action, args, context)
	} else {
		for _, cfg := range B.sources {
			handled, err = do(cfg, action, args, context)
			if err != nil {
				break
			}
		}
	}

	return handled, err
}
