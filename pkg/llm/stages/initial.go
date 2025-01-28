package stages

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/zachwalton/devoid/pkg/brain"
	"github.com/zachwalton/devoid/pkg/config"
	"github.com/zachwalton/devoid/pkg/errors"
	"github.com/zachwalton/devoid/pkg/tui"
)

// HandleInitial performs safety checks on ProjectLayout paths, so the scaffolding
// stage is able to create them confidently
func HandleInitial(payload *brain.StagePayload, cfg *config.Config) error {
	// Ensure ProjectPath is set
	if payload.Meta.ProjectPath == "" {
		return errors.ErrNoProjectPath
	}

	// Clean the ProjectPath to remove any redundant slashes or relative components
	projectPath, err := filepath.Abs(payload.Meta.ProjectPath)
	if err != nil {
		return fmt.Errorf("%s: %w", errors.ErrAbsolutePath, err)
	}

	// Iterate through ProjectLayout paths
	for _, path := range payload.Stages.Scaffolding.ProjectLayout {
		// Resolve the path relative to ProjectPath
		resolvedPath := filepath.Join(projectPath, path)

		// Get the absolute version of the resolved path
		absolutePath, err := filepath.Abs(resolvedPath)
		if err != nil {
			return fmt.Errorf("%s: %w", errors.ErrAbsolutePath, err)
		}

		// Check if the resolved path is within the ProjectPath
		if !strings.HasPrefix(absolutePath, projectPath) {
			return fmt.Errorf("%w: %s", errors.ErrProjectPathEscape, projectPath)
		}
	}

	if len(payload.Stages.Scaffolding.BootstrapCommands) > 0 {
		log.Warn(
			"review these commands carefully, as they will be run in the next stage, and make a selection below",
			"commands",
			payload.Stages.Scaffolding.BootstrapCommands,
		)
		choiceConfirm := "I have reviewed these commands and would like to run them"
		choiceRecover := "I don't want to run these commands, and want to describe how they should be fixed"
		choiceSkip := "Skipping interactive safety check for bootstrap commands"
		choiceExit := "I don't want to run these commands and want to exit"

		choice := ""
		if cfg.SkipInteractiveSafetyChecks {
			choice = choiceSkip
		} else {
			choice = tui.List([]string{
				choiceConfirm,
				choiceRecover,
				choiceExit,
			})
		}
		switch choice {
		case choiceConfirm:
			break
		case choiceSkip:
			log.Warn("interactive safety checks are disabled, moving ahead without prompting")
			return nil
		case choiceRecover:
			response := tui.Input("How do you want to change the commands?")
			if response == "" {
				return fmt.Errorf("%w: no guidance provided on how to address unwanted commands", errors.ErrBadCommands)
			}
			return fmt.Errorf("%w: %s: %s", errors.ErrRecoverable, errors.ErrBadCommands, response)
		case choiceExit:
			return errors.ErrBadCommands
		}
	}

	log.Info(
		"completed safety checks on project layout paths and bootstrap commands",
		"paths",
		payload.Stages.Scaffolding.ProjectLayout,
		"commands",
		payload.Stages.Scaffolding.BootstrapCommands,
	)
	// Return nil if all paths pass validation
	return nil
}
