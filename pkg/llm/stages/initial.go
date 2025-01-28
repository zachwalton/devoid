package stages

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/zachwalton/devoid/pkg/brain"
	"github.com/zachwalton/devoid/pkg/errors"
)

// HandleInitial performs safety checks on ProjectLayout paths, so the scaffolding
// stage is able to create them confidently
func HandleInitial(payload *brain.StagePayload) error {
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

	log.Info(
		"completed safety checks on project layout paths; either there are none or all are relative to the project dir",
		"paths",
		payload.Stages.Scaffolding.ProjectLayout,
	)
	// Return nil if all paths pass validation
	return nil
}
