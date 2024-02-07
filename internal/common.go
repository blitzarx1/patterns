package internal

import (
	"os/exec"
	"strings"
)

func MustGetGitVersion() string {
	tagCmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	commitCmd := exec.Command("git", "rev-parse", "--short", "HEAD")

	tagOutput, err := tagCmd.Output()
	if err != nil {
		panic(err)
	}

	commitOutput, err := commitCmd.Output()
	if err != nil {
		panic(err)
	}

	tag := strings.TrimSpace(string(tagOutput))
	commit := strings.TrimSpace(string(commitOutput))

	return tag + "-" + commit
}
