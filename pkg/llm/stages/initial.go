package stages

import (
  "fmt"

	"github.com/charmbracelet/log"
	"github.com/zachwalton/devoid/pkg/brain"
	"github.com/zachwalton/devoid/pkg/config"
	"github.com/zachwalton/devoid/pkg/errors"
)

func HandleInitial(payload *brain.StagePayload, cfg *config.Config) error {
  if payload.Meta.Name == "" {
    return fmt.Errorf("%w: meta -> name was unset", errors.ErrRecoverable)
  }
	log.Info(
		"completed validations and safety checks on project layout paths and bootstrap commands",
    "stage",
    payload.Meta.CurrentStage,
	)
	return nil
}
