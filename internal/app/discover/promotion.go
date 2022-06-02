package discover

import (
	"fmt"
	"strings"

	"golang.org/x/mod/semver"
)

type Promotion struct {
	Version     string
	Build       int
	DownloadURL string
	Latest      bool
	Major       bool
	MajorMinor  bool
	Canonical   bool
}

func (p Promotion) Semver() string {
	return fmt.Sprintf("%s+%d", semver.Canonical("v"+p.Version), p.Build)
}

func (p Promotion) DockerTags() string {
	var tags []string
	tags = append(tags, strings.ReplaceAll(p.Semver()[1:], "+", "-r"))

	if p.Canonical {
		tags = append(tags, semver.Canonical(p.Semver())[1:])
	}

	if p.MajorMinor {
		tags = append(tags, semver.MajorMinor(p.Semver())[1:])
	}

	if p.Major {
		tags = append(tags, semver.Major(p.Semver())[1:])
	}

	if p.Latest {
		tags = append(tags, "latest")
	}

	return strings.Join(tags, ",")
}
