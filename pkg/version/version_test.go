package version

import (
	"fmt"
	"testing"
)

const (
	v = "v1.1.1"
)

func TestMajorVersion(t *testing.T) {
	nextVersion, err := FindNextVersion(Major, v)
	if err != nil {
		t.Error(err)
	}

	if nextVersion != "v2.0.0" {
		t.Errorf("Expected v2.0.0 got %s", nextVersion)
	}
}

func TestMinorVersion(t *testing.T) {
	nextVersion, err := FindNextVersion(Minor, v)
	if err != nil {
		t.Error(err)
	}

	if nextVersion != "v1.2.0" {
		t.Errorf("Expected v1.2.0 got %s", nextVersion)
	}
}

func TestPatchVersion(t *testing.T) {
	nextVersion, err := FindNextVersion(Patch, v)
	if err != nil {
		t.Error(err)
	}

	if nextVersion != "v1.1.2" {
		t.Errorf("Expected v1.1.2 got %s", nextVersion)
	}
}

func TestPreVersion(t *testing.T) {
	nextVersion, err := FindNextVersion(Prerelease, v)
	if err != nil {
		t.Error(err)
	}

	if nextVersion != "v1.1.1-pre0" {
		t.Errorf("Expected v1.1.1-pre0 got %s", nextVersion)
	}
}

func TestIncreaseVersion(t *testing.T) {
	type Test struct {
		currentVersion      string
		versionIncreaseType VersionType
		nextVersion         string
	}
	tests := []Test{
		{
			currentVersion:      "v1.1.1",
			versionIncreaseType: Major,
			nextVersion:         "v2.0.0",
		},
		{
			currentVersion:      "v1.1.1",
			versionIncreaseType: Minor,
			nextVersion:         "v1.2.0",
		},
		{
			currentVersion:      "v1.1.1",
			versionIncreaseType: Patch,
			nextVersion:         "v1.1.2",
		},
		{
			currentVersion:      "v1.1.1",
			versionIncreaseType: Prerelease,
			nextVersion:         "v1.1.2-pre0",
		},
		{
			currentVersion:      "v1.1.1-pre0",
			versionIncreaseType: Prerelease,
			nextVersion:         "v1.1.1-pre1",
		},
	}
	for _, test := range tests {
		test := test
		testName := fmt.Sprintf("%s increase %s to %s", test.currentVersion, test.versionIncreaseType, test.nextVersion)
		t.Run(testName, func(t *testing.T) {
			nextVer, err := FindNextVersion(test.versionIncreaseType, test.currentVersion)
			if err != nil {
				t.Error(err)
			}

			if nextVer != test.nextVersion {
				t.Errorf("expected %s, got %s", test.nextVersion, nextVer)
			}
		})
	}
}
