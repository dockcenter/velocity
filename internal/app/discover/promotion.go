package discover

import (
	"fmt"
	"github.com/dockcenter/velocity/internal/pkg/utils/slices"
	"golang.org/x/mod/semver"
	"sort"
	"strings"
)

type Event int

const (
	Rebuild Event = iota
	Cron
)

type Promotion struct {
	Environment string
	DownloadURL string
	DockerTags  string
}

func BuildPromotions(builds []VersionFamilyBuild, existingTags []string, event Event, environment string) []Promotion {
	var notExistingTags []int
	sharedTags := make(map[string]int)

	// Mark builds
	for i, build := range builds {
		// Check if tag is existed
		if !slices.Contains(existingTags, GetUniqueTag(build.Version, build.Build)) {
			notExistingTags = append(notExistingTags, i)
		}

		// handle stable tag
		if build.Promoted {
			sharedTags["stable"] = i
		}

		// handle latest tag
		if i == len(builds)-1 {
			sharedTags["latest"] = i
		}

		// handle version shared tags
		sharedTags[build.Version] = i

		// handle semver tags
		if semver.Prerelease("v"+build.Version) == "" {
			sharedTags[semver.Major("v" + build.Version)[1:]] = i
			sharedTags[semver.MajorMinor("v" + build.Version)[1:]] = i
			sharedTags[semver.Canonical("v" + build.Version)[1:]] = i
		}
	}

	// Filter builds
	var promotedBuildIndex []int
	if event == Rebuild {
		// Build semver, latest snapshot, latest and stable
		var latestSnapshot string
		var latestSnapshotIndex int
		for tag, index := range sharedTags {
			if semver.Prerelease("v"+tag) == "-SNAPSHOT" {
				if semver.Compare("v"+tag, "v"+latestSnapshot) > 0 {
					latestSnapshot = tag
					latestSnapshotIndex = index
				}
				continue
			}

			if !slices.Contains(promotedBuildIndex, index) {
				promotedBuildIndex = append(promotedBuildIndex, index)
			}
		}

		if !slices.Contains(promotedBuildIndex, latestSnapshotIndex) {
			promotedBuildIndex = append(promotedBuildIndex, latestSnapshotIndex)
		}
	} else {
		promotedBuildIndex = append(promotedBuildIndex, notExistingTags...)
	}

	sort.Ints(promotedBuildIndex)

	// Build promotions
	var promotions []Promotion
	for _, index := range promotedBuildIndex {
		build := builds[index]

		// If environment is not specified, use version as environment name
		promotionEnvironment := environment
		if environment == "" {
			promotionEnvironment = build.Version
		}

		// Build docker tags
		var tags []string
		tags = append(tags, GetUniqueTag(build.Version, build.Build))
		for k, v := range sharedTags {
			if v == index {
				tags = append(tags, k)
			}
		}

		promotion := Promotion{
			Environment: promotionEnvironment,
			DownloadURL: fmt.Sprintf("https://api.papermc.io/v2/projects/%s/versions/%s/builds/%d/downloads/%s", PROJECT, build.Version, build.Build, build.Downloads[DownloadsKey].Name),
			DockerTags:  strings.Join(tags, ","),
		}
		promotions = append(promotions, promotion)
	}

	return promotions
}
