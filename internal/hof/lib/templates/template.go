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

package templates

import (
	"bytes"
	"text/template"
)

type Delims struct {
	LHS string
	RHS string
}

type Template struct {
	// Original inputs
	Name   string
	Source string
	Delims Delims

	// golang
	T *template.Template

	Buf *bytes.Buffer
}

func (T *Template) Render(data interface{}) ([]byte, error) {
	// endure we don't have nil, if so, there is a bug somewhere
	if T.T == nil {
		panic("template not set!")
	}

	var err error

	T.Buf.Reset()

	err = T.T.Execute(T.Buf, data)
	if err != nil {
		return nil, err
	}

	// we need to get a string
	// and then turn it into bytes
	// to work around a memory issue
	// with bytes.Buffer
	out := T.Buf.String()
	bs := []byte(out)

	return bs, nil
}

// Creates a hof Template struct, initializing the correct template system. The system will be inferred if left empty
func CreateFromString(name, content string, delims Delims) (t *Template, err error) {
	t = new(Template)
	t.Name = name
	t.Source = content

	// Golang wants helpers before parsing, and catches these errors early
	t.T = template.New(name)

	if delims.LHS != "" {
		t.T = t.T.Delims(delims.LHS, delims.RHS)
	}

	t.Buf = new(bytes.Buffer)

	t.AddGolangHelpers()

	t.T, err = t.T.Parse(content)

	return t, err
}
