package cmd

import (
	"context"
	"fmt"

	"github.com/zachwalton/devoid/pkg/config"
	"github.com/zachwalton/devoid/pkg/errors"
	"github.com/zachwalton/devoid/pkg/llm"

	"github.com/urfave/cli/v3"
)

var Cmd = &cli.Command{
	Name:      "main",
	ArgsUsage: "prompt: the prompt used to bootstrap the project",
	Usage:     "Generate a codebase from scratch interactively",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "project-path",
			Usage:    "Path to the directory where the project should be created. The directory should be empty, and will be created if it doesn't exist",
			Required: true,
		},
		&cli.BoolFlag{
			Name:  "skip-interactive-safety-checks",
			Usage: "When true, commands will be run without prompting. Use cautiously",
			Value: false,
		},
		&cli.StringFlag{
			Name:  "llm.model",
			Usage: "Name of the model to use, e.g. deepseek-r1:8b for Ollama",
			Value: "deepseek-r1:8b",
		},
		&cli.StringFlag{
			Name:  "llm.type",
			Usage: "LLM to use. Currently only 'ollama'",
			Value: "ollama",
		},
		&cli.FloatFlag{
			Name:  "llm.temperature",
			Usage: "Temperature to provide to the model",
			Value: .6,
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		cfg, err := setUpCfg(cmd)
		if err != nil {
			return err
		}

		var reasoner llm.Reasoner

		switch cfg.LLM.Type {
		case config.ReasonerOllama:
			if reasoner, err = llm.NewOllamaReasoner(cfg); err != nil {
				return fmt.Errorf("could not set up reasoner: %w", err)
			}
		default:
			return fmt.Errorf("%w: %s", errors.ErrUnknownType, cfg.LLM.Type)
		}
		<-llm.Start(
			ctx,
			reasoner,
			cfg.Prompt,
			cfg.ProjectPath,
			cfg,
		)
		return nil
	},
}

func setUpCfg(cmd *cli.Command) (*config.Config, error) {
	prompt := cmd.Args().Get(0)
	if prompt == "" {
		return nil, errors.ErrNoPrompt
	}
	return &config.Config{
		ProjectPath:                 cmd.String("project-path"),
		SkipInteractiveSafetyChecks: cmd.Bool("skip-interactive-safety-checks"),
		Prompt:                      prompt,
		LLM: config.LLM{
			Model:       cmd.String("llm.model"),
			Type:        config.Reasoner(cmd.String("llm.type")),
			Temperature: cmd.Float("llm.temperature"),
		},
	}, nil
}
