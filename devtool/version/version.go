package version

import (
	"strconv"
	"strings"
)

type Version struct {
	Patch int
	Minor int
	Major int
}

func (v *Version) String() string {
	return strconv.Itoa(v.Major) + "." + strconv.Itoa(v.Minor) + "." + strconv.Itoa(v.Patch)
}

func Parse(version string) (*Version, error) {
	parts := strings.Split(version, ".")
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}
	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, err
	}
	return &Version{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil
}
