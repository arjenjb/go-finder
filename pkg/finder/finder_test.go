package finder

import (
	"embed"
	"github.com/arjenjb/go-finder/internal/util"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testDir = "test"

//go:embed test/**/*
var assetsFS embed.FS

func Dump(items []Entry) string {
	sb := strings.Builder{}
	sb.WriteString("\n")
	for _, item := range items {
		sb.WriteString(strings.Repeat("  ", item.Depth()))
		sb.WriteString(item.Name())
		if item.IsDir() {
			sb.WriteString(" (dir)")
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func TestFinderFS(t *testing.T) {
	println(Dump(util.Must(NewFinder().InFS(assetsFS).Find())))
}

func TestFinderName(t *testing.T) {
	items := util.Must(NewFinder().In(testDir).Files().Name("*.txt").Find())
	assert.Equal(t, `
a.txt
b.txt
  x.txt
  y.txt
`, Dump(items))
}

func TestFinderDirectories(t *testing.T) {
	items := util.Must(NewFinder().In(testDir).Directories().Find())
	assert.Equal(t, `
.gitx (dir)
dir-a (dir)
  subdir-a (dir)
  subdir-b (dir)
`, Dump(items))
}

func TestFinderNotPath(t *testing.T) {
	items := util.Must(NewFinder().In(testDir).Directories().NotPath("ir-a/subdir-b").Find())
	assert.Equal(t, `
.gitx (dir)
dir-a (dir)
  subdir-a (dir)
`, Dump(items))
}

func TestFinderMaxDepth(t *testing.T) {
	items := util.Must(NewFinder().In(testDir).MaxDepth(1).Find())

	assert.Equal(
		t,
		`
.gitx (dir)
  file
.hidden
a.txt
b.txt
dir-a (dir)
  subdir-a (dir)
  subdir-b (dir)
  x.txt
  y.txt
`,
		Dump(items),
	)
}

func TestFinderMinDepth(t *testing.T) {
	items := util.Must(NewFinder().In(testDir).MinDepth(1).Find())

	assert.Equal(t, `
  file
  subdir-a (dir)
    README
  subdir-b (dir)
    README
  x.txt
  y.txt
`, Dump(items))
}
