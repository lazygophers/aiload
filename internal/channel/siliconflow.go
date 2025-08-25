package channel

import (
	"encoding/json"
	"errors"
	"io"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/lazygophers/aiload"
	"github.com/lazygophers/log"
)

var (
	ErrAuthTokenNil = errors.New("auth token is nil")
)

// SiliconFlowErrorRsp 表示来自 SiliconFlow API 的错误响应。
type SiliconFlowErrorRsp struct {
	Err struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Param   string `json:"param"`
		Code    string `json:"code"`
	} `json:"error"`
}

// Error 实现 error 接口，返回错误信息。
func (e *SiliconFlowErrorRsp) Error() string {
	return e.Err.Message
}

// SiliconFlow 是用于与 SiliconFlow API 交互的客户端。
type SiliconFlow struct {
	channel *aiload.ModelChannel
}

// NewSiliconFlow 创建一个新的 SiliconFlow 客户端实例。
func NewSiliconFlow(channel *aiload.ModelChannel) *SiliconFlow {
	if channel.BaseUrl == "" {
		channel.BaseUrl = "https://api.siliconflow.cn"
	}
	return &SiliconFlow{
		channel: channel,
	}
}

// GetRequestRequired 检查请求是否可以被认证。
func (p *SiliconFlow) GetRequestRequired() bool {
	return p.channel.Token != ""
}

// GetRequest 创建一个带有认证信息的新 resty 请求。
func (p *SiliconFlow) GetRequest() *resty.Request {
	return client.R().SetAuthToken(p.channel.Token)
}

// handleSiliconFlowError 检查并处理来自 SiliconFlow 的 API 错误。
func handleSiliconFlowError(resp *resty.Response, err error) error {
	if err != nil {
		log.Errorf("error: %s", err)
		return err
	}
	if !resp.IsError() {
		return nil
	}
	var errRsp SiliconFlowErrorRsp
	e := json.Unmarshal(resp.Body(), &errRsp)
	if e == nil && errRsp.Err.Message != "" {
		err = &errRsp
		log.Errorf("error: %s", err)
		return err
	}
	err = errors.New(resp.String())
	log.Errorf("error: %s", err)
	return err
}

// post 是一个辅助函数，用于发送带有 JSON 主体的 POST 请求。
func (p *SiliconFlow) post(url string, req, rsp interface{}) error {
	resp, err := p.GetRequest().
		SetBody(req).
		SetResult(rsp).
		Post(url)
	return handleSiliconFlowError(resp, err)
}

// SiliconFlowGetInfoListRsp 定义了获取用户信息响应的结构。
type SiliconFlowGetInfoListRsp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  bool   `json:"status"`
	Data    struct {
		Id            string `json:"id"`
		Name          string `json:"name"`
		Image         string `json:"image"`
		Email         string `json:"email"`
		IsAdmin       bool   `json:"isAdmin"`
		Balance       string `json:"balance"`
		Status        string `json:"status"`
		Introduction  string `json:"introduction"`
		Role          string `json:"role"`
		ChargeBalance string `json:"chargeBalance"`
		TotalBalance  string `json:"totalBalance"`
	} `json:"data"`
}

// GetInfoList 获取用户信息。
func (p *SiliconFlow) GetInfoList() (*SiliconFlowGetInfoListRsp, error) {
	if !p.GetRequestRequired() {
		err := ErrAuthTokenNil
		log.Errorf("error: %s", err)
		return nil, err
	}
	var rsp SiliconFlowGetInfoListRsp
	resp, err := p.GetRequest().SetResult(&rsp).Get(p.channel.BaseUrl + "/v1/user/info")
	err = handleSiliconFlowError(resp, err)
	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// SiliconFlowGetModelListReq defines the request for getting model list.
type SiliconFlowGetModelListReq struct {
	Type    string
	SubType string
}

// Model represents a model object returned by the SiliconFlow API.
type Model struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"` // 使用 int64 以兼容 Unix 时间戳
	OwnedBy string `json:"owned_by"`
}

// SiliconFlowGetModelListRsp defines the response for getting model list.
type SiliconFlowGetModelListRsp struct {
	Object string  `json:"object"`
	Data   []Model `json:"data"`
}

// GetModelList 获取模型列表。
func (p *SiliconFlow) GetModelList() (*SiliconFlowGetModelListRsp, error) {
	var rsp SiliconFlowGetModelListRsp
	resp, err := p.GetRequest().SetResult(&rsp).Get(p.channel.BaseUrl + "/v1/models")

	err = handleSiliconFlowError(resp, err)
	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// GetModel 获取指定 ID 的模型实例。
func (p *SiliconFlow) GetModel(modelId string) (*Model, error) {
	var rsp Model
	resp, err := p.GetRequest().
		SetResult(&rsp).
		SetPathParam("model", modelId).
		Get(p.channel.BaseUrl + "/v1/models/{model}")

	err = handleSiliconFlowError(resp, err)
	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// SiliconFlowCreateChatCompletionReq defines the request for creating chat completion.
type SiliconFlowCreateChatCompletionReq struct {
	Model            string                                   `json:"model"`
	Messages         []SiliconFlowCreateChatCompletionMessage `json:"messages"`
	Temperature      float64                                  `json:"temperature,omitempty"`
	TopP             float64                                  `json:"top_p,omitempty"`
	MaxTokens        int                                      `json:"max_tokens,omitempty"`
	Stream           bool                                     `json:"stream,omitempty"`
	Stop             []string                                 `json:"stop,omitempty"`
	PresencePenalty  float64                                  `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64                                  `json:"frequency_penalty,omitempty"`
}

// SiliconFlowCreateChatCompletionMessage defines a message in the chat completion request.
type SiliconFlowCreateChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// SiliconFlowCreateChatCompletionRsp defines the response for creating chat completion.
type SiliconFlowCreateChatCompletionRsp struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// SiliconFlowCreateAudioTranslationReq 定义了创建音频翻译的请求。
type SiliconFlowCreateAudioTranslationReq struct {
	File           io.Reader
	FileName       string
	Model          string
	Prompt         string
	ResponseFormat string
	Temperature    float64
}

// SiliconFlowCreateAudioTranslationRsp 定义了创建音频翻译的响应。
type SiliconFlowCreateAudioTranslationRsp struct {
	Text string `json:"text"`
}

// CreateAudioTranslation 创建音频翻译。
func (p *SiliconFlow) CreateAudioTranslation(req *SiliconFlowCreateAudioTranslationReq) (*SiliconFlowCreateAudioTranslationRsp, error) {
	if req.File == nil {
		err := errors.New("file reader is required")
		log.Errorf("error: %s", err)
		return nil, err
	}

	var rsp SiliconFlowCreateAudioTranslationRsp
	request := p.GetRequest()

	formData := make(map[string]string)
	formData["model"] = req.Model
	if req.Prompt != "" {
		formData["prompt"] = req.Prompt
	}
	if req.ResponseFormat != "" {
		formData["response_format"] = req.ResponseFormat
	}
	if req.Temperature > 0 {
		formData["temperature"] = strconv.FormatFloat(req.Temperature, 'f', -1, 64)
	}

	request.SetFormData(formData)
	request.SetFileReader("file", req.FileName, req.File)

	resp, err := request.
		SetResult(&rsp).
		Post(p.channel.BaseUrl + "/v1/audio/translations")

	err = handleSiliconFlowError(resp, err)
	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}

	return &rsp, nil
}

// SiliconFlowCreateAudioTranscriptionReq 定义了创建音频转录的请求。
type SiliconFlowCreateAudioTranscriptionReq struct {
	File           io.Reader
	FileName       string
	Model          string
	Prompt         string
	ResponseFormat string
	Temperature    float64
}

// SiliconFlowCreateAudioTranscriptionResp 定义了创建音频转录的响应。
type SiliconFlowCreateAudioTranscriptionResp struct {
	Text string `json:"text"`
}

// CreateAudioTranscription 创建音频转录。
func (p *SiliconFlow) CreateAudioTranscription(req *SiliconFlowCreateAudioTranscriptionReq) (*SiliconFlowCreateAudioTranscriptionResp, error) {
	if req.File == nil {
		err := errors.New("file reader is required")
		log.Errorf("error: %s", err)
		return nil, err
	}

	var rsp SiliconFlowCreateAudioTranscriptionResp
	request := p.GetRequest()

	formData := make(map[string]string)
	formData["model"] = req.Model
	if req.Prompt != "" {
		formData["prompt"] = req.Prompt
	}
	if req.ResponseFormat != "" {
		formData["response_format"] = req.ResponseFormat
	}
	if req.Temperature > 0 {
		formData["temperature"] = strconv.FormatFloat(req.Temperature, 'f', -1, 64)
	}

	request.SetFormData(formData)
	request.SetFileReader("file", req.FileName, req.File)

	resp, err := request.
		SetResult(&rsp).
		Post(p.channel.BaseUrl + "/v1/audio/transcriptions")

	err = handleSiliconFlowError(resp, err)
	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}

	return &rsp, nil
}

// CreateChatCompletion creates a chat completion.
func (p *SiliconFlow) CreateChatCompletion(req *SiliconFlowCreateChatCompletionReq) (*SiliconFlowCreateChatCompletionRsp, error) {
	var rsp SiliconFlowCreateChatCompletionRsp
	err := p.post(p.channel.BaseUrl+"/v1/chat/completions", req, &rsp)
	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// SiliconFlowCreateEmbeddingReq defines the request for creating embeddings.
type SiliconFlowCreateEmbeddingReq struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

// SiliconFlowCreateEmbeddingRsp defines the response for creating embeddings.
type SiliconFlowCreateEmbeddingRsp struct {
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

// CreateEmbedding creates embeddings.
func (p *SiliconFlow) CreateEmbedding(req *SiliconFlowCreateEmbeddingReq) (*SiliconFlowCreateEmbeddingRsp, error) {
	var rsp SiliconFlowCreateEmbeddingRsp
	err := p.post(p.channel.BaseUrl+"/v1/embeddings", req, &rsp)
	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// SiliconFlowCreateImageGenerationsReq defines the request for creating image generations.
type SiliconFlowCreateImageGenerationsReq struct {
	Prompt         string `json:"prompt"`
	Model          string `json:"model,omitempty"`
	N              int    `json:"n,omitempty"`
	Size           string `json:"size,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
	User           string `json:"user,omitempty"`
}

// SiliconFlowCreateImageGenerationsRsp defines the response for creating image generations.
type SiliconFlowCreateImageGenerationsRsp struct {
	Created int64 `json:"created"`
	Data    []struct {
		URL     string `json:"url,omitempty"`
		B64JSON string `json:"b64_json,omitempty"`
	} `json:"data"`
}

// CreateImageGenerations creates image generations.
func (p *SiliconFlow) CreateImageGenerations(req *SiliconFlowCreateImageGenerationsReq) (*SiliconFlowCreateImageGenerationsRsp, error) {
	var rsp SiliconFlowCreateImageGenerationsRsp
	err := p.post(p.channel.BaseUrl+"/v1/images/generations", req, &rsp)
	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// SiliconFlowCreateImageEditsReq defines the request for creating image edits.
type SiliconFlowCreateImageEditsReq struct {
	Prompt         string
	Image          io.Reader
	ImageFileName  string
	Mask           io.Reader // Optional
	MaskFileName   string    // Optional
	Model          string
	N              int
	Size           string
	ResponseFormat string
	User           string
}

// SiliconFlowCreateImageEditsRsp defines the response for creating image edits.
type SiliconFlowCreateImageEditsRsp = SiliconFlowCreateImageGenerationsRsp

// CreateImageEdits 根据提示词编辑图片。
func (p *SiliconFlow) CreateImageEdits(req *SiliconFlowCreateImageEditsReq) (*SiliconFlowCreateImageEditsRsp, error) {
	if req.Image == nil {
		err := errors.New("image reader is required")
		log.Errorf("error: %s", err)
		return nil, err
	}

	var rsp SiliconFlowCreateImageEditsRsp
	request := p.GetRequest()

	formData := make(map[string]string)
	formData["prompt"] = req.Prompt
	if req.Model != "" {
		formData["model"] = req.Model
	}
	if req.N > 0 {
		formData["n"] = strconv.Itoa(req.N)
	}
	if req.Size != "" {
		formData["size"] = req.Size
	}
	if req.ResponseFormat != "" {
		formData["response_format"] = req.ResponseFormat
	}
	if req.User != "" {
		formData["user"] = req.User
	}
	request.SetFormData(formData)
	request.SetFileReader("image", req.ImageFileName, req.Image)

	if req.Mask != nil {
		request.SetFileReader("mask", req.MaskFileName, req.Mask)
	}

	resp, err := request.
		SetResult(&rsp).
		Post(p.channel.BaseUrl + "/v1/images/edits")

	err = handleSiliconFlowError(resp, err)
	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}

	return &rsp, nil
}

// SiliconFlowCreateImageVariationsReq defines the request for creating image variations.
type SiliconFlowCreateImageVariationsReq struct {
	Image          io.Reader
	ImageFileName  string
	Model          string
	N              int
	Size           string
	ResponseFormat string
	User           string
}

// SiliconFlowCreateImageVariationsRsp defines the response for creating image variations.
type SiliconFlowCreateImageVariationsRsp = SiliconFlowCreateImageGenerationsRsp

// SiliconFlowFileObj defines the file object.
type SiliconFlowFileObj struct {
	ID            string `json:"id"`
	Bytes         int    `json:"bytes"`
	CreatedAt     int64  `json:"created_at"`
	Filename      string `json:"filename"`
	Object        string `json:"object"`
	Purpose       string `json:"purpose"`
	Status        string `json:"status"`
	StatusDetails any    `json:"status_details"`
}

// SiliconFlowListFilesRsp defines the response for listing files.
type SiliconFlowListFilesRsp struct {
	Object string               `json:"object"`
	Data   []SiliconFlowFileObj `json:"data"`
}

// SiliconFlowDeleteFileRsp corresponds to the response schema for deleting a file.
type SiliconFlowDeleteFileRsp struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

// SiliconFlowUploadFileReq defines the request for uploading a file.
type SiliconFlowUploadFileReq struct {
	File     io.Reader
	FileName string
	Purpose  string
}

// DeleteFile 删除指定 ID 的文件。
func (p *SiliconFlow) DeleteFile(fileID string) (*SiliconFlowDeleteFileRsp, error) {
	var rsp SiliconFlowDeleteFileRsp
	resp, err := p.GetRequest().
		SetResult(&rsp).
		SetPathParams(map[string]string{
			"file_id": fileID,
		}).
		Delete(p.channel.BaseUrl + "/v1/files/{file_id}")

	err = handleSiliconFlowError(resp, err)
	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// UploadFile 上传一个文件。
func (p *SiliconFlow) UploadFile(req *SiliconFlowUploadFileReq) (*SiliconFlowFileObj, error) {
	if req.File == nil {
		err := errors.New("file reader is required")
		log.Errorf("error: %s", err)
		return nil, err
	}

	var rsp SiliconFlowFileObj
	request := p.GetRequest()

	formData := make(map[string]string)
	formData["purpose"] = req.Purpose

	request.SetFormData(formData)
	request.SetFileReader("file", req.FileName, req.File)

	resp, err := request.
		SetResult(&rsp).
		Post(p.channel.BaseUrl + "/v1/files")

	err = handleSiliconFlowError(resp, err)
	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}

	return &rsp, nil
}

// ListFiles 列出用户上传的文件。
func (p *SiliconFlow) ListFiles() (*SiliconFlowListFilesRsp, error) {
	var rsp SiliconFlowListFilesRsp
	resp, err := p.GetRequest().
		SetResult(&rsp).
		Get(p.channel.BaseUrl + "/v1/files")

	err = handleSiliconFlowError(resp, err)
	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// RetrieveFile 获取指定文件的元数据。
func (p *SiliconFlow) RetrieveFile(fileID string) (*SiliconFlowFileObj, error) {
	var rsp SiliconFlowFileObj
	resp, err := p.GetRequest().
		SetResult(&rsp).
		SetPathParam("file_id", fileID).
		Get(p.channel.BaseUrl + "/v1/files/{file_id}")

	err = handleSiliconFlowError(resp, err)
	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// RetrieveFileContent 下载指定文件的内容。
func (p *SiliconFlow) RetrieveFileContent(fileID string) (io.ReadCloser, error) {
	resp, err := p.GetRequest().
		SetDoNotParseResponse(true). // 防止自动解析响应体
		SetPathParam("file_id", fileID).
		Get(p.channel.BaseUrl + "/v1/files/{file_id}/content")

	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}

	// 自定义错误处理，因为我们正在流式传输响应体
	if resp.IsError() {
		defer resp.RawBody().Close() // 确保在错误路径上关闭响应体
		body, readErr := io.ReadAll(resp.RawBody())
		if readErr != nil {
			log.Errorf("error: %s", readErr)
			return nil, readErr // 读取错误响应体时发生错误
		}

		var errRsp SiliconFlowErrorRsp
		e := json.Unmarshal(body, &errRsp)
		if e == nil && errRsp.Err.Message != "" {
			err = &errRsp
			log.Errorf("error: %s", err)
			return nil, err
		}
		err = errors.New(string(body))
		log.Errorf("error: %s", err)
		return nil, err
	}

	return resp.RawBody(), nil
}

// CreateImageVariations 创建图片的变体。
func (p *SiliconFlow) CreateImageVariations(req *SiliconFlowCreateImageVariationsReq) (*SiliconFlowCreateImageVariationsRsp, error) {
	if req.Image == nil {
		err := errors.New("image reader is required")
		log.Errorf("error: %s", err)
		return nil, err
	}

	var rsp SiliconFlowCreateImageVariationsRsp
	request := p.GetRequest()

	formData := make(map[string]string)
	if req.Model != "" {
		formData["model"] = req.Model
	}
	if req.N > 0 {
		formData["n"] = strconv.Itoa(req.N)
	}
	if req.Size != "" {
		formData["size"] = req.Size
	}
	if req.ResponseFormat != "" {
		formData["response_format"] = req.ResponseFormat
	}
	if req.User != "" {
		formData["user"] = req.User
	}
	request.SetFormData(formData)
	request.SetFileReader("image", req.ImageFileName, req.Image)

	resp, err := request.
		SetResult(&rsp).
		Post(p.channel.BaseUrl + "/v1/images/variations")

	err = handleSiliconFlowError(resp, err)
	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}

	return &rsp, nil
}

// SiliconFlowFineTuningJobRequest defines the request for creating a fine-tuning job.
type SiliconFlowFineTuningJobRequest struct {
	TrainingFile    string `json:"training_file"`
	ValidationFile  string `json:"validation_file,omitempty"`
	Model           string `json:"model,omitempty"`
	Hyperparameters struct {
		NEpochs int `json:"n_epochs,omitempty"`
	} `json:"hyperparameters,omitempty"`
	Suffix string `json:"suffix,omitempty"`
}

// SiliconFlowFineTuningJob defines the fine-tuning job object.
type SiliconFlowFineTuningJob struct {
	ID              string   `json:"id"`
	Object          string   `json:"object"`
	Model           string   `json:"model"`
	CreatedAt       int64    `json:"created_at"`
	FinishedAt      int64    `json:"finished_at"`
	FineTunedModel  string   `json:"fine_tuned_model"`
	OrganizationID  string   `json:"organization_id"`
	ResultFiles     []string `json:"result_files"`
	Status          string   `json:"status"`
	ValidationFile  string   `json:"validation_file"`
	TrainingFile    string   `json:"training_file"`
	Hyperparameters struct {
		NEpochs int `json:"n_epochs"`
	} `json:"hyperparameters"`
	TrainedTokens int `json:"trained_tokens"`
}

// SiliconFlowFineTuningJobList defines the response for listing fine-tuning jobs.
type SiliconFlowFineTuningJobList struct {
	Object  string                     `json:"object"`
	Data    []SiliconFlowFineTuningJob `json:"data"`
	HasMore bool                       `json:"has_more"`
}

// SiliconFlowFineTuningJobEvent defines a single event in a fine-tuning job.
type SiliconFlowFineTuningJobEvent struct {
	Object    string `json:"object"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	CreatedAt int64  `json:"created_at"`
}

// SiliconFlowFineTuningJobEventList defines the response for listing fine-tuning job events.
type SiliconFlowFineTuningJobEventList struct {
	Object string                          `json:"object"`
	Data   []SiliconFlowFineTuningJobEvent `json:"data"`
}

// CreateFineTuningJob 创建一个新的微调任务。
func (p *SiliconFlow) CreateFineTuningJob(req *SiliconFlowFineTuningJobRequest) (*SiliconFlowFineTuningJob, error) {
	var rsp SiliconFlowFineTuningJob
	err := p.post(p.channel.BaseUrl+"/v1/fine_tuning/jobs", req, &rsp)
	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// ListFineTuningJobs 列出现有的微调任务。
func (p *SiliconFlow) ListFineTuningJobs(limit int, after string) (*SiliconFlowFineTuningJobList, error) {
	var rsp SiliconFlowFineTuningJobList
	req := p.GetRequest().SetResult(&rsp)

	if limit > 0 {
		req.SetQueryParam("limit", strconv.Itoa(limit))
	}
	if after != "" {
		req.SetQueryParam("after", after)
	}

	resp, err := req.Get(p.channel.BaseUrl + "/v1/fine_tuning/jobs")

	err = handleSiliconFlowError(resp, err)
	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// RetrieveFineTuningJob 获取特定微调任务的信息。
func (p *SiliconFlow) RetrieveFineTuningJob(jobID string) (*SiliconFlowFineTuningJob, error) {
	var rsp SiliconFlowFineTuningJob
	resp, err := p.GetRequest().
		SetResult(&rsp).
		SetPathParam("job_id", jobID).
		Get(p.channel.BaseUrl + "/v1/fine_tuning/jobs/{job_id}")

	err = handleSiliconFlowError(resp, err)
	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// CancelFineTuningJob 取消一个微调任务。
func (p *SiliconFlow) CancelFineTuningJob(jobID string) (*SiliconFlowFineTuningJob, error) {
	var rsp SiliconFlowFineTuningJob
	resp, err := p.GetRequest().
		SetResult(&rsp).
		SetPathParam("job_id", jobID).
		Post(p.channel.BaseUrl + "/v1/fine_tuning/jobs/{job_id}/cancel")

	err = handleSiliconFlowError(resp, err)
	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}
	return &rsp, nil
}

// SiliconFlowListFineTuningJobEventsReq defines the request for listing fine-tuning job events.
type SiliconFlowListFineTuningJobEventsReq struct {
	FineTuningJobID string
	Limit           int
	After           string
}

// ListFineTuningJobEvents 获取微调任务的事件列表。
func (p *SiliconFlow) ListFineTuningJobEvents(req *SiliconFlowListFineTuningJobEventsReq) (*SiliconFlowFineTuningJobEventList, error) {
	var rsp SiliconFlowFineTuningJobEventList
	request := p.GetRequest().
		SetResult(&rsp).
		SetPathParam("job_id", req.FineTuningJobID)

	if req.Limit > 0 {
		request.SetQueryParam("limit", strconv.Itoa(req.Limit))
	}
	if req.After != "" {
		request.SetQueryParam("after", req.After)
	}

	resp, err := request.Get(p.channel.BaseUrl + "/v1/fine_tuning/jobs/{job_id}/events")

	err = handleSiliconFlowError(resp, err)
	if err != nil {
		log.Errorf("error: %s", err)
		return nil, err
	}
	return &rsp, nil
}
