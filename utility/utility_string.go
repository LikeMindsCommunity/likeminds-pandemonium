package utility

import "strings"

// JoinStrs returns joinchar concatenated string
func JoinStrs(joinStrsList []string, joinChar string) string {
	return strings.Join(joinStrsList, joinChar)
}
