package channel

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/lazygophers/aiload"
	"github.com/lazygophers/log"
)

// OpenAICompatible 实现了与 OpenAI API 兼容的渠道交互。
// 它封装了所有与特定模型渠道相关的 HTTP 请求和响应处理。
type OpenAICompatible struct {
	channel *aiload.ModelChannel
}

// NewOpenAICompatible 创建一个 OpenAICompatible 渠道的实例。
// 如果渠道配置中未指定 BaseUrl，则会使用默认的 OpenAI API 地址。
func NewOpenAICompatible(channel *aiload.ModelChannel) *OpenAICompatible {
	if channel.BaseUrl == "" {
		channel.BaseUrl = "https://api.openai.com"
	}
	return &OpenAICompatible{
		channel: channel,
	}
}

// OpenAICompatibleGetModelListRsp 定义了获取模型列表 API 的响应结构。
type OpenAICompatibleGetModelListRsp struct {
	Object string `json:"object"`
	Data   []struct {
		Id      string `json:"id"`
		Object  string `json:"object"`
		Created int    `json:"created"`
		OwnedBy string `json:"owned_by"`
	} `json:"data"`
}

// GetModelList 从渠道获取可用模型列表。
func (p *OpenAICompatible) GetModelList() (*OpenAICompatibleGetModelListRsp, error) {
	var rsp OpenAICompatibleGetModelListRsp
	resp, err := p.GetRequest().Get(p.channel.BaseUrl + "/v1/models")
	if err != nil {
		log.Errorf("failed to get model list, err: %s", err)
		return nil, err
	}
	if resp.IsError() {
		err = errors.New(resp.String())
		log.Errorf("get model list API returned error, err: %s", err)
		return nil, err
	}
	err = json.Unmarshal(resp.Body(), &rsp)
	if err != nil {
		log.Errorf("failed to parse model list response, err: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// =================================================================================
// ================================= Audio API =====================================
// =================================================================================

// OpenAICompatibleCreateSpeechReq 定义了创建语音合成请求的结构。
type OpenAICompatibleCreateSpeechReq struct {
	Model          string  `json:"model"`
	Input          string  `json:"input"`
	Voice          string  `json:"voice"`
	ResponseFormat string  `json:"response_format,omitempty"`
	Speed          float64 `json:"speed,omitempty"`
}

// CreateSpeech 通过文本输入生成音频。
func (p *OpenAICompatible) CreateSpeech(req *OpenAICompatibleCreateSpeechReq) (*resty.Response, error) {
	resp, err := p.GetRequest().
		SetBody(req).
		Post(p.channel.BaseUrl + "/v1/audio/speech")
	if err != nil {
		log.Errorf("failed to create speech request, err: %s", err)
		return nil, err
	}
	if resp.IsError() {
		err = errors.New(resp.String())
		log.Errorf("create speech API returned error, err: %s", err)
		return nil, err
	}
	return resp, nil
}

// OpenAICompatibleCreateTranscriptionRsp 定义了创建音频转录响应的结构。
type OpenAICompatibleCreateTranscriptionRsp struct {
	Text string `json:"text"`
}

// CreateTranscription 将音频文件转录为指定语言的文本。
func (p *OpenAICompatible) CreateTranscription(file, model, language, prompt, responseFormat string, temperature float64) (*OpenAICompatibleCreateTranscriptionRsp, error) {
	var rsp OpenAICompatibleCreateTranscriptionRsp
	req := p.GetRequest().
		SetFile("file", file).
		SetFormData(map[string]string{
			"model":           model,
			"language":        language,
			"prompt":          prompt,
			"response_format": responseFormat,
			"temperature":     fmt.Sprintf("%f", temperature),
		})

	resp, err := req.Post(p.channel.BaseUrl + "/v1/audio/transcriptions")
	if err != nil {
		log.Errorf("failed to create transcription request, err: %s", err)
		return nil, err
	}
	if resp.IsError() {
		err = errors.New(resp.String())
		log.Errorf("create transcription API returned error, err: %s", err)
		return nil, err
	}
	err = json.Unmarshal(resp.Body(), &rsp)
	if err != nil {
		log.Errorf("failed to parse transcription response, err: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// OpenAICompatibleCreateTranslationRsp 定义了创建音频翻译响应的结构。
type OpenAICompatibleCreateTranslationRsp struct {
	Text string `json:"text"`
}

// CreateTranslation 将音频文件翻译成英文。
func (p *OpenAICompatible) CreateTranslation(file, model, prompt, responseFormat string, temperature float64) (*OpenAICompatibleCreateTranslationRsp, error) {
	var rsp OpenAICompatibleCreateTranslationRsp
	req := p.GetRequest().
		SetFile("file", file).
		SetFormData(map[string]string{
			"model":           model,
			"prompt":          prompt,
			"response_format": responseFormat,
			"temperature":     fmt.Sprintf("%f", temperature),
		})

	resp, err := req.Post(p.channel.BaseUrl + "/v1/audio/translations")
	if err != nil {
		log.Errorf("failed to create translation request, err: %s", err)
		return nil, err
	}
	if resp.IsError() {
		err = errors.New(resp.String())
		log.Errorf("create translation API returned error, err: %s", err)
		return nil, err
	}
	err = json.Unmarshal(resp.Body(), &rsp)
	if err != nil {
		log.Errorf("failed to parse translation response, err: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// GetRequest 返回一个配置了认证信息的 resty.Request 实例。
// 注意：由于 resty 库的限制，BaseURL 需要在每个请求中单独拼接，而不是在客户端级别设置。
func (p *OpenAICompatible) GetRequest() *resty.Request {
	return client.R().
		SetAuthToken(p.channel.Token)
}

// =================================================================================
// ================================== Chat API =====================================
// =================================================================================

// OpenAICompatibleChatMessage 定义了聊天消息的结构。
type OpenAICompatibleChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAICompatibleCreateChatCompletionReq 定义了创建聊天补全请求的结构。
type OpenAICompatibleCreateChatCompletionReq struct {
	Model            string                        `json:"model"`
	Messages         []OpenAICompatibleChatMessage `json:"messages"`
	Temperature      float64                       `json:"temperature,omitempty"`
	TopP             float64                       `json:"top_p,omitempty"`
	N                int                           `json:"n,omitempty"`
	Stream           bool                          `json:"stream,omitempty"`
	Stop             []string                      `json:"stop,omitempty"`
	MaxTokens        int                           `json:"max_tokens,omitempty"`
	PresencePenalty  float64                       `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64                       `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]int                `json:"logit_bias,omitempty"`
	User             string                        `json:"user,omitempty"`
	ResponseFormat   map[string]string             `json:"response_format,omitempty"`
}

// OpenAICompatibleCreateChatCompletionRsp 定义了创建聊天补全响应的结构。
type OpenAICompatibleCreateChatCompletionRsp struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Choices []struct {
		Index        int                         `json:"index"`
		Message      OpenAICompatibleChatMessage `json:"message"`
		FinishReason string                      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// CreateChatCompletion 基于给定的消息创建一个模型生成的补全。
func (p *OpenAICompatible) CreateChatCompletion(req *OpenAICompatibleCreateChatCompletionReq) (*OpenAICompatibleCreateChatCompletionRsp, error) {
	var rsp OpenAICompatibleCreateChatCompletionRsp
	resp, err := p.GetRequest().
		SetBody(req).
		Post(p.channel.BaseUrl + "/v1/chat/completions")
	if err != nil {
		log.Errorf("failed to create chat completion request, err: %s", err)
		return nil, err
	}
	if resp.IsError() {
		err = errors.New(resp.String())
		log.Errorf("create chat completion API returned error, err: %s", err)
		return nil, err
	}
	err = json.Unmarshal(resp.Body(), &rsp)
	if err != nil {
		log.Errorf("failed to parse chat completion response, err: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// =================================================================================
// =============================== Completions API =================================
// =================================================================================

// OpenAICompatibleCreateCompletionReq 定义了创建文本补全请求的结构。
type OpenAICompatibleCreateCompletionReq struct {
	Model            string         `json:"model"`
	Prompt           any            `json:"prompt"` // string or []string
	BestOf           int            `json:"best_of,omitempty"`
	Echo             bool           `json:"echo,omitempty"`
	FrequencyPenalty float64        `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]int `json:"logit_bias,omitempty"`
	LogProbs         int            `json:"logprobs,omitempty"`
	MaxTokens        int            `json:"max_tokens,omitempty"`
	N                int            `json:"n,omitempty"`
	PresencePenalty  float64        `json:"presence_penalty,omitempty"`
	Seed             int            `json:"seed,omitempty"`
	Stop             any            `json:"stop,omitempty"` // string or []string
	Stream           bool           `json:"stream,omitempty"`
	Suffix           string         `json:"suffix,omitempty"`
	Temperature      float64        `json:"temperature,omitempty"`
	User             string         `json:"user,omitempty"`
	TopP             float64        `json:"top_p,omitempty"`
}

// OpenAICompatibleCreateCompletionRsp 定义了创建文本补全响应的结构。
type OpenAICompatibleCreateCompletionRsp struct {
	Id                string `json:"id"`
	Object            string `json:"object"`
	Created           int64  `json:"created"`
	Model             string `json:"model"`
	SystemFingerprint string `json:"system_fingerprint"`
	Choices           []struct {
		Text         string      `json:"text"`
		Index        int         `json:"index"`
		LogProbs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// CreateCompletion 为提供的提示和参数创建补全。
func (p *OpenAICompatible) CreateCompletion(req *OpenAICompatibleCreateCompletionReq) (*OpenAICompatibleCreateCompletionRsp, error) {
	var rsp OpenAICompatibleCreateCompletionRsp
	resp, err := p.GetRequest().
		SetBody(req).
		Post(p.channel.BaseUrl + "/v1/completions")
	if err != nil {
		log.Errorf("failed to create completion request, err: %s", err)
		return nil, err
	}
	if resp.IsError() {
		err = errors.New(resp.String())
		log.Errorf("create completion API returned error, err: %s", err)
		return nil, err
	}
	err = json.Unmarshal(resp.Body(), &rsp)
	if err != nil {
		log.Errorf("failed to parse completion response, err: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// =================================================================================
// =============================== Embeddings API ==================================
// =================================================================================

// OpenAICompatibleCreateEmbeddingReq 定义了创建嵌入请求的结构。
type OpenAICompatibleCreateEmbeddingReq struct {
	Model string `json:"model"`
	Input any    `json:"input"` // string or []string
}

// OpenAICompatibleCreateEmbeddingRsp 定义了创建嵌入响应的结构。
type OpenAICompatibleCreateEmbeddingRsp struct {
	Object string `json:"object"`
	Data   []struct {
		Object    string    `json:"object"`
		Embedding []float64 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		PromptTokens int `json:"prompt_tokens"`
		TotalTokens  int `json:"total_tokens"`
	} `json:"usage"`
}

// CreateEmbedding 为输入文本创建嵌入向量。
func (p *OpenAICompatible) CreateEmbedding(req *OpenAICompatibleCreateEmbeddingReq) (*OpenAICompatibleCreateEmbeddingRsp, error) {
	var rsp OpenAICompatibleCreateEmbeddingRsp
	resp, err := p.GetRequest().
		SetBody(req).
		Post(p.channel.BaseUrl + "/v1/embeddings")
	if err != nil {
		log.Errorf("failed to create embedding request, err: %s", err)
		return nil, err
	}
	if resp.IsError() {
		err = errors.New(resp.String())
		log.Errorf("create embedding API returned error, err: %s", err)
		return nil, err
	}
	err = json.Unmarshal(resp.Body(), &rsp)
	if err != nil {
		log.Errorf("failed to parse embedding response, err: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// =================================================================================
// ============================== Fine-tuning API ==================================
// =================================================================================

// Hyperparameters 定义了微调任务的超参数。
type Hyperparameters struct {
	BatchSize              string `json:"batch_size,omitempty"`
	LearningRateMultiplier string `json:"learning_rate_multiplier,omitempty"`
	NEpochs                string `json:"n_epochs,omitempty"`
}

// OpenAICompatibleCreateFineTuningJobReq 定义了创建微调任务请求的结构。
type OpenAICompatibleCreateFineTuningJobReq struct {
	TrainingFile    string           `json:"training_file"`
	Model           string           `json:"model"`
	Hyperparameters *Hyperparameters `json:"hyperparameters,omitempty"`
	Suffix          string           `json:"suffix,omitempty"`
	ValidationFile  string           `json:"validation_file,omitempty"`
}

// OpenAICompatibleFineTuningJob 定义了微调任务的详细信息。
type OpenAICompatibleFineTuningJob struct {
	Object          string   `json:"object"`
	Id              string   `json:"id"`
	Model           string   `json:"model"`
	CreatedAt       int64    `json:"created_at"`
	FinishedAt      int64    `json:"finished_at"`
	FineTunedModel  string   `json:"fine_tuned_model"`
	OrganizationId  string   `json:"organization_id"`
	ResultFiles     []string `json:"result_files"`
	Status          string   `json:"status"`
	ValidationFile  string   `json:"validation_file"`
	TrainingFile    string   `json:"training_file"`
	Hyperparameters struct {
		NEpochs int `json:"n_epochs"`
	} `json:"hyperparameters"`
	TrainedTokens int `json:"trained_tokens"`
}

// CreateFineTuningJob 创建一个新的微调任务。
func (p *OpenAICompatible) CreateFineTuningJob(req *OpenAICompatibleCreateFineTuningJobReq) (*OpenAICompatibleFineTuningJob, error) {
	var rsp OpenAICompatibleFineTuningJob
	resp, err := p.GetRequest().
		SetBody(req).
		Post(p.channel.BaseUrl + "/v1/fine_tuning/jobs")
	if err != nil {
		log.Errorf("failed to create fine-tuning job request, err: %s", err)
		return nil, err
	}
	if resp.IsError() {
		err = errors.New(resp.String())
		log.Errorf("create fine-tuning job API returned error, err: %s", err)
		return nil, err
	}
	err = json.Unmarshal(resp.Body(), &rsp)
	if err != nil {
		log.Errorf("failed to parse fine-tuning job response, err: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// OpenAICompatibleFineTuningJobEvent 定义了微调任务事件的结构。
type OpenAICompatibleFineTuningJobEvent struct {
	Object    string      `json:"object"`
	Id        string      `json:"id"`
	CreatedAt int64       `json:"created_at"`
	Level     string      `json:"level"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Type      string      `json:"type"`
}

// OpenAICompatibleListFineTuningJobsRsp 定义了列出微调任务响应的结构。
type OpenAICompatibleListFineTuningJobsRsp struct {
	Object  string                               `json:"object"`
	Data    []OpenAICompatibleFineTuningJobEvent `json:"data"`
	HasMore bool                                 `json:"has_more"`
}

// ListFineTuningJobs 列出组织的微调任务。
func (p *OpenAICompatible) ListFineTuningJobs(after string, limit int) (*OpenAICompatibleListFineTuningJobsRsp, error) {
	var rsp OpenAICompatibleListFineTuningJobsRsp
	req := p.GetRequest()
	if after != "" {
		req.SetQueryParam("after", after)
	}
	if limit > 0 {
		req.SetQueryParam("limit", fmt.Sprintf("%d", limit))
	}
	resp, err := req.Get(p.channel.BaseUrl + "/v1/fine_tuning/jobs")
	if err != nil {
		log.Errorf("failed to list fine-tuning jobs, err: %s", err)
		return nil, err
	}
	if resp.IsError() {
		err = errors.New(resp.String())
		log.Errorf("list fine-tuning jobs API returned error, err: %s", err)
		return nil, err
	}
	err = json.Unmarshal(resp.Body(), &rsp)
	if err != nil {
		log.Errorf("failed to parse list fine-tuning jobs response, err: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// RetrieveFineTuningJob 获取有关微调任务的信息。
func (p *OpenAICompatible) RetrieveFineTuningJob(jobId string) (*OpenAICompatibleFineTuningJob, error) {
	var rsp OpenAICompatibleFineTuningJob
	resp, err := p.GetRequest().
		Get(p.channel.BaseUrl + fmt.Sprintf("/v1/fine_tuning/jobs/%s", jobId))
	if err != nil {
		log.Errorf("failed to retrieve fine-tuning job, err: %s", err)
		return nil, err
	}
	if resp.IsError() {
		err = errors.New(resp.String())
		log.Errorf("retrieve fine-tuning job API returned error, err: %s", err)
		return nil, err
	}
	err = json.Unmarshal(resp.Body(), &rsp)
	if err != nil {
		log.Errorf("failed to parse retrieve fine-tuning job response, err: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// CancelFineTuningJob 立即取消一个微调任务。
func (p *OpenAICompatible) CancelFineTuningJob(jobId string) (*OpenAICompatibleFineTuningJob, error) {
	var rsp OpenAICompatibleFineTuningJob
	resp, err := p.GetRequest().
		Post(p.channel.BaseUrl + fmt.Sprintf("/v1/fine_tuning/jobs/%s/cancel", jobId))
	if err != nil {
		log.Errorf("failed to cancel fine-tuning job, err: %s", err)
		return nil, err
	}
	if resp.IsError() {
		err = errors.New(resp.String())
		log.Errorf("cancel fine-tuning job API returned error, err: %s", err)
		return nil, err
	}
	err = json.Unmarshal(resp.Body(), &rsp)
	if err != nil {
		log.Errorf("failed to parse cancel fine-tuning job response, err: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// OpenAICompatibleListFineTuningEventsRsp 定义了列出微调事件响应的结构。
type OpenAICompatibleListFineTuningEventsRsp struct {
	Object  string                               `json:"object"`
	Data    []OpenAICompatibleFineTuningJobEvent `json:"data"`
	HasMore bool                                 `json:"has_more"`
}

// ListFineTuningEvents 获取微调任务的精细事件。
func (p *OpenAICompatible) ListFineTuningEvents(jobId, after string, limit int) (*OpenAICompatibleListFineTuningEventsRsp, error) {
	var rsp OpenAICompatibleListFineTuningEventsRsp
	req := p.GetRequest()
	if after != "" {
		req.SetQueryParam("after", after)
	}
	if limit > 0 {
		req.SetQueryParam("limit", fmt.Sprintf("%d", limit))
	}
	resp, err := req.Get(p.channel.BaseUrl + fmt.Sprintf("/v1/fine_tuning/jobs/%s/events", jobId))
	if err != nil {
		log.Errorf("failed to list fine-tuning events, err: %s", err)
		return nil, err
	}
	if resp.IsError() {
		err = errors.New(resp.String())
		log.Errorf("list fine-tuning events API returned error, err: %s", err)
		return nil, err
	}
	err = json.Unmarshal(resp.Body(), &rsp)
	if err != nil {
		log.Errorf("failed to parse list fine-tuning events response, err: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// =================================================================================
// ================================= Images API ====================================
// =================================================================================

// OpenAICompatibleCreateImageReq 定义了创建图像请求的结构。
type OpenAICompatibleCreateImageReq struct {
	Prompt         string `json:"prompt"`
	Model          string `json:"model,omitempty"`
	N              int    `json:"n,omitempty"`
	Quality        string `json:"quality,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
	Style          string `json:"style,omitempty"`
	User           string `json:"user,omitempty"`
	Size           string `json:"size,omitempty"`
}

// OpenAICompatibleImageRsp 定义了图像 API 响应的通用结构。
type OpenAICompatibleImageRsp struct {
	Created int64 `json:"created"`
	Data    []struct {
		Url     string `json:"url"`
		B64Json string `json:"b64_json"`
	} `json:"data"`
}

// CreateImage 根据文本提示创建图像。
func (p *OpenAICompatible) CreateImage(req *OpenAICompatibleCreateImageReq) (*OpenAICompatibleImageRsp, error) {
	var rsp OpenAICompatibleImageRsp
	resp, err := p.GetRequest().
		SetBody(req).
		Post(p.channel.BaseUrl + "/v1/images/generations")
	if err != nil {
		log.Errorf("failed to create image request, err: %s", err)
		return nil, err
	}
	if resp.IsError() {
		err = errors.New(resp.String())
		log.Errorf("create image API returned error, err: %s", err)
		return nil, err
	}
	err = json.Unmarshal(resp.Body(), &rsp)
	if err != nil {
		log.Errorf("failed to parse create image response, err: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// CreateImageEdit 在给定图像的基础上根据提示创建编辑后的图像。
func (p *OpenAICompatible) CreateImageEdit(image, mask, prompt, n, size, responseFormat, user string) (*OpenAICompatibleImageRsp, error) {
	var rsp OpenAICompatibleImageRsp
	req := p.GetRequest().
		SetFile("image", image).
		SetFormData(map[string]string{
			"prompt":          prompt,
			"n":               n,
			"size":            size,
			"response_format": responseFormat,
			"user":            user,
		})
	if mask != "" {
		req.SetFile("mask", mask)
	}

	resp, err := req.Post(p.channel.BaseUrl + "/v1/images/edits")
	if err != nil {
		log.Errorf("failed to create image edit request, err: %s", err)
		return nil, err
	}
	if resp.IsError() {
		err = errors.New(resp.String())
		log.Errorf("create image edit API returned error, err: %s", err)
		return nil, err
	}
	err = json.Unmarshal(resp.Body(), &rsp)
	if err != nil {
		log.Errorf("failed to parse image edit response, err: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// CreateImageVariation 创建给定图像的变体。
func (p *OpenAICompatible) CreateImageVariation(image, n, size, responseFormat, user string) (*OpenAICompatibleImageRsp, error) {
	var rsp OpenAICompatibleImageRsp
	req := p.GetRequest().
		SetFile("image", image).
		SetFormData(map[string]string{
			"n":               n,
			"size":            size,
			"response_format": responseFormat,
			"user":            user,
		})

	resp, err := req.Post(p.channel.BaseUrl + "/v1/images/variations")
	if err != nil {
		log.Errorf("failed to create image variation request, err: %s", err)
		return nil, err
	}
	if resp.IsError() {
		err = errors.New(resp.String())
		log.Errorf("create image variation API returned error, err: %s", err)
		return nil, err
	}
	err = json.Unmarshal(resp.Body(), &rsp)
	if err != nil {
		log.Errorf("failed to parse image variation response, err: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// =================================================================================
// ================================= Models API ====================================
// =================================================================================

// OpenAICompatibleModel 定义了模型信息的结构。
type OpenAICompatibleModel struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

// RetrieveModel 检索并返回一个模型实例的详细信息。
func (p *OpenAICompatible) RetrieveModel(modelId string) (*OpenAICompatibleModel, error) {
	var rsp OpenAICompatibleModel
	resp, err := p.GetRequest().Get(p.channel.BaseUrl + fmt.Sprintf("/v1/models/%s", modelId))
	if err != nil {
		log.Errorf("failed to retrieve model request, err: %s", err)
		return nil, err
	}
	if resp.IsError() {
		err = errors.New(resp.String())
		log.Errorf("retrieve model API returned error, err: %s", err)
		return nil, err
	}
	err = json.Unmarshal(resp.Body(), &rsp)
	if err != nil {
		log.Errorf("failed to parse retrieve model response, err: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// OpenAICompatibleDeleteModelRsp 定义了删除模型响应的结构。
type OpenAICompatibleDeleteModelRsp struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

// DeleteModel 删除一个微调模型。
func (p *OpenAICompatible) DeleteModel(modelId string) (*OpenAICompatibleDeleteModelRsp, error) {
	var rsp OpenAICompatibleDeleteModelRsp
	resp, err := p.GetRequest().Delete(p.channel.BaseUrl + fmt.Sprintf("/v1/models/%s", modelId))
	if err != nil {
		log.Errorf("failed to delete model request, err: %s", err)
		return nil, err
	}
	if resp.IsError() {
		err = errors.New(resp.String())
		log.Errorf("delete model API returned error, err: %s", err)
		return nil, err
	}
	err = json.Unmarshal(resp.Body(), &rsp)
	if err != nil {
		log.Errorf("failed to parse delete model response, err: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// =================================================================================
// ================================== Files API ====================================
// =================================================================================

// OpenAICompatibleFile 定义了文件对象的结构。
type OpenAICompatibleFile struct {
	Id        string `json:"id"`
	Object    string `json:"object"`
	Bytes     int    `json:"bytes"`
	CreatedAt int64  `json:"created_at"`
	Filename  string `json:"filename"`
	Purpose   string `json:"purpose"`
}

// OpenAICompatibleListFilesRsp 定义了列出文件响应的结构。
type OpenAICompatibleListFilesRsp struct {
	Data   []OpenAICompatibleFile `json:"data"`
	Object string                 `json:"object"`
}

// ListFiles 返回属于用户组织的文件列表。
func (p *OpenAICompatible) ListFiles(purpose string) (*OpenAICompatibleListFilesRsp, error) {
	var rsp OpenAICompatibleListFilesRsp
	req := p.GetRequest()
	if purpose != "" {
		req.SetQueryParam("purpose", purpose)
	}
	resp, err := req.Get(p.channel.BaseUrl + "/v1/files")
	if err != nil {
		log.Errorf("failed to list files request, err: %s", err)
		return nil, err
	}
	if resp.IsError() {
		err = errors.New(resp.String())
		log.Errorf("list files API returned error, err: %s", err)
		return nil, err
	}
	err = json.Unmarshal(resp.Body(), &rsp)
	if err != nil {
		log.Errorf("failed to parse list files response, err: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// UploadFile 上传一个包含文档的文件，以便与“assistants”和“fine-tuning”等功能一起使用。
func (p *OpenAICompatible) UploadFile(file, purpose string) (*OpenAICompatibleFile, error) {
	var rsp OpenAICompatibleFile
	resp, err := p.GetRequest().
		SetFile("file", file).
		SetFormData(map[string]string{
			"purpose": purpose,
		}).
		Post(p.channel.BaseUrl + "/v1/files")
	if err != nil {
		log.Errorf("failed to upload file request, err: %s", err)
		return nil, err
	}
	if resp.IsError() {
		err = errors.New(resp.String())
		log.Errorf("upload file API returned error, err: %s", err)
		return nil, err
	}
	err = json.Unmarshal(resp.Body(), &rsp)
	if err != nil {
		log.Errorf("failed to parse upload file response, err: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// OpenAICompatibleDeleteFileRsp 定义了删除文件响应的结构。
type OpenAICompatibleDeleteFileRsp struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

// DeleteFile 删除一个文件。
func (p *OpenAICompatible) DeleteFile(fileId string) (*OpenAICompatibleDeleteFileRsp, error) {
	var rsp OpenAICompatibleDeleteFileRsp
	resp, err := p.GetRequest().Delete(p.channel.BaseUrl + fmt.Sprintf("/v1/files/%s", fileId))
	if err != nil {
		log.Errorf("failed to delete file request, err: %s", err)
		return nil, err
	}
	if resp.IsError() {
		err = errors.New(resp.String())
		log.Errorf("delete file API returned error, err: %s", err)
		return nil, err
	}
	err = json.Unmarshal(resp.Body(), &rsp)
	if err != nil {
		log.Errorf("failed to parse delete file response, err: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// RetrieveFile 检索并返回一个文件的详细信息。
func (p *OpenAICompatible) RetrieveFile(fileId string) (*OpenAICompatibleFile, error) {
	var rsp OpenAICompatibleFile
	resp, err := p.GetRequest().Get(p.channel.BaseUrl + fmt.Sprintf("/v1/files/%s", fileId))
	if err != nil {
		log.Errorf("failed to retrieve file request, err: %s", err)
		return nil, err
	}
	if resp.IsError() {
		err = errors.New(resp.String())
		log.Errorf("retrieve file API returned error, err: %s", err)
		return nil, err
	}
	err = json.Unmarshal(resp.Body(), &rsp)
	if err != nil {
		log.Errorf("failed to parse retrieve file response, err: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// RetrieveFileContent 返回指定文件的内容。
func (p *OpenAICompatible) RetrieveFileContent(fileId string) (string, error) {
	resp, err := p.GetRequest().Get(p.channel.BaseUrl + fmt.Sprintf("/v1/files/%s/content", fileId))
	if err != nil {
		log.Errorf("failed to retrieve file content, err: %s", err)
		return "", err
	}
	if resp.IsError() {
		err = errors.New(resp.String())
		log.Errorf("retrieve file content API returned error, err: %s", err)
		return "", err
	}
	return resp.String(), nil
}
