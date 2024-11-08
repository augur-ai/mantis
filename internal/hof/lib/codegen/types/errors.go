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
 * This work is based on code from https://github.com/gofireflyio/aiac, licensed under the MIT License.
 */
package types

import "errors"

var (
	// ErrNoSuchBackend is returned when the user provides a backend name that
	// does not exist in the configuration.
	ErrNoSuchBackend = errors.New("no such backend")

	// ErrNoDefaultBackend is returned when the user does not select a backend,
	// and the configuration file does not define a default backend.
	ErrNoDefaultBackend = errors.New("backend not selected and no default configured")

	// ErrNoDefaultModel is returned when the user does not select a model, and
	// the configuration file does not defined a default model.
	ErrNoDefaultModel = errors.New("model not selected and no default configured")

	// ErrNoResults is returned if the LLM provider API returned an empty
	// result. This should not generally happen.
	ErrNoResults = errors.New("no results returned from API")

	// ErrUnexpectedStatus is returned when the LLM provider API returned a
	// response with an unexpected status code.
	ErrUnexpectedStatus = errors.New("backend returned unexpected response")

	// ErrRequestFailed is returned when the LLM provider API returned an error
	// for the request.
	ErrRequestFailed = errors.New("request failed")
)
