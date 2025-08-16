package channel

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/lazygophers/aiload"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/anyx"
)

// GeminiError 封装了从 Gemini API 返回的错误信息。
type GeminiError struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

// Gemini 实现了与 Google Gemini 模型进行交互的 aiload.ModelChannel。
// 它处理 API 请求的构建、发送以及响应的解析。
type Gemini struct {
	channel *aiload.ModelChannel
	client  *resty.Client
}

// NewGemini 创建并返回一个新的 Gemini channel 实例。
// 它会初始化一个配置了 API token 的 resty 客户端。
func NewGemini(channel *aiload.ModelChannel) *Gemini {
	client := resty.New()
	client.SetQueryParam("key", channel.Token)
	client.SetError(&GeminiError{})
	return &Gemini{
		channel: channel,
		client:  client,
	}
}

// GeminiGetModelListReq 定义了获取模型列表的请求参数。
type GeminiGetModelListReq struct {
	// PageSize 指定每页返回的模型数量。
	PageSize int
	// PageToken 用于分页的页面令牌。
	PageToken string
}

// GeminiModel 代表一个具体的 Gemini 模型及其属性。
type GeminiModel struct {
	// Name 模型的唯一名称。
	Name string `json:"name"`
	// BaseModelID 基础模型的 ID。
	BaseModelID string `json:"baseModelId"`
	// Version 模型的版本号。
	Version string `json:"version"`
	// DisplayName 模型在 UI 中展示的名称。
	DisplayName string `json:"displayName"`
	// Description 模型的详细描述。
	Description string `json:"description"`
	// InputTokenLimit 模型支持的最大输入 token 数。
	InputTokenLimit int `json:"inputTokenLimit"`
	// OutputTokenLimit 模型支持的最大输出 token 数。
	OutputTokenLimit int `json:"outputTokenLimit"`
	// SupportedGenerationMethods 支持的生成方法列表。
	SupportedGenerationMethods []string `json:"supportedGenerationMethods"`
	// Thinking 指示模型是否处于“思考”状态。
	Thinking bool `json:"thinking"`
	// Temperature 模型的温度参数。
	Temperature int `json:"temperature"`
	// MaxTemperature 模型的最大温度参数。
	MaxTemperature int `json:"maxTemperature"`
	// TopP 模型的 Top-p 采样参数。
	TopP int `json:"topP"`
	// TopK 模型的 Top-k 采样参数。
	TopK int `json:"topK"`
}

// GeminiGetModelListRsp 定义了获取模型列表的响应结构。
type GeminiGetModelListRsp struct {
	// Models 返回的模型列表。
	Models []GeminiModel `json:"models"`
	// NextPageToken 用于获取下一页结果的令牌。
	NextPageToken string `json:"nextPageToken"`
}

// GetModelList 从 Gemini API 获取可用的模型列表。
// 它支持通过 PageSize 和 PageToken进行分页。
func (p *Gemini) GetModelList(ctx context.Context, req *GeminiGetModelListReq) (*GeminiGetModelListRsp, error) {
	var rsp GeminiGetModelListRsp
	var err error

	pageSize := anyx.ToString(req.PageSize)
	if req.PageSize <= 0 {
		pageSize = "50"
	}
	if req.PageSize > 1000 {
		pageSize = "1000"
	}
	resp, err := p.client.R().SetQueryParams(map[string]string{
		"pageSize":  pageSize,
		"pageToken": req.PageToken,
	}).Get(p.channel.BaseUrl + "/v1beta/models")
	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}
	if resp.IsError() {
		if geminiErr, ok := resp.Error().(*GeminiError); ok && geminiErr.Error.Message != "" {
			err = fmt.Errorf("API error: %s", geminiErr.Error.Message)
			log.Errorf("error: %s", err)
			return nil, err
		}
		err = errors.New(resp.String())
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

// GeminiGetModelReq 定义了获取单个模型信息的请求。
type GeminiGetModelReq struct {
	// Model 要获取的模型的名称。
	Model string
}

// GeminiGetModelRsp 定义了获取单个模型信息的响应。
type GeminiGetModelRsp struct {
	GeminiModel
}

// GetModel 根据模型名称获取单个模型的详细信息。
func (p *Gemini) GetModel(ctx context.Context, req *GeminiGetModelReq) (*GeminiGetModelRsp, error) {
	var rsp GeminiGetModelRsp
	var err error

	resp, err := p.client.R().Get(p.channel.BaseUrl + "/v1/models/" + req.Model)
	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}
	if resp.IsError() {
		if geminiErr, ok := resp.Error().(*GeminiError); ok && geminiErr.Error.Message != "" {
			err = fmt.Errorf("API error: %s", geminiErr.Error.Message)
			log.Errorf("error: %s", err)
			return nil, err
		}
		err = errors.New(resp.String())
		log.Errorf("error: %s", err)
		return nil, err
	}
	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// GeminiCommonPart represents a common part of a request or response.
type GeminiCommonPart struct {
	Text       string      `json:"text,omitempty"`
	InlineData *Blob       `json:"inlineData,omitempty"`
	FileData   *FileData   `json:"fileData,omitempty"`
	Function   *Function   `json:"function,omitempty"`
	ToolConfig *ToolConfig `json:"toolConfig,omitempty"`
}

// Blob represents binary data.
type Blob struct {
	MimeType string `json:"mimeType"`
	Data     string `json:"data"`
}

// FileData represents a file URI.
type FileData struct {
	MimeType string `json:"mimeType"`
	FileURI  string `json:"fileUri"`
}

// Function represents a function declaration.
type Function struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  any    `json:"parameters"`
}

// ToolConfig represents the configuration for tools.
type ToolConfig struct {
	FunctionCallingConfig *FunctionCallingConfig `json:"functionCallingConfig,omitempty"`
}

// FunctionCallingConfig represents the configuration for function calling.
type FunctionCallingConfig struct {
	Mode            string   `json:"mode,omitempty"`
	AllowedFunction []string `json:"allowedFunction,omitempty"`
}

// GeminiContent represents the content of a message.
type GeminiContent struct {
	Parts GeminiCommonPart `json:"parts"`
	Role  string           `json:"role,omitempty"`
}

// SafetySetting represents a safety setting.
type SafetySetting struct {
	Category  string `json:"category"`
	Threshold string `json:"threshold"`
}

// GenerationConfig represents the generation configuration.
type GenerationConfig struct {
	Temperature     float64  `json:"temperature,omitempty"`
	TopP            float64  `json:"topP,omitempty"`
	TopK            int      `json:"topK,omitempty"`
	CandidateCount  int      `json:"candidateCount,omitempty"`
	MaxOutputTokens int      `json:"maxOutputTokens,omitempty"`
	StopSequences   []string `json:"stopSequences,omitempty"`
}

// GeminiEmbedContentReq represents a request to embed content.
type GeminiEmbedContentReq struct {
	Model                string         `json:"-"`
	Content              *GeminiContent `json:"content"`
	TaskType             string         `json:"taskType,omitempty"`
	Title                string         `json:"title,omitempty"`
	OutputDimensionality int            `json:"outputDimensionality,omitempty"`
}

// ContentEmbedding represents the embedding of content.
type ContentEmbedding struct {
	Values []float32 `json:"values"`
}

// GeminiEmbedContentRsp represents a response from embedding content.
type GeminiEmbedContentRsp struct {
	Embedding ContentEmbedding `json:"embedding"`
}

// GeminiBatchEmbedContentsReq represents a request to batch embed contents.
type GeminiBatchEmbedContentsReq struct {
	Requests []*GeminiEmbedContentReq `json:"requests"`
}

// GeminiBatchEmbedContentsRsp represents a response from batch embedding contents.
type GeminiBatchEmbedContentsRsp struct {
	Embeddings []ContentEmbedding `json:"embeddings"`
}

// GeminiGenerateContentReq represents a request to generate content.
type GeminiGenerateContentReq struct {
	Contents         []*GeminiContent  `json:"contents"`
	SafetySettings   []*SafetySetting  `json:"safetySettings,omitempty"`
	GenerationConfig *GenerationConfig `json:"generationConfig,omitempty"`
	Tools            []*Tool           `json:"tools,omitempty"`
}

// Tool represents a tool that can be used by the model.
type Tool struct {
	FunctionDeclarations []*Function `json:"functionDeclarations,omitempty"`
}

// GeminiGenerateContentRsp represents a response from generating content.
type GeminiGenerateContentRsp struct {
	Candidates     []*Candidate    `json:"candidates"`
	PromptFeedback *PromptFeedback `json:"promptFeedback"`
	Error          *GeminiError    `json:"error,omitempty"`
}

// Candidate represents a candidate generated by the model.
type Candidate struct {
	Content       *GeminiContent  `json:"content"`
	FinishReason  string          `json:"finishReason"`
	Index         int             `json:"index"`
	SafetyRatings []*SafetyRating `json:"safetyRatings"`
	TokenCount    int             `json:"tokenCount"`
}

// SafetyRating represents a safety rating.
type SafetyRating struct {
	Category    string `json:"category"`
	Probability string `json:"probability"`
}

// PromptFeedback represents feedback on a prompt.
type PromptFeedback struct {
	BlockReason   string          `json:"blockReason"`
	SafetyRatings []*SafetyRating `json:"safetyRatings"`
}

// GeminiCountTokensReq represents a request to count tokens.
type GeminiCountTokensReq struct {
	// Model is the name of the model to use for counting tokens.
	Model string `json:"-"`
	// Contents is a list of contents to count tokens for.
	Contents []*GeminiContent `json:"contents"`
}

// GeminiCountTokensRsp represents a response from counting tokens.
type GeminiCountTokensRsp struct {
	// TotalTokens is the total number of tokens.
	TotalTokens int `json:"totalTokens"`
}

// CountTokens counts the number of tokens in a prompt.
func (p *Gemini) CountTokens(ctx context.Context, req *GeminiCountTokensReq) (*GeminiCountTokensRsp, error) {
	var rsp GeminiCountTokensRsp

	if req.Model == "" {
		err := errors.New("model name cannot be empty")
		log.Errorf("error: %s", err)
		return nil, err
	}

	resp, err := p.client.R().SetBody(req).Post(p.channel.BaseUrl + "/v1/models/" + req.Model + ":countTokens")
	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}

	if resp.IsError() {
		if geminiErr, ok := resp.Error().(*GeminiError); ok && geminiErr.Error.Message != "" {
			err := fmt.Errorf("API error: %s", geminiErr.Error.Message)
			log.Errorf("error: %s", err)
			return nil, err
		}
		err := errors.New(resp.String())
		log.Errorf("error: %s", err)
		return nil, err
	}

	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}

	return &rsp, nil
}

// EmbedContent embeds a single content.
func (p *Gemini) EmbedContent(ctx context.Context, req *GeminiEmbedContentReq) (*GeminiEmbedContentRsp, error) {
	var rsp GeminiEmbedContentRsp

	resp, err := p.client.R().SetBody(req).Post(p.channel.BaseUrl + "/v1/models/" + req.Model + ":embedContent")
	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}

	if resp.IsError() {
		if geminiErr, ok := resp.Error().(*GeminiError); ok && geminiErr.Error.Message != "" {
			err := fmt.Errorf("API error: %s", geminiErr.Error.Message)
			log.Errorf("error: %s", err)
			return nil, err
		}
		err := errors.New(resp.String())
		log.Errorf("error: %s", err)
		return nil, err
	}

	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}

	return &rsp, nil
}

// BatchEmbedContents embeds a batch of contents.
func (p *Gemini) BatchEmbedContents(ctx context.Context, req *GeminiBatchEmbedContentsReq) (*GeminiBatchEmbedContentsRsp, error) {
	var rsp GeminiBatchEmbedContentsRsp

	if len(req.Requests) == 0 {
		err := errors.New("requests cannot be empty")
		log.Errorf("error: %s", err)
		return nil, err
	}

	// 从第一个请求中获取模型名称，因为批量请求通常针对同一个模型
	model := req.Requests[0].Model
	if model == "" {
		err := errors.New("model name cannot be empty in the first request")
		log.Errorf("error: %s", err)
		return nil, err
	}

	resp, err := p.client.R().SetBody(req).Post(p.channel.BaseUrl + "/v1/models/" + model + ":batchEmbedContents")
	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}

	if resp.IsError() {
		if geminiErr, ok := resp.Error().(*GeminiError); ok && geminiErr.Error.Message != "" {
			err := fmt.Errorf("API error: %s", geminiErr.Error.Message)
			log.Errorf("error: %s", err)
			return nil, err
		}
		err := errors.New(resp.String())
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

// GenerateContent generates content from a prompt.
func (p *Gemini) GenerateContent(ctx context.Context, model string, req *GeminiGenerateContentReq) (*GeminiGenerateContentRsp, error) {
	var rsp GeminiGenerateContentRsp

	if model == "" {
		err := errors.New("model name cannot be empty")
		log.Errorf("error: %s", err)
		return nil, err
	}

	resp, err := p.client.R().SetBody(req).Post(p.channel.BaseUrl + "/v1/models/" + model + ":generateContent")
	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}

	if resp.IsError() {
		if geminiErr, ok := resp.Error().(*GeminiError); ok && geminiErr.Error.Message != "" {
			err := fmt.Errorf("API error: %s", geminiErr.Error.Message)
			log.Errorf("error: %s", err)
			return nil, err
		}
		err := errors.New(resp.String())
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

// StreamGenerateContent 以流式方式从 prompt 生成内容。
// 此函数是并发安全的，它启动一个 goroutine 来处理流式响应，并通过一个 channel 将结果返回。
//
// 工作流程:
// 1. 输入验证：检查 model 名称是否为空。
// 2. 通道创建：创建一个无缓冲的 channel，用于异步地将 GeminiGenerateContentRsp 推送给调用者。
// 3. 异步请求：启动一个新的 goroutine 来执行 HTTP 请求和流处理。
//   - defer close(ch): 确保在 goroutine 退出时，无论成功还是失败，都会关闭 channel，以防调用者死锁。
//
// 4. 请求构建：
//   - SetDoNotParseResponse(true): 关键设置！这告诉 resty 不要自动解析响应体，因为我们需要手动逐行读取流数据。
//
// 5. 请求执行：发送 POST 请求到 streamGenerateContent 端点。
// 6. 错误处理：
//   - 请求发送失败：记录错误并直接退出 goroutine。
//   - API 返回错误：尝试将响应体解析为 GeminiError。如果成功，记录具体的 API 错误信息；否则，记录原始响应体。然后退出 goroutine。
//
// 7. 流处理：
//   - defer resp.RawBody().Close(): 确保在处理完响应后关闭响应体。
//   - bufio.Scanner: 使用 scanner 高效地逐行读取响应流。
//   - Server-Sent Events (SSE) 解析:
//   - `strings.HasPrefix(line, "data: ")`: 检查每一行是否是 SSE 的数据事件。
//   - `strings.TrimPrefix`: 提取 "data: " 后面的 JSON 数据。
//   - JSON 解析：将提取的 JSON 数据解析到 GeminiGenerateContentRsp 结构体中。
//   - 流内错误检查：检查解析后的响应中是否包含错误字段。如果存在，记录错误并终止流处理。
//
// 8. 数据发送与上下文取消：
//   - `select` 语句：这是一个关键的并发控制模式。
//   - `case ch <- &rsp`: 尝试将成功解析的响应发送到 channel。这是一个阻塞操作，如果调用者没有准备好接收，它会等待。
//   - `case <-ctx.Done()`: 同时监听调用者上下文的取消信号。如果上层操作（例如，用户关闭了请求）被取消，`ctx.Done()` 会被关闭，此 case 会被选中，goroutine 会立即退出，从而避免了 goroutine 泄漏。
//
// 9. 扫描器错误检查：在循环结束后，检查 scanner 是否在读取过程中遇到了 IO 错误。
// 10. 返回值：函数立即返回 channel 和 nil error，调用者可以立即开始从 channel 中读取流式响应。
func (p *Gemini) StreamGenerateContent(ctx context.Context, model string, req *GeminiGenerateContentReq) (<-chan *GeminiGenerateContentRsp, error) {
	if model == "" {
		err := errors.New("model name cannot be empty")
		log.Errorf("error: %s", err)
		return nil, err
	}

	ch := make(chan *GeminiGenerateContentRsp)

	go func() {
		defer close(ch)

		request := p.client.R().SetBody(req)
		request.SetDoNotParseResponse(true)

		resp, err := request.Post(p.channel.BaseUrl + "/v1/models/" + model + ":streamGenerateContent")
		if err != nil {
			log.Errorf("failed to send stream request: %v", err)
			return
		}
		if resp.IsError() {
			var geminiErr GeminiError
			if err := json.Unmarshal(resp.Body(), &geminiErr); err == nil && geminiErr.Error.Message != "" {
				log.Errorf("received error response: API error: %s", geminiErr.Error.Message)
			} else {
				log.Errorf("received error response: %s", resp.String())
			}
			return
		}

		defer resp.RawBody().Close()
		scanner := bufio.NewScanner(resp.RawBody())
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "data: ") {
				jsonData := strings.TrimPrefix(line, "data: ")
				var rsp GeminiGenerateContentRsp
				if err := json.Unmarshal([]byte(jsonData), &rsp); err != nil {
					log.Errorf("failed to unmarshal stream data: %v", err)
					continue
				}

				if rsp.Error != nil && rsp.Error.Error.Message != "" {
					log.Errorf("received error in stream: %s", rsp.Error.Error.Message)
					return
				}

				select {
				case ch <- &rsp:
				case <-ctx.Done():
					return
				}
			}
		}

		if err := scanner.Err(); err != nil {
			log.Errorf("error reading stream: %v", err)
		}
	}()

	return ch, nil
}
