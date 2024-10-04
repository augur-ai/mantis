/*
 * Augur AI Proprietary
 * Copyright (c) 2024 Augur AI, Inc.
 *
 * This file is licensed under the Augur AI Proprietary License.
 *
 * Attribution:
 * This work is based on code from https://github.com/hofstadter-io/hof, licensed under the Apache License 2.0.
 */

package config

import (
	"fmt"
	//"os"
	//"strings"

	"cuelang.org/go/cue"
	//"github.com/opentofu/opentofu/internal/hof/cmd/hof/flags"
	//"github.com/opentofu/opentofu/internal/hof/gen/cuefig"
	//"github.com/opentofu/opentofu/internal/hof/lib/structural"
)

var rt *Runtime

func init() {
	rt = NewRuntime()
}

func Init() error {
	r := NewRuntime()
	err := r.Init()
	if err != nil {
		return err
	}

	rt = r

	return nil
}

func GetRuntime() *Runtime {
	return rt
}

// Runtime holds the app config/secrets
type Runtime struct {
	ContextType  string
	ContextValue cue.Value
	ConfigType   string
	ConfigValue  cue.Value
	SecretType   string
	SecretValue  cue.Value
}

func NewRuntime() *Runtime {
	return &Runtime{}
}

// TODO Load user/app config/secret

// We can safely ignore errors here. If the file exists, cue errors will be printed, otherwise up to the user
func (R *Runtime) Init() (err error) {
	// These are used to track if we found a file or not
	// contextFound, configFound, secretFound := false, false, false

	// First check config/secret flags, non-existence should err as user specified a flag
	//  if they exist, we load into local because we prefer that later
	//if flags.RootPflags.Context != "" {
	//val, err := cuefig.LoadContextConfig("", flags.RootPflags.Context)
	//if err != nil {
	//// Return early if they specify a file and we don't find it
	//return err
	//}
	//contextFound = true
	//R.ContextValue = val
	//R.ContextType = "custom-context"
	//}
	//if flags.RootPflags.Config != "" {
	//val, err := cuefig.LoadConfigConfig("", flags.RootPflags.Config)
	//if err != nil {
	//// Return early if they specify a file and we don't find it
	//return err
	//}
	//configFound = true
	//R.ConfigValue = val
	//R.ConfigType = "custom-config"
	//}
	//if flags.RootPflags.Secret != "" {
	//val, err := cuefig.LoadSecretConfig("", flags.RootPflags.Secret)
	//if err != nil {
	//// Return early if they specify a file and we don't find it
	//return err
	//}
	//secretFound = true
	//R.SecretValue = val
	//R.SecretType = "custom-secret"
	//}

	// Second, look for local config/secret
	//if !contextFound {
	//val, err := cuefig.LoadContextDefault()
	//// NOTE, we are doing the opposite of normal err checks here
	//if err == nil {
	//configFound = true
	//R.ContextValue = val
	//R.ContextType = "local-context"
	//}
	//}
	//if !configFound {
	//val, err := cuefig.LoadConfigDefault()
	//// NOTE, we are doing the opposite of normal err checks here
	//if err == nil {
	//configFound = true
	//R.ConfigValue = val
	//R.ConfigType = "local-config"
	//}
	//}
	//if !secretFound {
	//val, err := cuefig.LoadSecretDefault()
	//// NOTE, we are doing the opposite of normal err checks here
	//if err == nil {
	//secretFound = true
	//R.SecretValue = val
	//R.SecretType = "local-secret"
	//}
	//}

	// Finally, check for global config/secret
	//if !contextFound {
	//val, err := cuefig.LoadHofctxDefault()
	//// NOTE, we are doing the opposite of normal err checks here
	//if err == nil {
	//contextFound = true
	//R.ContextValue = val
	//R.ContextType = "global-context"
	//}
	//}
	//if !configFound {
	//val, err := cuefig.LoadHofcfgDefault()
	//// NOTE, we are doing the opposite of normal err checks here
	//if err == nil {
	//configFound = true
	//R.ConfigValue = val
	//R.ConfigType = "global-config"
	//}
	//}
	//if !secretFound {
	//val, err := cuefig.LoadHofshhDefault()
	//// NOTE, we are doing the opposite of normal err checks here
	//if err == nil {
	//secretFound = true
	//R.SecretValue = val
	//R.SecretType = "global-secret"
	//}
	//}

	return err
}

func (R *Runtime) PrintConfig() error {
	// Get top level struct from cuelang
	S, err := R.ConfigValue.Struct()
	if err != nil {
		return err
	}

	iter := S.Fields()
	for iter.Next() {

		label := iter.Label()
		value := iter.Value()
		fmt.Println("  -", label, value)
		for attrKey, attrVal := range value.Attributes() {
			fmt.Println("  --", attrKey)
			for i := 0; i < 5; i++ {
				str, err := attrVal.String(i)
				if err != nil {
					break
				}
				fmt.Println("  ---", str)
			}
		}
	}

	return nil
}

func (R *Runtime) PrintSecret() error {
	// Get top level struct from cuelang
	S, err := R.SecretValue.Struct()
	if err != nil {
		return err
	}

	iter := S.Fields()
	for iter.Next() {

		label := iter.Label()
		value := iter.Value()
		fmt.Println("  -", label, value)
		for attrKey, attrVal := range value.Attributes() {
			fmt.Println("  --", attrKey)
			for i := 0; i < 5; i++ {
				str, err := attrVal.String(i)
				if err != nil {
					break
				}
				fmt.Println("  ---", str)
			}
		}
	}

	return nil
}

func (R *Runtime) ContextGet(path string) (cue.Value, error) {
	var orig cue.Value
	// var err error
	//if flags.RootPflags.Context != "" {
	//orig, err = cuefig.LoadContextConfig("", flags.RootPflags.Context)
	//} else if flags.RootPflags.Local {
	//orig, err = cuefig.LoadContextConfig("", cuefig.ContextEntrypoint)
	//} else if flags.RootPflags.Global {
	//orig, err = cuefig.LoadHofctxDefault()
	//} else {
	//orig, err = cuefig.LoadContextDefault()
	//}

	//// now check for error
	//if err != nil {
	//return orig, err
	//}

	//if path == "" {
	//return orig, nil
	//}
	//paths := strings.Split(path, ".")
	val := orig.LookupPath(cue.ParsePath(path))
	return val, nil
}

func (R *Runtime) ConfigGet(path string) (cue.Value, error) {
	var orig cue.Value
	// var err error
	//if flags.RootPflags.Config != "" {
	//orig, err = cuefig.LoadConfigConfig("", flags.RootPflags.Config)
	////} else if flags.RootPflags.Local {
	////orig, err = cuefig.LoadConfigConfig("", cuefig.ConfigEntrypoint)
	////} else if flags.RootPflags.Global {
	////orig, err = cuefig.LoadHofcfgDefault()
	//} else {
	//orig, err = cuefig.LoadConfigDefault()
	//}

	//// now check for error
	//if err != nil {
	//return orig, err
	//}

	//if path == "" {
	//return orig, nil
	//}
	//paths := strings.Split(path, ".")
	val := orig.LookupPath(cue.ParsePath(path))
	return val, nil
}

func (R *Runtime) SecretGet(path string) (cue.Value, error) {
	var orig cue.Value
	// var err error
	//if flags.RootPflags.Secret != "" {
	//orig, err = cuefig.LoadSecretConfig("", flags.RootPflags.Secret)
	//} else if flags.RootPflags.Local {
	//orig, err = cuefig.LoadSecretConfig("", cuefig.SecretEntrypoint)
	//} else if flags.RootPflags.Global {
	//orig, err = cuefig.LoadHofshhDefault()
	//} else {
	//orig, err = cuefig.LoadSecretDefault()
	//}

	//// now check for error
	//if err != nil {
	//return orig, err
	//}

	//if path == "" {
	//return orig, nil
	//}
	//paths := strings.Split(path, ".")
	val := orig.LookupPath(cue.ParsePath(path))
	return val, nil
}

func (R *Runtime) ContextSet(expr string) error {
	// var orig cue.Value
	// var val cue.Value
	var err error

	// Check which config we want to work with
	//if flags.RootPflags.Context != "" {
	//orig, err = cuefig.LoadContextConfig("", flags.RootPflags.Context)
	//} else if flags.RootPflags.Local {
	//orig, err = cuefig.LoadContextConfig("", cuefig.ContextEntrypoint)
	//} else if flags.RootPflags.Global {
	//orig, err = cuefig.LoadHofctxDefault()
	//} else {
	//orig, err = cuefig.LoadContextDefault()
	//}

	//// now check for error from that config selection process
	//if err != nil {
	//if _, ok := err.(*os.PathError); !ok && (strings.Contains(err.Error(), "file does not exist") || strings.Contains(err.Error(), "no such file")) {
	//// error is worse than non-existent
	//return err
	//}
	//// file does not exist, so we should just set
	//var r cue.Runtime
	//inst, err := r.Compile("", expr)
	//if err != nil {
	//return err
	//}
	//val = inst.Value()
	//if val.Err() != nil {
	//return val.Err()
	//}

	//} else {
	//val, err = structural.Merge(orig, expr)
	//if err != nil {
	//return err
	//}
	//}

	//// Now save
	//if flags.RootPflags.Context != "" {
	//err = cuefig.SaveContextConfig("", flags.RootPflags.Context, val)
	//} else if flags.RootPflags.Local {
	//err = cuefig.SaveContextConfig("", cuefig.ContextEntrypoint, val)
	//} else if flags.RootPflags.Global {
	//err = cuefig.SaveHofctxDefault(val)
	//} else {
	//err = cuefig.SaveContextDefault(val)
	//}
	return err
}

func (R *Runtime) ConfigSet(expr string) error {
	//var orig cue.Value
	//var val cue.Value
	var err error

	// Check which config we want to work with
	//if flags.RootPflags.Config != "" {
	//orig, err = cuefig.LoadConfigConfig("", flags.RootPflags.Config)
	////} else if flags.RootPflags.Local {
	////orig, err = cuefig.LoadConfigConfig("", cuefig.ConfigEntrypoint)
	////} else if flags.RootPflags.Global {
	////orig, err = cuefig.LoadHofcfgDefault()
	//} else {
	//orig, err = cuefig.LoadConfigDefault()
	//}

	//// now check for error from that config selection process
	//if err != nil {
	//if _, ok := err.(*os.PathError); !ok && (strings.Contains(err.Error(), "file does not exist") || strings.Contains(err.Error(), "no such file")) {
	//// error is worse than non-existent
	//return err
	//}
	//// file does not exist, so we should just set
	//var r cue.Runtime
	//inst, err := r.Compile("", expr)
	//if err != nil {
	//return err
	//}
	//val = inst.Value()
	//if val.Err() != nil {
	//return val.Err()
	//}

	//} else {
	//val, err = structural.Merge(orig, expr)
	//if err != nil {
	//return err
	//}
	//}

	// Now save
	//if flags.RootPflags.Config != "" {
	//err = cuefig.SaveConfigConfig("", flags.RootPflags.Config, val)
	//} else if flags.RootPflags.Local {
	//err = cuefig.SaveConfigConfig("", cuefig.ConfigEntrypoint, val)
	//} else if flags.RootPflags.Global {
	//err = cuefig.SaveHofcfgDefault(val)
	//} else {
	//err = cuefig.SaveConfigDefault(val)
	//}
	return err
}

func (R *Runtime) SecretSet(expr string) error {
	//var orig cue.Value
	//var val cue.Value
	var err error

	// Check which config we want to work with
	//if flags.RootPflags.Secret != "" {
	//orig, err = cuefig.LoadSecretConfig("", flags.RootPflags.Secret)
	//} else if flags.RootPflags.Local {
	//orig, err = cuefig.LoadSecretConfig("", cuefig.SecretEntrypoint)
	//} else if flags.RootPflags.Global {
	//orig, err = cuefig.LoadHofshhDefault()
	//} else {
	//orig, err = cuefig.LoadSecretDefault()
	//}

	//// now check for error from that config selection process
	//if err != nil {
	//if _, ok := err.(*os.PathError); !ok && (strings.Contains(err.Error(), "file does not exist") || strings.Contains(err.Error(), "no such file")) {
	//// error is worse than non-existent
	//return err
	//}
	//// file does not exist, so we should just set
	//var r cue.Runtime
	//inst, err := r.Compile("", expr)
	//if err != nil {
	//return err
	//}
	//val = inst.Value()
	//if val.Err() != nil {
	//return val.Err()
	//}

	//} else {
	//val, err = structural.Merge(orig, expr)
	//if err != nil {
	//return err
	//}
	//}

	//// Now save
	//if flags.RootPflags.Secret != "" {
	//err = cuefig.SaveSecretConfig("", flags.RootPflags.Secret, val)
	//} else if flags.RootPflags.Local {
	//err = cuefig.SaveSecretConfig("", cuefig.SecretEntrypoint, val)
	//} else if flags.RootPflags.Global {
	//err = cuefig.SaveHofshhDefault(val)
	//} else {
	//err = cuefig.SaveSecretDefault(val)
	//}
	return err
}
