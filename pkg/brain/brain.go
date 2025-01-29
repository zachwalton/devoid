package brain

import (
	"bytes"
	"text/template"
)

type StagePayload struct {
	Meta         MetaPayload         `json:"meta"`
	StateMachine StateMachinePayload `json:"state_machine"`
}

type MetaPayload struct {
	Name         string `json:"name"`
	Language     string `json:"language"`
	Test         string `json:"test"`
	Framework    string `json:"framework"`
	Architecture string `json:"architecture"`
	Description  string `json:"description"`
	CurrentStage string `json:"-"`
	ProjectPath  string `json:"-"`
}

type StateMachinePayload struct {
	Next           string   `json:"next"`
	Final          bool     `json:"final"`
	ModifiedResult bool     `json:"modified_result"`
	Description    string   `json:"description"`
	Questions      []string `json:"questions"`
}

func (p *StagePayload) Markdown(stage, projectPath string) string {
	p.Meta.CurrentStage = stage
	p.Meta.ProjectPath = projectPath

	t, _ := template.New("markdown").Parse(`
# {{.Meta.Name}}

This is a summary of current status. When you're ready, press ` + "`" + `q` + "`" + ` and you'll be presented with some options.

{{ if .StateMachine.ModifiedResult }}
## Changes Applied

The changes requested have been applied as follows:

{{.StateMachine.Description}}

{{ end }}
{{ if eq .Meta.CurrentStage "initial" }}
## Core Design Choices

> {{.Meta.Description}}

The following are the high-level design choices for the "{{.Meta.Name}}" app. These may evolve in subsequent stages.

* *Language:* {{.Meta.Language}}
* *Test Strategy:* {{.Meta.Test}}
* *Framework(s):* {{.Meta.Framework}}
* *Architecture(s):* {{.Meta.Architecture}}
{{ end }}

{{ if gt (len .StateMachine.Questions) 0 }}
### Clarity Requested

devoid has some questions for you:

{{ range $question := .StateMachine.Questions }}
* {{$question}} 
{{ end }}
{{ end }}
    `,
	)
	var b bytes.Buffer
	t.Execute(&b, p)
	return b.String()
}
