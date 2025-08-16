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

type OpenAI struct {
	channel *aiload.ModelChannel
	baseURL string
}

func NewOpenAI(channel *aiload.ModelChannel) *OpenAI {
	baseURL := "https://api.openai.com"
	if channel.BaseUrl != "" {
		baseURL = channel.BaseUrl
	}
	return &OpenAI{
		channel: channel,
		baseURL: baseURL,
	}
}

type ModelListResponse struct {
	Object string      `json:"object"`
	Data   []ModelData `json:"data"`
}

type ModelData struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	OwnedBy string `json:"owned_by"`
}

func (o *OpenAI) GetModelList() (*ModelListResponse, error) {
	var rsp ModelListResponse
	resp, err := o.GetRequest().Get(o.baseURL + "/v1/models")
	if err != nil {
		log.Errorf("err:%s", err)
		return nil, err
	}
	if resp.IsError() {
		err := parseError(resp)
		log.Errorf("err:%s", err)
		return nil, err
	}
	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		log.Errorf("err:%s", err)
		return nil, err
	}
	return &rsp, nil
}

func (o *OpenAI) GetRequest() *resty.Request {
	// NOTE: The baseURL is prepended to the URL in each request method
	// because resty doesn't allow setting BaseURL on a per-request basis.
	return client.R().
		SetAuthToken(o.channel.Token)
}

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

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func (o *OpenAI) CreateChatCompletion(req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	var rsp ChatCompletionResponse
	resp, err := o.GetRequest().
		SetBody(req).
		Post(o.baseURL + "/v1/chat/completions")
	if err != nil {
		log.Errorf("err:%s", err)
		return nil, err
	}
	if resp.IsError() {
		err := parseError(resp)
		log.Errorf("err:%s", err)
		return nil, err
	}
	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		log.Errorf("err:%s", err)
		return nil, err
	}
	return &rsp, nil
}

type EmbeddingRequest struct {
	Input          []string `json:"input"`
	Model          string   `json:"model"`
	EncodingFormat string   `json:"encoding_format,omitempty"`
	User           string   `json:"user,omitempty"`
}

type EmbeddingResponse struct {
	Object string      `json:"object"`
	Data   []Embedding `json:"data"`
	Model  string      `json:"model"`
	Usage  Usage       `json:"usage"`
}

type Embedding struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}

func (o *OpenAI) CreateEmbedding(req *EmbeddingRequest) (*EmbeddingResponse, error) {
	var rsp EmbeddingResponse
	resp, err := o.GetRequest().
		SetBody(req).
		Post(o.baseURL + "/v1/embeddings")
	if err != nil {
		log.Errorf("err:%s", err)
		return nil, err
	}
	if resp.IsError() {
		err := parseError(resp)
		log.Errorf("err:%s", err)
		return nil, err
	}
	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		log.Errorf("err:%s", err)
		return nil, err
	}
	return &rsp, nil
}

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

type ImageGenerationResponse struct {
	Created int64               `json:"created"`
	Data    []ImageResponseData `json:"data"`
}

type ImageResponseData struct {
	URL           string `json:"url"`
	B64JSON       string `json:"b64_json"`
	RevisedPrompt string `json:"revised_prompt"`
}

func (o *OpenAI) CreateImage(req *ImageGenerationRequest) (*ImageGenerationResponse, error) {
	var rsp ImageGenerationResponse
	resp, err := o.GetRequest().
		SetBody(req).
		Post(o.baseURL + "/v1/images/generations")
	if err != nil {
		log.Errorf("err:%s", err)
		return nil, err
	}
	if resp.IsError() {
		err := parseError(resp)
		log.Errorf("err:%s", err)
		return nil, err
	}
	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		log.Errorf("err:%s", err)
		return nil, err
	}
	return &rsp, nil
}

type AudioTranscriptionRequest struct {
	File           io.Reader `json:"file"`
	FileName       string    `json:"file_name"`
	Model          string    `json:"model"`
	Language       string    `json:"language,omitempty"`
	Prompt         string    `json:"prompt,omitempty"`
	ResponseFormat string    `json:"response_format,omitempty"`
	Temperature    float64   `json:"temperature,omitempty"`
}

type AudioTranscriptionResponse struct {
	Text string `json:"text"`
}

func (o *OpenAI) CreateAudioTranscription(req *AudioTranscriptionRequest) (*AudioTranscriptionResponse, error) {
	if req.File == nil {
		err := errors.New("file reader is nil")
		log.Errorf("err:%s", err)
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
		Post(o.baseURL + "/v1/audio/transcriptions")
	if err != nil {
		log.Errorf("err:%s", err)
		return nil, err
	}
	if resp.IsError() {
		err := parseError(resp)
		log.Errorf("err:%s", err)
		return nil, err
	}
	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		log.Errorf("err:%s", err)
		return nil, err
	}
	return &rsp, nil
}

type APIError struct {
	StatusCode int    `json:"status_code"`
	Type       string `json:"type"`
	Message    string `json:"message"`
	Code       string `json:"code"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("api error: status_code: %d, type: %s, message: %s, code: %s", e.StatusCode, e.Type, e.Message, e.Code)
}

type ErrorResponse struct {
	Error APIError `json:"error"`
}

func parseError(resp *resty.Response) error {
	var errRsp ErrorResponse
	if err := json.Unmarshal(resp.Body(), &errRsp); err == nil {
		errRsp.Error.StatusCode = resp.StatusCode()
		err := &errRsp.Error
		log.Errorf("err:%s", err)
		return err
	}
	err := errors.New(resp.String())
	log.Errorf("err:%s", err)
	return err
}
