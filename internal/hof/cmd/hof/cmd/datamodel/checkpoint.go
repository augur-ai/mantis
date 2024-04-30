/*
 * Augur AI Proprietary
 * Copyright (c) 2024 Augur AI, Inc.
 *
 * This file is licensed under the Augur AI Proprietary License.
 *
 * Attribution:
 * This work is based on code from https://github.com/hofstadter-io/hof, licensed under the Apache License 2.0.
 */

package cmddatamodel

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/opentofu/opentofu/internal/hof/cmd/hof/flags"
	"github.com/opentofu/opentofu/internal/hof/cmd/hof/ga"

	"github.com/opentofu/opentofu/internal/hof/lib/datamodel/cmd"
)

var checkpointLong = `create a snapshot of the data model`

func init() {

	flags.SetupDatamodel__CheckpointFlags(CheckpointCmd.Flags(), &(flags.Datamodel__CheckpointFlags))

}

func CheckpointRun(args []string) (err error) {

	// you can safely comment this print out
	// fmt.Println("not implemented")

	err = cmd.Run("checkpoint", args, flags.RootPflags, flags.DatamodelPflags)

	return err
}

var CheckpointCmd = &cobra.Command{

	Use: "checkpoint",

	Aliases: []string{
		"cp",
		"x",
	},

	Short: "create a snapshot of the data model",

	Long: checkpointLong,

	Run: func(cmd *cobra.Command, args []string) {

		ga.SendCommandPath(cmd.CommandPath())

		var err error

		// Argument Parsing

		err = CheckpointRun(args)
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

	ohelp := CheckpointCmd.HelpFunc()
	ousage := CheckpointCmd.UsageFunc()

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
	CheckpointCmd.SetHelpFunc(thelp)
	CheckpointCmd.SetUsageFunc(tusage)

}
