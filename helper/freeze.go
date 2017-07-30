package helper

import (
	"fmt"
	"github.com/caojia/gip/log"
	"go/build"
	"os"
	"path/filepath"
	"strings"
	"sort"
)

type dir struct {
	name   string
	parent string
}

var wd string = ""
var packages map[string]bool = make(map[string]bool)
var repos map[string]Package = make(map[string]Package)

func init() {
	var err error = nil
	wd, err = os.Getwd()
	if err != nil {
		panic("fail to get current workspace: err=" + err.Error())
	}
}

func findRepo(src, p string) (Package, error) {
	folders := strings.Split(p, string(filepath.Separator))
	base := ""
	for _, f := range folders {
		base = filepath.Join(base, f)
		if r, ok := repos[base]; ok {
			return r, nil
		}
		if strings.HasSuffix(wd, base) {
			r := Package{Self: true}
			if f, _ := os.Stat(filepath.Join(src, base, "vendor")); f != nil && f.IsDir() {
				r.Vendor = true
			}
			r.Package = base
			repos[base] = r
			return r, nil
		}
		fullpath := filepath.Join(src, base)
		if isGitDir(fullpath) {
			remote, err := gitRemote(fullpath)
			if err != nil {
				log.Debug(fmt.Sprintf("[WARNING] unsupported cvs, pkg=%s, err=%s", p, err.Error()))
			}
			commit, err := gitCommit(fullpath)
			if err != nil {
				log.Debug("[WARNING] got error while getting commit hash: %s\n", err.Error())
			}
			r := Package{
				Package: base,
				Repo:    remote,
				Version: commit,
			}
			if f, _ := os.Stat(filepath.Join(src, base, "vendor")); f != nil && f.IsDir() {
				r.Vendor = true
			}
			if cachedR, ok := cfg.packagesMap[base]; ok {
				if len(cachedR.Repo) > 0 {
					r.Repo = cachedR.Repo
				}
				if len(cachedR.Version) > 0 {
					r.Version = cachedR.Version
				}
				r.Global = cachedR.Global
			}
			repos[base] = r
			return r, nil
		}
	}
	return Package{}, fmt.Errorf("ooops, unsupported cvs, pkg=%s, srcPath=%s", p, src)
}

// Traverse the packages and figure out the dependency recursively.
// if the package contains vendor, we assume it already solve the dependencies.
// so we don't look into it any more
func recursive(ctx build.Context, dir string, parent string) {
	p, err := ctx.ImportDir(dir, build.IgnoreVendor)
	if err != nil {
		return
	}
	for _, x := range p.Imports {
		if strings.Contains(x, ".") {
			packages[x] = true
			r, err := findRepo(p.SrcRoot, x)
			if err != nil {
				for _, src := range srcPaths {
					r, err = findRepo(src, x)
					if err == nil {
						break
					}
				}
			}
			if err != nil {
				panic(err)
			}
			if !r.Vendor {
				recursive(ctx, filepath.Join(p.SrcRoot, x), dir)
			}
		}
	}
}

func bfs(ctx build.Context, path string) {
	// BFS
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
			log.Debug("fail to open file %s, err = %s", base, err.Error())
			continue
		}
		recursive(ctx, base, "")
		subFiles, err := file.Readdirnames(-1)
		if err != nil {
			log.Debug("fail to get dir names %s, err = %s", base, err.Error())
			continue
		}
		file.Close()
		for _, sf := range subFiles {
			if fi, _ := os.Stat(filepath.Join(base, sf)); fi != nil && fi.IsDir() && !strings.HasPrefix(fi.Name(), ".") && fi.Name() != "vendor" {
				fs = append(fs, dir{sf, base})
			}
		}
	}
}

func Freeze() {
	bfs(build.Default, wd)
	keys := make([]string, len(repos))
	i := 0
	for k, _ := range repos {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := repos[k]
		if v.Self {
			continue
		}
		log.Info(v.String())
	}

}
