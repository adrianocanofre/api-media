package utils

import (
	"path/filepath"
	"regexp"
	"strings"
)

var invalidChars = regexp.MustCompile(`[^a-zA-Z0-9._-]`)

func SanitizeFilename(name string) string {
	// remove qualquer path enviado pelo client
	name = filepath.Base(name)

	// lowercase
	name = strings.ToLower(name)

	// troca espaços por "-"
	name = strings.ReplaceAll(name, " ", "-")

	// remove caracteres inválidos
	name = invalidChars.ReplaceAllString(name, "")

	return name
}
