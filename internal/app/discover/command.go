package discover

import (
	"fmt"
	"golang.org/x/mod/semver"
)

func BuildCommand(promotion Promotion, environment string) string {
	if environment == "" {
		environment = semver.Canonical(promotion.Semver())[1:]
	}

	// Calculate BASE_IMAGE_TAG
	var baseImageTag string
	if semver.Compare(promotion.Semver(), "v1.17.1") < 0 {
		baseImageTag = "16-jdk-alpine"
	} else {
		baseImageTag = "17-jre-alpine"
	}

	return fmt.Sprintf("drone build promote \"$DRONE_REPO\" \"$DRONE_BUILD_NUMBER\" %s --param=BASE_IMAGE_TAG=%s --param=DOWNLOAD_URL=%s --param=DOCKER_TAGS=%s", environment, baseImageTag, promotion.DownloadURL, promotion.DockerTags())
}
