package versioning

import (
	"strconv"
	"strings"
)

type VersionPart int

const (
	Major VersionPart = iota + 1 // update major version number still won't work
	Minor
	Patch
	RC
)

func IncreaseVersion(version string, part VersionPart) string {
	parts := strings.Split(version, ".")
	var major, minor, patch, rc string
	if len(parts) == 4 {
		major, minor, patch, rc = parts[0], parts[1], parts[2], parts[3]
	} else if len(parts) == 3 {
		major, minor, patch = parts[0], parts[1], parts[2]
	} else {
		return "Invalid version string format"
	}

	switch part {
	case Major:
		major = strconv.Itoa(toInt(major) + 1)
	case Minor:
		minor = strconv.Itoa(toInt(minor) + 1)
	case Patch:
		patch = strconv.Itoa(toInt(patch) + 1)
	case RC:
		rc = "rc." + strconv.Itoa(toInt(strings.Split(rc, ".")[1])+1)
	default:
		return "Invalid version part"
	}

	if len(rc) == 0 {
		return strings.Join([]string{major, minor, patch}, ".")
	}
	return strings.Join([]string{major, minor, patch, rc}, ".")
}

func GetVersion(s string) string {
	// Split the string by "/"
	split := strings.Split(s, "/")

	// Get the last element of the split slice
	version := split[len(split)-1]

	// Split the version string by "-"
	versionSplit := strings.Split(version, "-")

	// Return the first element of the version split slice
	return versionSplit[0]
}

func toInt(value string) int {
	i, err := strconv.Atoi(strings.TrimPrefix(value, "v"))
	if err != nil {
		return 0
	}
	return i
}
