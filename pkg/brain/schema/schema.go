package schema

import (
	"encoding/json"
)

var (
	stateMachineRequires = []string{"description", "next", "final", "questions"}
	metaRequires         = []string{"language", "test", "framework", "architecture", "description", "database"}
)

type Schema struct {
	Schema     string            `json:"$schema"`
	Title      string            `json:"title"`
	Type       string            `json:"type"`
	Required   []string          `json:"required"`
	Properties *SchemaProperties `json:"properties"`
}

type SchemaProperties struct {
	Meta         *MetaProperties         `json:"meta"`
	StateMachine *StateMachineProperties `json:"state_machine"`
}

type MetaProperties struct {
	Type        string            `json:"type"`
	Description string            `json:"description"`
	Required    []string          `json:"required"`
	Properties  *MetaPropertyList `json:"properties"`
}

type MetaPropertyList struct {
	Name         *Property `json:"name"`
	Description  *Property `json:"description"`
	Language     *Property `json:"language"`
	Test         *Property `json:"test"`
	Framework    *Property `json:"framework"`
	Database     *Property `json:"database"`
	Architecture *Property `json:"architecture"`
}

type StateMachineProperties struct {
	Type        string                    `json:"type"`
	Description string                    `json:"description"`
	Required    []string                  `json:"required"`
	Properties  *StateMachinePropertyList `json:"properties"`
}

type StateMachinePropertyList struct {
	Description *Property `json:"description"`
	Next        *Property `json:"next"`
	Final       *Property `json:"final"`
	Questions   *Property `json:"questions"`
}

type Property struct {
	Type        string      `json:"type"`
	Description string      `json:"description"`
	Items       *Item       `json:"items,omitempty"`
	Default     interface{} `json:"default"`
}

type Item struct {
	Type string `json:"type"`
}

func schemaDefault() *Schema {
	return &Schema{
		Schema:   "http://json-schema.org/draft-07/schema#",
		Title:    "Devoid",
		Type:     "object",
		Required: []string{"state_machine"},
		Properties: &SchemaProperties{
			StateMachine: &StateMachineProperties{
				Type:        "object",
				Description: "Properties that drive state machine behavior.",
				Required:    stateMachineRequires,
				Properties: &StateMachinePropertyList{
					Description: &Property{Type: "string", Description: "Description of changes made by the model for this inference", Default: ""},
					Next:        &Property{Type: "string", Description: "Next phase to execute with this output as the next stage's input", Default: "scaffolding"},
					Final:       &Property{Type: "boolean", Description: "True if this is the final stage for the project.", Default: true},
					Questions:   &Property{Type: "array", Items: &Item{Type: "string"}, Description: "Put any clarifying questions here if needed. This should only be used to satisfy 'unset' fields. The questions are you asking the user for project clarification, not random stuff like asking about what algorithms to use that the user does not know", Default: []string{}},
				},
			},
		},
	}
}

func metaDefault() *MetaProperties {
	return &MetaProperties{
		Type:        "object",
		Description: "General characteristics for the codebase.",
		Required:    metaRequires,
		Properties: &MetaPropertyList{
			Name:         &Property{Type: "string", Description: "Name of this project. Never use the default for this, you must always choose a name.", Default: "unset"},
			Description:  &Property{Type: "string", Description: "Description of the project. Be concise but thorough in describing the various high-level approaches.", Default: "unset"},
			Language:     &Property{Type: "string", Description: "The language that will be used for the codebase. Can be multiple languages depending on the project. When multiple languages are set, make it a comma-delimited list", Default: "unset"},
			Test:         &Property{Type: "string", Description: "Preferences related to testing (unit, integration, etc.). When adopting more than one testing approach, e.g. both unit and integration, make it a comma-delimited list", Default: "unset"},
			Framework:    &Property{Type: "string", Description: "Preferences related to frameworks (Ruby on Rails, Django, etc.). This is always a project-wide development framework and never something more specific like 'unittest'", Default: "unset"},
			Database:     &Property{Type: "string", Description: "Preferences related to database usage. 'database' encompasses all types of stateful data, so it's acceptable to put both things like 'cassandra' and 'static json' here. Default to simple solutions like flat files unless specifically requested by the user", Default: "unset"},
			Architecture: &Property{Type: "string", Description: "Preferences related to architecture (SOA, MVC, etc.).", Default: "unset"},
		},
	}
}

func (s Schema) JSON() string {
	j, _ := json.Marshal(s)
	return string(j)
}

func SchemaInitial() string {
	s := schemaDefault()
	s.Required = append(s.Required, "meta")
	s.Properties.Meta = metaDefault()
	return s.JSON()
}
