package channel

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/lazygophers/aiload"
	"time"
)

type Ollama struct {
	channel *aiload.ModelChannel
}

func NewOllama(channel *aiload.ModelChannel) *Ollama {
	return &Ollama{
		channel: channel,
	}
}

type OllamaModel struct {
	Name       string             `json:"name"`
	Model      string             `json:"model"`
	ModifiedAt time.Time          `json:"modified_at"`
	Size       int64              `json:"size"`
	Digest     string             `json:"digest"`
	Details    OllamaModelDetails `json:"details"`
}

type OllamaGetLocalModelListRsp struct {
	Models []OllamaModel `json:"models"`
}

func (p *Ollama) GetLocalModelList() (*OllamaGetLocalModelListRsp, error) {
	var rsp OllamaGetLocalModelListRsp
	resp, err := p.GetRequest().Get("http://localhost:11434/api/tags")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("api error: %s", resp.String())
	}
	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		return nil, err
	}
	return &rsp, nil
}

func (p *Ollama) GetRequest() *resty.Request {
	return client.R()
}

// OllamaGenerateReq represents the request for generating a completion.
type OllamaGenerateReq struct {
	Model    string                 `json:"model"`
	Prompt   string                 `json:"prompt"`
	Images   []string               `json:"images,omitempty"`
	Format   string                 `json:"format,omitempty"`
	Options  map[string]interface{} `json:"options,omitempty"`
	System   string                 `json:"system,omitempty"`
	Template string                 `json:"template,omitempty"`
	Context  []int                  `json:"context,omitempty"`
	Stream   bool                   `json:"stream,omitempty"`
	Raw      bool                   `json:"raw,omitempty"`
}

// OllamaGenerateRsp represents the response for a generated completion.
type OllamaGenerateRsp struct {
	Model           string    `json:"model"`
	CreatedAt       time.Time `json:"created_at"`
	Response        string    `json:"response"`
	Done            bool      `json:"done"`
	Context         []int     `json:"context,omitempty"`
	TotalDuration   int64     `json:"total_duration,omitempty"`
	LoadDuration    int64     `json:"load_duration,omitempty"`
	PromptEvalCount int       `json:"prompt_eval_count,omitempty"`
	EvalCount       int       `json:"eval_count,omitempty"`
	EvalDuration    int64     `json:"eval_duration,omitempty"`
}

// OllamaMessage represents a message in a chat completion request.
type OllamaMessage struct {
	Role    string   `json:"role"`
	Content string   `json:"content"`
	Images  []string `json:"images,omitempty"`
}

// OllamaChatReq represents the request for a chat completion.
type OllamaChatReq struct {
	Model    string                 `json:"model"`
	Messages []OllamaMessage        `json:"messages"`
	Format   string                 `json:"format,omitempty"`
	Options  map[string]interface{} `json:"options,omitempty"`
	Stream   bool                   `json:"stream,omitempty"`
}

// OllamaChatRsp represents the response for a chat completion.
type OllamaChatRsp struct {
	Model     string        `json:"model"`
	CreatedAt time.Time     `json:"created_at"`
	Message   OllamaMessage `json:"message"`
	Done      bool          `json:"done"`
}

// OllamaCreateModelReq represents the request to create a model.
type OllamaCreateModelReq struct {
	Name      string `json:"name"`
	Modelfile string `json:"modelfile"`
	Stream    bool   `json:"stream,omitempty"`
}

// OllamaCreateModelRsp represents the response for creating a model.
type OllamaCreateModelRsp struct {
	Status string `json:"status"`
}

// OllamaShowModelReq represents the request to show model information.
type OllamaShowModelReq struct {
	Name string `json:"name"`
}

// OllamaShowModelRsp represents the response with model information.
type OllamaShowModelRsp struct {
	License    string `json:"license"`
	Modelfile  string `json:"modelfile"`
	Parameters string `json:"parameters"`
	Template   string `json:"template"`
}

// OllamaCopyModelReq represents the request to copy a model.
type OllamaCopyModelReq struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

// OllamaDeleteModelReq represents the request to delete a model.
type OllamaDeleteModelReq struct {
	Name string `json:"name"`
}

// OllamaPullModelReq represents the request to pull a model.
type OllamaPullModelReq struct {
	Name     string `json:"name"`
	Insecure bool   `json:"insecure,omitempty"`
	Stream   bool   `json:"stream,omitempty"`
}

// OllamaPullModelRsp represents the response for pulling a model.
type OllamaPullModelRsp struct {
	Status    string `json:"status"`
	Digest    string `json:"digest,omitempty"`
	Total     int64  `json:"total,omitempty"`
	Completed int64  `json:"completed,omitempty"`
}

// OllamaPushModelReq represents the request to push a model.
type OllamaPushModelReq struct {
	Name     string `json:"name"`
	Insecure bool   `json:"insecure,omitempty"`
	Stream   bool   `json:"stream,omitempty"`
}

// OllamaPushModelRsp represents the response for pushing a model.
type OllamaPushModelRsp struct {
	Status    string `json:"status"`
	Digest    string `json:"digest,omitempty"`
	Total     int64  `json:"total,omitempty"`
	Completed int64  `json:"completed,omitempty"`
}

// OllamaEmbeddingsReq represents the request for generating embeddings.
type OllamaEmbeddingsReq struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// OllamaEmbeddingsRsp represents the response with the generated embeddings.
type OllamaEmbeddingsRsp struct {
	Embedding []float64 `json:"embedding"`
}

// OllamaRunningModel represents information about a running model.
type OllamaRunningModel struct {
	Name      string             `json:"name"`
	Model     string             `json:"model"`
	Size      int64              `json:"size"`
	Digest    string             `json:"digest"`
	Details   OllamaModelDetails `json:"details"`
	ExpiresAt time.Time          `json:"expires_at"`
	SizeVRAM  int64              `json:"size_vram"`
}

// OllamaModelDetails provides detailed information about a model.
type OllamaModelDetails struct {
	ParentModel       string   `json:"parent_model"`
	Format            string   `json:"format"`
	Family            string   `json:"family"`
	Families          []string `json:"families"`
	ParameterSize     string   `json:"parameter_size"`
	QuantizationLevel string   `json:"quantization_level"`
}

// OllamaGetRunningModelListRsp represents the response for listing running models.
type OllamaGetRunningModelListRsp struct {
	Models []OllamaRunningModel `json:"models"`
}

// OllamaVersionRsp represents the response for the version request.
type OllamaVersionRsp struct {
	Version string `json:"version"`
}

// GenerateCompletion sends a request to generate a completion.
func (p *Ollama) GenerateCompletion(req *OllamaGenerateReq) (*OllamaGenerateRsp, error) {
	var rsp OllamaGenerateRsp
	resp, err := p.GetRequest().SetBody(req).Post("http://localhost:11434/api/generate")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("api error: %s", resp.String())
	}
	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		return nil, err
	}
	return &rsp, nil
}

// CreateChatCompletion sends a request to create a chat completion.
func (p *Ollama) CreateChatCompletion(req *OllamaChatReq) (*OllamaChatRsp, error) {
	var rsp OllamaChatRsp
	resp, err := p.GetRequest().SetBody(req).Post("http://localhost:11434/api/chat")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("api error: %s", resp.String())
	}
	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		return nil, err
	}
	return &rsp, nil
}

// CreateModel sends a request to create a new model.
func (p *Ollama) CreateModel(req *OllamaCreateModelReq) (*OllamaCreateModelRsp, error) {
	var rsp OllamaCreateModelRsp
	resp, err := p.GetRequest().SetBody(req).Post("http://localhost:11434/api/create")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("api error: %s", resp.String())
	}
	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		return nil, err
	}
	return &rsp, nil
}

// ShowModelInfo sends a request to get information about a model.
func (p *Ollama) ShowModelInfo(req *OllamaShowModelReq) (*OllamaShowModelRsp, error) {
	var rsp OllamaShowModelRsp
	resp, err := p.GetRequest().SetBody(req).Post("http://localhost:11434/api/show")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("api error: %s", resp.String())
	}
	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		return nil, err
	}
	return &rsp, nil
}

// CopyModel sends a request to copy a model.
func (p *Ollama) CopyModel(req *OllamaCopyModelReq) error {
	resp, err := p.GetRequest().SetBody(req).Post("http://localhost:11434/api/copy")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("api error: %s", resp.String())
	}
	return nil
}

// DeleteModel sends a request to delete a model.
func (p *Ollama) DeleteModel(req *OllamaDeleteModelReq) error {
	resp, err := p.GetRequest().SetBody(req).Delete("http://localhost:11434/api/delete")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("api error: %s", resp.String())
	}
	return nil
}

// PullModel sends a request to pull a model from the registry.
func (p *Ollama) PullModel(req *OllamaPullModelReq) (*OllamaPullModelRsp, error) {
	var rsp OllamaPullModelRsp
	resp, err := p.GetRequest().SetBody(req).Post("http://localhost:11434/api/pull")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("api error: %s", resp.String())
	}
	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		return nil, err
	}
	return &rsp, nil
}

// PushModel sends a request to push a model to the registry.
func (p *Ollama) PushModel(req *OllamaPushModelReq) (*OllamaPushModelRsp, error) {
	var rsp OllamaPushModelRsp
	resp, err := p.GetRequest().SetBody(req).Post("http://localhost:11434/api/push")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("api error: %s", resp.String())
	}
	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		return nil, err
	}
	return &rsp, nil
}

// GenerateEmbeddings sends a request to generate embeddings for a prompt.
func (p *Ollama) GenerateEmbeddings(req *OllamaEmbeddingsReq) (*OllamaEmbeddingsRsp, error) {
	var rsp OllamaEmbeddingsRsp
	resp, err := p.GetRequest().SetBody(req).Post("http://localhost:11434/api/embeddings")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("api error: %s", resp.String())
	}
	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		return nil, err
	}
	return &rsp, nil
}

// ListRunningModels sends a request to list currently running models.
func (p *Ollama) ListRunningModels() (*OllamaGetRunningModelListRsp, error) {
	var rsp OllamaGetRunningModelListRsp
	resp, err := p.GetRequest().Get("http://localhost:11434/api/ps")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("api error: %s", resp.String())
	}
	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		return nil, err
	}
	return &rsp, nil
}

// GetVersion sends a request to get the version of Ollama.
func (p *Ollama) GetVersion() (*OllamaVersionRsp, error) {
	var rsp OllamaVersionRsp
	resp, err := p.GetRequest().Get("http://localhost:11434/api/version")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("api error: %s", resp.String())
	}
	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		return nil, err
	}
	return &rsp, nil
}

// CheckBlobExists sends a request to check if a blob exists.
func (p *Ollama) CheckBlobExists(digest string) (bool, error) {
	resp, err := p.GetRequest().Head("http://localhost:11434/api/blobs/" + digest)
	if err != nil {
		return false, err
	}
	if resp.IsError() {
		if resp.StatusCode() == 404 {
			return false, nil
		}
		return false, fmt.Errorf("api error: %s", resp.String())
	}
	return resp.StatusCode() == 200, nil
}