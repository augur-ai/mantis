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

	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/opentofu/opentofu/internal/hof/cmd/hof/cmd/fmt"
	"github.com/opentofu/opentofu/internal/hof/cmd/hof/flags"
	"github.com/opentofu/opentofu/internal/hof/cmd/hof/ga"

	hfmt "github.com/opentofu/opentofu/internal/hof/lib/fmt"
)

var fmtLong = `With hof fmt, you can
  1. format any language from a single tool
  2. run formatters as api servers for IDEs and hof
  3. manage the underlying formatter containers`

func init() {

	flags.SetupFmtFlags(FmtCmd.Flags(), &(flags.FmtFlags))

}

func FmtRun(files []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = hfmt.Run(files, flags.RootPflags, flags.FmtFlags)

	return err
}

var FmtCmd = &cobra.Command{

	Use: "fmt [filepaths or globs]",

	Short: "format any code and manage the formatters",

	Long: fmtLong,

	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		glob := toComplete + "*"
		matches, _ := filepath.Glob(glob)
		return matches, cobra.ShellCompDirectiveDefault
	},

	Run: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

		var err error

		// Argument Parsing

		if 0 >= len(args) {
			fmt.Println("missing required argument: 'files'")
			cmd.Usage()
			os.Exit(1)
		}

		var files []string

		if 0 < len(args) {

			files = args[0:]

		}

		err = FmtRun(files)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	extra := func(cmd *cobra.Command) bool {

		return false
	}

	ohelp := FmtCmd.HelpFunc()
	ousage := FmtCmd.UsageFunc()

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
	FmtCmd.SetHelpFunc(thelp)
	FmtCmd.SetUsageFunc(tusage)

	FmtCmd.AddCommand(cmdfmt.InfoCmd)
	FmtCmd.AddCommand(cmdfmt.PullCmd)
	FmtCmd.AddCommand(cmdfmt.StartCmd)
	FmtCmd.AddCommand(cmdfmt.TestCmd)
	FmtCmd.AddCommand(cmdfmt.StopCmd)

}
