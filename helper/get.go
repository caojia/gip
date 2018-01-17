package helper

import (
	"errors"
	"os"
	"os/exec"
	"github.com/caojia/gip/log"
	"path"
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

func gitStash(destPath string) error {
	c := exec.Command("git", "stash")
	c.Dir = destPath
	output, err := c.CombinedOutput()
	if err != nil {
		return errors.New(string(output))
	}
	return nil
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

func getRequirementsForPackage(pkg *Package, requirementFile string) ([]*Package, error) {
	realPath := pkg.RealPath()
	filepath := path.Join(realPath, requirementFile)
	if isFile(filepath) {
		return LoadPackagesFromFile(filepath)
	}
	return nil, nil
}

func GetRecursively(pkgs []*Package, requirementFile string) {
	allPackages := pkgs
	visitedPackage := map[string]bool{}
	for i := 0; i < len(allPackages); i++ {
		p := allPackages[i]
		if visitedPackage[p.Package] {
			continue
		}
		if err := Get(p); err != nil {
			log.Error("get %s failed: err=%s", p.Package, err.Error())
		}
		visitedPackage[p.Package] = true

		np, err := getRequirementsForPackage(p, requirementFile)
		if err != nil {
			log.Error("get requirements for %s failed: err=%s", p.Package, err.Error())
			continue
		}
		allPackages = append(allPackages, np...)
	}
}

func Get(pkg *Package) error {
	destPath := pkg.RealPath()
	log.Info("Getting %s, repo: %s, path: %s", pkg.Package, pkg.Repo, destPath)
	if isGitDir(destPath) {
		gitStash(destPath)
		needFetch := false
		// check is remote, if it is remote, always fetch origin
		remote, err := gitTryRemote(destPath, pkg.Version)
		if remote {
			if err != nil {
				return err
			}
			needFetch = true
		}

		// if it is not remote, check if the hash already checked out
		if !needFetch {
			err := gitCheckout(destPath, pkg.Version)
			if err != nil {
				needFetch = true
			}
		}

		// if need fetch origin, fetch it
		if needFetch {
			log.Debug("fetching origin: %s", destPath)
			// update the git
			c := exec.Command("git", "fetch", "origin")
			c.Dir = destPath
			c.Output()
		}
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
