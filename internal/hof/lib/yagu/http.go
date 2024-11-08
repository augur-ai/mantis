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

package yagu

import (
	"errors"
	"fmt"
	"strings"

	"github.com/parnurzeal/gorequest"
)

func BuildRequest(url string) *gorequest.SuperAgent {

	req := gorequest.New().Get(url)

	return req
}

const HTTP2_GOAWAY_CHECK = "http2: server sent GOAWAY and closed the connection"

func SimpleGet(url string) (string, error) {

	req := BuildRequest(url)
	resp, body, errs := req.End()

	if len(errs) != 0 && !strings.Contains(errs[0].Error(), HTTP2_GOAWAY_CHECK) {
		fmt.Println("errs:", errs)
		fmt.Println("resp:", resp)
		fmt.Println("body:", body)
		return body, errs[0]
	}

	if len(errs) != 0 || resp.StatusCode >= 500 {
		return body, errors.New("Internal Error: " + body)
	}
	if resp.StatusCode >= 400 {
		return body, errors.New("Bad Request: " + body)
	}

	return body, nil
}
