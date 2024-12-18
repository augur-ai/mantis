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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	cases := []struct {
		desc    string
		mod     string
		outKind Kind
	}{
		{
			desc:    "git github",
			mod:     "github.com/opentofu/opentofu/internal/hof",
			outKind: KindGit,
		},
		{
			desc:    "git github private",
			mod:     "github.com/andrewhare/env",
			outKind: KindGit,
		},
		{
			desc:    "git not github",
			mod:     "git.kernel.org/pub/scm/bluetooth/bluez.git",
			outKind: KindGit,
		},
		{
			desc:    "oci",
			mod:     "gcr.io/distroless/static-debian11",
			outKind: KindOCI,
		},
		{
			desc:    "oci",
			mod:     "us-central1-docker.pkg.dev/hof-io--develop/testing/test",
			outKind: KindOCI,
		},
	}

	for _, c := range cases {
		c := c

		t.Run(c.desc, func(t *testing.T) {
			t.Parallel()

			out, err := Parse(c.mod)
			require.NoError(t, err)
			assert.Equal(t, c.outKind, out.kind)
		})
	}
}
