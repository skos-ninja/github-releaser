package version

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/Masterminds/semver"
)

type VersionType string

const (
	Major      = VersionType("MAJOR")
	Minor      = VersionType("MINOR")
	Patch      = VersionType("PATCH")
	Prerelease = VersionType("PRERELEASE")
)

var (
	ErrNotValidVersion    = errors.New("Version provided is not a valid semver")
	ErrNotValidPrerelease = errors.New("Version provided is not a valid prerelease")
)

func FindNextVersion(t VersionType, v string) (string, error) {
	ver, err := semver.NewVersion(v)
	if err != nil {
		return "", ErrNotValidVersion
	}
	version := *ver

	switch t {
	case Major:
		version = version.IncMajor()
	case Minor:
		version = version.IncMinor()
	case Patch:
		version = version.IncPatch()
	case Prerelease:
		pre := version.Prerelease()
		if pre == "" {
			// As the current pre version doesn't match
			version, err = version.IncPatch().SetPrerelease("pre0")
			if err != nil {
				return "", err
			}
		} else {
			// We strip `pre` from the string before parsing
			preNum, err := strconv.ParseInt(pre[3:], 10, 64)
			if err != nil {
				return "", ErrNotValidPrerelease
			}

			version, err = version.SetPrerelease(fmt.Sprintf("pre%v", preNum+1))
			if err != nil {
				return "", err
			}
		}
	}

	if v[:1] == "v" {
		return "v" + version.String(), nil
	}
	return version.String(), nil
}

func GetVersionPriority(v VersionType) int {
	switch v {
	case Prerelease:
		return 1
	case Patch:
		return 2
	case Minor:
		return 3
	case Major:
		return 4
	}

	return 0
}

func ParseVersionType(versionType string) (VersionType, error) {
	t := VersionType(versionType)
	switch t {
	case Major, Minor, Patch, Prerelease:
		return t, nil
	}

	return VersionType(""), errors.New("Unknown version type")
}
