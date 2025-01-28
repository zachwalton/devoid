package llm

import (
	"context"
	"encoding/json"
	goerrors "errors"
	"fmt"
	"strings"

	"github.com/zachwalton/devoid/pkg/brain"
	"github.com/zachwalton/devoid/pkg/brain/schema"
	"github.com/zachwalton/devoid/pkg/brain/templates"
	"github.com/zachwalton/devoid/pkg/errors"
	stagepkg "github.com/zachwalton/devoid/pkg/llm/stages"
	"github.com/zachwalton/devoid/pkg/tui"

	"github.com/charmbracelet/log"
)

type (
	templateFunc func(string) string

	Reasoner interface {
		Generate(context.Context, string, string, string) error
		ResponseCh() <-chan Response
	}

	HandlerFunc func(payload *brain.StagePayload) error

	Prompt struct {
		Message        string
		SystemTemplate string
		Schema         string
	}

	Response struct {
		Response string
		Done     bool
	}

	Stage struct {
		Payload            *brain.StagePayload
		LLM                bool
		Description        string
		Next               string
		SystemTemplateFunc func(string) string
		Schema             string
		HandlerFunc        HandlerFunc
		Final              bool
	}
)

var (
	stages = map[string]Stage{
		"initial": {
			LLM:                true,
			Description:        "This stage analyzes the prompt and figures out things like language, frameworks, initial scaffolding, and so on.",
			SystemTemplateFunc: templates.SystemInitial,
			Schema:             schema.SchemaInitial(),
			Next:               "scaffolding",
			HandlerFunc:        stagepkg.HandleInitial,
		},
		"scaffolding": {
			LLM:         false,
			Description: "This stage applies the changes from the 'initial' stage to the project",
			HandlerFunc: stagepkg.HandleScaffolding,
			Next:        "run",
			Final:       true,
		},
	}
)

func Start(ctx context.Context, reasoner Reasoner, prompt, projectDir string) chan bool {
	doneCh := make(chan bool)
	jsonCh := make(chan string)
	stage := "initial"
	choice := ""
	choiceChanges := "Request changes in a live chat session"
	choiceExit := "Exit the program"
	choiceAnswers := "Answer some questions to help improve this result before proceeding"
	choiceTryAgain := "I just don't like the response. Try again"

	var jsonResponse string
	iteration := 1
	log.Info("creating project...", "path", projectDir)
	go func() {
		defer func() { doneCh <- true }()

		for {
			var payload brain.StagePayload
			log.Info("starting stage", "stage", stage, "description", stages[stage].Description, "iteration", iteration)
			go func() {
				var (
					resp Response
					sb   strings.Builder
				)
				for {
					select {
					case <-ctx.Done():
						return
					case resp = <-reasoner.ResponseCh():
						sb.WriteString(resp.Response)
					}
					if resp.Done {
						jsonCh <- sb.String()
						return
					}
				}
			}()

			if iteration > 1 && choice != choiceTryAgain {
				llmCtx, cancel := context.WithCancel(ctx)
				fmt.Println()
				spinner := tui.Spinner(llmCtx, cancel, "Working with the LLM on some changes...")
				defer spinner.Stop()

				if stages[stage].LLM {
					if err := reasoner.Generate(
						llmCtx, prompt,
						stages[stage].Schema, templates.SystemClarify(jsonResponse),
					); err != nil {
						spinner.Stop()
						log.Error("got an error during inference", "error", err)
						return
					}
					spinner.Stop()

					jsonResponse = <-jsonCh

					if err := json.Unmarshal([]byte(jsonResponse), &payload); err != nil {
						log.Error("got an error unmarshaling payload", "payload", jsonResponse, "error", err)
						return
					}
				}
			} else {
				llmCtx, cancel := context.WithCancel(ctx)
				fmt.Println()
				spinner := tui.Spinner(llmCtx, cancel, "Chatting with the LLM...")
				defer spinner.Stop()

				if stages[stage].LLM {
					if err := reasoner.Generate(
						llmCtx, prompt, stages[stage].Schema,
						stages[stage].SystemTemplateFunc(projectDir),
					); err != nil {
						spinner.Stop()
						log.Error("got an error during inference", "error", err)
						return
					}
					spinner.Stop()

					jsonResponse = <-jsonCh

					if err := json.Unmarshal([]byte(jsonResponse), &payload); err != nil {
						log.Error("got an error unmarshaling payload", "payload", jsonResponse, "error", err)
						return
					}
				}
			}

			payload.Meta.ProjectPath = projectDir
			err := stages[stage].HandlerFunc(&payload)
			if err != nil {
				switch {
				case goerrors.Is(err, errors.ErrRecoverable):
					prompt = stagepkg.UpdatePromptForErr(stage, err)
					iteration++
					continue
				default:
					log.Error("got an error handling stage", "stage", stage, "error", err)
					return
				}
			}
			log.Info("successfully applied stage", "stage", stage)
			s := stages[stage]
			s.Payload = &payload
			stages[stage] = s

			if stages[stage].LLM {
				if iteration > 1 {
					payload.StateMachine.ModifiedResult = true
				}
				tui.MarkdownView(payload.Markdown(stage, projectDir))
			}

			if stages[stage].Final {
				log.Info("All stages have been completed!")
				return
			}

			if !stages[stage].LLM {
				stage = stages[stage].Next
				prompt = stages[stage].Description
				iteration = 1
				continue
			}

			choiceMoveAhead := fmt.Sprintf("Move ahead to the '%s' stage", payload.StateMachine.Next)

			choices := []string{
				choiceMoveAhead,
				choiceChanges,
				choiceTryAgain,
				choiceExit,
			}
			if len(payload.StateMachine.Questions) != 0 {
				choices = append([]string{choiceAnswers}, choices...)
			}

			selected := false
			for !selected {
				choice := tui.List(choices)

				switch choice {
				case choiceMoveAhead:
					stage = payload.StateMachine.Next
					prompt = stages[stage].Description
					selected = true
					iteration = 1
				case choiceChanges:
					iteration++
					var addendum string
					for {
						addendum = tui.Input("What do you want to ask or tell the model?")
						if addendum == "" {
							log.Warn("You didn't enter any text! Try again...")
							continue
						}
						break
					}
					selected = true
					prompt = addendum
				case choiceExit:
					log.Info("exiting by user request...")
					return
				case choiceAnswers:
					iteration++
					p := strings.Builder{}
					for _, question := range payload.StateMachine.Questions {
						var answer string
						for {
							answer = tui.Input(question)
							if answer == "" {
								log.Warn("You didn't enter any text! Try again...")
								continue
							}
							break
						}
						p.WriteString("Question: ")
						p.WriteString(question + "\n")
						p.WriteString("Answer: ")
						p.WriteString(answer + "\n")
					}
					selected = true
					prompt = p.String()
				case choiceTryAgain:
					iteration++
					selected = true
				}
			}
		}
	}()
	return doneCh
}
