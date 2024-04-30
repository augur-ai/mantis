package script_test

import (
	"testing"

	"github.com/opentofu/opentofu/internal/hof/script/runtime"
)

func TestScriptBrowser(t *testing.T) {
	runtime.Run(t, runtime.Params{
		Dir:  "tests/browser",
		Glob: "*.hls",
	})
}

func TestScriptCmds(t *testing.T) {
	runtime.Run(t, runtime.Params{
		Dir:  "tests/cmds",
		Glob: "*.hls",
	})
}

func TestScriptHTTP(t *testing.T) {
	runtime.Run(t, runtime.Params{
		Dir:  "tests/http",
		Glob: "*.hls",
	})
}
