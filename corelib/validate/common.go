package validate


import (
	"regexp"
)

var semverRegex = regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)$`)

func isValidSemver(versionString string) bool {
    return semverRegex.MatchString(versionString)
}