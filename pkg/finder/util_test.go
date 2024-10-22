package finder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_asGlobRegexPattern(t *testing.T) {
	assert.Equal(t, "^.*\\.y.ml$", asGlobRegexPattern("*.y?ml", true))
	assert.Equal(t, "^.*\\.txt$", asGlobRegexPattern("*.txt", true))
	assert.Equal(t, "^.*\\.....$", asGlobRegexPattern("*.????", true))
}
