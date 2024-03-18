package github

import (
	"regexp"
)

var lre = regexp.MustCompile(`<([^>]*)>; *rel="([^"]*)"`)

// Parse Link headers, returning a map from Rel to URL.
// Only understands the URI and "rel" parameter. Very limited.
// See https://tools.ietf.org/html/rfc5988#section-5
func parseLinks(h string) map[string]string {
	links := map[string]string{}
	for _, m := range lre.FindAllStringSubmatch(h, 10) {
		if len(m) != 3 {
			continue
		}
		links[m[2]] = m[1]
	}
	return links
}
