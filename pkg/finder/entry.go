package finder

import "io/fs"

type Entry struct {
	path  string
	entry fs.DirEntry
	depth int
}

func (e Entry) Path() string {
	return e.path
}

func (e Entry) Depth() int {
	return e.depth
}

func (e Entry) Name() string {
	return e.entry.Name()
}

func (e Entry) IsDir() bool {
	return e.entry.IsDir()
}

func (e Entry) Type() fs.FileMode {
	return e.entry.Type()
}

func (e Entry) Info() (fs.FileInfo, error) {
	return e.entry.Info()
}
