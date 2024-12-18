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

package prompt

import (
	"fmt"

	"cuelang.org/go/cue"
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)

type handler func (Q map[string]any) (A any, err error)

var handlers map[string]handler

func init() {
	handlers = map[string]handler{
		// builtin handlers
		"input": handleInput,
		"multiline": handleMultiline,
		"password": handlePassword,
		"confirm": handleConfirm,
		"select": handleSelect,
		"multiselect": handleMultiselect,
		"subgroup": handleSubgroup,

		// custom handlers
	}
}

// this is for backwards compat, create wants things a little differently
func RunCreatePrompt(genVal cue.Value) (result cue.Value, err error) {
	create := genVal.LookupPath(cue.ParsePath("Create"))
	if create.Err() != nil {
		return genVal, create.Err()
	}

	questions := genVal.LookupPath(cue.ParsePath("Create.Questions"))
	if !questions.IsConcrete() || !questions.Exists() {
		// handle old create.Prompt
		prompt := genVal.LookupPath(cue.ParsePath("Create.Prompt"))
		if !prompt.IsConcrete() || !prompt.Exists() {
			// if no questions, just return the original
			return genVal, nil
		}
		// fill in Questions with Prompt if we get here
		create = create.FillPath(cue.ParsePath("Questions"), prompt)
	}

	r, err := RunPrompt(create)

	// fill the input back in with the output
	// TODO, change create to use Output
	genVal = genVal.FillPath(cue.ParsePath("Create"), r)
	outputVal := genVal.LookupPath(cue.ParsePath("Create.Output"))
	genVal = genVal.FillPath(cue.ParsePath("Create.Input"), outputVal)

	// finally return genVal and the prompt error if any
	return genVal, err
}

func RunPrompt(genVal cue.Value) (result cue.Value, err error) {
	// run while there are unanswered questions
	// we run in an extra loop to fill back answers
	// and recalculate the prompt questions each iteration

	// first, fill Output with Input
	inputVal := genVal.LookupPath(cue.ParsePath("Input"))
	if inputVal.Err() != nil {
		return genVal, inputVal.Err()
	}
	genVal = genVal.FillPath(cue.ParsePath("Output"), inputVal)

	done := false

	for !done {
		done = true


		outputVal := genVal.LookupPath(cue.ParsePath("Output"))
		if outputVal.Err() != nil {
			return genVal, outputVal.Err()
		}
		// fmt.Printf("outer loop input: %#v\n", inputVal)

		questions := genVal.LookupPath(cue.ParsePath("Questions"))
		if questions.Err() != nil {
			return genVal, questions.Err()
		}
		if !questions.IsConcrete() || !questions.Exists() {
			// to have a promptless generator, set it to the empty list
			return genVal, fmt.Errorf("missing Questions field")
		}

		// prompt should be an ordered list of questions
		iter, err := questions.List()
		if err != nil {
			return genVal, err
		}

		// loop over prompt questions, recursing as needed
		for iter.Next() {
			// todo, get label and check if input[label] is concrete
			value := iter.Value()

			Q := map[string]any{}
			err := value.Decode(&Q)
			if err != nil {
				return genVal, err
			}

			// fmt.Printf("%#v\n", Q)

			name := Q["Name"].(string)
			namePath := cue.ParsePath(name)

			// check if done already by inspececting in input
			i := outputVal.LookupPath(namePath)
			if i.Err() != nil {
				if i.Exists() {
					return genVal, i.Err()
				}
			}
			if i.Exists() && i.IsConcrete() {
				// question answer already exists in input
				// fmt.Println("continuing: ", name)
				continue
			}

			// there is a question to answer
			done = false

			// fmt.Println("q:", Q)
			// todo, extract Name
			A, err := handleQuestion(Q)
			if err != nil {
				if err == terminal.InterruptErr {
					return genVal, fmt.Errorf("user interrupt")
				}
				return genVal, err
			}

			// update input val
			outputVal = outputVal.FillPath(namePath, A)
			genVal = genVal.FillPath(cue.ParsePath("Output"), outputVal)

			// restart the prompt loop
			break
		}
	}

	return genVal, nil
}

func handleQuestion(Q map[string]any) (A any, err error) {
	// ask question until we get an answer or interrupt
	for {
		s, ok := Q["Type"]
		if !ok {
			panic("question type not set")
		}

		S, ok := s.(string)
		if !ok {
			panic("question 'Type' is not set using a string format")
		}

		h, ok := handlers[S]
		if !ok {
			panic("unknown question type: " + S)
		}

		a, err := h(Q)

		if err != nil {
			if err == terminal.InterruptErr {
				return nil, err
			}
			fmt.Println("error:", err)
			continue
		}

		A = a

		// we got an answer
		break
	}

	return A, nil
}

func handleInput(Q map[string]any) (A any, err error) {
	dval := ""
	if d, ok := Q["Default"]; ok {
		dval = d.(string)
	}
	prompt := &survey.Input {
		// todo, rename prompt to message in test and schema
		Message: Q["Prompt"].(string),
		Default: dval,
	}
	var a string
	err = survey.AskOne(prompt, &a)
	A = a

	return A, err
}

func handleMultiline(Q map[string]any) (A any, err error) {
	dval := ""
	if d, ok := Q["Default"]; ok {
		dval = d.(string)
	}
	prompt := &survey.Multiline {
		// todo, rename prompt to message in test and schema
		Message: Q["Prompt"].(string),
		Default: dval,
	}
	var a string
	err = survey.AskOne(prompt, &a)
	A = a

	return A, err
}

func handlePassword(Q map[string]any) (A any, err error) {
	prompt := &survey.Password {
		// todo, rename prompt to message in test and schema
		Message: Q["Prompt"].(string),
	}
	var a string
	err = survey.AskOne(prompt, &a)
	A = a

	return A, err
}

func handleConfirm(Q map[string]any) (A any, err error) {
	dval := false
	if d, ok := Q["Default"]; ok {
		dval = d.(bool)
	}
	prompt := &survey.Confirm {
		// todo, rename prompt to message in test and schema
		Message: Q["Prompt"].(string),
		Default: dval,
	}
	var a bool
	err = survey.AskOne(prompt, &a)
	if err != nil {
		if err == terminal.InterruptErr {
			return nil, err
		}
	}
	A = a

	if err != nil {
		return A, err
	}

	// possibly recurse if has own Questions
	QS, ok := Q["Questions"]
	if a && ok {
		A2 := map[string]any{}
		for _, Q2 := range QS.([]any) {
			q2 := Q2.(map[string]any)
			a2, e2 := handleQuestion(q2)
			// todo, think about if/how to handle this error
			if e2 != nil {
				return nil, e2
			}
			A2[q2["Name"].(string)] = a2
		}
		A = A2
	}

	return A, err
}

func handleSelect(Q map[string]any) (A any, err error) {
	opts := []string{}
	for _, o := range Q["Options"].([]any) {
		opts = append(opts, o.(string))
	}
	prompt := &survey.Select {
		// todo, rename prompt to message in test and schema
		Message: Q["Prompt"].(string),
		// todo, probably need to error handle options
		Options: opts,
		Default: Q["Default"],
	}
	var a string
	err = survey.AskOne(prompt, &a)
	A = a

	return A, err
}
	
func handleMultiselect(Q map[string]any) (A any, err error) {
	opts := []string{}
	for _, o := range Q["Options"].([]any) {
		opts = append(opts, o.(string))
	}
	prompt := &survey.MultiSelect {
		// todo, rename prompt to message in test and schema
		Message: Q["Prompt"].(string),
		// todo, probably need to error handle options
		Options: opts,
		Default: Q["Default"],
	}
	var a []string
	err = survey.AskOne(prompt, &a)
	A = a

	return A, err
}

func handleSubgroup(Q map[string]any) (A any, err error) {
	// get some strings
	name := Q["Name"].(string)
	msg := Q["Prompt"].(string)
	fmt.Println(msg)

	// get subgroup Questions
	QS, ok := Q["Questions"]
	if !ok {
		return nil, fmt.Errorf("subgroup prompt %q is missing 'Questions' field", name)
	}

	// gather nested answers
	A2 := map[string]any{}
	for _, Q2 := range QS.([]any) {
		q2 := Q2.(map[string]any)
		a2, e2 := handleQuestion(q2)
		// todo, think about if/how to handle this error
		if e2 != nil {
			return nil, e2
		}
		A2[q2["Name"].(string)] = a2
	}

	// set as nested value in A
	a := map[string]any{}
	a[name] = A2

	return a, err
}
