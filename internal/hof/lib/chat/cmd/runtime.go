package cmd

import (
	"github.com/opentofu/opentofu/internal/hof/cmd/hof/flags"
	"github.com/opentofu/opentofu/internal/hof/lib/chat"
	"github.com/opentofu/opentofu/internal/hof/lib/runtime"
)

func prepRuntime(args []string, rflags flags.RootPflagpole) (*runtime.Runtime, error) {

	// create our core runtime
	r, err := runtime.New(args, rflags)
	if err != nil {
		return nil, err
	}

	err = r.Load()
	if err != nil {
		return nil, err
	}

	err = r.EnrichChats(nil, EnrichChat)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func EnrichChat(R *runtime.Runtime, c *chat.Chat) error {

	// no-op
	return nil
}
