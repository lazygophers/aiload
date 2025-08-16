package channel

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"strings"

	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/lazygophers/aiload"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/anyx"
)

// GeminiError represents an error response from the Gemini API.
type GeminiError struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

// Gemini a channel for gemini models
type Gemini struct {
	channel *aiload.ModelChannel
	client  *resty.Client
}

// NewGemini creates a new gemini channel
func NewGemini(channel *aiload.ModelChannel) *Gemini {
	client := resty.New()
	client.SetQueryParam("key", channel.Token)
	client.SetError(&GeminiError{})

	return &Gemini{
		channel: channel,
		client:  client,
	}
}

// GeminiGetModelListReq represents a request to get a list of models.
type GeminiGetModelListReq struct {
	// PageSize is the number of models to return.
	PageSize int
	// PageToken is the page token to use for pagination.
	PageToken string
}

// GeminiModel represents a single model.
type GeminiModel struct {
	// Name is the name of the model.
	Name string `json:"name"`
	// BaseModelID is the base model ID.
	BaseModelID string `json:"baseModelId"`
	// Version is the version of the model.
	Version string `json:"version"`
	// DisplayName is the display name of the model.
	DisplayName string `json:"displayName"`
	// Description is the description of the model.
	Description string `json:"description"`
	// InputTokenLimit is the input token limit.
	InputTokenLimit int `json:"inputTokenLimit"`
	// OutputTokenLimit is the output token limit.
	OutputTokenLimit int `json:"outputTokenLimit"`
	// SupportedGenerationMethods is a list of supported generation methods.
	SupportedGenerationMethods []string `json:"supportedGenerationMethods"`
	// Thinking is whether the model is thinking.
	Thinking bool `json:"thinking"`
	// Temperature is the temperature of the model.
	Temperature int `json:"temperature"`
	// MaxTemperature is the maximum temperature of the model.
	MaxTemperature int `json:"maxTemperature"`
	// TopP is the top-p value of the model.
	TopP int `json:"topP"`
	// TopK is the top-k value of the model.
	TopK int `json:"topK"`
}

// GeminiGetModelListRsp represents a response from getting a list of models.
type GeminiGetModelListRsp struct {
	// Models is a list of models.
	Models []GeminiModel `json:"models"`
	// NextPageToken is the next page token.
	NextPageToken string `json:"nextPageToken"`
}

// GetModelList gets a list of models.
func (p *Gemini) GetModelList(ctx context.Context, req *GeminiGetModelListReq) (*GeminiGetModelListRsp, error) {
	var rsp GeminiGetModelListRsp
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
		log.Errorf("err:%s", err)
		return nil, err
	}

	if resp.IsError() {
		if geminiErr, ok := resp.Error().(*GeminiError); ok && geminiErr.Error.Message != "" {
			return nil, fmt.Errorf("API error: %s", geminiErr.Error.Message)
		}
		return nil, errors.New(resp.String())
	}

	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		return nil, err
	}

	return &rsp, nil
}

// GeminiGetModelReq represents a request to get a model.
type GeminiGetModelReq struct {
	// Model is the name of the model to get.
	Model string
}

// GeminiGetModelRsp represents a response from getting a model.
type GeminiGetModelRsp struct {
	GeminiModel
}

// GetModel gets a model.
func (p *Gemini) GetModel(ctx context.Context, req *GeminiGetModelReq) (*GeminiGetModelRsp, error) {
	var rsp GeminiGetModelRsp

	resp, err := p.client.R().Get(p.channel.BaseUrl + "/v1/models/" + req.Model)
	if err != nil {
		log.Errorf("err:%s", err)
		return nil, err
	}

	if resp.IsError() {
		if geminiErr, ok := resp.Error().(*GeminiError); ok && geminiErr.Error.Message != "" {
			return nil, fmt.Errorf("API error: %s", geminiErr.Error.Message)
		}
		return nil, errors.New(resp.String())
	}

	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
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
		return nil, errors.New("model name cannot be empty")
	}

	resp, err := p.client.R().SetBody(req).Post(p.channel.BaseUrl + "/v1/models/" + req.Model + ":countTokens")
	if err != nil {
		log.Errorf("err:%s", err)
		return nil, err
	}

	if resp.IsError() {
		if geminiErr, ok := resp.Error().(*GeminiError); ok && geminiErr.Error.Message != "" {
			return nil, fmt.Errorf("API error: %s", geminiErr.Error.Message)
		}
		return nil, errors.New(resp.String())
	}

	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		return nil, err
	}

	return &rsp, nil
}

// EmbedContent embeds a single content.
func (p *Gemini) EmbedContent(ctx context.Context, req *GeminiEmbedContentReq) (*GeminiEmbedContentRsp, error) {
	var rsp GeminiEmbedContentRsp

	resp, err := p.client.R().SetBody(req).Post(p.channel.BaseUrl + "/v1/models/" + req.Model + ":embedContent")
	if err != nil {
		log.Errorf("err:%s", err)
		return nil, err
	}

	if resp.IsError() {
		if geminiErr, ok := resp.Error().(*GeminiError); ok && geminiErr.Error.Message != "" {
			return nil, fmt.Errorf("API error: %s", geminiErr.Error.Message)
		}
		return nil, errors.New(resp.String())
	}

	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		return nil, err
	}

	return &rsp, nil
}

// BatchEmbedContents embeds a batch of contents.
func (p *Gemini) BatchEmbedContents(ctx context.Context, req *GeminiBatchEmbedContentsReq) (*GeminiBatchEmbedContentsRsp, error) {
	var rsp GeminiBatchEmbedContentsRsp

	if len(req.Requests) == 0 {
		return nil, errors.New("requests cannot be empty")
	}

	// 从第一个请求中获取模型名称，因为批量请求通常针对同一个模型
	model := req.Requests[0].Model
	if model == "" {
		return nil, errors.New("model name cannot be empty in the first request")
	}

	resp, err := p.client.R().SetBody(req).Post(p.channel.BaseUrl + "/v1/models/" + model + ":batchEmbedContents")
	if err != nil {
		log.Errorf("err:%s", err)
		return nil, err
	}

	if resp.IsError() {
		if geminiErr, ok := resp.Error().(*GeminiError); ok && geminiErr.Error.Message != "" {
			return nil, fmt.Errorf("API error: %s", geminiErr.Error.Message)
		}
		return nil, errors.New(resp.String())
	}

	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		return nil, err
	}

	return &rsp, nil
}

// GenerateContent generates content from a prompt.
func (p *Gemini) GenerateContent(ctx context.Context, model string, req *GeminiGenerateContentReq) (*GeminiGenerateContentRsp, error) {
	var rsp GeminiGenerateContentRsp

	if model == "" {
		return nil, errors.New("model name cannot be empty")
	}

	resp, err := p.client.R().SetBody(req).Post(p.channel.BaseUrl + "/v1/models/" + model + ":generateContent")
	if err != nil {
		log.Errorf("err:%s", err)
		return nil, err
	}

	if resp.IsError() {
		if geminiErr, ok := resp.Error().(*GeminiError); ok && geminiErr.Error.Message != "" {
			return nil, fmt.Errorf("API error: %s", geminiErr.Error.Message)
		}
		return nil, errors.New(resp.String())
	}

	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		return nil, err
	}

	return &rsp, nil
}

// StreamGenerateContent generates content from a prompt in a streaming fashion.
func (p *Gemini) StreamGenerateContent(ctx context.Context, model string, req *GeminiGenerateContentReq) (<-chan *GeminiGenerateContentRsp, error) {
	if model == "" {
		return nil, errors.New("model name cannot be empty")
	}

	// 创建一个通道用于流式返回响应
	ch := make(chan *GeminiGenerateContentRsp)

	// 使用 go aio.Run 异步执行请求
	go func() {
		defer close(ch) // 确保在函数退出时关闭通道

		// 设置请求，但不执行
		request := p.client.R().SetBody(req)
		request.SetDoNotParseResponse(true) // 必须设置，以便我们可以手动处理响应流

		// 执行请求
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

		// 手动处理响应流
		defer resp.RawBody().Close()
		scanner := bufio.NewScanner(resp.RawBody())
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "data: ") {
				jsonData := strings.TrimPrefix(line, "data: ")
				var rsp GeminiGenerateContentRsp
				if err := json.Unmarshal([]byte(jsonData), &rsp); err != nil {
					log.Errorf("failed to unmarshal stream data: %v", err)
					continue // 继续处理下一行
				}

				// 检查流中是否包含错误
				if rsp.Error != nil && rsp.Error.Error.Message != "" {
					log.Errorf("received error in stream: %s", rsp.Error.Error.Message)
					return // 遇到错误，停止处理并关闭通道
				}

				select {
				case ch <- &rsp:
				case <-ctx.Done(): // 如果上下文被取消，则停止发送
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
