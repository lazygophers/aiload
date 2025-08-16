package channel

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/go-resty/resty/v2"
	"github.com/lazygophers/aiload"
	"github.com/lazygophers/log"
)

// OpenAI 结构体封装了与 OpenAI API 交互所需的状态和方法。
type OpenAI struct {
	channel *aiload.ModelChannel
}

// NewOpenAI 创建并返回一个 OpenAI API 客户端实例。
// 如果 channel 中未指定 BaseUrl，则会使用默认的 OpenAI API 地址。
func NewOpenAI(channel *aiload.ModelChannel) *OpenAI {
	if channel.BaseUrl == "" {
		channel.BaseUrl = "https://api.openai.com"
	}
	return &OpenAI{
		channel: channel,
	}
}

// ModelListResponse 定义了模型列表 API 的响应结构。
type ModelListResponse struct {
	Object string      `json:"object"`
	Data   []ModelData `json:"data"`
}

// ModelData 定义了单个模型的数据结构。
type ModelData struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	OwnedBy string `json:"owned_by"`
}

// GetModelList 获取可用的模型列表。
func (o *OpenAI) GetModelList() (*ModelListResponse, error) {
	var rsp ModelListResponse
	resp, err := o.GetRequest().Get(o.channel.BaseUrl + "/v1/models")
	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}
	if resp.IsError() {
		err = parseError(resp)
		log.Errorf("error: %s", err)
		return nil, err
	}
	err = json.Unmarshal(resp.Body(), &rsp)
	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// GetRequest 创建并配置一个 resty 请求对象。
// 它会为请求设置认证令牌。
func (o *OpenAI) GetRequest() *resty.Request {
	// NOTE: The baseURL is prepended to the URL in each request method
	// because resty doesn't allow setting BaseURL on a per-request basis.
	return client.R().
		SetAuthToken(o.channel.Token)
}

// ChatCompletionRequest 定义了聊天补全 API 的请求结构。
type ChatCompletionRequest struct {
	Model            string    `json:"model"`
	Messages         []Message `json:"messages"`
	MaxTokens        int       `json:"max_tokens,omitempty"`
	Temperature      float64   `json:"temperature,omitempty"`
	TopP             float64   `json:"top_p,omitempty"`
	N                int       `json:"n,omitempty"`
	Stream           bool      `json:"stream,omitempty"`
	Stop             []string  `json:"stop,omitempty"`
	PresencePenalty  float64   `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64   `json:"frequency_penalty,omitempty"`
	User             string    `json:"user,omitempty"`
}

// Message 定义了聊天消息的结构。
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionResponse 定义了聊天补全 API 的响应结构。
type ChatCompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice 定义了返回的单个补全选项。
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage 定义了 API 调用中的 token 使用情况。
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// CreateChatCompletion 创建一个聊天补全请求。
func (o *OpenAI) CreateChatCompletion(req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	var rsp ChatCompletionResponse
	resp, err := o.GetRequest().
		SetBody(req).
		Post(o.channel.BaseUrl + "/v1/chat/completions")
	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}
	if resp.IsError() {
		err = parseError(resp)
		log.Errorf("error: %s", err)
		return nil, err
	}
	err = json.Unmarshal(resp.Body(), &rsp)
	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// EmbeddingRequest 定义了嵌入向量 API 的请求结构。
type EmbeddingRequest struct {
	Input          []string `json:"input"`
	Model          string   `json:"model"`
	EncodingFormat string   `json:"encoding_format,omitempty"`
	User           string   `json:"user,omitempty"`
}

// EmbeddingResponse 定义了嵌入向量 API 的响应结构。
type EmbeddingResponse struct {
	Object string      `json:"object"`
	Data   []Embedding `json:"data"`
	Model  string      `json:"model"`
	Usage  Usage       `json:"usage"`
}

// Embedding 定义了单个嵌入向量的数据结构。
type Embedding struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

// CreateEmbedding 创建一个嵌入向量请求。
func (o *OpenAI) CreateEmbedding(req *EmbeddingRequest) (*EmbeddingResponse, error) {
	var rsp EmbeddingResponse
	resp, err := o.GetRequest().
		SetBody(req).
		Post(o.channel.BaseUrl + "/v1/embeddings")
	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}
	if resp.IsError() {
		err = parseError(resp)
		log.Errorf("error: %s", err)
		return nil, err
	}
	err = json.Unmarshal(resp.Body(), &rsp)
	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// ImageGenerationRequest 定义了图像生成 API 的请求结构。
type ImageGenerationRequest struct {
	Prompt         string `json:"prompt"`
	Model          string `json:"model,omitempty"`
	N              int    `json:"n,omitempty"`
	Quality        string `json:"quality,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
	Size           string `json:"size,omitempty"`
	Style          string `json:"style,omitempty"`
	User           string `json:"user,omitempty"`
}

// ImageGenerationResponse 定义了图像生成 API 的响应结构。
type ImageGenerationResponse struct {
	Created int64               `json:"created"`
	Data    []ImageResponseData `json:"data"`
}

// ImageResponseData 定义了返回的单个图像数据。
type ImageResponseData struct {
	URL           string `json:"url"`
	B64JSON       string `json:"b64_json"`
	RevisedPrompt string `json:"revised_prompt"`
}

// CreateImage 创建一个图像生成请求。
func (o *OpenAI) CreateImage(req *ImageGenerationRequest) (*ImageGenerationResponse, error) {
	var rsp ImageGenerationResponse
	resp, err := o.GetRequest().
		SetBody(req).
		Post(o.channel.BaseUrl + "/v1/images/generations")
	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}
	if resp.IsError() {
		err = parseError(resp)
		log.Errorf("error: %s", err)
		return nil, err
	}
	err = json.Unmarshal(resp.Body(), &rsp)
	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// AudioTranscriptionRequest 定义了音频转录 API 的请求结构。
type AudioTranscriptionRequest struct {
	File           io.Reader `json:"file"`
	FileName       string    `json:"file_name"`
	Model          string    `json:"model"`
	Language       string    `json:"language,omitempty"`
	Prompt         string    `json:"prompt,omitempty"`
	ResponseFormat string    `json:"response_format,omitempty"`
	Temperature    float64   `json:"temperature,omitempty"`
}

// AudioTranscriptionResponse 定义了音频转录 API 的响应结构。
type AudioTranscriptionResponse struct {
	Text string `json:"text"`
}

// CreateAudioTranscription 创建一个音频转录请求。
func (o *OpenAI) CreateAudioTranscription(req *AudioTranscriptionRequest) (*AudioTranscriptionResponse, error) {
	if req.File == nil {
		err := errors.New("file reader is nil")
		log.Errorf("error: %s", err)
		return nil, err
	}
	var rsp AudioTranscriptionResponse

	formData := map[string]string{
		"model": req.Model,
	}
	if req.Language != "" {
		formData["language"] = req.Language
	}
	if req.Prompt != "" {
		formData["prompt"] = req.Prompt
	}
	if req.ResponseFormat != "" {
		formData["response_format"] = req.ResponseFormat
	}

	resp, err := o.GetRequest().
		SetFileReader("file", req.FileName, req.File).
		SetFormData(formData).
		Post(o.channel.BaseUrl + "/v1/audio/transcriptions")
	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}
	if resp.IsError() {
		err = parseError(resp)
		log.Errorf("error: %s", err)
		return nil, err
	}
	err = json.Unmarshal(resp.Body(), &rsp)
	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// APIError 定义了 OpenAI API 返回的标准错误结构。
type APIError struct {
	StatusCode int    `json:"status_code"`
	Type       string `json:"type"`
	Message    string `json:"message"`
	Code       string `json:"code"`
}

// Error 实现了 error 接口，返回格式化的错误信息。
func (e *APIError) Error() string {
	return fmt.Sprintf("api error: status_code: %d, type: %s, message: %s, code: %s", e.StatusCode, e.Type, e.Message, e.Code)
}

// ErrorResponse 定义了包含 APIError 的响应结构。
type ErrorResponse struct {
	Error APIError `json:"error"`
}

// parseError 从 resty.Response 中解析并返回一个 APIError。
// 它首先尝试将响应体 unmarshal 为 ErrorResponse 结构体，
// 如果失败，则将整个响应体作为错误信息。
func parseError(resp *resty.Response) error {
	var errRsp ErrorResponse
	err := json.Unmarshal(resp.Body(), &errRsp)
	if err == nil {
		errRsp.Error.StatusCode = resp.StatusCode()
		apiErr := &errRsp.Error
		log.Errorf("error: %s", apiErr)
		return apiErr
	}

	newErr := errors.New(resp.String())
	log.Errorf("error: %s", newErr)
	return newErr
}
