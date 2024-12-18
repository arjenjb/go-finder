package finder

import (
	"fmt"
	. "github.com/arjenjb/go-finder/internal/util"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type FileType int
type IgnoreMask int

const (
	Any FileType = iota
	File
	Directory
)

const (
	VcsFiles IgnoreMask = 1 << 0
	DotFiles            = 1 << 1
)

// Finder tracks the finders configuration prior to execution
type Finder struct {
	name []string

	matchName    []regexp.Regexp
	notMatchName []regexp.Regexp

	matchPath          []regexp.Regexp
	notMatchPath       []regexp.Regexp
	t                  FileType
	excludedDirs       []string
	directories        []string
	ignore             IgnoreMask
	maxDepth, minDepth int
	fs                 []fs.FS
	followSymlinks     bool

	//maxSize, minSize int // Not supported yet, but will be implemented later
}

// Name adds a name condition to the finder. The name of a file should match the given glob pattern.
// If you call this function multiple times, a file should match any of the given name patterns.
//
// Valid patterns are:
//   - *.txt -> matches file.txt
//   - README.* -> matches README.txt and README.md
//   - ?.x -> matches a.x, b.x etc
func (f Finder) Name(n string) Finder {
	nf := f
	nf.name = append(f.name, n)
	return nf
}

// Files configures the finder to only match files
func (f Finder) Files() Finder {
	nf := f
	nf.t = File
	return nf
}

// Directories configures the finder to only match directories
func (f Finder) Directories() Finder {
	nf := f
	nf.t = Directory
	return nf
}

func (f Finder) IgnoreVcsFiles(ignore bool) Finder {
	nf := f
	if ignore {
		nf.ignore |= VcsFiles
	} else {
		nf.ignore &^= VcsFiles
	}
	return nf
}

func (f Finder) IgnoreDotFiles(ignore bool) Finder {
	nf := f
	if ignore {
		nf.ignore |= DotFiles
	} else {
		nf.ignore &^= DotFiles
	}
	return nf
}

// In specifies in which directories to search, you may invoke this method multiple times to add additional directories
func (f Finder) In(directories ...string) Finder {
	nf := f
	nf.directories = append(nf.directories, directories...)
	return nf
}

func (f Finder) Exclude(directories ...string) Finder {
	nf := f
	nf.excludedDirs = append(nf.excludedDirs, directories...)
	return nf
}

func (f Finder) MaxDepth(d int) Finder {
	nf := f
	nf.maxDepth = d
	return nf
}

func (f Finder) MinDepth(d int) Finder {
	nf := f
	nf.minDepth = d
	return nf
}

func (f Finder) MustFind() []Entry {
	return Must(f.Find())
}

type linkEntry struct {
	info fs.FileInfo
}

func (l linkEntry) Name() string {
	return l.info.Name()
}

func (l linkEntry) IsDir() bool {
	return l.info.IsDir()
}

func (l linkEntry) Type() fs.FileMode {
	return l.info.Mode()
}

func (l linkEntry) Info() (fs.FileInfo, error) {
	return l.info, nil
}

var _ fs.DirEntry = linkEntry{}

func (f Finder) Find() ([]Entry, error) {
	var entries []Entry

	nameRegexes := append(f.matchName, Map(f.name, func(str string) regexp.Regexp {
		return asGlobRegex(str, true)
	})...)

	notDirNameRegexes := f.notDirNameRegexes()

	depth := 0
	depthPrefix := ""

	WalkFunc := func(path string, entry fs.DirEntry, err error, root string) error {
		// Figure out the relative dir and filename
		rdir, rfile := filepath.Split(path[Min(len(path), len(root)+1):])
		npath := strings.ReplaceAll(path[Min(len(path), len(root)+1):], string(os.PathSeparator), string('/'))

		// Skip the root directory
		if len(rfile) == 0 {
			return nil
		}

		if entry.IsDir() && len(notDirNameRegexes) > 0 {
			if AnySatisfy(notDirNameRegexes, func(r regexp.Regexp) bool {
				return r.MatchString(entry.Name())
			}) {
				return filepath.SkipDir
			}
		}

		// Figure out the current depth
		if rdir == depthPrefix {
			// Nothing the matter
		} else if len(depthPrefix) > 0 && strings.HasPrefix(rdir, depthPrefix) {
			depth++
			depthPrefix = rdir
		} else {
			depth = len(strings.Split(rdir, string(os.PathSeparator))) - 1
			depthPrefix = rdir
		}

		if f.minDepth != -1 && depth < f.minDepth {
			return nil
		}

		var ret error = nil
		if entry.IsDir() && f.maxDepth != -1 && depth == f.maxDepth {
			ret = filepath.SkipDir
		}

		// Resolve symlinks, if thats what we want
		if (entry.Type()&os.ModeSymlink) != 0 && f.followSymlinks {
			entry, err = resolveLink(path)
		}

		if f.t == File && !entry.Type().IsRegular() {
			return ret
		}
		if f.t == Directory && !entry.IsDir() {
			return nil
		}

		// Check names
		if len(f.notMatchName) > 0 {
			if AnySatisfy(f.notMatchName, func(r regexp.Regexp) bool {
				return r.MatchString(entry.Name())
			}) {
				return nil
			}
		}

		// Check names of files
		if !entry.IsDir() && len(nameRegexes) > 0 {
			if !AnySatisfy(nameRegexes, func(r regexp.Regexp) bool {
				return r.MatchString(entry.Name())
			}) {
				return nil
			}
		}

		// Check paths names
		if len(f.notMatchPath) > 0 {
			if AnySatisfy(f.notMatchPath, func(r regexp.Regexp) bool {
				return r.MatchString(npath)
			}) {
				return ret
			}
		}

		if len(f.matchPath) > 0 {
			if !AnySatisfy(f.matchPath, func(r regexp.Regexp) bool {
				return r.MatchString(npath)
			}) {
				return ret
			}
		}

		entries = append(entries, Entry{
			path:  path,
			entry: entry,
			depth: depth,
		})

		return ret
	}

	f.WalkDirectories(WalkFunc)
	f.WalkFS(WalkFunc)

	return entries, nil
}

func resolveLink(path string) (fs.DirEntry, error) {
	for {
		target, err := os.Readlink(path)
		if err != nil {
			return nil, err
		}

		if !filepath.IsAbs(target) {
			target, err = filepath.Abs(filepath.Join(filepath.Dir(path), target))
			if err != nil {
				return nil, err
			}
		}

		t, err := os.Stat(target)
		if t.Mode()&os.ModeSymlink != 0 {
			path = filepath.Join(filepath.Dir(path), target)
		} else {
			return linkEntry{
				info: t,
			}, nil
		}

	}
}

func asFullGlobRegex(str string) regexp.Regexp {
	return asGlobRegex(str, true)
}

type myWalkFunc func(path string, entry fs.DirEntry, err error, root string) error

func (f Finder) WalkDirectories(walkFunc myWalkFunc) {
	for _, d := range f.directories {
		filepath.WalkDir(d, func(path string, entry fs.DirEntry, err error) error {
			return walkFunc(path, entry, err, d)
		})
	}
}

func (f Finder) WalkFS(walkFunc myWalkFunc) error {
	for _, each := range f.fs {
		fs.WalkDir(each, ".", func(path string, entry fs.DirEntry, err error) error {
			return walkFunc(path, entry, err, ".")
		})
	}
	return nil
}

func (f Finder) notDirNameRegexes() []regexp.Regexp {
	var result []regexp.Regexp

	if (f.ignore & VcsFiles) > 0 {
		result = append(result, *regexp.MustCompile(fmt.Sprintf("^(%s)$", strings.Join(Map([]string{".svn", "_svn", "CVS", "_darcs", ".arch-params", ".monotone", ".bzr", ".git", ".hg"}, regexp.QuoteMeta), "|"))))
	}

	return result
}

func (f Finder) InFS(fs fs.FS) Finder {
	nf := f
	nf.fs = append(nf.fs, fs)
	return nf
}

func (f Finder) Path(path ...string) Finder {
	regexes := Map(path, func(each string) regexp.Regexp {
		return asGlobRegex(each, false)
	})

	nf := f
	nf.matchPath = append(nf.matchPath, regexes...)
	return nf
}

func (f Finder) NotPath(path ...string) Finder {
	nf := f
	nf.notMatchPath = append(nf.notMatchPath, Map(path, func(each string) regexp.Regexp {
		return asGlobRegex(each, false)
	})...)
	return nf
}

func (f Finder) NameRegex(r *regexp.Regexp) Finder {
	nf := f
	nf.matchName = append(nf.matchName, *r)
	return nf
}

func (f Finder) NotName(name string) Finder {
	nf := f
	nf.notMatchName = append(f.notMatchName, asGlobRegex(name, true))
	return nf
}

func (f Finder) FollowSymlinks() Finder {
	nf := f
	nf.followSymlinks = true
	return nf
}

// NewFinder creates a new finder for you
func NewFinder() Finder {
	return Finder{
		ignore:   DotFiles | VcsFiles,
		maxDepth: -1,
		minDepth: -1,
	}
}
