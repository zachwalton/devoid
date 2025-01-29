package templates

import (
	"bytes"
	"fmt"
	"text/template"
)

func SystemInitial(projectDirectory string) string {
	return systemPrompt(
		projectDirectory,
		"scaffolding",
		false,
	)
}

func SystemClarify(lastResponse string) string {
	return fmt.Sprintf(`
  Modify the JSON below the ^^^ line to incorporate the changes from the user prompt. For example, if the user requests to add unit tests, the meta -> test field must be updated in the returned JSON.

  Guidelines:
  ---
  - If answers are included in the user's prompt for questions in the "questions" array, remove those questions from the array.
  - Populate the state_machine -> description field with a description of changes made to the original result.
  - If the user requests changing e.g. the app name, update meta -> name appropriately. Be thorough about considering all fields that may require changes based on the user's prompt.
  ---
  ^^^
  %s
  `, lastResponse)
}

func systemPrompt(projectDirectory, next string, final bool) string {
	t, _ := template.New("prompt").Parse(`
    You are about to bootstrap a codebase from scratch as an expert software engineer. Please don't make grand claims about the codebase doing highly complex things (LLMs, databases) unless they are requested explicitly by the user.

    Project Directory: 
    ---
    {{.ProjectDirectory}}
    ---

    Guidelines:
    ---
    - "none" is NEVER a valid field value. To indicate that the field is unset, just use the default value from the schema or "unset" if it's a string.
    - Again, it is never valid for any field to have a value of "none". Do not set field values to "none". Period.
    - Make results pass the common sense test, e.g. "MVC" is not acceptable to suggest as an architecture for a CLI, similarly Django and Flask would not be appropriate frameworks for a CLI.
    - The "prompt" section should factor into field values. e.g. "test" should not be set if the user says they don't want tests
    ---

    Field Descriptions:
    ---
    {{range $key, $value := .FieldDescriptions}}
    - "{{ $key }}" should be evaluated as follows: {{ $value }}
    {{end}}
    ---
    `,
	)
	var b bytes.Buffer
	t.Execute(
		&b,
		&systemTemplate{
			ProjectDirectory: projectDirectory,
			FieldDescriptions: map[string]string{
				"meta.name":                         "Should not ever be empty when 'meta' is part of the provided schema",
				"meta.languages":                    "Usually one language, but can be a comma-delimited list of multiple languages; example would be if the user describes a Python service with a UI. However, you may choose to implement that whole example with Python if it feels appropriate.",
				"meta.description":                  "Should be descriptive but concise, encompassing all major implementation approaches (e.g. testing, frameworks, languages, etc.). If the user has described an app such as a python app with a UI, and you don't choose to use two languages (e.g. python and javascript), explain how the requested app can be created in a single language.",
				"meta.framework":                    `The "meta -> framework" key refers to a project development framework like Django or Rails, not things for specific parts of the codebase like "unittest". Can be a comma-delimited list of multiple frameworks when using multiple languages. Should pass the common sense test, e.g. don't suggest an MVC framework for a CLI but a CLI framework could be good`,
				"meta.architecture":                 `The "meta -> architecture" key refers to things like MVC or SOA. Must pass the common sense test, e.g. don't suggest MVC for a CLI`,
				"state_machine.next":                fmt.Sprintf("Should be the static string '%s'", next),
				"state_machine.final":               fmt.Sprintf("Should be %t", final),
				"state_machine.questions":           `If you ask questions, make sure they are about specific characteristics of the codebase, not things like "Should I proceed?"`,
			},
		},
	)
	return b.String()
}
