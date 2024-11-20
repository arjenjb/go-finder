package finder

import (
	"github.com/arjenjb/go-finder/internal/util"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Dir string

func (d Dir) Mkdir(s string) Dir {
	n := filepath.Join(string(d), s)
	os.Mkdir(n, 0775)
	return Dir(n)
}

func (d Dir) Touch(s string) {
	n := filepath.Join(string(d), s)
	fp, _ := os.OpenFile(n, os.O_CREATE, 0660)
	fp.Close()
}

func createTestDir(d string) Dir {
	root := Dir(d)
	root.Mkdir(".gitx").Touch("file")
	root.Touch(".hidden")
	root.Touch("a.txt")
	root.Touch("b.txt")

	a := root.Mkdir("dir-a")
	a.Touch("x.txt")
	a.Touch("y.txt")

	suba := a.Mkdir("subdir-a")
	suba.Touch("README")
	suba.Touch("z.txt")

	subb := a.Mkdir("subdir-b")
	subb.Touch("README")

	return root
}

func testDirDo(f func(dir string)) {
	tmp, _ := os.MkdirTemp("", "")
	defer os.RemoveAll(tmp)

	d := createTestDir(tmp)
	f(string(d))
}

func testFinderDo(f func(f Finder)) {
	testDirDo(func(dir string) {
		f(NewFinder().In(dir))
	})
}

func Dump(items []Entry) string {
	sb := strings.Builder{}
	sb.WriteString("\n")
	for _, item := range items {
		sb.WriteString(strings.Repeat("  ", item.Depth()))

		if item.IsDir() {
			sb.WriteString("[")
			sb.WriteString(item.Name())
			sb.WriteString("]")
		} else {
			sb.WriteString(item.Name())
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

//
//func TestFinderFS(t *testing.T) {
//	println(Dump(util.Must(NewFinder().InFS(assetsFS).Find())))
//}

func TestFinderName(t *testing.T) {
	testDirDo(func(dir string) {
		items := util.Must(NewFinder().In(dir).Files().FollowSymlinks().Name("*.txt").Find())
		assert.Equal(t, `
a.txt
b.txt
    z.txt
  x.txt
  y.txt
`, Dump(items))
	})

}

func TestFinderDirectories(t *testing.T) {
	testFinderDo(func(f Finder) {
		items := util.Must(f.Directories().Find())
		assert.Equal(t, `
[.gitx]
[dir-a]
  [subdir-a]
  [subdir-b]
`, Dump(items))
	})
}

func TestFinderNotPath(t *testing.T) {
	testFinderDo(func(f Finder) {
		items := util.Must(f.Directories().NotPath("ir-a/subdir-b").Find())
		assert.Equal(t, `
[.gitx]
[dir-a]
  [subdir-a]
`, Dump(items))
	})

}

func TestFinderMaxDepth(t *testing.T) {
	testFinderDo(func(f Finder) {
		items := util.Must(f.MaxDepth(1).Find())

		assert.Equal(
			t,
			`
[.gitx]
  file
.hidden
a.txt
b.txt
[dir-a]
  [subdir-a]
  [subdir-b]
  x.txt
  y.txt
`,
			Dump(items),
		)
	})

}

func TestFinderMinDepth(t *testing.T) {
	testFinderDo(func(f Finder) {
		items := util.Must(f.MinDepth(1).Find())

		assert.Equal(t, `
  file
  [subdir-a]
    README
    z.txt
  [subdir-b]
    README
  x.txt
  y.txt
`, Dump(items))

	})
}
