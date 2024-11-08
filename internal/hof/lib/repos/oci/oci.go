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

package oci

import (
	"archive/tar"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	// "path"
	"path/filepath"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/remote/transport"
	"github.com/google/go-containerregistry/pkg/v1/stream"
	"github.com/google/go-containerregistry/pkg/v1/types"
)

const (
	HofstadterSchema1Beta types.MediaType = "application/vnd.hofstadter.module.v1beta1+json"
	HofstadterModuleDeps  types.MediaType = "application/vnd.hofstadter.module.deps.tar.gz"
	HofstadterModuleCode  types.MediaType = "application/vnd.hofstadter.module.code.tar.gz"
)

var debug = false

func IsNetworkReachable(mod string) (bool, error) {
	_, err := crane.Manifest(mod, crane.WithAuthFromKeychain(authn.DefaultKeychain))

	var terr *transport.Error
	if errors.As(err, &terr) {
		if len(terr.Errors) != 1 {
			return false, fmt.Errorf("multiple transport errors: %w", terr)
		}

		switch c := terr.Errors[0].Code; c {
		case transport.ManifestUnknownErrorCode:
			return true, nil
		case transport.NameUnknownErrorCode:
			return false, errors.New("remote repo does not exist")
		default:
			return false, fmt.Errorf("unhandled transport code: %s", c)
		}
	}

	return err == nil, err
}

func ListTags(mod string) ([]string, error) {
	return crane.ListTags(mod, crane.WithAuthFromKeychain(authn.DefaultKeychain))
}

// Looks up a Ref and returns the hash it currently points at
// we recommend you setup a registry with immutable tags
func GetRefHash(url, ref string) (string, error) {
	if debug {
		fmt.Println("oci.GetRefHash:", url, ref)
	}
	p := url + ":" + ref
	r, err := name.ParseReference(p)
	if err != nil {
		return "", fmt.Errorf("whil parsing reference: %w", err)
	}

	img, err := remote.Image(r, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return "", fmt.Errorf("while finding remote image in oci.GetRefHash: %w", err)
	}

	hash, err := img.Digest()
	if err != nil {
		return "", fmt.Errorf("error getting hash in oci.GetRefHash: %s %s: %w", url, ref, err)
	}

	// trim hash algo off of front
	s := hash.String()
	pos := strings.Index(s, ":")
	s = s[pos+1:]

	return s, nil
}

func Pull(url, outPath string) error {
	if debug {
		fmt.Println("oci.Pull:", outPath, url)
	}
	p := strings.Index(url, "@")
	P := url[:p]
	fmt.Println("fetch'n:", P)
	ref, err := name.ParseReference(url)
	if err != nil {
		return fmt.Errorf("name parse reference: %w", err)
	}

	img, err := remote.Image(ref, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		return fmt.Errorf("remote image: %w", err)
	}

	r := mutate.Extract(img)
	defer r.Close()

	if err := untar(r, outPath); err != nil {
		return fmt.Errorf("untar: %w", err)
	}

	return nil
}

func untar(r io.Reader, target string) error {
	tr := tar.NewReader(r)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("tar reader next: %w", err)
		}

		var (
			p = filepath.Join(target, header.Name)
			i = header.FileInfo()
		)

		if i.IsDir() {
			if err = os.MkdirAll(p, 0755); err != nil {
				return fmt.Errorf("mkdir all: %w", err)
			}
			continue
		}

		// mkdir for file, in case we didn't get it first in the tar before the file
		if err = os.MkdirAll(filepath.Dir(p), 0755); err != nil {
			return fmt.Errorf("mkdir all: %w", err)
		}

		f, err := os.OpenFile(p, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("open file: %w", err)
		}

		defer f.Close()

		if _, err = io.Copy(f, tr); err != nil {
			return fmt.Errorf("io copy: %w", err)
		}
	}

	return nil
}

func Push(tag string, img v1.Image) error {
	ref, err := name.ParseReference(tag)
	if err != nil {
		return fmt.Errorf("name parse reference: %w", err)
	}

	fmt.Println("pushing...")
	if err = remote.Write(ref, img, remote.WithAuthFromKeychain(authn.DefaultKeychain)); err != nil {
		return fmt.Errorf("remote write: %w", err)
	}

	return nil
}

func Build(workingDir string, dirs []Dir) (v1.Image, error) {
	var layers []v1.Layer

	for _, d := range dirs {
		// todo, enable printing base on verbosity
		fmt.Println("layer:", d.mediaType)
		l, err := layer(workingDir, d)
		if err != nil {
			return nil, fmt.Errorf("layer: %w", err)
		}

		layers = append(layers, l)
	}

	e := mutate.MediaType(empty.Image, types.OCIManifestSchema1)
	e = mutate.ConfigMediaType(e, HofstadterSchema1Beta)

	img, err := mutate.AppendLayers(e, layers...)
	if err != nil {
		return nil, fmt.Errorf("append layers: %w", err)
	}

	return img, nil
}

func layer(wd string, d Dir) (v1.Layer, error) {
	var (
		buf bytes.Buffer
		w   = tar.NewWriter(&buf)
	)

	err := filepath.Walk(d.relPath, func(p string, i os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if d.Excluded(p) {
			return nil
		}
		// todo, print included filename based on verbosity
		fmt.Println(" ", p)

		h, err := tar.FileInfoHeader(i, "")
		if err != nil {
			return fmt.Errorf("tar file info header: %w", err)
		}

		h.Name = strings.ReplaceAll(p, wd, "")

		if err = w.WriteHeader(h); err != nil {
			return fmt.Errorf("tar write header: %w", err)
		}

		if i.IsDir() {
			return nil
		}

		f, err := os.Open(p)
		if err != nil {
			return fmt.Errorf("open %s: %w", p, err)
		}

		defer f.Close()

		if _, err = io.Copy(w, f); err != nil {
			return fmt.Errorf("copy %s: %w", p, err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("filepath walk: %w", err)
	}

	if err = w.Close(); err != nil {
		return nil, fmt.Errorf("tar writer close: %w", err)
	}

	var (
		rc = io.NopCloser(&buf)
		mt = stream.WithMediaType(d.mediaType)
	)

	return stream.NewLayer(rc, mt), nil
}
