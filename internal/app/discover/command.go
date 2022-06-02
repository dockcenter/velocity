package discover

import (
	"fmt"
)

func BuildCommand(promotion Promotion) string {
	return fmt.Sprintf("drone build promote \"$DRONE_REPO\" \"$DRONE_BUILD_NUMBER\" %s --param=DOWNLOAD_URL=%s --param=DOCKER_TAGS=%s", promotion.Environment, promotion.DownloadURL, promotion.DockerTags)
}
