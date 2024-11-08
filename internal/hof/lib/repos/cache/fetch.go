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

package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/osfs"
	gogit "github.com/go-git/go-git/v5"

	"github.com/opentofu/opentofu/internal/hof/lib/repos/remote"
	"github.com/opentofu/opentofu/internal/hof/lib/repos/utils"
)

func OpenRepoSource(path string) (*gogit.Repository, error) {
	if debug {
		fmt.Println("cache.OpenRepoSource:", path)
	}

	remote, owner, repo := utils.ParseModURL(path)
	dir := SourceOutdir(remote, owner, repo)
	return gogit.PlainOpen(dir)
}

func FetchRepoSource(mod, ver string) (billy.Filesystem, error) {
	if debug {
		fmt.Println("cache.FetchRepoSource:", mod)
	}

	rmt, err := remote.Parse(mod)
	if err != nil {
		return nil, fmt.Errorf("remote parse: %w", err)
	}

	dir := SourceOutdirParts(rmt.Host, rmt.Owner, rmt.Name)

	// TODO:
	//   * Use a passed-in context.
	//   * Choose a better timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// only fetch if we haven't already this run
	if _, ok := syncedRepos.Load(mod); !ok {
		if err := rmt.Pull(ctx, dir, ver); err != nil {
			return nil, fmt.Errorf("remote pull: %w", err)
		}

		syncedRepos.Store(mod, true)
	}

	return osfs.New(dir), nil
}

func FetchOCISource(mod, ver string) (billy.Filesystem, error) {
	if debug {
		fmt.Println("cache.FetchOCISource:", mod, ver)
	}

	// upgrade pseudo version
	s, err := UpgradePseudoVersion(mod, ver)
	if err != nil {
		return nil, err
	}
	ver = s

	if debug {
		fmt.Println("cache.FetchOCISource version resolve:", mod, ver)
	}



	rmt, err := remote.Parse(mod)
	if err != nil {
		return nil, fmt.Errorf("remote parse: %w", err)
	}

	dir := ModuleOutdir(rmt.Host, rmt.Owner, rmt.Name, ver)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// only fetch if we haven't already this run
	if _, ok := syncedRepos.Load(mod); !ok {
		if err := rmt.Pull(ctx, dir, ver); err != nil {
			return nil, fmt.Errorf("remote pull: %w", err)
		}

		syncedRepos.Store(mod, true)
	}

	return osfs.New(dir), nil
}
