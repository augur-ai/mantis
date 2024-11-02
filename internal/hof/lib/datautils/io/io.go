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

package io

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/clbanning/mxj"
	"github.com/naoina/toml"
	"gopkg.in/yaml.v3"
)

/*
Where's your docs doc?!
*/
func ReadAll(reader io.Reader, obj *interface{}) (contentType string, err error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	// the following error checks are opposite the usual
	// we try from most common to least common types
	// xml needs to come first because json also seems to read it

	mv, merr := mxj.NewMapXml(data)
	if merr == nil {
		*obj = map[string]interface{}(mv)
		return "xml", nil
	}

	err = json.Unmarshal(data, obj)
	if err == nil {
		return "json", nil
	}

	if bytes.Contains(data, []byte("---")) {
		ydata := bytes.Split(data, []byte("---"))

		var yslice []interface{}
		for _, yd := range ydata {
			var yobj interface{}
			err = yaml.Unmarshal(yd, &yobj)
			if err != nil {
				return "", err
			}
			if yobj == nil {
				continue
			}
			yslice = append(yslice, yobj)
		}

		*obj = yslice
		return "yaml", nil
	} else {
		err = yaml.Unmarshal(data, obj)
		if err == nil {
			return "yaml", nil
		}
	}

	err = toml.Unmarshal(data, obj)
	if err == nil {
		return "toml", nil
	}

	return "", errors.New("unknown content type")
}

/*
Where's your docs doc?!
*/
func ReadFile(filename string, obj *interface{}) (contentType string, err error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	ext := filepath.Ext(filename)[1:]
	switch ext {

	case "json":
		err = json.Unmarshal(data, obj)
		if err != nil {
			return "", err
		}
		return "json", nil

	case "toml":
		err = toml.Unmarshal(data, obj)
		if err != nil {
			return "", err
		}
		return "toml", nil

	case "xml":
		mv, merr := mxj.NewMapXml(data)
		if merr != nil {
			return "", merr
		}
		*obj = map[string]interface{}(mv)
		return "xml", nil

	case "yaml", "yml":
		if bytes.Contains(data, []byte("---")) {
			ydata := bytes.Split(data, []byte("---"))

			var yslice []interface{}
			for _, yd := range ydata {
				var yobj interface{}
				err = yaml.Unmarshal(yd, &yobj)
				if err != nil {
					return "", err
				}
				if yobj == nil {
					continue
				}
				yslice = append(yslice, yobj)
			}

			*obj = yslice
			return "yaml", nil
		} else {

			err = yaml.Unmarshal(data, obj)

			// yobj, err := yamlB.Read(bytes.NewReader(data))
			if err == nil {
				// *obj = yobj
				return "yaml", nil
			}
		}

	// TODO, CUE once we can create values and share between runtimes

	default:
		return InferDataContentType(data)
	}

	return "", errors.New("unknown content type")
}
