package discover

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"time"
)

func GetExistingTags(repository string) []DockerTag {
	client := resty.New()
	url := fmt.Sprintf("https://registry.hub.docker.com/v2/repositories/%s/tags", repository)

	var tags []DockerTag
	for ok := true; ok; ok = url != "" {
		var response DockerTagsResponse
		resp, err := client.R().Get(url)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(resp.Body(), &response)
		if err != nil {
			panic(err)
		}

		tags = append(tags, response.Results...)
		url = response.Next
	}

	return tags
}

func GetUniqueTag(version string, build int) string {
	return fmt.Sprintf("%s-%d", version, build)
}

type DockerTagsResponse struct {
	Count    int         `json:"count"`
	Next     string      `json:"next"`
	Previous string      `json:"previous"`
	Results  []DockerTag `json:"results"`
}

type DockerTag struct {
	Creator             int           `json:"creator"`
	FullSize            int           `json:"full_size"`
	Id                  int           `json:"id"`
	ImageId             int           `json:"image_id"`
	Images              []DockerImage `json:"images"`
	LastUpdated         time.Time     `json:"last_updated"`
	LastUpdater         int           `json:"last_updater"`
	LastUpdaterUsername string        `json:"last_updater_username"`
	Name                string        `json:"name"`
	Repository          int           `json:"repository"`
	TagLastPulled       time.Time     `json:"tag_last_pulled"`
	TagLastPushed       time.Time     `json:"tag_last_pushed"`
	TagStatus           string        `json:"tag_status"`
	V2                  bool          `json:"v2"`
}

type DockerImage struct {
	Architecture string    `json:"architecture"`
	Digest       string    `json:"digest"`
	Features     string    `json:"features"`
	LastPulled   time.Time `json:"last_pulled"`
	LastPushed   time.Time `json:"last_pushed"`
	Os           string    `json:"os"`
	OsFeatures   string    `json:"os_features"`
	OsVersion    string    `json:"os_version"`
	Size         int       `json:"size"`
	Status       string    `json:"status"`
	Variant      string    `json:"variant"`
}
