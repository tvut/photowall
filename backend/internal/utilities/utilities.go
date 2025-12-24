package utilities

import (
	"regexp"
	"strings"
)

var re = regexp.MustCompile("[^a-zA-Z0-9-]+")

// ToSlug converts a string to a URL-friendly slug
func ToSlug(title string) string {
	return re.ReplaceAllString(strings.ToLower(strings.ReplaceAll(title, " ", "-")), "")
}
