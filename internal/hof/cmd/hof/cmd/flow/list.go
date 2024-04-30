/*
 * Augur AI Proprietary
 * Copyright (c) 2024 Augur AI, Inc.
 *
 * This file is licensed under the Augur AI Proprietary License.
 *
 * Attribution:
 * This work is based on code from https://github.com/hofstadter-io/hof, licensed under the Apache License 2.0.
 */

package cmdflow

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/opentofu/opentofu/internal/hof/cmd/hof/ga"

	"github.com/opentofu/opentofu/internal/hof/cmd/hof/flags"
	"github.com/opentofu/opentofu/internal/hof/flow/cmd"
)

var listLong = `print available flows`

func ListRun(args []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = cmd.List(args, flags.RootPflags, flags.FlowPflags)

	return err
}

var ListCmd = &cobra.Command{

	Use: "list",

	Short: "print available flows",

	Long: listLong,

	Run: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

		var err error

		// Argument Parsing

		err = ListRun(args)
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

	ohelp := ListCmd.HelpFunc()
	ousage := ListCmd.UsageFunc()

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
	ListCmd.SetHelpFunc(thelp)
	ListCmd.SetUsageFunc(tusage)

}
