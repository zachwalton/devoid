package brain

import (
	"bytes"
	"text/template"
)

type StagePayload struct {
	Meta         MetaPayload         `json:"meta"`
	StateMachine StateMachinePayload `json:"state_machine"`
	Stages       StagesPayload       `json:"stages"`
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

type StagesPayload struct {
	Scaffolding StageScaffoldingPayload `json:"scaffolding"`
}

type StageScaffoldingPayload struct {
	BootstrapCommands []string `json:"bootstrap_commands"`
	ProjectLayout     []string `json:"project_layout"`
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

## Upcoming Changes

The following changes will be made in the subsequent "{{.StateMachine.Next}}" phase.

{{ if eq .Meta.CurrentStage "initial" }}
{{ if gt (len .Stages.Scaffolding.BootstrapCommands) 0 }}
### Bootstrap Commands

The following commands will be run on your system to set up the local project. If they're not installed, you'll be prompted to install them.

{{ range $command := .Stages.Scaffolding.BootstrapCommands }}
* ` + "`" + `{{$command}}` + "`" + ` 
{{ end }}
{{ end }}

{{ if gt (len .Stages.Scaffolding.ProjectLayout) 0 }}
### Directory Layout

The following directories and files will be created in the ` + "`" + `{{.Meta.ProjectPath}}` + "`" + ` directory. If this step fails, you'll be prompted to correct permissions.

{{ range $path := .Stages.Scaffolding.ProjectLayout }}
* ` + "`" + `{{$path}}` + "`" + ` 
{{ end }}
{{ end }}
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
