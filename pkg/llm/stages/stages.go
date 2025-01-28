package stages

import (
	"fmt"

	"github.com/charmbracelet/log"
)

func UpdatePromptForErr(stage string, err error) string {
	log.Warn("got a recoverable error handling stage. going to try again", "stage", stage, "error", err)
	return fmt.Sprintf(
		"I got this error, please address it in the respons JSON: %s",
		err,
	)
}
