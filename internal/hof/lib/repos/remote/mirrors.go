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

package remote

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/opentofu/opentofu/internal/hof/lib/repos/git"
	"github.com/opentofu/opentofu/internal/hof/lib/repos/oci"
)

const debug = 0

var (
	mirrorsGit = []string{
		"github.com",
		"gitlab.com",
		"bitbucket.org",
	}
	mirrorsOCI = []string{
		"ghcr.io",
	}

	MirrorsSingleton *Mirrors
)

func init() {
	// TODO, this should be a singleton for the application
	// right now, we read the file every time we parse a mod path
	var err error
	MirrorsSingleton, err = NewMirrors()
	if err != nil {
		panic(err)
		// return nil, fmt.Errorf("new mirrors: %w", err)
	}

}

const (
	hofDir             = "hof"
	mirrorsFileName    = "mirrors.json"
	mirrorsFileNameEnv = "HOF_MOD_MIRRORFILE"
)

// Holds the mirror mappings
// todo, make the value more interesting
// we probably want to know if pub/priv & auth settings
// and what to mirror where, this only says what repo kind it is
// prehaps we should make the primary key the url prefix,
// and then look up details about [reg-type,reg-auth,reg-url,mods-mirrored]
type Mirrors struct {
	valuesMu sync.RWMutex
	values   map[Kind][]string
}

// tbd
type Mirror struct {
	Kind Kind
	URL  string	
	Auth any
	Prefixes []string
}

func mirrorsFilePath() (string, error) {
	p := os.Getenv(mirrorsFileNameEnv)
	if p != "" {
		return p, nil
	}

	d, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("user cache dir: %w", err)
	}

	return filepath.Join(d, hofDir, mirrorsFileName), nil
}

func NewMirrors() (*Mirrors, error) {
	p, err := mirrorsFilePath()
	if err != nil {
		return nil, fmt.Errorf("mirrors file path: %w", err)
	}

	info, err := os.Stat(p)
	if errors.Is(err, os.ErrNotExist) || info.Size() == 0 {
		return &Mirrors{values: make(map[Kind][]string)}, nil
	}

	f, err := os.Open(p)
	if err != nil {
		return nil, fmt.Errorf("os open %s: %w", p, err)
	}

	defer f.Close()

	var m Mirrors
	if err := json.NewDecoder(f).Decode(&m.values); err != nil {
		return nil, fmt.Errorf("json decode %s: %w", p, err)
	}

	return &m, nil
}

func (m *Mirrors) Is(ctx context.Context, k Kind, mod string) (bool, error) {
	if debug > 0 {
		fmt.Println("mirrors.Is", k, mod)	
	}
	var (
		mirrors  []string
		netCheck func(context.Context, string) (bool, error)
	)

	switch k {
	case KindGit:
		mirrors = mirrorsGit
		netCheck = m.netCheckGit
	case KindOCI:
		mirrors = mirrorsOCI
		netCheck = m.netCheckOCI
	default:
		return false, fmt.Errorf("unknow kind: %s", k)
	}

	for _, ss := range mirrors {
		if strings.HasPrefix(mod, ss) {
			return true, nil
		}
	}

	if m.hasValue(k, mod) {
		return true, nil
	}

	// TODO, think through conditions here
	// the error was taking priority over false
	is, err := netCheck(ctx, mod)
	if !is && err != nil {
		return false, err
	}

	if is {
		m.valuesMu.Lock()
		m.values[k] = append(m.values[k], mod)
		m.valuesMu.Unlock()
	}

	return is, nil
}

func (m *Mirrors) hasValue(k Kind, mod string) bool {
	m.valuesMu.RLock()
	defer m.valuesMu.RUnlock()

	if vals, ok := m.values[k]; ok {
		for _, v := range vals {
			if strings.HasPrefix(mod, v) {
				return true
			}
		}
	}

	return false
}

func (m *Mirrors) netCheckGit(ctx context.Context, mod string) (bool, error) {
	return git.IsNetworkReachable(ctx, mod)
}

func (m *Mirrors) netCheckOCI(ctx context.Context, mod string) (bool, error) {
	return oci.IsNetworkReachable(mod)
}

func (m *Mirrors) Set(k Kind, s string) {
	m.valuesMu.Lock()
	defer m.valuesMu.Unlock()

	if vals, ok := m.values[k]; ok {
		vals = append(vals, s)
		m.values[k] = vals
	}
}

func (m *Mirrors) Close() error {
	m.valuesMu.Lock()
	defer m.valuesMu.Unlock()

	if len(m.values) == 0 {
		return nil
	}

	p, err := mirrorsFilePath()
	if err != nil {
		return fmt.Errorf("mirrors file path: %w", err)
	}

	dir := filepath.Dir(p)
	if err = os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("mkdir %s: %w", dir, err)
	}

	f, err := os.OpenFile(p, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("open %s: %w", p, err)
	}
	defer f.Close()

	e := json.NewEncoder(f)
	e.SetIndent("", "  ")

	if err := e.Encode(&m.values); err != nil {
		return fmt.Errorf("json encode %s: %w", p, err)
	}

	return nil
}
