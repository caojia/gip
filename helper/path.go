package helper

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func isGitDir(path string) bool {
	return isDir(filepath.Join(path, ".git"))
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

func isDir(path string) bool {
	f, _ := os.Stat(path)
	return f != nil && f.IsDir()
}
