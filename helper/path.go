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
	c := exec.Command("git", "rev-parse", "HEAD")
	c.Dir = path
	commit, err := c.Output()
	return strings.TrimSpace(string(commit)), err
}

func gitRemote(path string) (string, error) {
	c := exec.Command("git", "remote", "get-url", "origin")
	c.Dir = path
	remote, err := c.Output()
	return strings.TrimSpace(string(remote)), err
}


