package helper

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/caojia/gip/log"
)

func gitTryRemote(destPath, version string) (bool, error) {
	c := exec.Command("git", "show-ref", "-q", "--verify", "refs/remotes/origin/"+version)
	c.Dir = destPath
	_, err := c.Output()
	if err == nil {
		log.Debug("%s reset to origin/%s", destPath, version)
		c := exec.Command("git", "reset", "--hard", "origin/"+version)
		c.Dir = destPath
		output, err := c.CombinedOutput()
		if err != nil {
			return true, errors.New(string(output))
		}
		return true, nil
	}
	return false, nil
}

func gitCheckout(destPath, version string) error {
	log.Debug("%s checking out %s", destPath, version)
	c := exec.Command("git", "checkout", version)
	c.Dir = destPath
	output, err := c.CombinedOutput()
	if err != nil {
		return errors.New(string(output))
	}
	return nil
}

func Get(pkg *Package) error {
	src := srcPath
	if pkg.Global {
		src = globalSrcPath
	}
	destPath := filepath.Join(src, pkg.Package)
	if isGitDir(destPath) {
		log.Debug("fetching origin: %s", destPath)
		// update the git
		c := exec.Command("git", "fetch", "origin")
		c.Dir = destPath
		c.Output()
	} else {
		log.Debug("cloning: %s", destPath)
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
