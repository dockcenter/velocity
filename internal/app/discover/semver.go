package discover

import (
	"golang.org/x/mod/semver"
	"strconv"
)

func MarkSemver(promotions []Promotion) {
	// Build semverMap
	semverMap := make(map[string]string)
	for _, promotion := range promotions {
		// Record largest semver information
		majorMinor := semver.MajorMinor(promotion.Semver())
		canonical := semver.Canonical(promotion.Semver())

		// Compare latest
		ver, ok := semverMap["latest"]
		if ok {
			compare := semver.Compare(promotion.Semver(), ver)
			promotionBuild, _ := strconv.Atoi(semver.Build(promotion.Semver()))
			verBuild, _ := strconv.Atoi(semver.Build(ver))
			if compare > 0 || (compare == 0 && promotionBuild > verBuild) {
				semverMap["latest"] = promotion.Semver()
			}
		} else {
			semverMap["latest"] = promotion.Semver()
		}

		// Compare MajorMinor
		ver, ok = semverMap[majorMinor]
		if ok {
			compare := semver.Compare(promotion.Semver(), ver)
			promotionBuild, _ := strconv.Atoi(semver.Build(promotion.Semver()))
			verBuild, _ := strconv.Atoi(semver.Build(ver))
			if compare > 0 || (compare == 0 && promotionBuild > verBuild) {
				semverMap[majorMinor] = promotion.Semver()
			}
		} else {
			semverMap[majorMinor] = promotion.Semver()
		}

		// Compare Canonical
		ver, ok = semverMap[canonical]
		if ok {
			promotionBuild, _ := strconv.Atoi(semver.Build(promotion.Semver()))
			verBuild, _ := strconv.Atoi(semver.Build(ver))
			if promotionBuild > verBuild {
				semverMap[canonical] = promotion.Semver()
			}
		} else {
			semverMap[canonical] = promotion.Semver()
		}
	}

	// Mark semver
	for i, promotion := range promotions {
		majorMinor := semver.MajorMinor(promotion.Semver())
		canonical := semver.Canonical(promotion.Semver())

		if promotion.Semver() == semverMap["latest"] {
			promotions[i].Latest = true
			promotions[i].Major = true
		}
		if promotion.Semver() == semverMap[majorMinor] {
			promotions[i].MajorMinor = true
		}
		if promotion.Semver() == semverMap[canonical] {
			promotions[i].Canonical = true
		}
	}
}
