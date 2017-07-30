// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"strings"
	"path/filepath"
	"os"
	"os/exec"
	"log"
	"go/build"
	"github.com/caojia/gip/helper"
)

// freezeCmd represents the freeze command
var freezeCmd = &cobra.Command{
	Use:   "freeze",
	Short: "Output installed packages in requirements format.",
	Long: `Output installed packages in requirements format.
Only packages installed and depenced directly or indirectly by current package will be output.
packages are listed in a case-insensitive sorted order.

Edit gip.yml if the git repo or version of a package need customized.

Example of gip.yml:

	imports:
		- package: golang.org/x/net
		  repo: https://github.com/golang/net
		  version: master

`,

	Run: func(cmd *cobra.Command, args []string) {
		helper.Freeze()
	},
}

type dir struct {
	name   string
	parent string
}

type Package struct {
	// package name
	Package string
	// repo url
	Repo    string
	// commit hash
	Version string
	// contains vendor dir
	Vendor  bool
	// isSelf
	Self    bool
}

var wd string = ""
var packages map[string]bool = make(map[string]bool)
var repos map[string]Package = make(map[string]Package)

func findRepo(srcPath, p string) Package {
	folders := strings.Split(p, string(filepath.Separator))
	base := ""
	for _, f := range folders {
		base = filepath.Join(base, f)
		if r, ok := repos[base]; ok {
			return r
		}
		if strings.HasSuffix(wd, base) {
			r := Package{ Self: true, }
			if f, _ := os.Stat(filepath.Join(srcPath, base, "vendor")); f != nil && f.IsDir() {
				r.Vendor = true
			}
			repos[base] = r
			return r
		}
		if f, _ := os.Stat(filepath.Join(srcPath, base, ".git")); f != nil && f.IsDir() {
			output, err := exec.Command("git", "-C", filepath.Join(srcPath, base), "remote", "get-url", "origin").Output()
			if err != nil {
				log.Printf(fmt.Sprintf("[WARNING] unsupported cvs, pkg=%s, err=%s", p, err.Error()))
			}
			commit, err := exec.Command("git", "-C", filepath.Join(srcPath, base), "rev-parse", "HEAD").Output()
			if err != nil {
				log.Printf("[WARNING] got error while getting commit hash: %s\n", err.Error())
			}
			r := Package{
				Repo:    strings.Trim(string(output), "\n"),
				Version: strings.Trim(string(commit), "\n"),
			}
			if f, _ := os.Stat(filepath.Join(srcPath, base, "vendor")); f != nil && f.IsDir() {
				r.Vendor = true
			}
			repos[base] = r
			return r
		}
	}
	panic(fmt.Sprintf("ooops, unsupported cvs, pkg=%s", p))
}

// Traverse the packages and figure out the dependency recursively.
// if the package contains vendor, we assume it already solve the dependencies.
// so we don't look into it any more
func recursive(ctx build.Context, dir string, parent string) {
	p, err := ctx.ImportDir(dir, build.IgnoreVendor)
	if err != nil {
		//log.Printf("parent=%s, pkg=%s %s\n", parent, dir, err.Error())
		return
	}
	for _, x := range p.Imports {
		if strings.Contains(x, ".") {
			packages[x] = true
			r := findRepo(p.SrcRoot, x)
			if !r.Vendor {
				recursive(ctx, filepath.Join(p.SrcRoot, x), dir)
			}
		}
	}
}

func dfs(ctx build.Context, path string) {
	// DFS
	var fs []dir
	fs = append(fs, dir{"", path})
	var i = 0
	for {
		if len(fs) <= i {
			break
		}
		f := fs[i]
		i += 1
		base := filepath.Join(f.parent, f.name)
		file, err := os.OpenFile(base, os.O_RDONLY, os.ModeDir)
		if err != nil {
			log.Printf(err.Error())
			continue
		}
		recursive(ctx, base, "")
		subFiles, err := file.Readdirnames(-1)
		if err != nil {
			log.Printf(err.Error())
			continue
		}
		for _, sf := range subFiles {
			if fi, _ := os.Stat(filepath.Join(base, sf)); fi != nil && fi.IsDir() && !strings.HasPrefix(fi.Name(), ".") && fi.Name() != "vendor" {
				fs = append(fs, dir{sf, base})
			}
		}
	}

}

func init() {
	RootCmd.AddCommand(freezeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// freezeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// freezeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
