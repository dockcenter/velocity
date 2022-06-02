package main

import (
	"encoding/json"
	"fmt"
	. "github.com/dockcenter/velocity/internal/app/discover"
	"github.com/dockcenter/velocity/internal/pkg/utils/slices"
	"github.com/go-resty/resty/v2"
	"os"
	"time"
)

func main() {
	client := resty.New()
	const PROJECT string = "paper"
	const DockerRepository string = "dockcenter/paper"
	const SupportedVersionGroup int = 2
	const DownloadsKey string = "application"

	// Parse environment variables
	event := os.Getenv("DRONE_BUILD_EVENT")
	branch := os.Getenv("DRONE_BRANCH")
	duration := os.Getenv("DURATION")
	environment := os.Getenv("ENVIRONMENT")
	fmt.Println("Trigger event:", event)
	fmt.Println("Branch:", branch)
	fmt.Println("Duration:", duration)

	// Get paper versions
	var project ProjectResponse
	url := fmt.Sprintf("https://api.papermc.io/v2/projects/%s", PROJECT)
	resp, err := client.R().Get(url)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(resp.Body(), &project)
	if err != nil {
		panic(err)
	}

	// When pushing to main, promote all supported versions' latest build
	// Otherwise, get all builds in duration in supported vrsion groups
	var promotions []Promotion
	if event == "push" && branch == "main" {
		// Iterate all version groups
		var versions []string
		for _, versionGroup := range project.VersionGroups[len(project.VersionGroups)-SupportedVersionGroup:] {
			// Get all versions in version group
			var versionFamily VersionFamilyResponse
			url := fmt.Sprintf("https://api.papermc.io/v2/projects/%s/version_group/%s", PROJECT, versionGroup)
			resp, err := client.R().Get(url)
			if err != nil {
				panic(err)
			}
			err = json.Unmarshal(resp.Body(), &versionFamily)
			if err != nil {
				panic(err)
			}

			versions = append(versions, versionFamily.Versions...)
		}

		// Get all builds of versions
		for _, version := range versions {
			// Get all builds for specific version
			var builds BuildsResponse
			url := fmt.Sprintf("https://api.papermc.io/v2/projects/%s/versions/%s/builds", PROJECT, version)
			resp, err := client.R().Get(url)
			if err != nil {
				panic(err)
			}
			err = json.Unmarshal(resp.Body(), &builds)
			if err != nil {
				panic(err)
			}

			// Build promotion
			var promotion Promotion
			promotion.Version = builds.Version

			// Select latest build in each version
			build := builds.Builds[len(builds.Builds)-1]
			promotion.Build = build.Build
			promotion.DownloadURL = fmt.Sprintf("https://api.papermc.io/v2/projects/%s/versions/%s/builds/%d/downloads/%s", PROJECT, promotion.Version, promotion.Build, build.Downloads[DownloadsKey].Name)
			promotions = append(promotions, promotion)
		}
	} else {
		// If it's not cron and environment is not specified, set it to development
		if event != "cron" && environment == "" {
			environment = "development"
		}

		// Get all docker tags when duration is not set
		var tagNames []string
		if duration == "" {
			tags := GetAllTags(DockerRepository)
			for _, tag := range tags {
				tagNames = append(tagNames, tag.Name)
			}
		}

		// Get all builds for supported version groups
		for _, versionGroup := range project.VersionGroups[len(project.VersionGroups)-SupportedVersionGroup:] {
			var versionFamilyBuilds VersionFamilyBuildsResponse
			url := fmt.Sprintf("https://api.papermc.io/v2/projects/%s/version_group/%s/builds", PROJECT, versionGroup)
			resp, err := client.R().Get(url)
			if err != nil {
				panic(err)
			}
			err = json.Unmarshal(resp.Body(), &versionFamilyBuilds)
			if err != nil {
				panic(err)
			}

			builds := versionFamilyBuilds.Builds
			var versionGroupPromotions []Promotion
			for i := len(builds) - 1; i >= 0; i-- {
				build := builds[i]

				// Filter builds
				if duration != "" {
					duration, err := time.ParseDuration(duration)
					if err != nil {
						panic(err)
					}

					// Filter out builds that are longer than duration
					if time.Since(build.Time) > duration {
						continue
					}
				} else {
					// Filter out build and following existed in Docker Hub
					if slices.ContainsString(tagNames, GetTagName(build.Version, build.Build)) {
						break
					}
				}

				// Build promotion and append to promotions
				var promotion Promotion
				promotion.Version = build.Version
				promotion.Build = build.Build
				promotion.DownloadURL = fmt.Sprintf("https://api.papermc.io/v2/projects/%s/versions/%s/builds/%d/downloads/%s", PROJECT, promotion.Version, promotion.Build, build.Downloads[DownloadsKey].Name)
				versionGroupPromotions = append(versionGroupPromotions, promotion)
			}

			// Reverse versionGroupPromotions and append to promotions
			slices.Reverse(versionGroupPromotions)
			promotions = append(promotions, versionGroupPromotions...)
		}
	}

	MarkSemver(promotions)

	// Print promotion environment and tags to promote
	fmt.Println("Promotion environment:", environment)
	fmt.Println("\nTags to promote:")
	for _, promotion := range promotions {
		fmt.Println(promotion.DockerTags())
	}

	// Build promote commands and write to scripts/promote.sh
	cmd := "#!/bin/sh\n\n"
	for _, promotion := range promotions {
		cmd += BuildCommand(promotion, environment) + "\n"
	}

	// Write to scripts/promote.sh
	err = os.MkdirAll("scripts", 0700)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("scripts/promote.sh", []byte(cmd), 0700)
	if err != nil {
		panic(err)
	}

	fmt.Println("\nShell script has been generated to scripts/promote.sh")
}
