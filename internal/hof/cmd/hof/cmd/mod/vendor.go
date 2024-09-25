/*
 * Augur AI Proprietary
 * Copyright (c) 2024 Augur AI, Inc.
 *
 * This file is licensed under the Augur AI Proprietary License.
 *
 * Attribution:
 * This work is based on code from https://github.com/hofstadter-io/hof, licensed under the Apache License 2.0.
 */

package cmdmod

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/opentofu/opentofu/internal/hof/lib/mod"

	"github.com/opentofu/opentofu/internal/hof/cmd/hof/flags"

	"github.com/opentofu/opentofu/internal/hof/cmd/hof/ga"
)

var vendorLong = `copy dependencies to cue.mod/pkg`

func VendorRun(args []string) (err error) {

	err = mod.Vendor(flags.RootPflags)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return err
}

var VendorCmd = &cobra.Command{

	Use: "vendor",

	Short: "copy dependencies to cue.mod/pkg",

	Long: vendorLong,

	Run: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

		var err error

		// Argument Parsing

		err = VendorRun(args)
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

	ohelp := VendorCmd.HelpFunc()
	ousage := VendorCmd.UsageFunc()

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
	VendorCmd.SetHelpFunc(thelp)
	VendorCmd.SetUsageFunc(tusage)

}