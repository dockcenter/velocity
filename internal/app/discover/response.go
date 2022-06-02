package discover

import "time"

type ProjectBase struct {
	ProjectID   string `json:"project_id"`
	ProjectName string `json:"project_name"`
}

type ProjectResponse struct {
	ProjectBase
	VersionGroups []string `json:"version_groups"`
	Versions      []string `json:"versions"`
}

type VersionFamilyResponse struct {
	ProjectBase
	VersionGroup string   `json:"version_group"`
	Versions     []string `json:"versions"`
}

type BuildsResponse struct {
	ProjectBase
	Version string         `json:"version"`
	Builds  []VersionBuild `json:"builds"`
}

type VersionBuild struct {
	Build     int                 `json:"build"`
	Time      time.Time           `json:"time"`
	Channel   string              `json:"channel"`
	Promoted  bool                `json:"promoted"`
	Changes   []Change            `json:"change"`
	Downloads map[string]Download `json:"downloads"`
}

type Change struct {
	Commit  string `json:"commit"`
	Summary string `json:"summary"`
	Message string `json:"message"`
}

type Download struct {
	Name   string `json:"name"`
	Sha256 string `json:"sha256"`
}

type VersionFamilyBuildsResponse struct {
	VersionFamilyResponse
	Builds []VersionFamilyBuild `json:"builds"`
}

type VersionFamilyBuild struct {
	VersionBuild
	Version string `json:"version"`
}
