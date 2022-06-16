package discover

import (
	"encoding/json"
	"fmt"
)

func BuildCommand(workflow string, promotion Promotion) string {
	// Convert promotion to json
	jsonBytes, err := json.Marshal(promotion)
	if err != nil {
		panic(err)
	}
	jsonStr := string(jsonBytes)

	return fmt.Sprintf("echo '%s' | gh workflow run %s --json", jsonStr, workflow)
}
