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

package container

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Image struct {
	ID       string
	Repository string
	Tag string
	RepoTags []string
	Names []string
}

type Container struct {
	Ports  PortList
	Image  string
	State  string
	Status string
	Names  NameList
}

type PortList []int

type NameList []string

type structuredPort struct {
	HostPort int `json:"host_port"`
}

var portExp = regexp.MustCompile(`:(\d+)->`)

func (l *PortList) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return nil
	}

	var ll []int

	switch b[0] {
	case '[':
		var hps []structuredPort
		if err := json.Unmarshal(b, &hps); err != nil {
			return fmt.Errorf("json unmarshal port list structured: %w", err)
		}

		for _, hp := range hps {
			ll = append(ll, hp.HostPort)
		}
	default:
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return fmt.Errorf("json unmarshal port list string: %w", err)
		}

		parts := strings.Split(s, ",")
		for _, p := range parts {
			pp := portExp.FindStringSubmatch(p)
			if len(pp) != 2 {
				continue
			}

			i, err := strconv.Atoi(pp[1])
			if err != nil {
				continue
			}

			ll = append(ll, i)
		}
	}

	*l = ll

	return nil
}

func (l *NameList) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		return nil
	}

	var ll []string
	switch b[0] {
	case '[':
		if err := json.Unmarshal(b, &ll); err != nil {
			return fmt.Errorf("json unmarshal name list array: %w", err)
		}
	default:
		var s string
		if err := json.Unmarshal(b, &s); err != nil {
			return fmt.Errorf("unmarshal as string: %w", err)
		}

		ll = strings.Split(s, ",")
	}

	*l = ll

	return nil
}
