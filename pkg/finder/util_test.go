package finder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_asGlobRegexPattern(t *testing.T) {
	assert.Equal(t, "^.*\\.y.ml$", asGlobRegexPattern("*.y?ml"))
	assert.Equal(t, "^.*\\.txt$", asGlobRegexPattern("*.txt"))
	assert.Equal(t, "^.*\\.....$", asGlobRegexPattern("*.????"))
}
