package channel

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/lazygophers/aiload"
	"github.com/lazygophers/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupSiliconFlowTest a helper function to setup tests
func setupSiliconFlowTest(t *testing.T) *SiliconFlow {
	t.Helper()

	mockChannel := &aiload.ModelChannel{
		Token: "test-token",
	}
	p := NewSiliconFlow(mockChannel)

	httpmock.ActivateNonDefault(client.GetClient())
	t.Cleanup(httpmock.DeactivateAndReset)

	return p
}

func TestMain(m *testing.M) {
	// 禁用测试期间的日志输出，保持测试结果的清洁
	log.SetOutput(io.Discard)
	// 运行包中的所有测试
	exitCode := m.Run()
	// 恢复日志输出到标准错误
	log.SetOutput(os.Stderr)
	// 退出并返回测试结果的退出码
	os.Exit(exitCode)
}

func TestSiliconFlowChat(t *testing.T) {
	p := setupSiliconFlowTest(t)

	mockSuccessResponse := SiliconFlowCreateChatCompletionRsp{
		ID:      "chatcmpl-mock-id",
		Object:  "chat.completion",
		Created: 1677652288,
		Model:   "deepseek-coder",
		Choices: []struct {
			Index   int `json:"index"`
			Message struct {
				Role    string `json:"role"`
				Content string `json:"content"`
			} `json:"message"`
			FinishReason string `json:"finish_reason"`
		}{
			{
				Index: 0,
				Message: struct {
					Role    string `json:"role"`
					Content string `json:"content"`
				}{
					Role:    "assistant",
					Content: "\n\nHello there, how may I assist you today?",
				},
				FinishReason: "stop",
			},
		},
		Usage: struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		}{
			PromptTokens:     9,
			CompletionTokens: 12,
			TotalTokens:      21,
		},
	}
	successRespBody, err := json.Marshal(mockSuccessResponse)
	require.NoError(t, err)

	tests := []struct {
		name          string
		giveRequest   *SiliconFlowCreateChatCompletionReq
		mockResponder httpmock.Responder
		checkResponse func(t *testing.T, rsp *SiliconFlowCreateChatCompletionRsp, err error)
	}{
		{
			name: "success",
			giveRequest: &SiliconFlowCreateChatCompletionReq{
				Model: "deepseek-coder",
				Messages: []SiliconFlowCreateChatCompletionMessage{
					{Role: "user", Content: "Hello"},
				},
			},
			mockResponder: newMockResponder(http.StatusOK, string(successRespBody)),
			checkResponse: func(t *testing.T, rsp *SiliconFlowCreateChatCompletionRsp, err error) {
				assert.NoError(t, err)
				require.NotNil(t, rsp)
				assert.Equal(t, "chatcmpl-mock-id", rsp.ID)
				require.NotEmpty(t, rsp.Choices)
				assert.Equal(t, "assistant", rsp.Choices[0].Message.Role)
				assert.Equal(t, "\n\nHello there, how may I assist you today?", rsp.Choices[0].Message.Content)
			},
		},
		{
			name: "api error",
			giveRequest: &SiliconFlowCreateChatCompletionReq{
				Model: "deepseek-coder",
				Messages: []SiliconFlowCreateChatCompletionMessage{
					{Role: "user", Content: "Hello"},
				},
			},
			mockResponder: newMockResponder(http.StatusInternalServerError, `{"error":{"message":"Internal server error"}}`),
			checkResponse: func(t *testing.T, rsp *SiliconFlowCreateChatCompletionRsp, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
				checkAPIError(t, err, "Internal server error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.RegisterResponder("POST", "https://api.siliconflow.cn/v1/chat/completions", tt.mockResponder)

			rsp, err := p.CreateChatCompletion(tt.giveRequest)

			tt.checkResponse(t, rsp, err)
		})
	}
}

func TestSiliconFlowEmbedding(t *testing.T) {
	p := setupSiliconFlowTest(t)

	mockSuccessResponse := SiliconFlowCreateEmbeddingRsp{
		Object: "list",
		Data: []struct {
			Object    string    `json:"object"`
			Embedding []float64 `json:"embedding"`
			Index     int       `json:"index"`
		}{
			{
				Object:    "embedding",
				Embedding: []float64{0.1, 0.2, 0.3},
				Index:     0,
			},
		},
		Model: "bge-large-zh-v1.5",
		Usage: struct {
			PromptTokens int `json:"prompt_tokens"`
			TotalTokens  int `json:"total_tokens"`
		}{
			PromptTokens: 10,
			TotalTokens:  10,
		},
	}
	successRespBody, err := json.Marshal(mockSuccessResponse)
	require.NoError(t, err)

	tests := []struct {
		name          string
		giveRequest   *SiliconFlowCreateEmbeddingReq
		mockResponder httpmock.Responder
		checkResponse func(t *testing.T, rsp *SiliconFlowCreateEmbeddingRsp, err error)
	}{
		{
			name: "success",
			giveRequest: &SiliconFlowCreateEmbeddingReq{
				Model: "bge-large-zh-v1.5",
				Input: []string{"hello"},
			},
			mockResponder: newMockResponder(http.StatusOK, string(successRespBody)),
			checkResponse: func(t *testing.T, rsp *SiliconFlowCreateEmbeddingRsp, err error) {
				assert.NoError(t, err)
				require.NotNil(t, rsp)
				assert.Equal(t, "list", rsp.Object)
				require.Len(t, rsp.Data, 1)
				assert.Equal(t, []float64{0.1, 0.2, 0.3}, rsp.Data[0].Embedding)
			},
		},
		{
			name: "api error",
			giveRequest: &SiliconFlowCreateEmbeddingReq{
				Model: "bge-large-zh-v1.5",
				Input: []string{"hello"},
			},
			mockResponder: newMockResponder(http.StatusInternalServerError, `{"error":{"message":"Internal server error"}}`),
			checkResponse: func(t *testing.T, rsp *SiliconFlowCreateEmbeddingRsp, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
				checkAPIError(t, err, "Internal server error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.RegisterResponder("POST", "https://api.siliconflow.cn/v1/embeddings", tt.mockResponder)

			rsp, err := p.CreateEmbedding(tt.giveRequest)

			tt.checkResponse(t, rsp, err)
		})
	}
}

func TestSiliconFlow_Image(t *testing.T) {
	p := setupSiliconFlowTest(t)

	// Mocks for successful responses
	mockGenerationsSuccessResponse := SiliconFlowCreateImageGenerationsRsp{
		Created: 1677652288,
		Data: []struct {
			URL     string `json:"url,omitempty"`
			B64JSON string `json:"b64_json,omitempty"`
		}{{URL: "https://example.com/image.png"}},
	}
	generationsSuccessBody, err := json.Marshal(mockGenerationsSuccessResponse)
	require.NoError(t, err)

	mockEditsSuccessResponse := SiliconFlowCreateImageEditsRsp{
		Created: 1677652288,
		Data: []struct {
			URL     string `json:"url,omitempty"`
			B64JSON string `json:"b64_json,omitempty"`
		}{{URL: "https://example.com/edited-image.png"}},
	}
	editsSuccessBody, err := json.Marshal(mockEditsSuccessResponse)
	require.NoError(t, err)

	mockVariationsSuccessResponse := SiliconFlowCreateImageVariationsRsp{
		Created: 1677652288,
		Data: []struct {
			URL     string `json:"url,omitempty"`
			B64JSON string `json:"b64_json,omitempty"`
		}{{URL: "https://example.com/variation-image.png"}},
	}
	variationsSuccessBody, err := json.Marshal(mockVariationsSuccessResponse)
	require.NoError(t, err)

	tests := []struct {
		name          string
		endpoint      string
		giveRequest   any
		mockResponder httpmock.Responder
		apiCall       func(p *SiliconFlow, req any) (any, error)
		checkResponse func(t *testing.T, rsp any, err error)
	}{
		// Generations
		{
			name:          "generations/success",
			endpoint:      "https://api.siliconflow.cn/v1/images/generations",
			giveRequest:   &SiliconFlowCreateImageGenerationsReq{Prompt: "a cat"},
			mockResponder: newMockResponder(http.StatusOK, string(generationsSuccessBody)),
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.CreateImageGenerations(req.(*SiliconFlowCreateImageGenerationsReq))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.NoError(t, err)
				r, ok := rsp.(*SiliconFlowCreateImageGenerationsRsp)
				require.True(t, ok)
				require.NotNil(t, r)
				require.Len(t, r.Data, 1)
				assert.Equal(t, "https://example.com/image.png", r.Data[0].URL)
			},
		},
		{
			name:          "generations/api_error",
			endpoint:      "https://api.siliconflow.cn/v1/images/generations",
			giveRequest:   &SiliconFlowCreateImageGenerationsReq{Prompt: "a cat"},
			mockResponder: newMockResponder(http.StatusInternalServerError, `{"error":{"message":"Internal server error"}}`),
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.CreateImageGenerations(req.(*SiliconFlowCreateImageGenerationsReq))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
				checkAPIError(t, err, "Internal server error")
			},
		},
		{
			name:          "generations/network_error",
			endpoint:      "https://api.siliconflow.cn/v1/images/generations",
			giveRequest:   &SiliconFlowCreateImageGenerationsReq{Prompt: "a cat"},
			mockResponder: httpmock.NewErrorResponder(errors.New("network error")),
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.CreateImageGenerations(req.(*SiliconFlowCreateImageGenerationsReq))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
				assert.Contains(t, err.Error(), "network error")
			},
		},
		// Edits
		{
			name:     "edits/success",
			endpoint: "https://api.siliconflow.cn/v1/images/edits",
			giveRequest: &SiliconFlowCreateImageEditsReq{
				Prompt:        "a cute cat",
				Image:         strings.NewReader("fake-image-data"),
				ImageFileName: "test-image.png",
				Model:         "stable-diffusion-xl-1024-v1-0",
			},
			mockResponder: newMockResponder(http.StatusOK, string(editsSuccessBody)),
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.CreateImageEdits(req.(*SiliconFlowCreateImageEditsReq))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.NoError(t, err)
				r, ok := rsp.(*SiliconFlowCreateImageEditsRsp)
				require.True(t, ok)
				require.NotNil(t, r)
				require.Len(t, r.Data, 1)
				assert.Equal(t, "https://example.com/edited-image.png", r.Data[0].URL)
			},
		},
		{
			name:     "edits/missing_image",
			endpoint: "https://api.siliconflow.cn/v1/images/edits",
			giveRequest: &SiliconFlowCreateImageEditsReq{
				Prompt: "a cute cat",
			},
			mockResponder: newMockResponder(http.StatusBadRequest, ""), // Not called
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.CreateImageEdits(req.(*SiliconFlowCreateImageEditsReq))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
				assert.Contains(t, err.Error(), "image reader is required")
			},
		},
		{
			name:     "edits/api_error",
			endpoint: "https://api.siliconflow.cn/v1/images/edits",
			giveRequest: &SiliconFlowCreateImageEditsReq{
				Prompt:        "a cute cat",
				Image:         strings.NewReader("fake-image-data"),
				ImageFileName: "test-image.png",
			},
			mockResponder: newMockResponder(http.StatusInternalServerError, `{"error":{"message":"Internal server error"}}`),
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.CreateImageEdits(req.(*SiliconFlowCreateImageEditsReq))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
				checkAPIError(t, err, "Internal server error")
			},
		},
		{
			name:     "edits/network_error",
			endpoint: "https://api.siliconflow.cn/v1/images/edits",
			giveRequest: &SiliconFlowCreateImageEditsReq{
				Prompt:        "a cute cat",
				Image:         strings.NewReader("fake-image-data"),
				ImageFileName: "test-image.png",
			},
			mockResponder: httpmock.NewErrorResponder(errors.New("network error")),
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.CreateImageEdits(req.(*SiliconFlowCreateImageEditsReq))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
				assert.Contains(t, err.Error(), "network error")
			},
		},
		{
			name:     "edits/success_with_all_fields",
			endpoint: "https://api.siliconflow.cn/v1/images/edits",
			giveRequest: &SiliconFlowCreateImageEditsReq{
				Prompt:         "a cute cat with a hat",
				Image:          strings.NewReader("fake-image-data-full"),
				ImageFileName:  "test-image-full.png",
				Mask:           strings.NewReader("fake-mask-data"),
				MaskFileName:   "mask.png",
				Model:          "stable-diffusion-xl-1024-v1-0",
				N:              2,
				Size:           "1024x1024",
				ResponseFormat: "b64_json",
				User:           "test-user",
			},
			mockResponder: newMockResponder(http.StatusOK, string(editsSuccessBody)),
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.CreateImageEdits(req.(*SiliconFlowCreateImageEditsReq))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.NoError(t, err)
				r, ok := rsp.(*SiliconFlowCreateImageEditsRsp)
				require.True(t, ok)
				require.NotNil(t, r)
				require.Len(t, r.Data, 1)
				assert.Equal(t, "https://example.com/edited-image.png", r.Data[0].URL)
			},
		},
		// Variations
		{
			name:     "variations/success",
			endpoint: "https://api.siliconflow.cn/v1/images/variations",
			giveRequest: &SiliconFlowCreateImageVariationsReq{
				Image:         strings.NewReader("fake-variation-image-data"),
				ImageFileName: "test-variation-image.png",
				Model:         "stable-diffusion-xl-1024-v1-0",
			},
			mockResponder: newMockResponder(http.StatusOK, string(variationsSuccessBody)),
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.CreateImageVariations(req.(*SiliconFlowCreateImageVariationsReq))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.NoError(t, err)
				r, ok := rsp.(*SiliconFlowCreateImageVariationsRsp)
				require.True(t, ok)
				require.NotNil(t, r)
				require.Len(t, r.Data, 1)
				assert.Equal(t, "https://example.com/variation-image.png", r.Data[0].URL)
			},
		},
		{
			name:          "variations/missing_image",
			endpoint:      "https://api.siliconflow.cn/v1/images/variations",
			giveRequest:   &SiliconFlowCreateImageVariationsReq{},
			mockResponder: newMockResponder(http.StatusBadRequest, ""), // Not called
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.CreateImageVariations(req.(*SiliconFlowCreateImageVariationsReq))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
				assert.Equal(t, "image reader is required", err.Error())
			},
		},
		{
			name:     "variations/network_error",
			endpoint: "https://api.siliconflow.cn/v1/images/variations",
			giveRequest: &SiliconFlowCreateImageVariationsReq{
				Image:         strings.NewReader("fake-variation-image-data"),
				ImageFileName: "test-variation-image.png",
			},
			mockResponder: httpmock.NewErrorResponder(errors.New("network error")),
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.CreateImageVariations(req.(*SiliconFlowCreateImageVariationsReq))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
				assert.Contains(t, err.Error(), "network error")
			},
		},
		{
			name:     "variations/success_with_all_fields",
			endpoint: "https://api.siliconflow.cn/v1/images/variations",
			giveRequest: &SiliconFlowCreateImageVariationsReq{
				Image:          strings.NewReader("fake-variation-image-data-full"),
				ImageFileName:  "test-variation-image-full.png",
				Model:          "stable-diffusion-xl-1024-v1-0",
				N:              2,
				Size:           "1024x1024",
				ResponseFormat: "url",
				User:           "test-user-variations",
			},
			mockResponder: newMockResponder(http.StatusOK, string(variationsSuccessBody)),
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.CreateImageVariations(req.(*SiliconFlowCreateImageVariationsReq))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.NoError(t, err)
				r, ok := rsp.(*SiliconFlowCreateImageVariationsRsp)
				require.True(t, ok)
				require.NotNil(t, r)
				require.Len(t, r.Data, 1)
				assert.Equal(t, "https://example.com/variation-image.png", r.Data[0].URL)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.RegisterResponder("POST", tt.endpoint, tt.mockResponder)

			rsp, err := tt.apiCall(p, tt.giveRequest)

			tt.checkResponse(t, rsp, err)
		})
	}
}

func TestSiliconFlow_Audio(t *testing.T) {
	p := setupSiliconFlowTest(t)

	// Mocks for successful responses
	mockTranslationSuccessResponse := SiliconFlowCreateAudioTranslationRsp{
		Text: "Hello, this is a test.",
	}
	translationSuccessBody, err := json.Marshal(mockTranslationSuccessResponse)
	require.NoError(t, err)

	mockTranscriptionSuccessResponse := SiliconFlowCreateAudioTranscriptionResp{
		Text: "This is a test transcription.",
	}
	transcriptionSuccessBody, err := json.Marshal(mockTranscriptionSuccessResponse)
	require.NoError(t, err)

	tests := []struct {
		name          string
		endpoint      string
		giveRequest   any
		mockResponder httpmock.Responder
		apiCall       func(p *SiliconFlow, req any) (any, error)
		checkResponse func(t *testing.T, rsp any, err error)
	}{
		// Translations
		{
			name:     "translations/success",
			endpoint: "https://api.siliconflow.cn/v1/audio/translations",
			giveRequest: &SiliconFlowCreateAudioTranslationReq{
				File:     strings.NewReader("fake-audio-data"),
				FileName: "test-audio.mp3",
				Model:    "whisper-large-v3",
			},
			mockResponder: newMockResponder(http.StatusOK, string(translationSuccessBody)),
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.CreateAudioTranslation(req.(*SiliconFlowCreateAudioTranslationReq))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.NoError(t, err)
				r, ok := rsp.(*SiliconFlowCreateAudioTranslationRsp)
				require.True(t, ok)
				require.NotNil(t, r)
				assert.Equal(t, "Hello, this is a test.", r.Text)
			},
		},
		{
			name:     "translations/api_error",
			endpoint: "https://api.siliconflow.cn/v1/audio/translations",
			giveRequest: &SiliconFlowCreateAudioTranslationReq{
				File:     strings.NewReader("fake-audio-data"),
				FileName: "test-audio.mp3",
				Model:    "whisper-large-v3",
			},
			mockResponder: newMockResponder(http.StatusInternalServerError, `{"error":{"message":"Internal server error"}}`),
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.CreateAudioTranslation(req.(*SiliconFlowCreateAudioTranslationReq))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
				checkAPIError(t, err, "Internal server error")
			},
		},
		{
			name:          "translations/missing_file",
			endpoint:      "https://api.siliconflow.cn/v1/audio/translations",
			giveRequest:   &SiliconFlowCreateAudioTranslationReq{Model: "whisper-large-v3"},
			mockResponder: newMockResponder(http.StatusBadRequest, ""), // Not called
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.CreateAudioTranslation(req.(*SiliconFlowCreateAudioTranslationReq))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
				assert.Equal(t, "file reader is required", err.Error())
			},
		},
		{
			name:     "translations/network_error",
			endpoint: "https://api.siliconflow.cn/v1/audio/translations",
			giveRequest: &SiliconFlowCreateAudioTranslationReq{
				File:     strings.NewReader("fake-audio-data"),
				FileName: "test-audio.mp3",
				Model:    "whisper-large-v3",
			},
			mockResponder: httpmock.NewErrorResponder(errors.New("network error")),
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.CreateAudioTranslation(req.(*SiliconFlowCreateAudioTranslationReq))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
				assert.Contains(t, err.Error(), "network error")
			},
		},
		// Transcriptions
		{
			name:     "transcriptions/success",
			endpoint: "https://api.siliconflow.cn/v1/audio/transcriptions",
			giveRequest: &SiliconFlowCreateAudioTranscriptionReq{
				File:     strings.NewReader("fake-audio-data-for-transcription"),
				FileName: "test-audio.mp3",
				Model:    "whisper-large-v3",
			},
			mockResponder: newMockResponder(http.StatusOK, string(transcriptionSuccessBody)),
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.CreateAudioTranscription(req.(*SiliconFlowCreateAudioTranscriptionReq))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.NoError(t, err)
				r, ok := rsp.(*SiliconFlowCreateAudioTranscriptionResp)
				require.True(t, ok)
				require.NotNil(t, r)
				assert.Equal(t, "This is a test transcription.", r.Text)
			},
		},
		{
			name:     "transcriptions/api_error",
			endpoint: "https://api.siliconflow.cn/v1/audio/transcriptions",
			giveRequest: &SiliconFlowCreateAudioTranscriptionReq{
				File:     strings.NewReader("fake-audio-data-for-transcription"),
				FileName: "test-audio.mp3",
				Model:    "whisper-large-v3",
			},
			mockResponder: newMockResponder(http.StatusInternalServerError, `{"error":{"message":"Internal server error"}}`),
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.CreateAudioTranscription(req.(*SiliconFlowCreateAudioTranscriptionReq))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
				checkAPIError(t, err, "Internal server error")
			},
		},
		{
			name:          "transcriptions/missing_file",
			endpoint:      "https://api.siliconflow.cn/v1/audio/transcriptions",
			giveRequest:   &SiliconFlowCreateAudioTranscriptionReq{Model: "whisper-large-v3"},
			mockResponder: newMockResponder(http.StatusBadRequest, ""), // Not called
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.CreateAudioTranscription(req.(*SiliconFlowCreateAudioTranscriptionReq))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
				assert.Equal(t, "file reader is required", err.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.RegisterResponder("POST", tt.endpoint, tt.mockResponder)

			rsp, err := tt.apiCall(p, tt.giveRequest)

			tt.checkResponse(t, rsp, err)
		})
	}
}

func TestSiliconflowFile(t *testing.T) {
	p := setupSiliconFlowTest(t)

	// Mocks for successful responses
	mockUploadFileSuccessResponse := SiliconFlowFileObj{
		ID: "file-mock-id", Bytes: 140, CreatedAt: 1677652288, Filename: "test-file.jsonl", Object: "file", Purpose: "fine-tune", Status: "processed",
	}
	uploadFileSuccessBody, err := json.Marshal(mockUploadFileSuccessResponse)
	require.NoError(t, err)

	mockListFilesSuccessResponse := SiliconFlowListFilesRsp{
		Object: "list",
		Data: []SiliconFlowFileObj{
			{ID: "file-mock-id-1", Filename: "test-file-1.jsonl"},
			{ID: "file-mock-id-2", Filename: "test-file-2.jsonl"},
		},
	}
	listFilesSuccessBody, err := json.Marshal(mockListFilesSuccessResponse)
	require.NoError(t, err)

	mockDeleteFileSuccessResponse := SiliconFlowDeleteFileRsp{
		ID: "file-mock-id-to-delete", Object: "file", Deleted: true,
	}
	deleteFileSuccessBody, err := json.Marshal(mockDeleteFileSuccessResponse)
	require.NoError(t, err)

	mockRetrieveFileSuccessResponse := SiliconFlowFileObj{
		ID: "file-mock-id-1", Filename: "test-file-1.jsonl",
	}
	retrieveFileSuccessBody, err := json.Marshal(mockRetrieveFileSuccessResponse)
	require.NoError(t, err)

	tests := []struct {
		name          string
		method        string
		endpoint      string
		giveRequest   any
		mockResponder httpmock.Responder
		apiCall       func(p *SiliconFlow, req any) (any, error)
		checkResponse func(t *testing.T, rsp any, err error)
	}{
		// UploadFile
		{
			name:     "upload/success",
			method:   "POST",
			endpoint: "https://api.siliconflow.cn/v1/files",
			giveRequest: &SiliconFlowUploadFileReq{
				File: strings.NewReader("fake-file-data"), FileName: "test-file.jsonl", Purpose: "fine-tune",
			},
			mockResponder: newMockResponder(http.StatusOK, string(uploadFileSuccessBody)),
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.UploadFile(req.(*SiliconFlowUploadFileReq))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.NoError(t, err)
				r, ok := rsp.(*SiliconFlowFileObj)
				require.True(t, ok)
				assert.Equal(t, "file-mock-id", r.ID)
			},
		},
		{
			name:     "upload/missing_file",
			method:   "POST",
			endpoint: "https://api.siliconflow.cn/v1/files",
			giveRequest: &SiliconFlowUploadFileReq{
				Purpose: "fine-tune",
			},
			mockResponder: newMockResponder(http.StatusBadRequest, ""), // Not called
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.UploadFile(req.(*SiliconFlowUploadFileReq))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
				assert.Equal(t, "file reader is required", err.Error())
			},
		},
		{
			name:          "upload/network_error",
			method:        "POST",
			endpoint:      "https://api.siliconflow.cn/v1/files",
			giveRequest:   &SiliconFlowUploadFileReq{File: strings.NewReader("d"), FileName: "d", Purpose: "fine-tune"},
			mockResponder: httpmock.NewErrorResponder(errors.New("network error")),
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.UploadFile(req.(*SiliconFlowUploadFileReq))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
				assert.Contains(t, err.Error(), "network error")
			},
		},
		// ListFiles
		{
			name:          "list/success",
			method:        "GET",
			endpoint:      "https://api.siliconflow.cn/v1/files",
			giveRequest:   nil,
			mockResponder: newMockResponder(http.StatusOK, string(listFilesSuccessBody)),
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.ListFiles()
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.NoError(t, err)
				r, ok := rsp.(*SiliconFlowListFilesRsp)
				require.True(t, ok)
				require.Len(t, r.Data, 2)
				assert.Equal(t, "file-mock-id-1", r.Data[0].ID)
			},
		},
		// DeleteFile
		{
			name:          "delete/success",
			method:        "DELETE",
			endpoint:      "https://api.siliconflow.cn/v1/files/file-mock-id-to-delete",
			giveRequest:   "file-mock-id-to-delete",
			mockResponder: newMockResponder(http.StatusOK, string(deleteFileSuccessBody)),
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.DeleteFile(req.(string))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.NoError(t, err)
				r, ok := rsp.(*SiliconFlowDeleteFileRsp)
				require.True(t, ok)
				assert.True(t, r.Deleted)
			},
		},
		// RetrieveFile
		{
			name:          "retrieve/success",
			method:        "GET",
			endpoint:      "https://api.siliconflow.cn/v1/files/file-mock-id-1",
			giveRequest:   "file-mock-id-1",
			mockResponder: newMockResponder(http.StatusOK, string(retrieveFileSuccessBody)),
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.RetrieveFile(req.(string))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.NoError(t, err)
				r, ok := rsp.(*SiliconFlowFileObj)
				require.True(t, ok)
				assert.Equal(t, "file-mock-id-1", r.ID)
			},
		},
		// RetrieveFileContent
		{
			name:          "retrieve_content/success",
			method:        "GET",
			endpoint:      "https://api.siliconflow.cn/v1/files/file-id-success/content",
			giveRequest:   "file-id-success",
			mockResponder: newMockResponder(http.StatusOK, "file content"),
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.RetrieveFileContent(req.(string))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.NoError(t, err)
				body, ok := rsp.(io.ReadCloser)
				require.True(t, ok)
				defer body.Close()
				content, _ := io.ReadAll(body)
				assert.Equal(t, "file content", string(content))
			},
		},
		{
			name:          "retrieve_content/api_error",
			method:        "GET",
			endpoint:      "https://api.siliconflow.cn/v1/files/file-id-api-error/content",
			giveRequest:   "file-id-api-error",
			mockResponder: newMockResponder(http.StatusNotFound, `{"error":{"message":"File not found"}}`),
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.RetrieveFileContent(req.(string))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.Nil(t, rsp)
				checkAPIError(t, err, "File not found")
			},
		},
		{
			name:          "retrieve_content/unstructured_error",
			method:        "GET",
			endpoint:      "https://api.siliconflow.cn/v1/files/file-id-unstructured-error/content",
			giveRequest:   "file-id-unstructured-error",
			mockResponder: newMockResponder(http.StatusInternalServerError, "internal server error"),
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.RetrieveFileContent(req.(string))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
				assert.Contains(t, err.Error(), "internal server error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.RegisterResponder(tt.method, tt.endpoint, tt.mockResponder)
			rsp, err := tt.apiCall(p, tt.giveRequest)
			tt.checkResponse(t, rsp, err)
		})
	}
}
func TestSiliconflowFineTuning(t *testing.T) {
	p := setupSiliconFlowTest(t)

	// Mocks for successful responses
	mockCreateJobSuccessResponse := SiliconFlowFineTuningJob{
		ID: "ft-job-mock-id", Status: "succeeded", FineTunedModel: "ft:gpt-3.5-turbo:my-org:custom-model-name:1",
	}
	createJobSuccessBody, err := json.Marshal(mockCreateJobSuccessResponse)
	require.NoError(t, err)

	mockListJobsSuccessResponse := SiliconFlowFineTuningJobList{
		Object: "list", Data: []SiliconFlowFineTuningJob{{ID: "ft-job-mock-id-1"}}, HasMore: false,
	}
	listJobsSuccessBody, err := json.Marshal(mockListJobsSuccessResponse)
	require.NoError(t, err)

	mockRetrieveJobSuccessResponse := SiliconFlowFineTuningJob{
		ID: "ft-job-mock-id-success", Status: "succeeded", FineTunedModel: "ft:gpt-3.5-turbo:my-org:custom-model-name:1",
	}
	retrieveJobSuccessBody, err := json.Marshal(mockRetrieveJobSuccessResponse)
	require.NoError(t, err)

	mockCancelJobSuccessResponse := SiliconFlowFineTuningJob{
		ID: "ft-job-mock-id-success", Status: "cancelled",
	}
	cancelJobSuccessBody, err := json.Marshal(mockCancelJobSuccessResponse)
	require.NoError(t, err)

	mockListEventsSuccessResponse := SiliconFlowFineTuningJobEventList{
		Object: "list", Data: []SiliconFlowFineTuningJobEvent{{Object: "fine_tuning.job.event", Message: "Job started"}},
	}
	listEventsSuccessBody, err := json.Marshal(mockListEventsSuccessResponse)
	require.NoError(t, err)

	tests := []struct {
		name          string
		method        string
		endpoint      string
		giveRequest   any
		mockResponder httpmock.Responder
		apiCall       func(p *SiliconFlow, req any) (any, error)
		checkResponse func(t *testing.T, rsp any, err error)
	}{
		// CreateFineTuningJob
		{
			name:     "create_job/success",
			method:   "POST",
			endpoint: "https://api.siliconflow.cn/v1/fine_tuning/jobs",
			giveRequest: &SiliconFlowFineTuningJobRequest{
				TrainingFile: "file-train-mock-id", Model: "gpt-3.5-turbo",
			},
			mockResponder: newMockResponder(http.StatusOK, string(createJobSuccessBody)),
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.CreateFineTuningJob(req.(*SiliconFlowFineTuningJobRequest))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.NoError(t, err)
				r, ok := rsp.(*SiliconFlowFineTuningJob)
				require.True(t, ok)
				assert.Equal(t, "ft-job-mock-id", r.ID)
			},
		},
		// ListFineTuningJobs
		{
			name:     "list_jobs/success",
			method:   "GET",
			endpoint: "https://api.siliconflow.cn/v1/fine_tuning/jobs",
			giveRequest: struct {
				limit int
				after string
			}{limit: 10, after: "job-id-123"},
			mockResponder: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, "10", req.URL.Query().Get("limit"))
				assert.Equal(t, "job-id-123", req.URL.Query().Get("after"))
				return newMockResponder(http.StatusOK, string(listJobsSuccessBody))(req)
			},
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				r := req.(struct {
					limit int
					after string
				})
				return p.ListFineTuningJobs(r.limit, r.after)
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.NoError(t, err)
				r, ok := rsp.(*SiliconFlowFineTuningJobList)
				require.True(t, ok)
				require.Len(t, r.Data, 1)
				assert.Equal(t, "ft-job-mock-id-1", r.Data[0].ID)
			},
		},
		{
			name:     "list_jobs/success_no_params",
			method:   "GET",
			endpoint: "https://api.siliconflow.cn/v1/fine_tuning/jobs",
			giveRequest: struct {
				limit int
				after string
			}{limit: 0, after: ""},
			mockResponder: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, "", req.URL.Query().Get("limit"))
				assert.Equal(t, "", req.URL.Query().Get("after"))
				return newMockResponder(http.StatusOK, string(listJobsSuccessBody))(req)
			},
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				r := req.(struct {
					limit int
					after string
				})
				return p.ListFineTuningJobs(r.limit, r.after)
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.NoError(t, err)
				r, ok := rsp.(*SiliconFlowFineTuningJobList)
				require.True(t, ok)
				require.Len(t, r.Data, 1)
				assert.Equal(t, "ft-job-mock-id-1", r.Data[0].ID)
			},
		},
		// RetrieveFineTuningJob
		{
			name:          "retrieve_job/success",
			method:        "GET",
			endpoint:      "https://api.siliconflow.cn/v1/fine_tuning/jobs/ft-job-mock-id-success",
			giveRequest:   "ft-job-mock-id-success",
			mockResponder: newMockResponder(http.StatusOK, string(retrieveJobSuccessBody)),
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.RetrieveFineTuningJob(req.(string))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.NoError(t, err)
				r, ok := rsp.(*SiliconFlowFineTuningJob)
				require.True(t, ok)
				assert.Equal(t, "ft-job-mock-id-success", r.ID)
			},
		},
		// CancelFineTuningJob
		{
			name:          "cancel_job/success",
			method:        "POST",
			endpoint:      "https://api.siliconflow.cn/v1/fine_tuning/jobs/ft-job-mock-id-success/cancel",
			giveRequest:   "ft-job-mock-id-success",
			mockResponder: newMockResponder(http.StatusOK, string(cancelJobSuccessBody)),
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.CancelFineTuningJob(req.(string))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.NoError(t, err)
				r, ok := rsp.(*SiliconFlowFineTuningJob)
				require.True(t, ok)
				assert.Equal(t, "cancelled", r.Status)
			},
		},
		// ListFineTuningJobEvents
		{
			name:     "list_events/success",
			method:   "GET",
			endpoint: "https://api.siliconflow.cn/v1/fine_tuning/jobs/ft-job-123/events",
			giveRequest: &SiliconFlowListFineTuningJobEventsReq{
				FineTuningJobID: "ft-job-123", Limit: 5, After: "event-id-0",
			},
			mockResponder: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, "5", req.URL.Query().Get("limit"))
				assert.Equal(t, "event-id-0", req.URL.Query().Get("after"))
				return newMockResponder(http.StatusOK, string(listEventsSuccessBody))(req)
			},
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.ListFineTuningJobEvents(req.(*SiliconFlowListFineTuningJobEventsReq))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.NoError(t, err)
				r, ok := rsp.(*SiliconFlowFineTuningJobEventList)
				require.True(t, ok)
				require.Len(t, r.Data, 1)
				assert.Equal(t, "Job started", r.Data[0].Message)
			},
		},
		{
			name:     "list_events/success_no_params",
			method:   "GET",
			endpoint: "https://api.siliconflow.cn/v1/fine_tuning/jobs/ft-job-123/events",
			giveRequest: &SiliconFlowListFineTuningJobEventsReq{
				FineTuningJobID: "ft-job-123",
			},
			mockResponder: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, "", req.URL.Query().Get("limit"))
				assert.Equal(t, "", req.URL.Query().Get("after"))
				return newMockResponder(http.StatusOK, string(listEventsSuccessBody))(req)
			},
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.ListFineTuningJobEvents(req.(*SiliconFlowListFineTuningJobEventsReq))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.NoError(t, err)
				r, ok := rsp.(*SiliconFlowFineTuningJobEventList)
				require.True(t, ok)
				require.Len(t, r.Data, 1)
				assert.Equal(t, "Job started", r.Data[0].Message)
			},
		},
		{
			name:     "list_events/api_error",
			method:   "GET",
			endpoint: "https://api.siliconflow.cn/v1/fine_tuning/jobs/ft-job-456/events",
			giveRequest: &SiliconFlowListFineTuningJobEventsReq{
				FineTuningJobID: "ft-job-456",
			},
			mockResponder: newMockResponder(http.StatusNotFound, `{"error":{"message":"Job not found"}}`),
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.ListFineTuningJobEvents(req.(*SiliconFlowListFineTuningJobEventsReq))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.Nil(t, rsp)
				checkAPIError(t, err, "Job not found")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.RegisterResponder(tt.method, tt.endpoint, tt.mockResponder)
			rsp, err := tt.apiCall(p, tt.giveRequest)
			tt.checkResponse(t, rsp, err)
		})
	}
}

func TestSiliconflowMetadataAPI(t *testing.T) {
	p := setupSiliconFlowTest(t)

	// Mocks for GetModel
	mockGetModelSuccessResponse := Model{
		Id: "deepseek-coder", Object: "model", Created: 1696931010, OwnedBy: "deepseek",
	}
	getModelSuccessBody, err := json.Marshal(mockGetModelSuccessResponse)
	require.NoError(t, err)

	// Mocks for GetModelList
	mockGetModelListFilteredResponse := SiliconFlowGetModelListRsp{
		Object: "list",
		Data:   []Model{{Id: "deepseek-coder", Object: "model", Created: 1696931010, OwnedBy: "deepseek"}},
	}
	getModelListFilteredBody, err := json.Marshal(mockGetModelListFilteredResponse)
	require.NoError(t, err)

	mockGetModelListUnfilteredResponse := SiliconFlowGetModelListRsp{
		Object: "list",
		Data:   []Model{{Id: "model-1"}, {Id: "model-2"}},
	}
	getModelListUnfilteredBody, err := json.Marshal(mockGetModelListUnfilteredResponse)
	require.NoError(t, err)

	// Mocks for GetInfoList
	mockGetInfoListSuccessResponse := SiliconFlowGetInfoListRsp{
		Code: 200, Message: "Success", Status: true,
		Data: struct {
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
		}{Id: "user-123", Name: "Test User", Email: "test@example.com", Balance: "100.00"},
	}
	getInfoListSuccessBody, err := json.Marshal(mockGetInfoListSuccessResponse)
	require.NoError(t, err)

	tests := []struct {
		name          string
		method        string
		endpoint      string
		giveRequest   any
		mockResponder httpmock.Responder
		apiCall       func(p *SiliconFlow, req any) (any, error)
		checkResponse func(t *testing.T, rsp any, err error)
	}{
		// GetModel tests
		{
			name:          "get_model/success",
			method:        "GET",
			endpoint:      "https://api.siliconflow.cn/v1/models/deepseek-coder",
			giveRequest:   "deepseek-coder",
			mockResponder: newMockResponder(http.StatusOK, string(getModelSuccessBody)),
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.GetModel(req.(string))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.NoError(t, err)
				r, ok := rsp.(*Model)
				require.True(t, ok)
				assert.Equal(t, "deepseek-coder", r.Id)
			},
		},
		{
			name:          "get_model/not_found",
			method:        "GET",
			endpoint:      "https://api.siliconflow.cn/v1/models/non-existent-model",
			giveRequest:   "non-existent-model",
			mockResponder: newMockResponder(http.StatusNotFound, `{"error":{"message":"model not found"}}`),
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.GetModel(req.(string))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.Nil(t, rsp)
				checkAPIError(t, err, "model not found")
			},
		},
		// GetModelList tests
		{
			name:     "get_model_list/with_filter",
			method:   "GET",
			endpoint: "https://api.siliconflow.cn/v1/models",
			giveRequest: &SiliconFlowGetModelListReq{
				Type: "chat", SubType: "open-source",
			},
			mockResponder: func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, "chat", req.URL.Query().Get("type"))
				assert.Equal(t, "open-source", req.URL.Query().Get("sub_type"))
				return newMockResponder(http.StatusOK, string(getModelListFilteredBody))(req)
			},
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.GetModelList(req.(*SiliconFlowGetModelListReq))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.NoError(t, err)
				r, ok := rsp.(*SiliconFlowGetModelListRsp)
				require.True(t, ok)
				require.Len(t, r.Data, 1)
				assert.Equal(t, "deepseek-coder", r.Data[0].Id)
			},
		},
		{
			name:          "get_model_list/no_filter",
			method:        "GET",
			endpoint:      "https://api.siliconflow.cn/v1/models",
			giveRequest:   &SiliconFlowGetModelListReq{},
			mockResponder: newMockResponder(http.StatusOK, string(getModelListUnfilteredBody)),
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.GetModelList(req.(*SiliconFlowGetModelListReq))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.NoError(t, err)
				r, ok := rsp.(*SiliconFlowGetModelListRsp)
				require.True(t, ok)
				assert.Len(t, r.Data, 2)
			},
		},
		// GetInfoList tests
		{
			name:          "get_info_list/success",
			method:        "GET",
			endpoint:      "https://api.siliconflow.cn/v1/user/info",
			giveRequest:   &aiload.ModelChannel{Token: "test-token"},
			mockResponder: newMockResponder(http.StatusOK, string(getInfoListSuccessBody)),
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				p.channel = req.(*aiload.ModelChannel)
				return p.GetInfoList()
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.NoError(t, err)
				r, ok := rsp.(*SiliconFlowGetInfoListRsp)
				require.True(t, ok)
				assert.True(t, r.Status)
				assert.Equal(t, "user-123", r.Data.Id)
			},
		},
		{
			name:          "get_info_list/no_token",
			method:        "GET",
			endpoint:      "https://api.siliconflow.cn/v1/user/info",
			giveRequest:   &aiload.ModelChannel{Token: ""},
			mockResponder: newMockResponder(http.StatusOK, ""), // Not called
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				p.channel = req.(*aiload.ModelChannel)
				return p.GetInfoList()
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.Nil(t, rsp)
				assert.True(t, errors.Is(err, ErrAuthTokenNil))
			},
		},
		{
			name:          "get_info_list/api_unstructured_error",
			method:        "GET",
			endpoint:      "https://api.siliconflow.cn/v1/user/info",
			giveRequest:   &aiload.ModelChannel{Token: "test-token"},
			mockResponder: newMockResponder(http.StatusInternalServerError, "internal server error"),
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				p.channel = req.(*aiload.ModelChannel)
				return p.GetInfoList()
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
				assert.Contains(t, err.Error(), "internal server error")
			},
		},
		{
			name:          "get_model_list/api_error",
			method:        "GET",
			endpoint:      "https://api.siliconflow.cn/v1/models",
			giveRequest:   &SiliconFlowGetModelListReq{},
			mockResponder: newMockResponder(http.StatusInternalServerError, `{"error":{"message":"Internal server error"}}`),
			apiCall: func(p *SiliconFlow, req any) (any, error) {
				return p.GetModelList(req.(*SiliconFlowGetModelListReq))
			},
			checkResponse: func(t *testing.T, rsp any, err error) {
				assert.Error(t, err)
				assert.Nil(t, rsp)
				checkAPIError(t, err, "Internal server error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpmock.RegisterResponder(tt.method, tt.endpoint, tt.mockResponder)
			rsp, err := tt.apiCall(p, tt.giveRequest)
			tt.checkResponse(t, rsp, err)
		})
	}
}

func TestSiliconFlow_RequestHandling(t *testing.T) {
	tests := []struct {
		name        string
		giveChannel *aiload.ModelChannel
		checkFunc   func(t *testing.T, p *SiliconFlow)
	}{
		{
			name:        "GetRequestRequired returns true when token is present",
			giveChannel: &aiload.ModelChannel{Token: "some-token"},
			checkFunc: func(t *testing.T, p *SiliconFlow) {
				assert.True(t, p.GetRequestRequired())
			},
		},
		{
			name:        "GetRequestRequired returns false when token is absent",
			giveChannel: &aiload.ModelChannel{Token: ""},
			checkFunc: func(t *testing.T, p *SiliconFlow) {
				assert.False(t, p.GetRequestRequired())
			},
		},
		{
			name:        "GetRequest sets auth token correctly",
			giveChannel: &aiload.ModelChannel{Token: "my-secret-token"},
			checkFunc: func(t *testing.T, p *SiliconFlow) {
				req := p.GetRequest()
				assert.Equal(t, "my-secret-token", req.Token)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewSiliconFlow(tt.giveChannel)
			tt.checkFunc(t, p)
		})
	}
}

// newMockResponder creates a new httpmock.Responder with a given status code and body.
// It automatically sets the Content-Type header to application/json.
func newMockResponder(statusCode int, body string) httpmock.Responder {
	resp := httpmock.NewStringResponse(statusCode, body)
	resp.Header.Set("Content-Type", "application/json")
	return httpmock.ResponderFromResponse(resp)
}

// checkAPIError checks if the given error is a *SiliconFlowErrorRsp and if its message matches the expected one.
func checkAPIError(t *testing.T, err error, expectedMessage string) {
	t.Helper()
	var siliconFlowError *SiliconFlowErrorRsp
	require.ErrorAs(t, err, &siliconFlowError, "error should be of type SiliconFlowErrorRsp")
	assert.Equal(t, expectedMessage, siliconFlowError.Err.Message, "error message should match")
}