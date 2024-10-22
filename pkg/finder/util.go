package finder

import (
	"github.com/arjenjb/go-finder/internal/util"
	"regexp"
	"strings"
)

func asGlobRegexPattern(glob string, full bool) string {
	re := regexp.MustCompile("[*?]|[^*?]+")
	result := re.FindAllStringSubmatch(glob, -1)
	var pattern strings.Builder

	if full {
		pattern.WriteString("^")
	}

	for _, each := range result {
		switch each[0] {
		case "?":
			pattern.WriteString(".")
		case "*":
			pattern.WriteString(".*")
		default:
			pattern.WriteString(regexp.QuoteMeta(each[0]))
		}
	}

	if full {
		pattern.WriteString("$")
	}

	return pattern.String()
}

func asGlobRegex(glob string, full bool) regexp.Regexp {
	return *util.Must(regexp.Compile(asGlobRegexPattern(glob, full)))
}
