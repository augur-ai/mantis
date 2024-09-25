/*
 * Augur AI Proprietary
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
	"os"
	"strconv"

	"github.com/gdamore/tcell/v2"

	"github.com/opentofu/opentofu/internal/hof/lib/tui"
	"github.com/opentofu/opentofu/internal/hof/lib/tui/components/panel"
	"github.com/opentofu/opentofu/internal/hof/lib/tui/components/widget"
	"github.com/opentofu/opentofu/internal/hof/lib/tui/components/cue/browser"
	// "github.com/opentofu/opentofu/internal/hof/lib/tui/components/cue/flower"
	"github.com/opentofu/opentofu/internal/hof/lib/tui/components/cue/helpers"
	"github.com/opentofu/opentofu/internal/hof/lib/tui/components/cue/playground"
)

// used for debugging panel CRUD & KEYS
var panel_debug = false

func init() {
	if v := os.Getenv("HOF_TUI_PANEL_DEBUG"); v != "" {
		vb, _ := strconv.ParseBool(v)
		if vb {
			panel_debug = true
		}
	}

	if !panel_debug {
		setupCreator()
	}
}

var itemCreator *panel.Factory

func setupCreator() {
	f := panel.NewFactory()

	f.Register("default", helpItem)
	f.Register("help", helpItem)
	f.Register("play", playItem)
	f.Register("view", viewItem)
	// f.Register("flow", flowItem)

	itemCreator = f
}

func (E *Eval) setThinking(thinking bool) {
	c := tcell.Color42
	if thinking {
		c = tcell.ColorViolet
	}

	E.SetBorderColor(c)
	go tui.Draw()
}

// this function is responsable for creating the components that fill slots in the panel
// these are the widgets that make up the application and should have their own operation
func (E *Eval) creator(context panel.ItemContext, parent *panel.Panel) (panel.PanelItem, error) {
	tui.Log("extra", fmt.Sprintf("Eval.creator: %v", context ))

	E.setThinking(true)
	defer E.setThinking(false)

	// short-circuit for developer mode (first, before user custom)
	if panel_debug {
		t := widget.NewTextView()
		i := panel.NewBaseItem(parent)
		i.SetWidget(t)
		return i, nil
	}

	// set default item
	if _, ok := context["item"]; !ok {
		context["item"] = "help"
	}

	i, e := itemCreator.Creator(context, parent)
	// todo, better error handling
	if e == nil {
		i.SetBorder(E.showOther)
	}
	return i, e
}

func helpItem(context panel.ItemContext, parent *panel.Panel) (panel.PanelItem, error) {
	// tui.Log("extra", fmt.Sprintf("new helpItem %v", context ))
	I := panel.NewBaseItem(parent)

	txt := widget.NewTextView()
	txt.SetBorderPadding(0,0,1,1)	
	fmt.Fprint(txt, EvalHelpText)

	I.SetWidget(txt)

	return I, nil
}

func playItem(context panel.ItemContext, parent *panel.Panel) (panel.PanelItem, error) {
	tui.Log("extra", fmt.Sprintf("Eval.playItem.context: %v", context ))

	args := []string{}
	if _args, ok := context["args"]; ok {
		args = _args.([]string)
	}

	play := playground.New("")
	play.HandleAction("create", args, context)

	I := panel.NewBaseItem(parent)
	I.SetWidget(play)

	return I, nil
}

func viewItem(context panel.ItemContext, parent *panel.Panel) (panel.PanelItem, error) {
	// tui.Log("extra", fmt.Sprintf("new viewItem %v", context ))

	args := []string{}
	if _args, ok := context["args"]; ok {
		args = _args.([]string)
	}

	// get source, defaults to runtime
	source := "runtime"
	if _source, ok := context["source"]; ok {
		source = _source.(string)
	}

	cfg := &helpers.SourceConfig{
		Source: helpers.EvalSource(source),
		Args: args,
	}

	b := browser.New()
	b.AddSourceConfig(cfg)
	b.SetTitle(fmt.Sprintf("  %v  ", args)).SetBorder(true)
	b.RebuildValue()
	b.Rebuild()

	I := panel.NewBaseItem(parent)
	I.SetWidget(b)

	return I, nil
}


//func flowItem(context panel.ItemContext, parent *panel.Panel) (panel.PanelItem, error) {
//  tui.Log("extra", fmt.Sprintf("Eval.flowItem.context: %v", context ))

//  args := []string{}
//  if _args, ok := context["args"]; ok {
//    args = _args.([]string)
//  }

//  flow := flower.New()
//  flow.HandleAction("update", args, context)
//  flow.Rebuild()

//  I := panel.NewBaseItem(parent)
//  I.SetWidget(flow)

//  return I, nil
//}