package helper

import (
	"path/filepath"
	"os/exec"
	"os"
	"errors"
)

func gitTryRemote(destPath, version string) (bool, error) {
	_, err := exec.Command("git", "-C", destPath, "show-ref", "-q", "--verify", "refs/remote/origin/" + version).Output()
	if err == nil {
		output, err := exec.Command("git", "-C", destPath, "reset", "--hard", "origin/" + version).CombinedOutput()
		if err != nil {
			return true, errors.New(string(output))
		}
		return true, nil
	}
	return false, nil
}

func gitCheckout(destPath, version string) error {
	output, err := exec.Command("git", "-C", destPath, "checkout", version).CombinedOutput()
	if err != nil {
		return errors.New(string(output))
	}
	return nil
}

func Get(pkg *Package) error {
	destPath := filepath.Join(srcPath, pkg.Package)
	if isGitDir(destPath) {
		// update the git
		exec.Command("git", "-C", destPath, "fetch", "origin").Output()
	} else {
		os.MkdirAll(destPath, os.ModePerm)
		output, err := exec.Command("git", "clone", pkg.Repo, destPath).CombinedOutput()
		if err != nil {
			return errors.New(string(output))
		}
	}
	if len(pkg.Version) > 0 {
		// the version is a remote branch
		remote, err := gitTryRemote(destPath, pkg.Version)
		if remote {
			return err
		}
		return gitCheckout(destPath, pkg.Version)
	}
	return nil
}
