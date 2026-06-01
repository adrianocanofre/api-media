package utils

import (
	"path/filepath"
	"regexp"
	"strings"
)

var invalidChars = regexp.MustCompile(`[^a-zA-Z0-9._-]`)

func SanitizeFilename(name string) string {
	name = filepath.Base(name)

	name = strings.ToLower(name)

	name = strings.ReplaceAll(name, " ", "-")

	name = invalidChars.ReplaceAllString(name, "")

	return name
}
