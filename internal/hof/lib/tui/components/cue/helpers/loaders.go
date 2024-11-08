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

package helpers

import (
	"fmt"
	"os"
	"net/url"
	"strings"

	"cuelang.org/go/cue"
	"github.com/spf13/pflag"

	"github.com/opentofu/opentofu/internal/hof/cmd/hof/flags"
	"github.com/opentofu/opentofu/internal/hof/lib/cuetils"
	"github.com/opentofu/opentofu/internal/hof/lib/runtime"
	"github.com/opentofu/opentofu/internal/hof/lib/singletons"
	"github.com/opentofu/opentofu/internal/hof/lib/tui"
	"github.com/opentofu/opentofu/internal/hof/lib/yagu"
)

func LoadRuntime(args []string) (*runtime.Runtime, error) {
	// tui.Log("trace", fmt.Sprintf("Panel.loadRuntime.inputs: %v", args))

	// build eval args & flags from the input args
	var (
		rflags flags.RootPflagpole
		cflags flags.EvalFlagpole
	)
	fset := pflag.NewFlagSet("panel", pflag.ContinueOnError)
	flags.SetupRootPflags(fset, &rflags)
	flags.SetupEvalFlags(fset, &cflags)
	fset.Parse(args)
	args = fset.Args()

	// tui.Log("trace", fmt.Sprintf("Panel.loadRuntime.parsed: %v %v", args, rflags))

	R, err := runtime.New(args, rflags)
	if err != nil {
		tui.Log("error", cuetils.ExpandCueError(err))
		return R, err
	}

	err = R.Load()
	if err != nil {
		tui.Log("error", cuetils.ExpandCueError(err))
		return R, err
	}

	return R, nil
}

func LoadFromText(content string) (string, cue.Value, error) {

	ctx := singletons.CueContext()
	v := ctx.CompileString(content, cue.Filename("SourceConfig.Text"))

	return content, v, nil
}

func LoadFromFile(filename string) (string, cue.Value, error) {

	ctx := singletons.CueContext()
	b, err := os.ReadFile(filename)
	if err != nil {
		return string(b), singletons.EmptyValue(), err
	}
	v := ctx.CompileBytes(b, cue.Filename(filename))

	return string(b), v, nil
}

func LoadFromHttp(fullurl string) (string, cue.Value, error) {
	// tui.Log("trace", fmt.Sprintf("Panel.loadHttpValue: %s %s", mode, from))

	// rework any cue/play links
	f := fullurl
	if strings.Contains(fullurl, "cuelang.org/play") {
		u, err := url.Parse(fullurl)
		if err != nil {
			tui.Log("error", err)
			return "", singletons.EmptyValue(), err
		}
		q, err := url.ParseQuery(u.RawQuery)
		if err != nil {
			tui.Log("error", err)
			return "", singletons.EmptyValue(), err
		}
		id := q["id"][0]
		f = fmt.Sprintf("https://%s/.netlify/functions/snippets?id=%s", u.Host, id)
	}

	// fetch content
	header := "// from: " + fullurl + "\n\n"
	content, err := yagu.SimpleGet(f)
	content = header + content

	if err != nil {
		return content, singletons.EmptyValue(), fmt.Errorf("%s -- %w", header, err)
	}


	// rebuild, TODO, if scope, use that value and scope.Context() here
	ctx := singletons.CueContext()
	v := ctx.CompileString(content, cue.InferBuiltins(true))

	return content, v, nil
}

func LoadFromBash(args []string) (string, cue.Value, error) {

	wd, err := os.Getwd()
	if err != nil {
		return "", singletons.EmptyValue(), err
	}

	script := strings.Join(args, " ")
	out, err := yagu.Bash(script, wd)
	if err != nil {
		return "", singletons.EmptyValue(), err
	}

	// TODO, infer output type, support yaml too

	header := "// bash " + strings.Join(args, " ") + "\n\n"
	out = header + out 

	// compile CUE (json, but all json is CUE, which is why we can add a comment)
	ctx := singletons.CueContext()
	v := ctx.CompileString(out, cue.InferBuiltins(true))

	return out, v, nil
}
