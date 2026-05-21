package validate


import (
	"regexp"
)

var semverRegex = regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)$`)

func isValidSemver(versionString string) bool {
    return semverRegex.MatchString(versionString)
}

func isValidDateVer(versionString string) bool {
    return semverRegex.MatchString(versionString)
}



func isValidVersioningSchema(versioning string) bool {
    switch versioning {
    case "semver", "date":
        return true
    default:
        return false
    }
}

func isValidVersion(versioning string, versionString string) bool {

	switch versioning {
    case "semver":
        return isValidSemver(versionString)
	case "date":
		return isValidDateVer(versionString)
		
    default:
        return false
    }

}