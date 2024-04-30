package cache

import (
	"os"

	"golang.org/x/mod/sumdb/dirhash"

	"github.com/opentofu/opentofu/internal/hof/lib/repos/utils"
)

func Checksum(mod, ver string) (string, error) {
	remote, owner, repo := utils.ParseModURL(mod)
	tag := ver

	dir := ModuleOutdir(remote, owner, repo, tag)

	_, err := os.Lstat(dir)
	if err != nil {
		return "", err
	}

	h, err := dirhash.HashDir(dir, mod, dirhash.Hash1)

	return h, err
}
