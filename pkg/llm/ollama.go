package llm

import (
	"context"
	"encoding/json"

	ollama "github.com/ollama/ollama/api"

	"github.com/zachwalton/devoid/pkg/config"
)

type OllamaReasoner struct {
	cfg        *config.Config
	client     *ollama.Client
	responseCh chan Response
}

func (r *OllamaReasoner) Generate(ctx context.Context, prompt, format string, systemTemplate string) error {
	req := &ollama.GenerateRequest{
		Model:  r.cfg.LLM.Model,
		Prompt: prompt,
		System: systemTemplate,
		Options: map[string]interface{}{
			"temperature": r.cfg.LLM.Temperature,
		},
	}
	if format != "" {
		req.Format = json.RawMessage(format)
	}
	return r.client.Generate(
		ctx,
		req,
		r.generateResponseFunc,
	)
}

func (r OllamaReasoner) ResponseCh() <-chan Response {
	return r.responseCh
}

func (r *OllamaReasoner) generateResponseFunc(response ollama.GenerateResponse) error {
	r.responseCh <- Response{
		Response: response.Response,
		Done:     response.Done,
	}
	return nil
}

func NewOllamaReasoner(cfg *config.Config) (*OllamaReasoner, error) {
	client, err := ollama.ClientFromEnvironment()
	return &OllamaReasoner{
		cfg:        cfg,
		client:     client,
		responseCh: make(chan Response),
	}, err
}
