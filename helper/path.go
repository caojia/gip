package helper

import (
	"os"
	"path/filepath"
	"os/exec"
	"strings"
)

func isGitDir(path string) bool {
	f, _ := os.Stat(filepath.Join(path, ".git"))
	return f != nil && f.IsDir()
}

func gitCommit(path string) (string, error) {
	commit, err := exec.Command("git", "-C", path, "rev-parse", "HEAD").Output()
	return strings.TrimSpace(string(commit)), err
}

func gitRemote(path string) (string, error) {
	remote, err := exec.Command("git", "-C", path, "remote", "get-url", "origin").Output()
	return strings.TrimSpace(string(remote)), err
}


