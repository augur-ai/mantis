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

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/opentofu/opentofu/internal/hof/cmd/hof/flags"

	"github.com/opentofu/opentofu/internal/hof/cmd/hof/ga"

	"github.com/opentofu/opentofu/internal/hof/lib/cuecmd"
)

var exportLong = `output data in a standard format`

func init() {

	flags.SetupExportFlags(ExportCmd.Flags(), &(flags.ExportFlags))

}

func ExportRun(args []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = cuecmd.Export(args, flags.RootPflags, flags.ExportFlags)

	return err
}

var ExportCmd = &cobra.Command{

	Use: "export",

	Short: "output data in a standard format",

	Long: exportLong,

	Run: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

		var err error

		// Argument Parsing

		err = ExportRun(args)
		if err != nil {
			// fmt.Println(err)
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	extra := func(cmd *cobra.Command) bool {

		return false
	}

	ohelp := ExportCmd.HelpFunc()
	ousage := ExportCmd.UsageFunc()

	help := func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath() + " help")

		if extra(cmd) {
			return
		}
		ohelp(cmd, args)
	}
	usage := func(cmd *cobra.Command) error {
		if extra(cmd) {
			return nil
		}
		return ousage(cmd)
	}

	thelp := func(cmd *cobra.Command, args []string) {
		help(cmd, args)
	}
	tusage := func(cmd *cobra.Command) error {
		return usage(cmd)
	}
	ExportCmd.SetHelpFunc(thelp)
	ExportCmd.SetUsageFunc(tusage)

}
