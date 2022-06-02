package main

import (
	"encoding/json"
	"fmt"
	. "github.com/dockcenter/velocity/internal/app/discover"
	"github.com/go-resty/resty/v2"
	"os"
	"strconv"
)

func main() {
	// Parse environment variables
	event := os.Getenv("DRONE_BUILD_EVENT")
	branch := os.Getenv("DRONE_BRANCH")
	environment := os.Getenv("ENVIRONMENT")
	dryRunStr := os.Getenv("DRY_RUN")
	dryRun, err := strconv.ParseBool(dryRunStr)
	if err != nil {
		dryRun = false
	}
	fmt.Println("Trigger event:", event)
	fmt.Println("Branch:", branch)
	fmt.Println("Dry run:", dryRun)

	client := resty.New()

	// Get velocity builds
	var versionFamilyBuilds VersionFamilyBuildsResponse
	url := fmt.Sprintf("https://api.papermc.io/v2/projects/%s/version_group/%s/builds", PROJECT, SupportedVersionGroup)
	resp, err := client.R().Get(url)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(resp.Body(), &versionFamilyBuilds)
	if err != nil {
		panic(err)
	}

	// Get pushed tags
	var tags []string
	dockerTags := GetExistingTags(DockerRepository)
	for _, dockerTag := range dockerTags {
		tags = append(tags, dockerTag.Name)
	}

	// build promotions
	var eventForPromotions Event
	if event == "push" && branch == "main" {
		eventForPromotions = Rebuild
	} else {
		eventForPromotions = Cron
	}
	promotions := BuildPromotions(versionFamilyBuilds.Builds, tags, eventForPromotions, environment)

	// Print tags to promotion
	fmt.Println("\nTags to promote:")
	for _, promotion := range promotions {
		fmt.Println(promotion.DockerTags)
	}

	// Build promote commands and write to scripts/promote.sh
	cmd := "#!/bin/sh\n\n"
	for _, promotion := range promotions {
		cmd += BuildCommand(promotion) + "\n"
	}

	// Create scripts folder
	err = os.MkdirAll("scripts", 0700)
	if err != nil {
		panic(err)
	}

	if dryRun {
		// Create empty scripts/promote.sh
		err := os.WriteFile("scripts/promote.sh", []byte("#!/bin/sh\n"), 0700)
		if err != nil {
			panic(err)
		}

		// print scripts content
		fmt.Println("\nThis is a dry run")
		fmt.Println("We generate the following script but not write to scripts/promote.sh")
		fmt.Println("\n" + cmd)
	} else {
		// Write to scripts/promote.sh
		err = os.WriteFile("scripts/promote.sh", []byte(cmd), 0700)
		if err != nil {
			panic(err)
		}

		fmt.Println("\nShell script has been generated to scripts/promote.sh")
	}
}
