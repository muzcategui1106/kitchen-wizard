package format

import (
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	extraSpaceRegex = regexp.MustCompile(`\s+`)
	caser           = cases.Title(language.AmericanEnglish)
)

func ComplyAsTitle(s string) string {
	return strings.TrimSpace(extraSpaceRegex.ReplaceAllString(caser.String(s), " "))
}
