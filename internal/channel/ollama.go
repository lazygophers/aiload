package channel

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/lazygophers/aiload"
	"github.com/lazygophers/log"
)

// Ollama 结构体封装了与 Ollama 服务交互所需的状态。
type Ollama struct {
	channel *aiload.ModelChannel
}

// NewOllama 创建一个新的 Ollama 客户端实例。
// 如果 channel 中没有指定 BaseUrl，则使用默认的 "http://localhost:11434"。
func NewOllama(channel *aiload.ModelChannel) *Ollama {
	if channel.BaseUrl == "" {
		channel.BaseUrl = "http://localhost:11434"
	}
	return &Ollama{
		channel: channel,
	}
}

// OllamaModel 代表一个 Ollama 模型的信息。
type OllamaModel struct {
	Name       string             `json:"name"`
	Model      string             `json:"model"`
	ModifiedAt time.Time          `json:"modified_at"`
	Size       int64              `json:"size"`
	Digest     string             `json:"digest"`
	Details    OllamaModelDetails `json:"details"`
}

// OllamaGetLocalModelListRsp 代表获取本地模型列表的响应。
type OllamaGetLocalModelListRsp struct {
	Models []OllamaModel `json:"models"`
}

// GetLocalModelList 获取本地可用的模型列表。
func (p *Ollama) GetLocalModelList() (*OllamaGetLocalModelListRsp, error) {
	var rsp OllamaGetLocalModelListRsp
	var err error

	resp, err := p.GetRequest().Get(p.channel.BaseUrl + "/api/tags")
	if err != nil {
		log.Errorf("failed to request ollama api: %s", err)
		return nil, err
	}

	if resp.IsError() {
		err = fmt.Errorf("ollama api error: %s", resp.String())
		log.Errorf("error detail: %s", err)
		return nil, err
	}

	err = json.Unmarshal(resp.Body(), &rsp)
	if err != nil {
		log.Errorf("failed to unmarshal response json: %s", err)
		return nil, err
	}

	return &rsp, nil
}

// GetRequest 创建并返回一个配置了基本认证和超时的 resty 请求。
func (p *Ollama) GetRequest() *resty.Request {
	return client.R()
}

// OllamaGenerateReq 代表生成补全的请求。
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

// OllamaGenerateRsp 代表生成补全的响应。
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

// OllamaMessage 代表聊天补全请求中的一条消息。
type OllamaMessage struct {
	Role    string   `json:"role"`
	Content string   `json:"content"`
	Images  []string `json:"images,omitempty"`
}

// OllamaChatReq 代表聊天补全的请求。
type OllamaChatReq struct {
	Model    string                 `json:"model"`
	Messages []OllamaMessage        `json:"messages"`
	Format   string                 `json:"format,omitempty"`
	Options  map[string]interface{} `json:"options,omitempty"`
	Stream   bool                   `json:"stream,omitempty"`
}

// OllamaChatRsp 代表聊天补全的响应。
type OllamaChatRsp struct {
	Model     string        `json:"model"`
	CreatedAt time.Time     `json:"created_at"`
	Message   OllamaMessage `json:"message"`
	Done      bool          `json:"done"`
}

// OllamaCreateModelReq 代表创建模型的请求。
type OllamaCreateModelReq struct {
	Name      string `json:"name"`
	Modelfile string `json:"modelfile"`
	Stream    bool   `json:"stream,omitempty"`
}

// OllamaCreateModelRsp 代表创建模型的响应。
type OllamaCreateModelRsp struct {
	Status string `json:"status"`
}

// OllamaShowModelReq 代表显示模型信息的请求。
type OllamaShowModelReq struct {
	Name string `json:"name"`
}

// OllamaShowModelRsp 代表包含模型信息的响应。
type OllamaShowModelRsp struct {
	License    string `json:"license"`
	Modelfile  string `json:"modelfile"`
	Parameters string `json:"parameters"`
	Template   string `json:"template"`
}

// OllamaCopyModelReq 代表复制模型的请求。
type OllamaCopyModelReq struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

// OllamaDeleteModelReq 代表删除模型的请求。
type OllamaDeleteModelReq struct {
	Name string `json:"name"`
}

// OllamaPullModelReq 代表拉取模型的请求。
type OllamaPullModelReq struct {
	Name     string `json:"name"`
	Insecure bool   `json:"insecure,omitempty"`
	Stream   bool   `json:"stream,omitempty"`
}

// OllamaPullModelRsp 代表拉取模型的响应。
type OllamaPullModelRsp struct {
	Status    string `json:"status"`
	Digest    string `json:"digest,omitempty"`
	Total     int64  `json:"total,omitempty"`
	Completed int64  `json:"completed,omitempty"`
}

// OllamaPushModelReq 代表推送模型的请求。
type OllamaPushModelReq struct {
	Name     string `json:"name"`
	Insecure bool   `json:"insecure,omitempty"`
	Stream   bool   `json:"stream,omitempty"`
}

// OllamaPushModelRsp 代表推送模型的响应。
type OllamaPushModelRsp struct {
	Status    string `json:"status"`
	Digest    string `json:"digest,omitempty"`
	Total     int64  `json:"total,omitempty"`
	Completed int64  `json:"completed,omitempty"`
}

// OllamaEmbeddingsReq 代表生成嵌入的请求。
type OllamaEmbeddingsReq struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// OllamaEmbeddingsRsp 代表包含生成嵌入的响应。
type OllamaEmbeddingsRsp struct {
	Embedding []float64 `json:"embedding"`
}

// OllamaRunningModel 代表有关正在运行模型的信息。
type OllamaRunningModel struct {
	Name      string             `json:"name"`
	Model     string             `json:"model"`
	Size      int64              `json:"size"`
	Digest    string             `json:"digest"`
	Details   OllamaModelDetails `json:"details"`
	ExpiresAt time.Time          `json:"expires_at"`
	SizeVRAM  int64              `json:"size_vram"`
}

// OllamaModelDetails 提供有关模型的详细信息。
type OllamaModelDetails struct {
	ParentModel       string   `json:"parent_model"`
	Format            string   `json:"format"`
	Family            string   `json:"family"`
	Families          []string `json:"families"`
	ParameterSize     string   `json:"parameter_size"`
	QuantizationLevel string   `json:"quantization_level"`
}

// OllamaGetRunningModelListRsp 代表列出正在运行模型的响应。
type OllamaGetRunningModelListRsp struct {
	Models []OllamaRunningModel `json:"models"`
}

// OllamaVersionRsp 代表版本请求的响应。
type OllamaVersionRsp struct {
	Version string `json:"version"`
}

// GenerateCompletion 发送生成补全的请求。
func (p *Ollama) GenerateCompletion(req *OllamaGenerateReq) (*OllamaGenerateRsp, error) {
	var rsp OllamaGenerateRsp
	var err error

	resp, err := p.GetRequest().SetBody(req).Post(p.channel.BaseUrl + "/api/generate")
	if err != nil {
		log.Errorf("failed to request ollama api: %s", err)
		return nil, err
	}

	if resp.IsError() {
		err = fmt.Errorf("ollama api error: %s", resp.String())
		log.Errorf("error detail: %s", err)
		return nil, err
	}

	err = json.Unmarshal(resp.Body(), &rsp)
	if err != nil {
		log.Errorf("failed to unmarshal response json: %s", err)
		return nil, err
	}

	return &rsp, nil
}

// CreateChatCompletion 发送创建聊天补全的请求。
func (p *Ollama) CreateChatCompletion(req *OllamaChatReq) (*OllamaChatRsp, error) {
	var rsp OllamaChatRsp
	var err error

	resp, err := p.GetRequest().SetBody(req).Post(p.channel.BaseUrl + "/api/chat")
	if err != nil {
		log.Errorf("failed to request ollama api: %s", err)
		return nil, err
	}

	if resp.IsError() {
		err = fmt.Errorf("ollama api error: %s", resp.String())
		log.Errorf("error detail: %s", err)
		return nil, err
	}

	err = json.Unmarshal(resp.Body(), &rsp)
	if err != nil {
		log.Errorf("failed to unmarshal response json: %s", err)
		return nil, err
	}

	return &rsp, nil
}

// CreateModel 发送创建新模型的请求。
func (p *Ollama) CreateModel(req *OllamaCreateModelReq) (*OllamaCreateModelRsp, error) {
	var rsp OllamaCreateModelRsp
	var err error

	resp, err := p.GetRequest().SetBody(req).Post(p.channel.BaseUrl + "/api/create")
	if err != nil {
		log.Errorf("failed to request ollama api: %s", err)
		return nil, err
	}

	if resp.IsError() {
		err = fmt.Errorf("ollama api error: %s", resp.String())
		log.Errorf("error detail: %s", err)
		return nil, err
	}

	err = json.Unmarshal(resp.Body(), &rsp)
	if err != nil {
		log.Errorf("failed to unmarshal response json: %s", err)
		return nil, err
	}

	return &rsp, nil
}

// ShowModelInfo 发送获取模型信息的请求。
func (p *Ollama) ShowModelInfo(req *OllamaShowModelReq) (*OllamaShowModelRsp, error) {
	var rsp OllamaShowModelRsp
	var err error

	resp, err := p.GetRequest().SetBody(req).Post(p.channel.BaseUrl + "/api/show")
	if err != nil {
		log.Errorf("failed to request ollama api: %s", err)
		return nil, err
	}

	if resp.IsError() {
		err = fmt.Errorf("ollama api error: %s", resp.String())
		log.Errorf("error detail: %s", err)
		return nil, err
	}

	err = json.Unmarshal(resp.Body(), &rsp)
	if err != nil {
		log.Errorf("failed to unmarshal response json: %s", err)
		return nil, err
	}

	return &rsp, nil
}

// CopyModel 发送复制模型的请求。
func (p *Ollama) CopyModel(req *OllamaCopyModelReq) error {
	var err error

	resp, err := p.GetRequest().SetBody(req).Post(p.channel.BaseUrl + "/api/copy")
	if err != nil {
		log.Errorf("failed to request ollama api: %s", err)
		return err
	}

	if resp.IsError() {
		err = fmt.Errorf("ollama api error: %s", resp.String())
		log.Errorf("error detail: %s", err)
		return err
	}

	return nil
}

// DeleteModel 发送删除模型的请求。
func (p *Ollama) DeleteModel(req *OllamaDeleteModelReq) error {
	var err error

	resp, err := p.GetRequest().SetBody(req).Delete(p.channel.BaseUrl + "/api/delete")
	if err != nil {
		log.Errorf("failed to request ollama api: %s", err)
		return err
	}

	if resp.IsError() {
		err = fmt.Errorf("ollama api error: %s", resp.String())
		log.Errorf("error detail: %s", err)
		return err
	}

	return nil
}

// PullModel 发送从注册表拉取模型的请求。
func (p *Ollama) PullModel(req *OllamaPullModelReq) (*OllamaPullModelRsp, error) {
	var rsp OllamaPullModelRsp
	var err error

	resp, err := p.GetRequest().SetBody(req).Post(p.channel.BaseUrl + "/api/pull")
	if err != nil {
		log.Errorf("failed to request ollama api: %s", err)
		return nil, err
	}

	if resp.IsError() {
		err = fmt.Errorf("ollama api error: %s", resp.String())
		log.Errorf("error detail: %s", err)
		return nil, err
	}

	err = json.Unmarshal(resp.Body(), &rsp)
	if err != nil {
		log.Errorf("failed to unmarshal response json: %s", err)
		return nil, err
	}

	return &rsp, nil
}

// PushModel 发送将模型推送到注册表的请求。
func (p *Ollama) PushModel(req *OllamaPushModelReq) (*OllamaPushModelRsp, error) {
	var rsp OllamaPushModelRsp
	var err error

	resp, err := p.GetRequest().SetBody(req).Post(p.channel.BaseUrl + "/api/push")
	if err != nil {
		log.Errorf("failed to request ollama api: %s", err)
		return nil, err
	}

	if resp.IsError() {
		err = fmt.Errorf("ollama api error: %s", resp.String())
		log.Errorf("error detail: %s", err)
		return nil, err
	}

	err = json.Unmarshal(resp.Body(), &rsp)
	if err != nil {
		log.Errorf("failed to unmarshal response json: %s", err)
		return nil, err
	}

	return &rsp, nil
}

// GenerateEmbeddings 发送为提示生成嵌入的请求。
func (p *Ollama) GenerateEmbeddings(req *OllamaEmbeddingsReq) (*OllamaEmbeddingsRsp, error) {
	var rsp OllamaEmbeddingsRsp
	var err error

	resp, err := p.GetRequest().SetBody(req).Post(p.channel.BaseUrl + "/api/embeddings")
	if err != nil {
		log.Errorf("failed to request ollama api: %s", err)
		return nil, err
	}

	if resp.IsError() {
		err = fmt.Errorf("ollama api error: %s", resp.String())
		log.Errorf("error detail: %s", err)
		return nil, err
	}

	err = json.Unmarshal(resp.Body(), &rsp)
	if err != nil {
		log.Errorf("failed to unmarshal response json: %s", err)
		return nil, err
	}

	return &rsp, nil
}

// ListRunningModels 发送列出当前正在运行模型的请求。
func (p *Ollama) ListRunningModels() (*OllamaGetRunningModelListRsp, error) {
	var rsp OllamaGetRunningModelListRsp
	var err error

	resp, err := p.GetRequest().Get(p.channel.BaseUrl + "/api/ps")
	if err != nil {
		log.Errorf("failed to request ollama api: %s", err)
		return nil, err
	}

	if resp.IsError() {
		err = fmt.Errorf("ollama api error: %s", resp.String())
		log.Errorf("error detail: %s", err)
		return nil, err
	}

	err = json.Unmarshal(resp.Body(), &rsp)
	if err != nil {
		log.Errorf("failed to unmarshal response json: %s", err)
		return nil, err
	}

	return &rsp, nil
}

// GetVersion 发送获取 Ollama 版本的请求。
func (p *Ollama) GetVersion() (*OllamaVersionRsp, error) {
	var rsp OllamaVersionRsp
	var err error

	resp, err := p.GetRequest().Get(p.channel.BaseUrl + "/api/version")
	if err != nil {
		log.Errorf("failed to request ollama api: %s", err)
		return nil, err
	}

	if resp.IsError() {
		err = fmt.Errorf("ollama api error: %s", resp.String())
		log.Errorf("error detail: %s", err)
		return nil, err
	}

	err = json.Unmarshal(resp.Body(), &rsp)
	if err != nil {
		log.Errorf("failed to unmarshal response json: %s", err)
		return nil, err
	}

	return &rsp, nil
}

// CheckBlobExists 发送检查 Blob 是否存在的请求。
func (p *Ollama) CheckBlobExists(digest string) (bool, error) {
	var err error

	resp, err := p.GetRequest().Head(p.channel.BaseUrl + "/api/blobs/" + digest)
	if err != nil {
		log.Errorf("failed to request ollama api: %s", err)
		return false, err
	}

	if resp.IsError() {
		if resp.StatusCode() == 404 {
			return false, nil
		}
		err = fmt.Errorf("ollama api error: %s", resp.String())
		log.Errorf("error detail: %s", err)
		return false, err
	}

	return resp.StatusCode() == 200, nil
}
