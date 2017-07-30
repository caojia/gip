package helper

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/caojia/gip/log"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

const configFile = "gip.yml"

var cfg config
var srcPath string
var globalSrcPath string
var srcPaths []string = make([]string, 0)

// Package
// define the structure of a package
type Package struct {
	// package name
	Package string `yaml:"package"`
	// repo url
	Repo string `yaml:"repo"`
	// commit hash
	Version string `yaml:"version"`
	// isGlobal, install the package in global path
	Global bool `yaml:"global"`
	// contains vendor dir
	Vendor bool
	// isSelf
	Self bool
}

func (p Package) String() string {
	return fmt.Sprintf("%s#%s,%s", p.Repo, p.Version, p.Package)
}

func LoadPackage(line string) (*Package, error) {
	splits := strings.SplitN(strings.TrimSpace(line), ",", 2)
	if len(splits) < 2 {
		return nil, errors.New("Invalid line: " + line)
	}

	repos := strings.SplitN(splits[0], "#", 2)
	repo := repos[0]
	version := ""
	if len(repos) == 2 {
		version = repos[1]
	}
	packageName := splits[1]

	return &Package{
		Package: packageName,
		Repo:    repo,
		Version: version,
	}, nil
}

func LoadPackagesFromFile(filename string) ([]*Package, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")
	pkgs := make([]*Package, 0)
	for _, l := range lines {
		p, err := LoadPackage(l)
		if err != nil {
			continue
		}
		pkgs = append(pkgs, p)
	}
	return pkgs, nil
}

func IgnoreGlobal() {
	globalSrcPath = srcPath
}

// Config
// define the structure of a config file: gip.yaml
type config struct {
	Imports     []Package `yaml:imports`
	packagesMap map[string]Package
}

func init() {
	srcPaths = strings.Split(build.Default.GOPATH, ":")
	for i, s := range srcPaths {
		srcPaths[i] = filepath.Join(s, "src")
	}
	srcPath = srcPaths[0]
	globalSrcPath = srcPaths[0]
	if len(srcPaths) > 1 {
		globalSrcPath = srcPaths[1]
	}
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Debug("fail to read the file, err=%s", err.Error())
		return
	}
	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		log.Debug("fail to unmarshal the file, err=%s", err.Error())
		return
	}
	cfg.packagesMap = make(map[string]Package)
	for _, pkg := range cfg.Imports {
		cfg.packagesMap[pkg.Package] = pkg
	}
}
