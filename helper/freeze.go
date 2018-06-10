package helper

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/caojia/gip/log"
	"time"
)

type dir struct {
	name   string
	parent string
}

var wd string = ""
var repos map[string]Package = make(map[string]Package)
var visitedPackage map[string]bool = make(map[string]bool)

func init() {
	var err error = nil
	wd, err = os.Getwd()
	if err != nil {
		panic("fail to get current workspace: err=" + err.Error())
	}
}

func findRepo(src, p string, ignoreCache bool) (Package, error) {
	if ignoreCache && isDir(filepath.Join(src, p)) {
		return Package{}, nil
	}

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
			r.SrcPath = wd[0:len(wd) - len(base)]
			if !ignoreCache {
				repos[base] = r
			}

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
			r.SrcPath = src
			if !ignoreCache {
				repos[base] = r
			}
			return r, nil
		}
	}
	return Package{}, fmt.Errorf("ooops, unsupported cvs, pkg=%s, srcPath=%s", p, src)
}

// Traverse the packages and figure out the dependency recursively.
// if the package contains vendor, we assume it already solve the dependencies.
// so we don't look into it any more
func recursive(ctx build.Context, dir string, parent string, lastPackage *Package) {
	ts := time.Now()
	defer func() {
		elapsed := time.Now().UnixNano() - ts.UnixNano()
		log.Debug("dir=%s, parent=%s, es=%d", dir, parent, elapsed)
	}()
	p, err := ctx.ImportDir(dir, build.IgnoreVendor)
	if err != nil {
		if !strings.Contains(err.Error(), "no buildable Go source files") {
			log.Debug("got err: %s, dir=%s, parent=%s", err.Error(), dir, parent)
		}
		return
	}
	imports := p.Imports
	if parent == "" {
		imports = append(imports, p.TestImports...)
		imports = append(imports, p.XTestImports...)
	}
	for _, x := range imports {
		if strings.Contains(x, ".") {
			if visitedPackage[x] {
				continue
			}
			srcRoot := p.SrcRoot
			var r Package
			// Parent is vendor, means current package is under vendor directory
			parentIsVendor := lastPackage != nil && lastPackage.Vendor && !strings.HasPrefix(x, lastPackage.Package)
			if parentIsVendor {
				r, err = findRepo(filepath.Join(lastPackage.SrcPath, lastPackage.Package, "vendor"), x, true)
				if err == nil {
					continue
				}
				log.Debug("fail to get repo: %s, dir=%s err=%s, lastPackage=%v, srcRoot=%s", x, dir, err.Error(), lastPackage, srcRoot)
			}
			for _, src := range srcPaths {
				r, err = findRepo(src, x, false)
				if err == nil {
					srcRoot = r.SrcPath
					break
				}
			}

			if err != nil {
				if !parentIsVendor {
					panic(err)
				}
				continue
			}

			visitedPackage[x] = true

			if !(parentIsVendor && r.Vendor) {
				recursive(ctx, filepath.Join(srcRoot, x), dir, &r)
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
		recursive(ctx, base, "", nil)
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
