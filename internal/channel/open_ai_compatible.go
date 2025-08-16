package channel

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/lazygophers/aiload"
)

type OpenAICompatible struct {
	channel *aiload.ModelChannel
	baseURL string
}

func NewOpenAICompatible(channel *aiload.ModelChannel) *OpenAICompatible {
	baseURL := "https://api.OpenAICompatible.com"
	if channel.BaseUrl != "" {
		baseURL = channel.BaseUrl
	}
	return &OpenAICompatible{
		channel: channel,
		baseURL: baseURL,
	}
}

type OpenAICompatibleGetModelListRsp struct {
	Object string `json:"object"`
	Data   []struct {
		Id      string `json:"id"`
		Object  string `json:"object"`
		Created int    `json:"created"`
		OwnedBy string `json:"owned_by"`
	} `json:"data"`
}

func (p *OpenAICompatible) GetModelList() (*OpenAICompatibleGetModelListRsp, error) {
	var rsp OpenAICompatibleGetModelListRsp
	resp, err := p.GetRequest().Get(p.baseURL + "/v1/models")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}
	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		return nil, err
	}
	return &rsp, nil
}

// ================================= Audio API =================================

type OpenAICompatibleCreateSpeechReq struct {
	Model          string  `json:"model"`
	Input          string  `json:"input"`
	Voice          string  `json:"voice"`
	ResponseFormat string  `json:"response_format,omitempty"`
	Speed          float64 `json:"speed,omitempty"`
}

// CreateSpeech generates audio from the input text.
func (p *OpenAICompatible) CreateSpeech(req *OpenAICompatibleCreateSpeechReq) (*resty.Response, error) {
	resp, err := p.GetRequest().
		SetBody(req).
		Post(p.baseURL + "/v1/audio/speech")

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}
	return resp, nil
}

type OpenAICompatibleCreateTranscriptionRsp struct {
	Text string `json:"text"`
}

// CreateTranscription transcribes audio into the input language.
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

	resp, err := req.Post(p.baseURL + "/v1/audio/transcriptions")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}
	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		return nil, err
	}
	return &rsp, nil
}

type OpenAICompatibleCreateTranslationRsp struct {
	Text string `json:"text"`
}

// CreateTranslation translates audio into English.
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

	resp, err := req.Post(p.baseURL + "/v1/audio/translations")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}
	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		return nil, err
	}
	return &rsp, nil
}

func (p *OpenAICompatible) GetRequest() *resty.Request {
	// NOTE: The baseURL is prepended to the URL in each request method
	// because resty doesn't allow setting BaseURL on a per-request basis.
	return client.R().
		SetAuthToken(p.channel.Token)
}

// ================================= Chat API =================================

type OpenAICompatibleChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

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

func (p *OpenAICompatible) CreateChatCompletion(req *OpenAICompatibleCreateChatCompletionReq) (*OpenAICompatibleCreateChatCompletionRsp, error) {
	var rsp OpenAICompatibleCreateChatCompletionRsp
	resp, err := p.GetRequest().
		SetBody(req).
		Post(p.baseURL + "/v1/chat/completions")

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}
	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		return nil, err
	}
	return &rsp, nil
}

// ================================= Completions API =================================

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

func (p *OpenAICompatible) CreateCompletion(req *OpenAICompatibleCreateCompletionReq) (*OpenAICompatibleCreateCompletionRsp, error) {
	var rsp OpenAICompatibleCreateCompletionRsp
	resp, err := p.GetRequest().
		SetBody(req).
		Post(p.baseURL + "/v1/completions")

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}
	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		return nil, err
	}
	return &rsp, nil
}

// ================================= Embeddings API =================================

type OpenAICompatibleCreateEmbeddingReq struct {
	Model string `json:"model"`
	Input any    `json:"input"` // string or []string
}

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

func (p *OpenAICompatible) CreateEmbedding(req *OpenAICompatibleCreateEmbeddingReq) (*OpenAICompatibleCreateEmbeddingRsp, error) {
	var rsp OpenAICompatibleCreateEmbeddingRsp
	resp, err := p.GetRequest().
		SetBody(req).
		Post(p.baseURL + "/v1/embeddings")

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}
	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		return nil, err
	}
	return &rsp, nil
}

// ================================= Fine-tuning API =================================

type Hyperparameters struct {
	BatchSize              string `json:"batch_size,omitempty"`
	LearningRateMultiplier string `json:"learning_rate_multiplier,omitempty"`
	NEpochs                string `json:"n_epochs,omitempty"`
}

type OpenAICompatibleCreateFineTuningJobReq struct {
	TrainingFile    string           `json:"training_file"`
	Model           string           `json:"model"`
	Hyperparameters *Hyperparameters `json:"hyperparameters,omitempty"`
	Suffix          string           `json:"suffix,omitempty"`
	ValidationFile  string           `json:"validation_file,omitempty"`
}

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

func (p *OpenAICompatible) CreateFineTuningJob(req *OpenAICompatibleCreateFineTuningJobReq) (*OpenAICompatibleFineTuningJob, error) {
	var rsp OpenAICompatibleFineTuningJob
	resp, err := p.GetRequest().
		SetBody(req).
		Post(p.baseURL + "/v1/fine_tuning/jobs")

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}
	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		return nil, err
	}
	return &rsp, nil
}

type OpenAICompatibleFineTuningJobEvent struct {
	Object    string      `json:"object"`
	Id        string      `json:"id"`
	CreatedAt int64       `json:"created_at"`
	Level     string      `json:"level"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Type      string      `json:"type"`
}

type OpenAICompatibleListFineTuningJobsRsp struct {
	Object  string                               `json:"object"`
	Data    []OpenAICompatibleFineTuningJobEvent `json:"data"`
	HasMore bool                                 `json:"has_more"`
}

func (p *OpenAICompatible) ListFineTuningJobs(after string, limit int) (*OpenAICompatibleListFineTuningJobsRsp, error) {
	var rsp OpenAICompatibleListFineTuningJobsRsp
	req := p.GetRequest()
	if after != "" {
		req.SetQueryParam("after", after)
	}
	if limit > 0 {
		req.SetQueryParam("limit", fmt.Sprintf("%d", limit))
	}
	resp, err := req.Get(p.baseURL + "/v1/fine_tuning/jobs")

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}
	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		return nil, err
	}
	return &rsp, nil
}

func (p *OpenAICompatible) RetrieveFineTuningJob(jobId string) (*OpenAICompatibleFineTuningJob, error) {
	var rsp OpenAICompatibleFineTuningJob
	resp, err := p.GetRequest().
		Get(p.baseURL + fmt.Sprintf("/v1/fine_tuning/jobs/%s", jobId))

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}
	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		return nil, err
	}
	return &rsp, nil
}

func (p *OpenAICompatible) CancelFineTuningJob(jobId string) (*OpenAICompatibleFineTuningJob, error) {
	var rsp OpenAICompatibleFineTuningJob
	resp, err := p.GetRequest().
		Post(p.baseURL + fmt.Sprintf("/v1/fine_tuning/jobs/%s/cancel", jobId))

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}
	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		return nil, err
	}
	return &rsp, nil
}

type OpenAICompatibleListFineTuningEventsRsp struct {
	Object  string                               `json:"object"`
	Data    []OpenAICompatibleFineTuningJobEvent `json:"data"`
	HasMore bool                                 `json:"has_more"`
}

func (p *OpenAICompatible) ListFineTuningEvents(jobId, after string, limit int) (*OpenAICompatibleListFineTuningEventsRsp, error) {
	var rsp OpenAICompatibleListFineTuningEventsRsp
	req := p.GetRequest()
	if after != "" {
		req.SetQueryParam("after", after)
	}
	if limit > 0 {
		req.SetQueryParam("limit", fmt.Sprintf("%d", limit))
	}
	resp, err := req.Get(p.baseURL + fmt.Sprintf("/v1/fine_tuning/jobs/%s/events", jobId))

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}
	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		return nil, err
	}
	return &rsp, nil
}

// ================================= Images API =================================

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

type OpenAICompatibleImageRsp struct {
	Created int64 `json:"created"`
	Data    []struct {
		Url     string `json:"url"`
		B64Json string `json:"b64_json"`
	} `json:"data"`
}

func (p *OpenAICompatible) CreateImage(req *OpenAICompatibleCreateImageReq) (*OpenAICompatibleImageRsp, error) {
	var rsp OpenAICompatibleImageRsp
	resp, err := p.GetRequest().
		SetBody(req).
		Post(p.baseURL + "/v1/images/generations")

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}
	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		return nil, err
	}
	return &rsp, nil
}

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

	resp, err := req.Post(p.baseURL + "/v1/images/edits")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}
	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		return nil, err
	}
	return &rsp, nil
}

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

	resp, err := req.Post(p.baseURL + "/v1/images/variations")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}
	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		return nil, err
	}
	return &rsp, nil
}

// ================================= Models API =================================

type OpenAICompatibleModel struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	OwnedBy string `json:"owned_by"`
}

func (p *OpenAICompatible) RetrieveModel(modelId string) (*OpenAICompatibleModel, error) {
	var rsp OpenAICompatibleModel
	resp, err := p.GetRequest().Get(p.baseURL + fmt.Sprintf("/v1/models/%s", modelId))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}
	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		return nil, err
	}
	return &rsp, nil
}

type OpenAICompatibleDeleteModelRsp struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

func (p *OpenAICompatible) DeleteModel(modelId string) (*OpenAICompatibleDeleteModelRsp, error) {
	var rsp OpenAICompatibleDeleteModelRsp
	resp, err := p.GetRequest().Delete(p.baseURL + fmt.Sprintf("/v1/models/%s", modelId))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}
	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		return nil, err
	}
	return &rsp, nil
}

// ================================= Files API =================================

type OpenAICompatibleFile struct {
	Id        string `json:"id"`
	Object    string `json:"object"`
	Bytes     int    `json:"bytes"`
	CreatedAt int64  `json:"created_at"`
	Filename  string `json:"filename"`
	Purpose   string `json:"purpose"`
}

type OpenAICompatibleListFilesRsp struct {
	Data   []OpenAICompatibleFile `json:"data"`
	Object string                 `json:"object"`
}

func (p *OpenAICompatible) ListFiles(purpose string) (*OpenAICompatibleListFilesRsp, error) {
	var rsp OpenAICompatibleListFilesRsp
	req := p.GetRequest()
	if purpose != "" {
		req.SetQueryParam("purpose", purpose)
	}
	resp, err := req.Get(p.baseURL + "/v1/files")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}
	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		return nil, err
	}
	return &rsp, nil
}

func (p *OpenAICompatible) UploadFile(file, purpose string) (*OpenAICompatibleFile, error) {
	var rsp OpenAICompatibleFile
	resp, err := p.GetRequest().
		SetFile("file", file).
		SetFormData(map[string]string{
			"purpose": purpose,
		}).
		Post(p.baseURL + "/v1/files")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}
	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		return nil, err
	}
	return &rsp, nil
}

type OpenAICompatibleDeleteFileRsp struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

func (p *OpenAICompatible) DeleteFile(fileId string) (*OpenAICompatibleDeleteFileRsp, error) {
	var rsp OpenAICompatibleDeleteFileRsp
	resp, err := p.GetRequest().Delete(p.baseURL + fmt.Sprintf("/v1/files/%s", fileId))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}
	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		return nil, err
	}
	return &rsp, nil
}

func (p *OpenAICompatible) RetrieveFile(fileId string) (*OpenAICompatibleFile, error) {
	var rsp OpenAICompatibleFile
	resp, err := p.GetRequest().Get(p.baseURL + fmt.Sprintf("/v1/files/%s", fileId))
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(resp.String())
	}
	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		return nil, err
	}
	return &rsp, nil
}

func (p *OpenAICompatible) RetrieveFileContent(fileId string) (string, error) {
	resp, err := p.GetRequest().Get(p.baseURL + fmt.Sprintf("/v1/files/%s/content", fileId))
	if err != nil {
		return "", err
	}
	if resp.IsError() {
		return "", errors.New(resp.String())
	}
	return resp.String(), nil
}
