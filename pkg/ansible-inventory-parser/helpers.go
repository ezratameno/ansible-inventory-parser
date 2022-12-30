package ansibleinventoryparser

import "strings"

func parseName(s string) string {
	// remove indent.
	s = strings.TrimSpace(s)

	s = strings.ReplaceAll(s, ":", "")
	// remove inline comment.
	s = strings.Split(s, "#")[0]

	return s

}
